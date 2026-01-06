/*
Copyright 2025 Ross Golder.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/crossplane/crossplane-runtime/pkg/feature"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/rossigee/provider-hostinger/apis"
	"github.com/rossigee/provider-hostinger/internal/controller"
)

func main() {
	var (
		metricsAddr          = os.Getenv("METRICS_BIND_ADDRESS")
		enableLeaderElection = os.Getenv("LEADER_ELECT") == "true"
		webhookCertDir       = os.Getenv("WEBHOOK_TLS_CERT_DIR")
		syncPeriod           = os.Getenv("SYNC_PERIOD")
		pollInterval         = os.Getenv("POLL_INTERVAL")
		maxReconcileRate     = os.Getenv("MAX_RECONCILE_RATE")
	)

	// Set default metrics address if not provided
	if metricsAddr == "" {
		metricsAddr = ":8080"
	}

	// Set default webhook cert dir if not provided (use environment variable)
	if webhookCertDir == "" {
		webhookCertDir = filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs")
	}

	logger, err := logging.NewDefaultLogger()
	if err != nil {
		panic(err)
	}

	ctrl.SetLogger(logger)

	logger.Info("Starting provider",
		"syncPeriod", syncPeriod,
		"pollInterval", pollInterval,
		"maxReconcileRate", maxReconcileRate,
	)

	cfg, err := ctrl.GetConfig()
	if err != nil {
		logger.Error(err, "Unable to get kubeconfig")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: apis.Scheme,
		WebhookServer: webhook.NewServer(webhook.Options{
			CertDir: webhookCertDir,
			Port:    9443,
		}),
		MetricsBindAddress:     metricsAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "provider-hostinger",
		SyncPeriod:             parseDuration(syncPeriod),
		ClientDisableCacheFor:  []interface{}{},
	})
	if err != nil {
		logger.Error(err, "Unable to create manager")
		os.Exit(1)
	}

	// Register APIs
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		logger.Error(err, "Unable to add APIs to scheme")
		os.Exit(1)
	}

	// Setup webhook server with TLS configuration from environment
	mgr.GetWebhookServer().Port = 9443
	if webhookCertDir != "" {
		mgr.GetWebhookServer().CertDir = webhookCertDir
	}

	// Register controllers
	if err := controller.Setup(mgr, logger, ratelimiter.NewTypedDefaultingRateLimiter[interface{}](nil)); err != nil {
		logger.Error(err, "Unable to setup controller")
		os.Exit(1)
	}

	// Setup feature flags
	if err := feature.Initialize(mgr.GetConfig()); err != nil {
		logger.Error(err, "Unable to initialize feature flags")
		os.Exit(1)
	}

	logger.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(err, "Problem running manager")
		os.Exit(1)
	}
}

// parseDuration parses a duration string or returns default
func parseDuration(s string) *time.Duration {
	if s == "" {
		return nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return nil
	}
	return &d
}

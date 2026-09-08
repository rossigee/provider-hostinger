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
	"context"
	"os"
	"path/filepath"
	"runtime"

	xpcontroller "github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/feature"
	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/v2/pkg/statemetrics"
	"github.com/rossigee/provider-hostinger/apis"
	backupv1beta1 "github.com/rossigee/provider-hostinger/apis/backup/v1beta1"
	firewallv1beta1 "github.com/rossigee/provider-hostinger/apis/firewall/v1beta1"
	instancev1beta1 "github.com/rossigee/provider-hostinger/apis/instance/v1beta1"
	sshkeyv1beta1 "github.com/rossigee/provider-hostinger/apis/sshkey/v1beta1"
	"github.com/rossigee/provider-hostinger/internal/controller"
	"github.com/rossigee/provider-hostinger/internal/features"
	"github.com/rossigee/provider-hostinger/internal/tracing"
	"github.com/rossigee/provider-hostinger/internal/version"
	"gopkg.in/alecthomas/kingpin.v2"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	metricserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func main() {
	var (
		app                      = kingpin.New(filepath.Base(os.Args[0]), "Hostinger VPS support for Crossplane.").DefaultEnvars()
		debug                    = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncPeriod               = app.Flag("sync", "Controller manager sync period such as 300ms, 1.5h, or 2h45m").Short('s').Default("1h").Duration()
		leaderElection           = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()
		pollInterval             = app.Flag("poll", "Poll interval controls how often an individual resource should be checked for drift.").Default("1m").Duration()
		maxReconcileRate         = app.Flag("max-reconcile-rate", "The global maximum rate per second at which resources may be checked for drift.").Default("10").Int()
		pollStateMetricInterval  = app.Flag("poll-state-metric", "State metric recording interval").Default("5s").Duration()
		metricsBindAddress       = app.Flag("metrics-bind-address", "The address the metrics endpoint binds to.").Default(":8080").String()
		enableManagementPolicies = app.Flag("enable-management-policies", "Enable support for management policies.").Default("true").Envar("ENABLE_MANAGEMENT_POLICIES").Bool()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-hostinger"))

	shutdownTracing := tracing.Init("provider-hostinger")
	defer shutdownTracing(context.Background())

	// Always set the controller-runtime logger to prevent logging errors
	ctrl.SetLogger(zl)

	log.Info("Provider starting up",
		"provider", "provider-hostinger",
		"version", version.Version,
		"go-version", runtime.Version(),
		"platform", runtime.GOOS+"/"+runtime.GOARCH,
		"sync-period", syncPeriod.String(),
		"poll-interval", pollInterval.String(),
		"max-reconcile-rate", *maxReconcileRate,
		"leader-election", *leaderElection,
		"leader-election-id", "crossplane-leader-election-provider-hostinger",
		"management-policies", *enableManagementPolicies,
		"debug-mode", *debug)

	s := apimachineryruntime.NewScheme()
	kingpin.FatalIfError(scheme.AddToScheme(s), "Cannot add k8s types to scheme")
	kingpin.FatalIfError(apis.AddToScheme(s), "Cannot add Hostinger APIs to scheme")

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:           s,
		LeaderElection:   *leaderElection,
		LeaderElectionID: "crossplane-leader-election-provider-hostinger",
		Metrics: metricserver.Options{
			BindAddress: *metricsBindAddress,
		},
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")

	mrStateMetrics := statemetrics.NewMRStateMetrics()
	metrics.Registry.MustRegister(mrStateMetrics)

	mo := xpcontroller.MetricOptions{
		PollStateMetricInterval: *pollStateMetricInterval,
		MRStateMetrics:          mrStateMetrics,
	}

	featureFlags := &feature.Flags{}
	if *enableManagementPolicies {
		featureFlags.Enable(features.EnableAlphaManagementPolicies)
		log.Info("Alpha feature enabled", "flag", features.EnableAlphaManagementPolicies)
	}

	o := xpcontroller.Options{
		Logger:                  log,
		MaxConcurrentReconciles: *maxReconcileRate,
		PollInterval:            *pollInterval,
		GlobalRateLimiter:       ratelimiter.NewGlobal(*maxReconcileRate),
		Features:                featureFlags,
		MetricOptions:           &mo,
	}

	kingpin.FatalIfError(controller.Setup(mgr, o), "Cannot setup Hostinger controllers")

	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &instancev1beta1.InstanceList{}, o.MetricOptions.PollStateMetricInterval)), "Cannot register state metrics for Instance")
	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &sshkeyv1beta1.SSHKeyList{}, o.MetricOptions.PollStateMetricInterval)), "Cannot register state metrics for SSHKey")
	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &backupv1beta1.BackupList{}, o.MetricOptions.PollStateMetricInterval)), "Cannot register state metrics for Backup")
	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), o.Logger, o.MetricOptions.MRStateMetrics, &firewallv1beta1.FirewallRuleList{}, o.MetricOptions.PollStateMetricInterval)), "Cannot register state metrics for FirewallRule")

	kingpin.FatalIfError(mgr.AddHealthzCheck("healthz", healthz.Ping), "Cannot add health check")
	kingpin.FatalIfError(mgr.AddReadyzCheck("readyz", healthz.Ping), "Cannot add ready check")

	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

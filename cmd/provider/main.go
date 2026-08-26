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

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/crossplane/crossplane-runtime/v2/pkg/statemetrics"
	"github.com/rossigee/provider-hostinger/apis"
	backupv1beta1 "github.com/rossigee/provider-hostinger/apis/backup/v1beta1"
	firewallv1beta1 "github.com/rossigee/provider-hostinger/apis/firewall/v1beta1"
	instancev1beta1 "github.com/rossigee/provider-hostinger/apis/instance/v1beta1"
	sshkeyv1beta1 "github.com/rossigee/provider-hostinger/apis/sshkey/v1beta1"
	"github.com/rossigee/provider-hostinger/internal/controller"
	"github.com/rossigee/provider-hostinger/internal/tracing"
	"github.com/rossigee/provider-hostinger/internal/version"
	"gopkg.in/alecthomas/kingpin.v2"
	apimachineryruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	metricserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func main() {
	var (
		app                     = kingpin.New(filepath.Base(os.Args[0]), "Hostinger VPS support for Crossplane.").DefaultEnvars()
		debug                   = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncPeriod              = app.Flag("sync", "Controller manager sync period such as 300ms, 1.5h, or 2h45m").Short('s').Default("1h").Duration()
		leaderElection          = app.Flag("leader-election", "Use leader election for the controller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()
		pollStateMetricInterval = app.Flag("poll-state-metric", "State metric recording interval").Default("5s").Duration()
		metricsBindAddress      = app.Flag("metrics-bind-address", "The address the metrics endpoint binds to.").Default(":8080").String()
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
		"leader-election", *leaderElection,
		"leader-election-id", "crossplane-leader-election-provider-hostinger",
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

	rl := workqueue.DefaultTypedControllerRateLimiter[any]()
	kingpin.FatalIfError(controller.Setup(mgr, log, rl), "Cannot setup Hostinger controllers")

	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), log, mrStateMetrics, &instancev1beta1.InstanceList{}, *pollStateMetricInterval)), "Cannot register state metrics for Instance")
	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), log, mrStateMetrics, &sshkeyv1beta1.SSHKeyList{}, *pollStateMetricInterval)), "Cannot register state metrics for SSHKey")
	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), log, mrStateMetrics, &backupv1beta1.BackupList{}, *pollStateMetricInterval)), "Cannot register state metrics for Backup")
	kingpin.FatalIfError(mgr.Add(statemetrics.NewMRStateRecorder(mgr.GetClient(), log, mrStateMetrics, &firewallv1beta1.FirewallRuleList{}, *pollStateMetricInterval)), "Cannot register state metrics for FirewallRule")

	kingpin.FatalIfError(mgr.AddHealthzCheck("healthz", healthz.Ping), "Cannot add health check")
	kingpin.FatalIfError(mgr.AddReadyzCheck("readyz", healthz.Ping), "Cannot add ready check")

	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

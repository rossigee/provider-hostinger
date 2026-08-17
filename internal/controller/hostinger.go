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

package controller

import (
	"context"
	"os"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	"github.com/rossigee/provider-hostinger/internal/controller/instance"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Setup registers all Hostinger provider controllers with the manager
func Setup(mgr ctrl.Manager, l logging.Logger, wl workqueue.TypedRateLimiter[any]) error {
	// Providers MUST manage their RBAC during install/upgrade (and optionally removal).
	// This eliminates the need for static per-revision bootstrap hacks in gitops
	// (the ownerReferences pinning workarounds that caused double-controller conflicts
	// and live patches).
	if err := setupRBAC(mgr.GetClient(), l); err != nil {
		l.Info("RBAC setup warning (may be transient)", "error", err)
	}
	return instance.Setup(mgr, l, wl)
}

// setupRBAC ensures the provider's system, edit, and view ClusterRoles have the
// full rules for all its resources. The provider manages this itself.
func setupRBAC(c client.Client, l logging.Logger) error {
	ctx := context.Background()

	// Full rules for hostinger provider's resources.
	// These are the rules previously maintained in static bootstrap files.
	rules := []rbacv1.PolicyRule{
		{
			APIGroups: []string{"hostinger.crossplane.io"},
			Resources: []string{"providerconfigs", "providerconfigs/status"},
			Verbs:     []string{"get", "list", "watch", "update", "patch", "create"},
		},
		{
			APIGroups: []string{"sshkey.m.hostinger.crossplane.io"},
			Resources: []string{"sshkeys", "sshkeys/status"},
			Verbs:     []string{"get", "list", "watch", "update", "patch", "create"},
		},
		{
			APIGroups: []string{"instance.m.hostinger.crossplane.io"},
			Resources: []string{"instances", "instances/status"},
			Verbs:     []string{"get", "list", "watch", "update", "patch", "create"},
		},
		{
			APIGroups: []string{"firewall.m.hostinger.crossplane.io"},
			Resources: []string{"firewallrules", "firewallrules/status"},
			Verbs:     []string{"get", "list", "watch", "update", "patch", "create"},
		},
		{
			APIGroups: []string{"backup.m.hostinger.crossplane.io"},
			Resources: []string{"backups", "backups/status"},
			Verbs:     []string{"get", "list", "watch", "update", "patch", "create"},
		},
		{
			APIGroups: []string{"hostinger.crossplane.io", "sshkey.m.hostinger.crossplane.io", "instance.m.hostinger.crossplane.io", "firewall.m.hostinger.crossplane.io", "backup.m.hostinger.crossplane.io"},
			Resources: []string{"*/finalizers"},
			Verbs:     []string{"update"},
		},
		{
			APIGroups: []string{"", "coordination.k8s.io"},
			Resources: []string{"secrets", "configmaps", "events", "leases"},
			Verbs:     []string{"*"},
		},
	}

	// Apply the system role (bound to provider SA).
	system := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crossplane:provider:provider-hostinger:system", // stable name; revision specific can aggregate or provider can update per rev if needed
			Labels: map[string]string{
				"rbac.crossplane.io/system": "provider-hostinger",
			},
		},
		Rules: rules,
	}
	if err := c.Create(ctx, system); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	if err := c.Update(ctx, system); err != nil {
		l.Info("system role update", "err", err)
	}

	// Ensure binding so the SA gets the permissions from the system role we manage.
	binding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crossplane:provider:provider-hostinger:system",
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "crossplane:provider:provider-hostinger:system",
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      os.Getenv("REVISION_NAME"),
			Namespace: "crossplane-system",
		}},
	}
	if err := c.Create(ctx, binding); err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	if err := c.Update(ctx, binding); err != nil {
		l.Info("system binding update", "err", err)
	}

	// Aggregate edit
	edit := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crossplane:provider:provider-hostinger:aggregate-to-edit",
			Labels: map[string]string{
				"rbac.crossplane.io/aggregate-to-edit":       "true",
				"rbac.crossplane.io/aggregate-to-admin":      "true",
				"rbac.crossplane.io/aggregate-to-crossplane": "true",
				"rbac.crossplane.io/system":                  "provider-hostinger",
			},
		},
		Rules: withVerbs(rules, []string{"*"}),
	}
	if err := c.Create(ctx, edit); err != nil && !errors.IsAlreadyExists(err) {
		l.Info("aggregate-to-edit create warning (non-fatal)", "err", err)
	}
	_ = c.Update(ctx, edit)

	// Aggregate view
	view := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crossplane:provider:provider-hostinger:aggregate-to-view",
			Labels: map[string]string{
				"rbac.crossplane.io/aggregate-to-view": "true",
				"rbac.crossplane.io/system":            "provider-hostinger",
			},
		},
		Rules: withVerbs(rules, []string{"get", "list", "watch"}),
	}
	if err := c.Create(ctx, view); err != nil && !errors.IsAlreadyExists(err) {
		l.Info("aggregate-to-view create warning (non-fatal)", "err", err)
	}
	_ = c.Update(ctx, view)

	l.Info("provider self-managed RBAC roles ensured")
	return nil
}

func withVerbs(r []rbacv1.PolicyRule, verbs []string) []rbacv1.PolicyRule {
	out := make([]rbacv1.PolicyRule, len(r))
	for i := range r {
		out[i] = r[i]
		out[i].Verbs = verbs
	}
	return out
}

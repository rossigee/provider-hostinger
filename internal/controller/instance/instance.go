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

package instance

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	v1beta1 "github.com/rossigee/provider-hostinger/apis/instance/v1beta1"
	providerv1beta1 "github.com/rossigee/provider-hostinger/apis/v1beta1"
	"github.com/rossigee/provider-hostinger/internal/clients"
	instanceclient "github.com/rossigee/provider-hostinger/internal/clients/instance"
)

const (
	errNotInstance = "managed resource is not a Instance custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errNewClient    = "cannot create new Hostinger client"
)

// Setup adds a controller that reconciles Instance managed resources.
func Setup(mgr ctrl.Manager, l logging.Logger, wl workqueue.TypedRateLimiter[any]) error {
	name := managed.ControllerName(v1beta1.InstanceGroupKind)

	o := controller.Options{
		RateLimiter: nil, // Use default rate limiter
		MaxConcurrentReconciles: 5,
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1beta1.InstanceGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:        mgr.GetClient(),
			newClientFn: clients.NewClientFactory,
		}),
		managed.WithLogger(l.WithValues("controller", name)),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithPollInterval(5*time.Minute),
		managed.WithInitializers(),
	)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&v1beta1.Instance{}).
		Complete(r)
}

// A connector is expected to produce typed ExternalClient for the managed
// resource it is supposed to manage.
type connector struct {
	kube        client.Client
	newClientFn func(client.Client, clients.HTTPClientConfig) *clients.ClientFactory
}

// Connect typically produces an ExternalClient by dialing for the provider
// configured in ProviderConfig and using this Provider as an authentication
// mechanism.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1beta1.Instance)
	if !ok {
		return nil, errors.New(errNotInstance)
	}

	// Get the ProviderConfig referenced by this Instance
	pc := &providerv1beta1.ProviderConfig{}
	if err := c.kube.Get(ctx, client.ObjectKey{Namespace: cr.GetNamespace(), Name: cr.Spec.ProviderConfigReference.Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	// Create the Hostinger client factory with default config
	clientFactory := c.newClientFn(c.kube, clients.DefaultHTTPClientConfig())

	// Create the Hostinger client
	hc, err := clientFactory.CreateHostingerClient(ctx, pc)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	// Create the instance client
	instanceClient := instanceclient.NewInstanceClient(hc)

	return &external{client: instanceClient}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	client instanceclient.Client
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1beta1.Instance)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotInstance)
	}

	// Get the external resource ID from the resource annotation
	externalName := meta.GetExternalName(cr)
	if externalName == "" {
		// Resource hasn't been created yet
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	// Fetch the current state of the instance
	instance, err := e.client.Get(ctx, externalName)
	if err != nil {
		if clients.IsNotFound(err) {
			return managed.ExternalObservation{ResourceExists: false}, nil
		}
		return managed.ExternalObservation{}, err
	}

	// Update the observation status
	cr.Status.AtProvider = *e.client.GetObservation(instance)

	// Check if the instance is up-to-date
	upToDate := e.client.UpToDate(instance, &cr.Spec.ForProvider)

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: upToDate,
	}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1beta1.Instance)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotInstance)
	}

	// Create the instance
	instance, err := e.client.Create(ctx, &cr.Spec.ForProvider)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, "failed to create instance")
	}

	// Set the external name annotation (Crossplane uses this as the resource ID)
	meta.SetExternalName(cr, instance.ID)

	// Perform late initialization
	if e.client.LateInitialize(instance, &cr.Spec.ForProvider) {
		// Update the spec if late initialization changed anything
		cr.Spec.ForProvider = *&cr.Spec.ForProvider
	}

	return managed.ExternalCreation{
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1beta1.Instance)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotInstance)
	}

	externalName := meta.GetExternalName(cr)
	if externalName == "" {
		return managed.ExternalUpdate{}, errors.New("external name not set")
	}

	// Update the instance
	if err := e.client.Update(ctx, externalName, &cr.Spec.ForProvider); err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, "failed to update instance")
	}

	return managed.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1beta1.Instance)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotInstance)
	}

	externalName := meta.GetExternalName(cr)
	if externalName == "" {
		return managed.ExternalDelete{}, nil // Already deleted or never created
	}

	// Delete the instance
	if err := e.client.Delete(ctx, externalName); err != nil {
		return managed.ExternalDelete{}, errors.Wrap(err, "failed to delete instance")
	}

	return managed.ExternalDelete{}, nil
}

// Disconnect closes the connection to the external service.
func (e *external) Disconnect(ctx context.Context) error {
	// No cleanup needed for Hostinger client
	return nil
}

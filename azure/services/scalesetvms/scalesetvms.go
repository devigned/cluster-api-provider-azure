/*
Copyright 2020 The Kubernetes Authors.

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

package scalesetvms

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1alpha3"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

type (
	// ScaleSetVMScope defines the scope interface for a scale sets service.
	ScaleSetVMScope interface {
		logr.Logger
		azure.ClusterDescriber
		InstanceID() string
		ScaleSetName() string
		SetVMSSVM(vmssvm *infrav1exp.VMSSVM)
		GetLongRunningOperationState() *infrav1.Future
	}

	// Service provides operations on azure resources
	Service struct {
		Client client
		Scope  ScaleSetVMScope
	}
)

// NewService creates a new service.
func NewService(scope ScaleSetVMScope) *Service {
	return &Service{
		Client: newClient(scope),
		Scope:  scope,
	}
}

// Reconcile idempotently gets, creates, and updates a scale set.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, span := tele.Tracer().Start(ctx, "scalesets.Service.Reconcile")
	defer span.End()

	var (
		resourceGroup = s.Scope.ResourceGroup()
		vmssName      = s.Scope.ScaleSetName()
		instanceID    = s.Scope.InstanceID()
	)

	// fetch the latest data about the instance -- model mutations are handled by the AzureMachinePoolReconciler
	instance, err := s.Client.Get(ctx, resourceGroup, vmssName, instanceID)
	if err != nil {
		if azure.ResourceNotFound(err) {
			return azure.WithTransientError(errors.New("instance does not exist yet"), 30*time.Second)
		}
		return errors.Wrap(err, "failed getting instance")
	}

	s.Scope.SetVMSSVM(converters.SDKToVMSSVM(instance))
	return nil
}

// Delete deletes a scaleset instance asynchronously returning a future which encapsulates the long running operation.
func (s *Service) Delete(ctx context.Context) (*infrav1.Future, error) {
	ctx, span := tele.Tracer().Start(ctx, "scalesets.Service.Delete")
	defer span.End()

	var (
		resourceGroup = s.Scope.ResourceGroup()
		vmssName      = s.Scope.ScaleSetName()
		instanceID    = s.Scope.InstanceID()
	)

	log := s.Scope.WithValues("resourceGroup", resourceGroup, "scaleset", vmssName, "instanceID", instanceID)
	log.Info("entering delete")
	future := s.Scope.GetLongRunningOperationState()
	if future != nil {
		if future.Type != DeleteFuture {
			return future, azure.WithTransientError(errors.New("attempting to delete, non-delete operation in progress"), 30*time.Second)
		}

		log.Info("checking if the instance is done deleting")
		if _, err := s.Client.GetResultIfDone(ctx, future); err != nil {
			// fetch instance to update status
			if instance, err := s.Client.Get(ctx, resourceGroup, vmssName, instanceID); err != nil {
				s.Scope.SetVMSSVM(converters.SDKToVMSSVM(instance))
			}
			return future, errors.Wrap(err, "failed to get result of long running operation")
		}

		// there was no error in fetching the result, the future has been completed
		log.Info("successfully deleted the instance")
		return nil, nil
	}

	// since the future was nil, there is no ongoing activity; start deleting the instance
	future, err := s.Client.DeleteAsync(ctx, resourceGroup, vmssName, instanceID)
	if err != nil {
		if azure.ResourceNotFound(err) {
			// already deleted
			return nil, nil
		}
		return nil, errors.Wrapf(err, "failed to delete instance %s/%s", vmssName, instanceID)
	}

	// fetch instance to update status
	if instance, err := s.Client.Get(ctx, resourceGroup, vmssName, instanceID); err != nil {
		s.Scope.SetVMSSVM(converters.SDKToVMSSVM(instance))
	}

	log.V(2).Info("successfully started deleting the instance")
	return future, nil
}
// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Package controller is used to provide the core functionalities of machine-controller-manager
package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"

	"k8s.io/klog/v2"

	"github.com/gardener/machine-controller-manager/pkg/apis/machine"
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machineutils"
)

func (c *controller) machineToMachineClassAdd(obj interface{}) {
	machine, ok := obj.(*v1alpha1.Machine)
	if machine == nil || !ok {
		klog.Warningf("Couldn't get machine from object: %+v", obj)
		return
	}
	if machine.Spec.Class.Kind == machineutils.MachineClassKind {
		c.machineClassQueue.Add(machine.Spec.Class.Name)
	}
}

func (c *controller) machineToMachineClassDelete(obj interface{}) {
	c.machineToMachineClassAdd(obj)
}

func (c *controller) machineClassAdd(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	c.machineClassQueue.Add(key)
}

func (c *controller) machineClassUpdate(oldObj, newObj interface{}) {
	old, ok := oldObj.(*v1alpha1.MachineClass)
	if old == nil || !ok {
		return
	}
	new, ok := newObj.(*v1alpha1.MachineClass)
	if new == nil || !ok {
		return
	}

	c.machineClassAdd(newObj)
}

func (c *controller) machineClassDelete(obj interface{}) {
	c.machineClassAdd(obj)
}

// reconcileClusterMachineClassKey reconciles an machineClass due to controller resync
// or an event on the machineClass.
func (c *controller) reconcileClusterMachineClassKey(key string) error {
	ctx := context.Background()
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	class, err := c.machineClassLister.MachineClasses(c.namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.Infof("%s %q: Not doing work because it has been deleted", machineutils.MachineClassKind, key)
		return nil
	}
	if err != nil {
		klog.Infof("%s %q: Unable to retrieve object from store: %v", machineutils.MachineClassKind, key, err)
		return err
	}

	err = c.reconcileClusterMachineClass(ctx, class)
	if err != nil {
		// Re-enqueue after a ShortRetry window
		c.enqueueMachineClassAfter(class, time.Duration(machineutils.ShortRetry))
	} else {
		// Re-enqueue periodically to avoid missing of events
		// TODO: Get ride of this logic
		c.enqueueMachineClassAfter(class, time.Duration(machineutils.LongRetry))
	}

	return nil
}

func (c *controller) reconcileClusterMachineClass(ctx context.Context, class *v1alpha1.MachineClass) error {
	klog.V(4).Info("Start Reconciling machineclass: ", class.Name)
	defer klog.V(4).Info("Stop Reconciling machineclass: ", class.Name)

	// Validate internal to external scheme conversion
	internalClass := &machine.MachineClass{}
	err := c.internalExternalScheme.Convert(class, internalClass, nil)
	if err != nil {
		return err
	}

	// Fetch all machines referring the machineClass
	machines, err := c.findMachinesForClass(machineutils.MachineClassKind, class.Name)
	if err != nil {
		return err
	}

	if class.DeletionTimestamp == nil && len(machines) > 0 {
		// If deletionTimestamp is not set and at least one machine is referring this machineClass

		if finalizers := sets.NewString(class.Finalizers...); !finalizers.Has(MCMFinalizerName) {
			// Add machineClassFinalizer as if doesn't exist
			err = c.addMCMFinalizerToMachineClass(ctx, class)
			if err != nil {
				return err
			}

			// Enqueue all machines once finalizer is added to machineClass
			// This is to allow processing of such machines
			for _, machine := range machines {
				c.enqueueMachine(machine, "finalizer placed on machineClass")
			}
		}

		return nil
	}

	if len(machines) > 0 {
		// Machines are still referring the machine class, please wait before deletion
		klog.V(3).Infof("Cannot remove finalizer on %s because still (%d) machines are referencing it", class.Name, len(machines))

		for _, machine := range machines {
			c.addMachine(machine)
		}

		return fmt.Errorf("Retry as machine objects are still referring the machineclass")
	}

	if finalizers := sets.NewString(class.Finalizers...); finalizers.Has(MCMFinalizerName) {
		// Delete finalizer if exists on machineClass
		return c.deleteMCMFinalizerFromMachineClass(ctx, class)
	}

	return nil
}

/*
	SECTION
	Manipulate Finalizers
*/

func (c *controller) addMCMFinalizerToMachineClass(ctx context.Context, class *v1alpha1.MachineClass) error {
	finalizers := sets.NewString(class.Finalizers...)
	finalizers.Insert(MCMFinalizerName)
	return c.updateMachineClassFinalizers(ctx, class, finalizers.List(), true)
}

func (c *controller) deleteMCMFinalizerFromMachineClass(ctx context.Context, class *v1alpha1.MachineClass) error {
	finalizers := sets.NewString(class.Finalizers...)
	finalizers.Delete(MCMFinalizerName)
	return c.updateMachineClassFinalizers(ctx, class, finalizers.List(), false)
}

func (c *controller) updateMachineClassFinalizers(ctx context.Context, class *v1alpha1.MachineClass, finalizers []string, addFinalizers bool) error {
	// Get the latest version of the class so that we can avoid conflicts
	class, err := c.controlMachineClient.MachineClasses(class.Namespace).Get(ctx, class.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	clone := class.DeepCopy()
	clone.Finalizers = finalizers
	_, err = c.controlMachineClient.MachineClasses(class.Namespace).Update(ctx, clone, metav1.UpdateOptions{})
	if err != nil {
		klog.Warning("Updating machineClass failed, retrying. ", class.Name, err)
		return err
	}
	if addFinalizers {
		klog.V(3).Infof("Successfully added finalizer on the machineclass %q", class.Name)
	} else {
		klog.V(3).Infof("Successfully removed finalizer on the machineclass %q", class.Name)
	}
	return err
}

func (c *controller) enqueueMachineClassAfter(obj interface{}, after time.Duration) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	c.machineClassQueue.AddAfter(key, after)
}

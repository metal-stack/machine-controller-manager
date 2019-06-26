/*
Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved.

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

// Package driver contains the cloud provider specific implementations to manage machines
package driver

import (
	"fmt"
	"strings"

	v1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	corev1 "k8s.io/api/core/v1"

	"github.com/golang/glog"
	metalgo "github.com/metal-pod/metal-go"
)

// MetalDriver is the driver struct for holding Metal machine information
type MetalDriver struct {
	MetalMachineClass *v1alpha1.MetalMachineClass
	CloudConfig       *corev1.Secret
	UserData          string
	MachineID         string
	MachineName       string
}

// NewMetalDriver returns an empty MetalDriver object
func NewMetalDriver(create func() (string, error), delete func() error, existing func() (string, error)) Driver {
	return &MetalDriver{}
}

// Create method is used to create a Metal machine
func (d *MetalDriver) Create() (string, string, error) {

	svc, err := d.createSVC()
	if err != nil {
		return "", "", err
	}
	createRequest := &metalgo.MachineCreateRequest{
		Description:   d.MachineName + " created by Gardener.",
		Name:          d.MachineName,
		Hostname:      d.MachineName,
		UserData:      d.UserData,
		Size:          d.MetalMachineClass.Spec.Size,
		Project:       d.MetalMachineClass.Spec.Project,
		Tenant:        d.MetalMachineClass.Spec.Tenant,
		Partition:     d.MetalMachineClass.Spec.Partition,
		Image:         d.MetalMachineClass.Spec.Image,
		Tags:          d.MetalMachineClass.Spec.Tags,
		SSHPublicKeys: d.MetalMachineClass.Spec.SSHKeys,
	}

	mcr, err := svc.MachineCreate(createRequest)
	if err != nil {
		glog.Errorf("Could not create machine: %v", err)
		return "", "", err
	}
	return d.encodeMachineID(*mcr.Machine.Partition.ID, *mcr.Machine.ID), *mcr.Machine.Allocation.Name, nil
}

// Delete method is used to delete a Machine machine
func (d *MetalDriver) Delete() error {

	svc, err := d.createSVC()
	if err != nil {
		return err
	}
	machineID := d.decodeMachineID(d.MachineID)
	_, err = svc.MachineDelete(machineID)
	if err != nil {
		glog.Errorf("Could not terminate machine %s: %v", d.MachineID, err)
		return err
	}
	return nil
}

// GetExisting method is used to get machineID for existing Metal machine
func (d *MetalDriver) GetExisting() (string, error) {
	return d.MachineID, nil
}

// GetVMs returns a machine matching the machineID
// If machineID is an empty string then it returns all matching instances
func (d *MetalDriver) GetVMs(machineID string) (VMs, error) {
	listOfVMs := make(map[string]string)

	clusterName := ""
	nodeRole := ""

	for _, key := range d.MetalMachineClass.Spec.Tags {
		if strings.Contains(key, "kubernetes.io/cluster/") {
			clusterName = key
		} else if strings.Contains(key, "kubernetes.io/role/") {
			nodeRole = key
		}
	}

	if clusterName == "" || nodeRole == "" {
		return listOfVMs, nil
	}

	svc, err := d.createSVC()
	if err != nil {
		return nil, err
	}
	if machineID == "" {
		listRequest := &metalgo.MachineListRequest{
			Project:   &d.MetalMachineClass.Spec.Project,
			Partition: &d.MetalMachineClass.Spec.Partition,
		}
		mlr, err := svc.MachineList(listRequest)
		if err != nil {
			glog.Errorf("Could not list machines for project %s in partition:%s: %v", d.MetalMachineClass.Spec.Project, d.MetalMachineClass.Spec.Partition, err)
			return nil, err
		}
		for _, m := range mlr.Machines {
			matchedCluster := false
			matchedRole := false
			for _, tag := range m.Tags {
				switch tag {
				case clusterName:
					matchedCluster = true
				case nodeRole:
					matchedRole = true
				}
			}
			if matchedCluster && matchedRole {
				listOfVMs[*m.ID] = *m.Allocation.Hostname
			}
		}
	} else {
		machineID = d.decodeMachineID(machineID)
		mgr, err := svc.MachineGet(machineID)
		if err != nil {
			glog.Errorf("Could not get machine %s: %v", machineID, err)
			return nil, err
		}
		listOfVMs[machineID] = *mgr.Machine.Allocation.Hostname
	}
	return listOfVMs, nil
}

// Helper function to create SVC
func (d *MetalDriver) createSVC() (*metalgo.Driver, error) {
	token := strings.TrimSpace(string(d.CloudConfig.Data[v1alpha1.MetalAPIKey]))
	hmac := strings.TrimSpace(string(d.CloudConfig.Data[v1alpha1.MetalAPIHMac]))

	u, ok := d.CloudConfig.Data[v1alpha1.MetalAPIURL]
	if !ok {
		return nil, fmt.Errorf("missing %s in secret", v1alpha1.MetalAPIURL)
	}
	url := strings.TrimSpace(string(u))

	return metalgo.NewDriver(url, token, hmac)
}

func (d *MetalDriver) encodeMachineID(partition, machineID string) string {
	return fmt.Sprintf("metal:///%s/%s", partition, machineID)
}

func (d *MetalDriver) decodeMachineID(id string) string {
	splitProviderID := strings.Split(id, "/")
	return splitProviderID[len(splitProviderID)-1]
}

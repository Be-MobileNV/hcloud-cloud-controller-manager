/*
Copyright 2018 Hetzner Cloud GmbH.

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

package hcloud

import (
	"errors"
	"strconv"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type instances struct {
	client *HetznerClient
}

func newInstances(client *HetznerClient) *instances {
	return &instances{client}
}

func (i *instances) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	server, err := i.client.GetServerByProviderID(providerID)
	if err != nil {
		return nil, err
	}
	return nodeAddresses(server.Name, server.Ipv4), nil
	// return nil, cloudprovider.InstanceNotFound
}

func (i *instances) NodeAddresses(nodeName types.NodeName) ([]v1.NodeAddress, error) {
	server, err := i.client.GetServerByName(string(nodeName))
	if err != nil {
		return nil, err
	}
	return nodeAddresses(server.Name, server.Ipv4), nil
}

func (i *instances) ExternalID(nodeName types.NodeName) (string, error) {
	return i.InstanceID(nodeName)
}

func (i *instances) InstanceID(nodeName types.NodeName) (string, error) {
	server, err := i.client.GetServerByName(string(nodeName))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(server.ID), nil
}

func (i *instances) InstanceType(nodeName types.NodeName) (string, error) {
	server, err := i.client.GetServerByName(string(nodeName))
	if err != nil {
		return "", err
	}
	return server.InstanceType, nil
}

func (i *instances) InstanceTypeByProviderID(providerID string) (string, error) {
	server, err := i.client.GetServerByProviderID(providerID)
	if err != nil {
		return "", err
	}
	return server.InstanceType, nil
}

func (i *instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return errors.New("not implemented")
}

func (i *instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i instances) InstanceExistsByProviderID(providerID string) (exists bool, err error) {
	var server *Server
	server, err = i.client.GetServerByProviderID(providerID)
	if err != nil {
		return
	}
	exists = server != nil
	if !exists {
		// Node does not exists, but should not be removed, hybrid cloud, see https://github.com/kubernetes/kubernetes/pull/73171
		exists = true
	}
	return
}

func nodeAddresses(name, ipv4 string) []v1.NodeAddress {
	var addresses []v1.NodeAddress
	addresses = append(
		addresses,
		v1.NodeAddress{Type: v1.NodeHostName, Address: name},
		v1.NodeAddress{Type: v1.NodeExternalIP, Address: ipv4},
	)
	return addresses
}

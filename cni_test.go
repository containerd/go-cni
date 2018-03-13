/*
   Copyright The containerd Authors.

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

package cni

import (
	"net"
	"testing"

	cnilibrary "github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/020"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestLibCNIType020 tests the cni version 2.0 plugin
// config and parses the result into structured data
func TestLibCNIType020(t *testing.T) {
	// Get the default CNI config
	l := defaultCNIConfig()

	// Create a fake cni config directory and file
	cniDir, confDir := makeFakeCNIConfig(t)
	defer tearDownCNIConfig(t, cniDir)
	l.pluginConfDir = confDir
	// Set the minimum network count as 2 for this test
	l.networkCount = 2
	err := l.populateNetworkConfig()
	assert.NoError(t, err)

	mockCNI := &MockCNI{}
	l.networks[0].cni = mockCNI
	expectedRT := &cnilibrary.RuntimeConf{ContainerID: "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth0",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[0].config, expectedRT).Return(&types020.Result{
		IP4: &types020.IPConfig{
			IP: net.IPNet{
				IP: []byte{10, 0, 0, 1},
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[0].config, expectedRT).Return(nil)

	l.networks[1].cni = mockCNI
	expectedRT = &cnilibrary.RuntimeConf{ContainerID: "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth1",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[1].config, expectedRT).Return(&types020.Result{
		IP4: &types020.IPConfig{
			IP: net.IPNet{
				IP: []byte{10, 0, 0, 2},
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[1].config, expectedRT).Return(nil)

	r, err := l.Setup("container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
	assert.Contains(t, r.Interfaces, "eth0")
	assert.NotNil(t, r.Interfaces["eth0"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth0"].IPConfigs[0].IP.String(), "10.0.0.1")
	assert.Equal(t, r.Interfaces["eth0"].IPConfigs[1].IP.String(), "10.0.0.2")
	err = l.Remove("container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
}

// TestLibCNITypeCurrent tests the cni version 3.x plugin
// config and parses the result into structured data
func TestLibCNITypeCurrent(t *testing.T) {
	// Get the default CNI config
	l := defaultCNIConfig()
	// Create a fake cni config directory and file
	cniDir, confDir := makeFakeCNIConfig(t)
	defer tearDownCNIConfig(t, cniDir)
	l.pluginConfDir = confDir
	// Set the minimum network count as 2 for this test
	l.networkCount = 2
	err := l.populateNetworkConfig()
	assert.NoError(t, err)

	mockCNI := &MockCNI{}
	l.networks[0].cni = mockCNI
	l.networks[1].cni = mockCNI
	ipv4, err := types.ParseCIDR("10.0.0.1/24")
	expectedRT := &cnilibrary.RuntimeConf{
		ContainerID:    "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth0",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[0].config, expectedRT).Return(&current.Result{
		CNIVersion: "0.3.1",
		Interfaces: []*current.Interface{
			{
				Name: "eth0",
			},
		},
		IPs: []*current.IPConfig{
			{
				Version:   "4",
				Interface: current.Int(0),
				Address:   *ipv4,
				Gateway:   net.ParseIP("10.0.0.255"),
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[0].config, expectedRT).Return(nil)

	ipv4, err = types.ParseCIDR("10.0.0.2/24")
	l.networks[1].cni = mockCNI
	expectedRT = &cnilibrary.RuntimeConf{
		ContainerID:    "container-id1",
		NetNS:          "/proc/12345/ns/net",
		IfName:         "eth1",
		Args:           [][2]string(nil),
		CapabilityArgs: map[string]interface{}{},
	}
	mockCNI.On("AddNetworkList", l.networks[1].config, expectedRT).Return(&current.Result{
		CNIVersion: "0.3.1",
		Interfaces: []*current.Interface{
			{
				Name: "eth1",
			},
		},
		IPs: []*current.IPConfig{
			{
				Version:   "4",
				Interface: current.Int(0),
				Address:   *ipv4,
				Gateway:   net.ParseIP("10.0.0.2"),
			},
		},
	}, nil)
	mockCNI.On("DelNetworkList", l.networks[1].config, expectedRT).Return(nil)

	r, err := l.Setup("container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
	assert.Contains(t, r.Interfaces, "eth0")
	assert.NotNil(t, r.Interfaces["eth0"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth0"].IPConfigs[0].IP.String(), "10.0.0.1")
	assert.Contains(t, r.Interfaces, "eth1")
	assert.NotNil(t, r.Interfaces["eth1"].IPConfigs)
	assert.Equal(t, r.Interfaces["eth1"].IPConfigs[0].IP.String(), "10.0.0.2")
	err = l.Remove("container-id1", "/proc/12345/ns/net")
	assert.NoError(t, err)
}

type MockCNI struct {
	mock.Mock
}

func (m *MockCNI) AddNetwork(net *cnilibrary.NetworkConfig, rt *cnilibrary.RuntimeConf) (types.Result, error) {
	args := m.Called(net, rt)
	return args.Get(0).(types.Result), args.Error(1)
}

func (m *MockCNI) DelNetwork(net *cnilibrary.NetworkConfig, rt *cnilibrary.RuntimeConf) error {
	args := m.Called(net, rt)
	return args.Error(0)
}

func (m *MockCNI) DelNetworkList(net *cnilibrary.NetworkConfigList, rt *cnilibrary.RuntimeConf) error {
	args := m.Called(net, rt)
	return args.Error(0)
}

func (m *MockCNI) AddNetworkList(net *cnilibrary.NetworkConfigList, rt *cnilibrary.RuntimeConf) (types.Result, error) {
	args := m.Called(net, rt)
	return args.Get(0).(types.Result), args.Error(1)
}

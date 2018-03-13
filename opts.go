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
	cnilibrary "github.com/containernetworking/cni/libcni"
)

type ConfigOptions func(c *libcni) error

// WithInterfacePrefix sets the prefix for network interfaces
// e.g. eth or wlan
func WithInterfacePrefix(prefix string) ConfigOptions {
	return func(c *libcni) error {
		c.prefix = prefix
		return nil
	}
}

func WithPluginDir(dirs []string) ConfigOptions {
	return func(c *libcni) error {
		c.pluginDirs = dirs
		c.cniConfig = &cnilibrary.CNIConfig{Path: dirs}
		return nil
	}
}

func WithPluginConfDir(dir string) ConfigOptions {
	return func(c *libcni) error {
		c.pluginConfDir = dir
		return nil
	}
}

func WithLoNetwork() ConfigOptions {
	return func(c *libcni) error {
		loConfig, _ := cnilibrary.ConfListFromBytes([]byte(`{
"cniVersion": "0.3.1",
"name": "cni-loopback",
"plugins": [{
  "type": "loopback"
}]
}`))
		c.networks = append(c.networks, &Network{
			cni:    c.cniConfig,
			config: loConfig,
			ifName: "lo",
		})
		return nil
	}
}

func WithMinNetworkCount(count int) ConfigOptions {
	return func(c *libcni) error {
		c.networkCount = count
		return nil
	}
}

//TODO: Should we support direct network configs?
/*
func WithConf(byte []bytes) ConfigOptions {
	return func(c *config) error {
			c.networks=
	}
}

func WithConfFile(fileName string) ConfigOptions {
	return func(c *config) error {

	}
}

func WithConfList(byte []bytes) ConfigOptions {
	return func(c *config) error {
      c.
	}
}
func WithConfListFile(files []string) ConfigOptions {
	return func(c *config) error {

	}
}
*/

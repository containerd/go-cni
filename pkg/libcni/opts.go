package libcni

import (
	cnilibrary "github.com/containernetworking/cni/libcni"
)

type ConfigOptions func(c *libcni) error

func WithDefaultIfName(ifName string) ConfigOptions {
	return func(c *libcni) error {
		c.defaultIfName = ifName
		return nil
	}
}

func WithPluginDir(dirs []string) ConfigOptions {
	return func(c *libcni) error {
		c.pluginDirs = append(c.pluginDirs, dirs...)
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
		c.networks["cni-loopback"] = loConfig
		return nil
	}
}

//TODO: Should we support direct network configs?
/*
func WithConf(byte []bytes) ConfigOptions {
	return func(c *config) error {

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

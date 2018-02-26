package libcni

import (
	"fmt"
	"sort"
	"strings"

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

func WithDefaultConfig() ConfigOptions {
	return func(c *libcni) error {
		files, err := cnilibrary.ConfFiles(c.pluginConfDir, []string{".conf", ".conflist", ".json"})
		switch {
		case err != nil:
			return err
		case len(files) == 0:
			return fmt.Errorf("No network config found in %s", c.pluginConfDir)
		}
		// files contains the network config files associated with cni network.
		// Use lexicographical way as a defined order for network config files.
		sort.Strings(files)

		for _, confFile := range files {
			var confList *cnilibrary.NetworkConfigList
			if strings.HasSuffix(confFile, ".conflist") {
				confList, err = cnilibrary.ConfListFromFile(confFile)
				if err != nil {
					fmt.Errorf("Error loading CNI config list file %s: %v", confFile, err)
					continue
				}
			} else {
				conf, err := cnilibrary.ConfFromFile(confFile)
				if err != nil {
					fmt.Errorf("Error loading CNI config file %s: %v", confFile, err)
					continue
				}
				// Ensure the config has a "type" so we know what plugin to run.
				// Also catches the case where somebody put a conflist into a conf file.
				if conf.Network.Type == "" {
					fmt.Errorf("Error loading CNI config file %s: no 'type'; perhaps this is a .conflist?", confFile)
					continue
				}

				confList, err = cnilibrary.ConfListFromConf(conf)
				if err != nil {
					fmt.Errorf("Error converting CNI config file %s to list: %v", confFile, err)
					continue
				}
			}
			if len(confList.Plugins) == 0 {
				fmt.Errorf("CNI config list %s has no networks, skipping", confFile)
				continue
			}
			c.networks = append(c.networks, &Network{
				cni:    c.cniConfig,
				config: confList,
			})
		}

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

func WithConfig(b []byte) ConfigOptions {
	return func(c *libcni) error {
		cfg, err := cnilibrary.ConfListFromBytes(b)
		if err != nil {
			return err
		}

		c.networks = append(c.networks, &Network{
			cni:    c.cniConfig,
			config: cfg,
		})

		return nil
	}
}

//TODO: Should we support direct network configs?
/*
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

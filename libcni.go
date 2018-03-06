package libcni

import (
	"fmt"
	"sort"
	"strings"

	cnilibrary "github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

type CNI interface {
	// Status returns whether the cni plugin is ready.
	Status() error
	// Setup setups the networking for the container.
	Setup(id string, path string, opts ...NamespaceOpts) (*CNIResult, error)
	// Remove tears down the network of the container.
	Remove(id string, path string, opts ...NamespaceOpts) error
}

type libcni struct {
	config

	cniConfig    *cnilibrary.CNIConfig
	networkCount int // minimum network plugin configurations needed to initialize cni
	networks     []*Network
}

func defaultCNIConfig() *libcni {
	return &libcni{
		config: config{
			pluginDirs:    []string{DefaultCNIDir},
			pluginConfDir: DefaultNetDir,
			prefix:        DefaultPrefix,
		},
	}
}

func New(config ...ConfigOptions) CNI {
	cni := defaultCNIConfig()
	cni.cniConfig = &cnilibrary.CNIConfig{Path: cni.pluginDirs}
	for _, c := range config {
		c(cni)
	}
	cni.populateNetworkConfig()
	return cni
}

func (c *libcni) Status() error {
	// TODO this logic changes when CNI Supports
	// Dynamic network updates
	if len(c.networks) < c.networkCount {
		err := c.populateNetworkConfig()
		if err != nil {
			return fmt.Errorf("cni config not intialized: %v", err)
		}
	}
	return nil
}

func (c *libcni) populateNetworkConfig() error {
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
	// Since the CNI spec does not specify a way to detect default networks,
	// the convention chosen is - the first network configuration in the sorted
	// list of network conf files as the default network and choose the default
	// interface provided during init as the network interface for this default
	// network. For every other network use a generated interface id.
	i := 0
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
			ifName: getIfName(c.prefix, i),
		})
		i++
	}
	if len(c.networks) == 0 {
		return fmt.Errorf("No valid networks found in %s", c.pluginDirs)
	}
	return nil
}

func (c *libcni) Setup(id string, path string, opts ...NamespaceOpts) (*CNIResult, error) {
	ns, err := newNamespace(id, path, opts...)
	if err != nil {
		return nil, err
	}
	var results []*current.Result
	for _, network := range c.networks {
		r, err := network.Attach(ns)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return c.GetCNIResultFromResults(results)
}

func (c *libcni) Remove(id string, path string, opts ...NamespaceOpts) error {
	ns, err := newNamespace(id, path, opts...)
	if err != nil {
		return err
	}
	for _, network := range c.networks {
		if err := network.Remove(ns); err != nil {
			return err
		}
	}
	return nil
}

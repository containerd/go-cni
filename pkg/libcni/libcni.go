package libcni

import (
	"fmt"
	"sort"
	"strings"

	cnilibrary "github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

type libcni struct {
	config
	cniConfig *cnilibrary.CNIConfig
	networks  map[string]*cnilibrary.NetworkConfigList
}

type CNI interface {
	// PluginStatus returns whether the cni plugin is ready.
	PluginStatus() error
	// SetupNetworkContainer setups the network Container.
	Setup(ID string, netNS string, opts ...ContainerOptions) ([]*current.Result, error)
	// RemoveNetworkContainer removes the network Container.
	Remove(ID string, netNS string, opts ...ContainerOptions) error
	// ContainerStatus returns the network status of the Container
	ContainerStatus() *Container
}

func defaultCNIConfig() *libcni {
	return &libcni{
		config: config{
			pluginDirs:    []string{DefaultCNIDir},
			pluginConfDir: DefaultNetDir,
			defaultIfName: DefaultIfName,
		},
	}
}

func New(config ...ConfigOptions) CNI {
	cni := defaultCNIConfig()
	cni.cniConfig = &cnilibrary.CNIConfig{Path: cni.pluginDirs}
	cni.networks = make(map[string]*cnilibrary.NetworkConfigList)
	for _, c := range config {
		c(cni)
	}
	cni.populateNetworkConfig()
	return cni
}

func (c *libcni) PluginStatus() error {
	// TODO this logic changes when CNI Support
	// Dynamic network updates
	if len(c.networks) == 0 {
		c.populateNetworkConfig()
	}

	if len(c.networks) == 0 {
		return fmt.Errorf("cni config not intialized")
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
		c.networks[confList.Name] = confList
	}
	if len(c.networks) == 0 {
		return fmt.Errorf("No valid networks found in %s", c.pluginDirs)
	}

	return nil
}

func (c *libcni) Setup(ID string, netNS string, opts ...ContainerOptions) ([]*current.Result, error) {
	container, err := NewContainer(ID, netNS, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Container for %s: %v", ID, err)
	}
	r := container.constructRuntimeConf()
	cninet := &cnilibrary.CNIConfig{
		Path: c.pluginDirs,
	}
	results := []*current.Result{}
	//By default attach container to all networks
	for _, n := range c.networks {
		result, err := container.addNetworks(r, n, cninet)
		if err != nil {
			return nil, fmt.Errorf("failed to attach container %s to network %s", ID, n.Name)
		}
		results = append(results, result)
	}
	return results, nil
}

func (c *libcni) Remove(ID string, netNS string, opts ...ContainerOptions) error {
	container, err := NewContainer(ID, netNS, opts...)
	if err != nil {
		return fmt.Errorf("failed to remove Container for %s: %v", ID, err)
	}
	r := container.constructRuntimeConf()
	cninet := &cnilibrary.CNIConfig{
		Path: c.pluginDirs,
	}
	//By default detach container from all networks
	for _, n := range c.networks {
		err := container.deleteNetworks(r, n, cninet)
		if err != nil {
			return fmt.Errorf("failed to detach container %s from network %s", ID, n.Name)
		}
	}
	return nil
}

func (c *libcni) ContainerStatus() *Container {
	return nil
}

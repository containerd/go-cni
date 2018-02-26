package libcni

import (
	"fmt"

	cnilibrary "github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

type CNI interface {
	// Status returns whether the cni plugin is ready.
	Status() error
	// Setup setups the networking for the container.
	Setup(id string, path string, opts ...NamespaceOpts) ([]*current.Result, error)
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
			defaultIfName: DefaultIfName,
		},
	}
}

func New(config ...ConfigOptions) (CNI, error) {
	cni := defaultCNIConfig()
	cni.cniConfig = &cnilibrary.CNIConfig{Path: cni.pluginDirs}
	for _, c := range config {
		if err := c(cni); err != nil {
			return nil, err
		}
	}
	return cni, nil
}

func (c *libcni) Status() error {
	// TODO this logic changes when CNI Supports
	// Dynamic network updates
	if len(c.networks) < c.networkCount {
		return fmt.Errorf("cni config not intialized")
	}
	return nil
}

func (c *libcni) Setup(id string, path string, opts ...NamespaceOpts) ([]*current.Result, error) {
	ns, err := newNamespace(id, path, c.defaultIfName, opts...)
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
	return results, nil
}

func (c *libcni) Remove(id string, path string, opts ...NamespaceOpts) error {
	ns, err := newNamespace(id, path, c.defaultIfName, opts...)
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

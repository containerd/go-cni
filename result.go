package libcni

import (
	"fmt"
	"net"

	"github.com/containernetworking/cni/pkg/types"
	"github.com/containernetworking/cni/pkg/types/current"
)

type IPConfig struct {
	IP      net.IP
	Gateway net.IP
}

type CNIResult struct {
	Interfaces map[string]*Config
	DNS        []types.DNS
	Routes     []*types.Route
}

type Config struct {
	IPConfigs []*IPConfig
	Mac       string
	Sandbox   string
}

// GetCNIResultFromResults returns a structured data containing the
// interface configuration for each of the interfaces created in the namespace.
// Conforms with
// Result:
// a) Interfaces list. Depending on the plugin, this can include the sandbox
// (eg, container or hypervisor) interface name and/or the host interface
// name, the hardware addresses of each interface, and details about the
// sandbox (if any) the interface is in.
// b) IP configuration assigned to each  interface. The IPv4 and/or IPv6 addresses,
// gateways, and routes assigned to sandbox and/or host interfaces.
// c) DNS information. Dictionary that includes DNS information for nameservers,
// domain, search domains and options.
func (c *libcni) GetCNIResultFromResults(results []*current.Result) (*CNIResult, error) {
	r := &CNIResult{
		Interfaces: make(map[string]*Config),
	}

	// Plugins may not need to return Interfaces in result if
	// if there are no multiple interfaces created. In that case
	// all configs should be applied against default interface
	r.Interfaces[defaultInterface(c.prefix)] = &Config{}

	// Walk through all the results
	for _, result := range results {
		// Walk through all the interface in each result
		for _, intf := range result.Interfaces {
			r.Interfaces[intf.Name] = &Config{
				Mac:     intf.Mac,
				Sandbox: intf.Sandbox,
			}
		}
		// Walk through all the IPs in the result and attach it to corresponding
		// interfaces
		for _, ipConf := range result.IPs {
			if err := validateInterfaceConfig(ipConf, len(result.Interfaces)); err != nil {
				return nil, err
			}
			name := c.getInterfaceName(result.Interfaces, ipConf)
			r.Interfaces[name].IPConfigs = append(r.Interfaces[name].IPConfigs,
				&IPConfig{IP: ipConf.Address.IP, Gateway: ipConf.Gateway})
		}
		r.DNS = append(r.DNS, result.DNS)
		r.Routes = append(r.Routes, result.Routes...)
	}
	if _, ok := r.Interfaces[defaultInterface(c.prefix)]; !ok {
		return nil, fmt.Errorf("namespace not intialized with defualt network")
	}
	return r, nil
}

// getInterfaceName returns the interface name if the plugins
// return the result with associated interfaces. If interface
// is not present then default interface name is used
func (c *libcni) getInterfaceName(interfaces []*current.Interface,
	ipConf *current.IPConfig) string {
	if ipConf.Interface != nil {
		return interfaces[*ipConf.Interface].Name
	}
	return defaultInterface(c.prefix)
}

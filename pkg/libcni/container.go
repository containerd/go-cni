package libcni

import (
	"fmt"

	cnilibrary "github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

// Container holds the info associated with the Container setup request
type Container struct {
	ID          string // ID that uniquely identifies the Container
	NetNS       string // Network Namespace associated with the Container
	IPRanges    []IPRanges
	PortMapping []PortMapping
	IfName      string
	Labels      map[string]string
	IPs         []string
}

func NewContainer(ID string, netNS string, opts ...ContainerOptions) (*Container, error) {
	c := new(Container)
	c.ID = ID
	c.NetNS = netNS
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Container) constructRuntimeConf() *cnilibrary.RuntimeConf {
	r := new(cnilibrary.RuntimeConf)
	r.ContainerID = c.ID
	r.NetNS = c.NetNS
	r.IfName = c.IfName
	for k, v := range c.Labels {
		r.Args = append(r.Args, [2]string{k, v})
	}
	r.CapabilityArgs = make(map[string]interface{})
	if len(c.IPRanges) > 0 {
		r.CapabilityArgs["ipRanges"] = c.IPRanges
	}
	if len(c.PortMapping) > 0 {
		r.CapabilityArgs["portMappings"] = c.PortMapping
	}
	return r
}

func (c *Container) addNetworks(
	r *cnilibrary.RuntimeConf,
	n *cnilibrary.NetworkConfigList,
	cniConfig *cnilibrary.CNIConfig,
) (*current.Result, error) {
	res, err := cniConfig.AddNetworkList(n, r)
	if err != nil {
		fmt.Errorf("Error adding network: %v", err)
		return nil, err
	}
	curRes, err := current.NewResultFromResult(res)
	if err != nil {
		fmt.Errorf("Error translating result: %v", err)
		return nil, err
	}
	return curRes, nil
}

func (c *Container) deleteNetworks(
	r *cnilibrary.RuntimeConf,
	n *cnilibrary.NetworkConfigList,
	cniConfig *cnilibrary.CNIConfig,
) error {
	if err := cniConfig.DelNetworkList(n, r); err != nil {
		return fmt.Errorf("Error adding network: %v", err)
	}
	return nil
}

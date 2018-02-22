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

func NewContainer(ID string, netNS string, ifName string, opts ...ContainerOptions) (*Container, error) {
	c := &Container{
		ID:     ID,
		NetNS:  netNS,
		IfName: ifName,
	}
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Container) constructRuntimeConf() *cnilibrary.RuntimeConf {
	r := &cnilibrary.RuntimeConf{
		ContainerID: c.ID,
		NetNS:       c.NetNS,
		IfName:      c.IfName,
	}
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
		return nil, fmt.Errorf("failed adding network: %v", err)
	}
	newRes, err := current.NewResultFromResult(res)
	if err != nil {
		return nil, fmt.Errorf("failed translating result: %v", err)
	}
	return newRes, nil
}

func (c *Container) deleteNetworks(
	r *cnilibrary.RuntimeConf,
	n *cnilibrary.NetworkConfigList,
	cniConfig *cnilibrary.CNIConfig,
) error {
	if err := cniConfig.DelNetworkList(n, r); err != nil {
		return fmt.Errorf("failed deleting network: %v", err)
	}
	return nil
}

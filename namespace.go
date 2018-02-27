package libcni

import (
	cnilibrary "github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/types/current"
)

type Network struct {
	cni    *cnilibrary.CNIConfig
	config *cnilibrary.NetworkConfigList
	ifName string
}

func (n *Network) Attach(ns *Namespace) (*current.Result, error) {
	r, err := n.cni.AddNetworkList(n.config, ns.config(n.ifName))
	if err != nil {
		return nil, err
	}
	return current.NewResultFromResult(r)
}

func (n *Network) Remove(ns *Namespace) error {
	return n.cni.DelNetworkList(n.config, ns.config(n.ifName))
}

type Namespace struct {
	ID   string
	Path string

	IPRanges    []IPRanges
	PortMapping []PortMapping
	Labels      map[string]string
	IPs         []string
}

func newNamespace(id, path string, opts ...NamespaceOpts) (*Namespace, error) {
	ns := &Namespace{
		ID:   id,
		Path: path,
	}
	for _, o := range opts {
		if err := o(ns); err != nil {
			return nil, err
		}
	}
	return ns, nil
}

func (ns *Namespace) config(ifName string) *cnilibrary.RuntimeConf {
	c := &cnilibrary.RuntimeConf{
		ContainerID: ns.ID,
		NetNS:       ns.Path,
		IfName:      ifName,
	}
	for k, v := range ns.Labels {
		c.Args = append(c.Args, [2]string{k, v})
	}
	c.CapabilityArgs = make(map[string]interface{})
	if len(ns.IPRanges) > 0 {
		c.CapabilityArgs["ipRanges"] = ns.IPRanges
	}
	if len(ns.PortMapping) > 0 {
		c.CapabilityArgs["portMappings"] = ns.PortMapping
	}
	return c
}

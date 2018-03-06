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
	id             string
	path           string
	capabilityArgs map[string]interface{}
	args           map[string]string
}

func newNamespace(id, path string, opts ...NamespaceOpts) (*Namespace, error) {
	ns := &Namespace{
		id:             id,
		path:           path,
		capabilityArgs: make(map[string]interface{}),
		args:           make(map[string]string),
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
		ContainerID: ns.id,
		NetNS:       ns.path,
		IfName:      ifName,
	}
	for k, v := range ns.args {
		c.Args = append(c.Args, [2]string{k, v})
	}
	c.CapabilityArgs = ns.capabilityArgs
	return c
}

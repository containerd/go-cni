package libcni

type ContainerOptions func(s *Container) error

// Non default interface name
func WithContainerIfName(ifname string) ContainerOptions {
	return func(c *Container) error {
		c.IfName = ifname
		return nil
	}
}

// Capabilities
func WithCapabilityPortMap(portMapping []PortMapping) ContainerOptions {
	return func(c *Container) error {
		for _, pmap := range portMapping {
			c.PortMapping = append(c.PortMapping, pmap)
		}
		return nil
	}
}

func WithCapabilityIPRanges(ipRanges []IPRanges) ContainerOptions {
	return func(c *Container) error {
		for _, ipr := range ipRanges {
			c.IPRanges = append(c.IPRanges, ipr)
		}
		return nil
	}
}

// Args
func WithLabels(labels map[string]string) ContainerOptions {
	return func(c *Container) error {
		c.Labels = make(map[string]string)
		for k, v := range labels {
			c.Labels[k] = v
		}
		return nil
	}
}

func WithIPs(ips []string) ContainerOptions {
	return func(c *Container) error {
		c.IPs = append(c.IPs, ips...)
		return nil
	}
}

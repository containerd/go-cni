package libcni

type NamespaceOpts func(s *Namespace) error

// Capabilities
func WithCapabilityPortMap(portMapping []PortMapping) NamespaceOpts {
	return func(c *Namespace) error {
		for _, pmap := range portMapping {
			c.PortMapping = append(c.PortMapping, pmap)
		}
		return nil
	}
}

func WithCapabilityIPRanges(ipRanges []IPRanges) NamespaceOpts {
	return func(c *Namespace) error {
		for _, ipr := range ipRanges {
			c.IPRanges = append(c.IPRanges, ipr)
		}
		return nil
	}
}

// Args
func WithLabels(labels map[string]string) NamespaceOpts {
	return func(c *Namespace) error {
		c.Labels = make(map[string]string)
		for k, v := range labels {
			c.Labels[k] = v
		}
		return nil
	}
}

func WithIPs(ips []string) NamespaceOpts {
	return func(c *Namespace) error {
		c.IPs = append(c.IPs, ips...)
		return nil
	}
}

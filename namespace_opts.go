package libcni

type NamespaceOpts func(s *Namespace) error

// Capabilities
func WithCapabilityPortMap(portMapping []PortMapping) NamespaceOpts {
	return func(c *Namespace) error {
		c.capabilityArgs["portMappings"] = portMapping
		return nil
	}
}

func WithCapabilityIPRanges(ipRanges []IPRanges) NamespaceOpts {
	return func(c *Namespace) error {
		c.capabilityArgs["ipRanges"] = ipRanges
		return nil
	}
}

func WithCapability(name string, capability interface{}) NamespaceOpts {
	return func(c *Namespace) error {
		c.capabilityArgs[name] = capability
		return nil
	}
}

// Args
func WithLabels(labels map[string]string) NamespaceOpts {
	return func(c *Namespace) error {
		for k, v := range labels {
			c.args[k] = v
		}
		return nil
	}
}

func WithArgs(k, v string) NamespaceOpts {
	return func(c *Namespace) error {
		c.args[k] = v
		return nil
	}
}

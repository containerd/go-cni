package libcni

import (
	"net"
)

const (
	CNIPluginName        = "cni"
	DefaultNetDir        = "/etc/cni/net.d"
	DefaultCNIDir        = "/opt/cni/bin"
	VendorCNIDirTemplate = "%s/opt/%s/bin"
	DefaultIfName        = "eth0"
)

type config struct {
	pluginDirs    []string
	pluginConfDir string
	defaultIfName string
}

type Interface struct {
	IP     []net.IP
	Mac    net.HardwareAddr
	Routes string
}

type PortMapping struct {
	HostPort      int
	ContainerPort int
	Protocol      string
}

type IPRanges struct {
	Subnet     string
	RangeStart string
	RangeEnd   string
	Gateway    string
}

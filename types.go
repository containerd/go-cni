package libcni

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

type PortMapping struct {
	HostPort      int32
	ContainerPort int32
	Protocol      string
	HostIP        string
}

type IPRanges struct {
	Subnet     string
	RangeStart string
	RangeEnd   string
	Gateway    string
}

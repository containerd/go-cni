package libcni

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/types/current"
)

func validateInterfaceConfig(ipConf *current.IPConfig, ifs int) error {
	if ipConf == nil {
		return fmt.Errorf("Invalid IP configuration")
	}
	if ipConf.Interface != nil && *ipConf.Interface > ifs {
		return fmt.Errorf("Invalid IP configuration with invalid interface")
	}
	return nil
}

func getIfName(prefix string, i int) string {
	return fmt.Sprintf("%s%d", prefix, i)
}

func defaultInterface(prefix string) string {
	return getIfName(prefix, 0)
}

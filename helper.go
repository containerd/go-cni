package libcni

import (
	"fmt"

	"github.com/containernetworking/cni/pkg/types/current"
)

func validateInterfaceConfig(ipConf *current.IPConfig, ifs int) error {
	if ipConf == nil {
		return fmt.Errorf("invalid IP configuration")
	}
	if ipConf.Interface != nil && *ipConf.Interface > ifs {
		return fmt.Errorf("invalid IP configuration with invalid interface %d", *ipConf.Interface)
	}
	return nil
}

func getIfName(prefix string, i int) string {
	return fmt.Sprintf("%s%d", prefix, i)
}

func defaultInterface(prefix string) string {
	return getIfName(prefix, 0)
}

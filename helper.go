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

func getIfName(ifName string, index int) string {
	if index > 0 {
		return fmt.Sprintf("%s%d", ifName, index)
	}
	return ifName
}

# go-cni

[![PkgGoDev](https://pkg.go.dev/badge/github.com/containerd/go-cni)](https://pkg.go.dev/github.com/containerd/go-cni)
[![Build Status](https://github.com/containerd/go-cni/workflows/CI/badge.svg)](https://github.com/containerd/go-cni/actions?query=workflow%3ACI)
[![codecov](https://codecov.io/gh/containerd/go-cni/branch/main/graph/badge.svg)](https://codecov.io/gh/containerd/go-cni)
[![Go Report Card](https://goreportcard.com/badge/github.com/containerd/go-cni)](https://goreportcard.com/report/github.com/containerd/go-cni)

A generic CNI library to provide APIs for CNI plugin interactions. The library provides APIs to:

- Load CNI network config from different sources  
- Setup networks for container namespace
- Remove networks from container namespace
- Query status of CNI network plugin initialization
- Check verifies the network is still in desired state

go-cni aims to support plugins that implement [Container Network Interface](https://github.com/containernetworking/cni)

## Usage
#### 1. Confirm that the cni plugin exists

```shell
➜  ~ ls /opt/cni/bin
bandwidth  bridge  dhcp  firewall  flannel  host-device  host-local  ipvlan  loopback  macvlan  portmap  ptp  sbr  static  tuning  vlan  vrf
```

#### 2. create `example-ns-1` netns

```shell
➜  ~ ip netns add example-ns-1
```

#### 3. Add cni config

```shell
echo '{"cniVersion":"0.4.0","name":"myptp","type":"ptp","ipMasq":true,"ipam":{"type":"host-local","subnet":"172.16.29.0/24","routes":[{"dst":"0.0.0.0/0"}]}}' | sudo tee /etc/cni/net.d/10-myptp.conf
```

#### 4. Build code
```go
package main

import (
	"context"
	"fmt"
	"log"
	
	gocni "github.com/containerd/go-cni"
)

func main() {
	id := "example"
	netns := "/var/run/netns/example-ns-1"

	// CNI allows multiple CNI configurations and the network interface
	// will be named by eth0, eth1, ..., ethN.
	ifPrefixName := "eth"
	defaultIfName := "eth0"

	// Initializes library
	l, err := gocni.New(
		// one for loopback network interface
		gocni.WithMinNetworkCount(2),
		gocni.WithPluginConfDir("/etc/cni/net.d"),
		gocni.WithPluginDir([]string{"/opt/cni/bin"}),
		// Sets the prefix for network interfaces, eth by default
		gocni.WithInterfacePrefix(ifPrefixName))
	if err != nil {
		log.Fatalf("failed to initialize cni library: %v", err)
	}

	// Load the cni configuration
	if err := l.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
		log.Fatalf("failed to load cni configuration: %v", err)
	}

	// Setup network for namespace.
	labels := map[string]string{
		"K8S_POD_NAMESPACE":          "namespace1",
		"K8S_POD_NAME":               "pod1",
		"K8S_POD_INFRA_CONTAINER_ID": id,
		// Plugin tolerates all Args embedded by unknown labels, like
		// K8S_POD_NAMESPACE/NAME/INFRA_CONTAINER_ID...
		"IgnoreUnknown": "1",
	}

	ctx := context.Background()

	// Teardown network
	defer func() {
		if err := l.Remove(ctx, id, netns, gocni.WithLabels(labels)); err != nil {
			log.Fatalf("failed to teardown network: %v", err)
		}
	}()

	// Setup network
	result, err := l.Setup(ctx, id, netns, gocni.WithLabels(labels))
	if err != nil {
		log.Fatalf("failed to setup network for namespace: %v", err)
	}

	// Get IP of the default interface
	IP := result.Interfaces[defaultIfName].IPConfigs[0].IP.String()
	fmt.Printf("IP of the default interface %s:%s", defaultIfName, IP)
}
```
```shell
➜  cni go build cni-demo.go
➜  cni ./cni-demo
IP of the default interface eth0:172.16.29.4#
```

## Project details

The go-cni is a containerd sub-project, licensed under the [Apache 2.0 license](./LICENSE).
As a containerd sub-project, you will find the:

 * [Project governance](https://github.com/containerd/project/blob/main/GOVERNANCE.md),
 * [Maintainers](https://github.com/containerd/project/blob/main/MAINTAINERS),
 * and [Contributing guidelines](https://github.com/containerd/project/blob/main/CONTRIBUTING.md)

information in our [`containerd/project`](https://github.com/containerd/project) repository.

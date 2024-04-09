module github.com/containerd/go-cni/integration

go 1.21

toolchain go1.21.5

require (
	github.com/containerd/continuity v0.2.2
	github.com/containerd/go-cni v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.2
)

require (
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/containernetworking/cni v1.1.2 // indirect
	github.com/containernetworking/plugins v1.4.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/containerd/go-cni => ../

module github.com/containerd/go-cni

go 1.19

require (
	github.com/containernetworking/cni v1.1.2
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/containernetworking/cni => github.com/MikeZappa87/cni v1.0.2-0.20240226173106-330a4a70d3ab

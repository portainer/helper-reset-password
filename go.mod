module github.com/portainer/helper-reset-password

go 1.13

require (
	github.com/containerd/containerd v1.6.18 // indirect
	github.com/portainer/portainer/api v0.0.0-20230226132524-e0481f69b18a
	github.com/sethvargo/go-password v0.1.3
	golang.org/x/crypto v0.1.0
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200204220554-5f6d6f3f2203 // indirect

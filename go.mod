module github.com/portainer/helper-reset-password

go 1.13

require (
	github.com/portainer/portainer/api v0.0.0-20220725230233-8c1977e0aaba
	github.com/sethvargo/go-password v0.1.3
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200204220554-5f6d6f3f2203 // indirect

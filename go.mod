module github.com/portainer/helper-reset-password

go 1.13

require (
	github.com/boltdb/bolt v1.3.1
	github.com/portainer/portainer/api v0.0.0-20200602235039-38066ece3384
	github.com/sethvargo/go-password v0.1.3
	golang.org/x/crypto v0.0.0-20200602180216-279210d13fed
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200204220554-5f6d6f3f2203 // indirect

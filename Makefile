MAIN=cmd/helper-reset-password/main.go

build: compile image

run:
	go run $(MAIN)

image:
	docker build -t portainer/helper-reset-password .

compile:
	CGO_ENABLED=0 go build -o bin/helper-reset-password $(MAIN)

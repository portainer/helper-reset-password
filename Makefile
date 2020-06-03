MAIN=cmd/helper-reset-password/main.go

build: compile image

run:
	go run $(MAIN)

image:
	docker build -t portainer/helper-reset-password .

compile:
	go build -o bin/helper-reset-password $(MAIN)

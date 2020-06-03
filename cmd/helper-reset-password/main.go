package main

import (
	helper_reset_password "github.com/portainer/helper-reset-password"
	"github.com/portainer/helper-reset-password/bcrypt"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/bolt"
	"github.com/portainer/portainer/api/filesystem"

	"github.com/portainer/helper-reset-password/password"
	"log"
)

func main() {
	fileService := initFileService(helper_reset_password.DataStorePath)

	store, err := createBoltStore(helper_reset_password.DataStorePath, fileService)
	if err != nil {
		log.Fatalf("Unable to create datastore object, err: %s", err)
	}

	err = store.Open()
	if err != nil {
		log.Fatalf("Unable to open the database, err: %s", err)
	}
	defer store.Close()

	user, err := store.User().User(portainer.UserID(1))
	if err != nil {
		log.Fatalf("Unable to retrieve user with ID 1, err: %s", err)
	}

	newPassword, err := password.GeneratePlainTextPassword()
	if err != nil {
		log.Fatalf("An error occured during password generation, err: %s", err)
	}

	hash, err := bcrypt.HashPassword(newPassword)
	if err != nil {
		log.Fatalf("Unable to hash password, err: %s", err)
	}

	user.Password = hash

	err = store.User().UpdateUser(user.ID, user)
	if err != nil {
		log.Fatalf("Unable to persist password changes inside the database, err: %s", err)
	}

	log.Printf("Password succesfully updated for user: %s", user.Username)
	log.Printf("Use the following password to login: %s", newPassword)
}

func createBoltStore(dataStorePath string, fileService portainer.FileService) (portainer.DataStore, error) {
	return bolt.NewStore(dataStorePath, fileService)
}

func initFileService(dataStorePath string) portainer.FileService {
	fileService, err := filesystem.NewService(dataStorePath, "")
	if err != nil {
		log.Fatal(err)
	}
	return fileService
}

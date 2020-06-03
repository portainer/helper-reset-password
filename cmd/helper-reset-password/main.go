package main

import (
	helper_reset_password "github.com/portainer/helper-reset-password"
	"github.com/portainer/helper-reset-password/bolt"
	"github.com/portainer/helper-reset-password/password"
	"log"
	"os"
)

func main() {
	_, err := os.Stat(helper_reset_password.DatabaseFilePath)
	if err != nil {
		log.Fatalf("Unable to locate database file at %s, err: %s", helper_reset_password.DatabaseFilePath, err)
	}

	store, err := createBoltStore()
	if err != nil {
		log.Fatalf("Unable to create datastore object, err: %s", err)
	}

	err = store.Open()
	if err != nil {
		log.Fatalf("Unable to open the database, err: %s", err)
	}
	defer store.Close()

	user, err := store.User().User(1)
	if err != nil {
		log.Fatalf("Unable to retrieve user with ID 1, err: %s", err)
	}

	newPassword, err := password.GeneratePassword()
	if err != nil {
		log.Fatalf("An error occured during password generation, err: %s", err)
	}

	user.Password = newPassword

	err = store.User().UpdateUser(user.ID, user)
	if err != nil {
		log.Fatalf("Unable to persist password changes inside the database, err: %s", err)
	}

	log.Printf("Password succesfully updated for user: %s", user.Username)
	log.Printf("Use the following password to login: %s", user.Password)
}

func createBoltStore() (helper_reset_password.DataStore, error) {
	return bolt.NewStore(helper_reset_password.DatabaseFilePath)
}

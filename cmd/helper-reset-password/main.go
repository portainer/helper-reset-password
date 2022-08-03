package main

import (
	"fmt"
	"log"
	"os"
	"path"

	helper_reset_password "github.com/portainer/helper-reset-password"
	"github.com/portainer/helper-reset-password/bcrypt"
	"github.com/portainer/helper-reset-password/password"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/bolt"
	"github.com/portainer/portainer/api/filesystem"
)

func main() {
	if _, err := os.Stat(path.Join(helper_reset_password.DataStorePath, "portainer.db")); err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("Unable to locate /data/portainer.db on disk")
		}
		log.Fatalf("Unable to verify database file existence, err: %s", err)
	}

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

	createAdmin := false
	adminName := "admin"

	user, err := store.User().User(portainer.UserID(1))
	if err != nil {
		log.Printf("[WARN] Unable to retrieve user with ID 1, will try to create, err: %s", err)

		createAdmin = true

		adminUser, err := store.User().UserByUsername(adminName)

		if err != nil {
			log.Fatalf("Unable to query user %s from database, err: %s", adminName, err)
		} else if adminUser != nil {
			adminName, err = password.GenerateRandomString()
			if err != nil {
				log.Fatalf("Unable to generate random admin user name, err: %s", err)
			}
			adminName = fmt.Sprintf("admin-%s", adminName)
		}
	}

	newPassword, err := password.GeneratePlainTextPassword()
	if err != nil {
		log.Fatalf("An error occured during password generation, err: %s", err)
	}

	hash, err := bcrypt.HashPassword(newPassword)
	if err != nil {
		log.Fatalf("Unable to hash password, err: %s", err)
	}

	if createAdmin {
		if err := store.User().CreateUser(&portainer.User{
			Username: adminName,
			Role:     portainer.AdministratorRole,
			Password: hash,
		}); err != nil {
			log.Fatalf("Unable to create admin user %s inside the database, err: %s", user.Username, err)
		}

		log.Printf("Admin user %s successfully created", user.Username)
	} else {
		user.Password = hash

		err = store.User().UpdateUser(user.ID, user)
		if err != nil {
			log.Fatalf("Unable to persist password changes inside the database, err: %s", err)
		}

		log.Printf("Password succesfully updated for user: %s", user.Username)
	}

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

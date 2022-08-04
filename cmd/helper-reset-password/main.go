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
	"github.com/portainer/portainer/api/database"
	"github.com/portainer/portainer/api/datastore"
	"github.com/portainer/portainer/api/filesystem"
)

func main() {
	// try to locate the db file
	if _, err := os.Stat(path.Join(helper_reset_password.DataStorePath, "portainer.db")); err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("Unable to locate /data/portainer.db on disk")
		}
		log.Fatalf("Unable to verify database file existence, err: %s", err)
	}

	// db init bolt store from teh db file
	// note: encrypted db isn't supported ATM
	fileService := initFileService(helper_reset_password.DataStorePath)
	store, err := createBoltStore(helper_reset_password.DataStorePath, fileService)
	isNew, err := store.Open()
	if err != nil {
		log.Fatalf("Unable to open the database, err: %v", err)
	} else if isNew {
		log.Fatalf("Data store not found at %s", helper_reset_password.DataStorePath)
	}
	defer store.Close()

	createAdmin := false
	// default user1 name
	adminName := "admin"

	// try to find user1
	user, err := store.User().User(portainer.UserID(1))
	if err != nil {
		// if user1 doesn't exist, will create later
		log.Printf("[WARN] Unable to retrieve user with ID 1, will try to create, err: %s", err)

		createAdmin = true

		// if there is already a user named admin, will randomize the name with suffix (admin-asdkjfh123)
		adminUser, err := store.User().UserByUsername(adminName)

		if adminUser != nil {
			adminName, err = password.GenerateRandomString()
			if err != nil {
				log.Fatalf("Unable to generate random admin user name, err: %s", err)
			}
			adminName = fmt.Sprintf("admin-%s", adminName)
		}
	}

	// generate the new password
	newPassword, err := password.GeneratePlainTextPassword()
	if err != nil {
		log.Fatalf("An error occured during password generation, err: %s", err)
	}

	// hash the password
	hash, err := bcrypt.HashPassword(newPassword)
	if err != nil {
		log.Fatalf("Unable to hash password, err: %s", err)
	}

	if createAdmin {
		// create user1 when needed
		if err := store.GetConnection().CreateObjectWithId(
			store.User().BucketName(),
			1,
			&portainer.User{
				ID:       1,
				Username: adminName,
				Role:     portainer.AdministratorRole,
				Password: hash,
			},
		); err != nil {
			log.Fatalf("Unable to create admin user %s inside the database, err: %s", adminName, err)
		}

		// try to make sure the bolt db user busket sequence is > 1, 10 attempts
		seq := store.GetConnection().GetNextIdentifier(store.User().BucketName())
		for i := 1; i <= 10 && seq <= 1; i++ {
			seq = store.GetConnection().GetNextIdentifier(store.User().BucketName())
		}
		// if the bucket sequence is still less than 1, exit gracefully. in theory this should not happen.
		if seq <= 1 {
			log.Fatalf("Unable to increase the data store key sequence after 10 attempts. Current seq: %d", seq)
		}

		log.Printf("Admin user %s successfully created", adminName)
	} else {
		// update user1 with the generated password
		user.Password = hash

		err = store.User().UpdateUser(user.ID, user)
		if err != nil {
			log.Fatalf("Unable to persist password changes inside the database, err: %s", err)
		}

		log.Printf("Password succesfully updated for user: %s", user.Username)
	}

	log.Printf("Use the following password to login: %s", newPassword)
}

func createBoltStore(dataStorePath string, fileService portainer.FileService) (datastore.Store, error) {
	connection, err := database.NewDatabase("boltdb", dataStorePath, nil)
	if err != nil {
		log.Fatalf("failed creating database connection: %s", err)
	}

	store := datastore.NewStore(dataStorePath, fileService, connection)

	return *store, nil
}

func initFileService(dataStorePath string) portainer.FileService {
	fileService, err := filesystem.NewService(dataStorePath, "")
	if err != nil {
		log.Fatal(err)
	}
	return fileService
}

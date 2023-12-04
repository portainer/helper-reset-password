package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	helper_reset_password "github.com/portainer/helper-reset-password"
	"github.com/portainer/helper-reset-password/password"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/crypto"
	"github.com/portainer/portainer/api/database"
	"github.com/portainer/portainer/api/datastore"
	"github.com/portainer/portainer/api/filesystem"
)

func parseCommandLineArguments() (string, string, error) {
	var (
		password     string
		passwordHash string
		err          error
	)

	flag.StringVar(&password, "password", "", "The new admin password")
	flag.StringVar(&passwordHash, "password-hash", "", "The new admin password hash")

	flag.Parse()

	if password != "" && passwordHash != "" {
		err = errors.New("You cannot use the 'password' and 'password-hash' arguments at the same time")
	}

	return password, passwordHash, err
}

func main() {
	// parse CLI arguments
	cliPassword, cliPasswordHash, err := parseCommandLineArguments()
	if err != nil {
		log.Fatalf("Invalid CLI usage! err: %s", err)
	}
	// try to locate the db file
	if _, err := os.Stat(path.Join(helper_reset_password.DataStorePath, "portainer.db")); err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("Unable to locate /data/portainer.db on disk")
		}
		log.Fatalf("Unable to verify database file existence, err: %s", err)
	}

	// db init bolt store from the db file
	// note: encrypted db isn't supported ATM
	fileService := initFileService(helper_reset_password.DataStorePath)
	store, err := createBoltStore(helper_reset_password.DataStorePath, fileService)
	if err != nil {
		log.Fatalf("Unable to create boltdb store, err: %v", err)
	}

	isNew, err := store.Open()
	if err != nil {
		log.Fatalf("Unable to open the database, err: %v", err)
	} else if isNew {
		log.Fatalf("Data store not found at %s", helper_reset_password.DataStorePath)
	}
	defer store.Close()

	cryptoService := crypto.Service{}

	createAdmin := false
	// default user1 name
	adminName := "admin"

	settings, err := store.Settings().Settings()
	if err != nil {
		log.Fatalf("Unable to retrieve settings, err: %s", err)
	}

	if settings.IsDockerDesktopExtension {
		log.Fatalf("Database from a Docker Desktop Portainer instance detected - exiting without resetting")
	}

	// try to find user1
	user, err := store.User().User(portainer.UserID(1))
	if err != nil {
		// if user1 doesn't exist, will create later
		log.Printf("[WARN] Unable to retrieve user with ID 1, will try to create, err: %s", err)

		createAdmin = true

		// if there is already a user named admin, will randomize the name with suffix (admin-asdkjfh123)
		adminUser, _ := store.User().UserByUsername(adminName)

		if adminUser != nil {
			adminName, err = password.GenerateRandomString()
			if err != nil {
				log.Fatalf("Unable to generate random admin user name, err: %s", err)
			}
			adminName = fmt.Sprintf("admin-%s", adminName)
		}
	}

	// If password is used for docker extension. It won't return an error if passwords are same.
	if err := cryptoService.CompareHashAndData(user.Password, "K7yJPP5qNK4hf1QsRnfV"); err == nil {
		log.Fatalf("Database from a Docker Desktop Portainer instance detected - exiting without resetting")
	}

	// generate the new password if not given via CLI
	var newPassword string
	if cliPassword == "" {
		newPassword, err = password.GeneratePlainTextPassword()
		if err != nil {
			log.Fatalf("An error occurred during password generation, err: %s", err)
		}
	} else {
		log.Printf("Using password provided via CLI")
		newPassword = cliPassword
	}

	// hash the password if not given via CLI
	var hash string
	if cliPasswordHash == "" {
		hash, err = cryptoService.Hash(newPassword)
		if err != nil {
			log.Fatalf("Unable to hash password, err: %s", err)
		}
	} else {
		log.Printf("Using password hash provided via CLI")
		hash = cliPasswordHash
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

		// try to make sure the bolt db user bucket sequence is > 1, 10 attempts
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

		log.Printf("Password successfully updated for user: %s", user.Username)
	}

	if cliPasswordHash == "" {
		log.Printf("Use the following password to login: %s", newPassword)
	} else {
		log.Printf("Use the password from your provided hash to login")
	}
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

package cmdline

import (
	"errors"
	"flag"

	helper_reset_password "github.com/portainer/helper-reset-password"
)

type Arguments struct {
	Password     string
	PasswordHash string
	DataPath     string
}

func ParseArguments() (Arguments, error) {
	var (
		password     string
		passwordHash string
		dataPath     string
	)

	flag.StringVar(&password, "password", "", "The new admin password")
	flag.StringVar(&passwordHash, "password-hash", "", "The new admin password hash")
	flag.StringVar(&dataPath, "data-path", helper_reset_password.DataStorePath, "The path to the Portainer data store")

	flag.Parse()

	if password != "" && passwordHash != "" {
		return Arguments{}, errors.New("you cannot use the 'password' and 'password-hash' arguments at the same time")
	}

	// set default data path if empty (flags does not give default value when empty string is provided)
	if dataPath == "" {
		dataPath = helper_reset_password.DataStorePath
	}

	return Arguments{
		Password:     password,
		PasswordHash: passwordHash,
		DataPath:     dataPath,
	}, nil
}

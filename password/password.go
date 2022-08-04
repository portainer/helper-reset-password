package password

import "github.com/sethvargo/go-password/password"

func GeneratePlainTextPassword() (string, error) {
	pass, err := password.Generate(32, 8, 8, false, false)
	if err != nil {
		return "", err
	}

	return pass, nil
}

func GenerateRandomString() (string, error) {
	pass, err := password.Generate(16, 8, 0, true, true)
	if err != nil {
		return "", err
	}

	return pass, nil
}

package helper_reset_password

type (
	DataStore interface {
		Open() error
		Close() error
		User() UserService
	}

	UserService interface {
		User(ID int) (*User, error)
		UpdateUser(ID int, user *User) error
	}

	User struct {
		ID       int    `json:"Id"`
		Username string `json:"Username"`
		Password string `json:"Password"`
	}
)

const (
	DatabaseFilePath = "/data/portainer.db"
)

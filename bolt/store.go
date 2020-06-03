package bolt

import (
	"github.com/boltdb/bolt"
	helper_reset_password "github.com/portainer/helper-reset-password"
	"github.com/portainer/helper-reset-password/bolt/user"
	"time"
)

type Store struct {
	path string
	db   *bolt.DB

	userService *user.Service
}

// NewStore initializes a new Store and the associated services
func NewStore(storePath string) (*Store, error) {
	store := &Store{
		path: storePath,
	}

	return store, nil
}

// Open opens and initializes the BoltDB database.
func (store *Store) Open() error {
	db, err := bolt.Open(store.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	store.db = db
	store.userService = user.NewService(db)
	return nil
}

// Close closes the BoltDB database.
func (store *Store) Close() error {
	if store.db != nil {
		return store.db.Close()
	}
	return nil
}

// User gives access to the User data management layer
func (store *Store) User() helper_reset_password.UserService {
	return store.userService
}

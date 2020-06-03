package internal

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
)

// Itob returns an 8-byte big endian representation of v.
// This function is typically used for encoding integer IDs to byte slices
// so that they can be used as BoltDB keys.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// GetObject is a generic function used to retrieve an unmarshalled object from a bolt database.
func GetObject(db *bolt.DB, bucketName string, key []byte, object interface{}) error {
	var data []byte

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		value := bucket.Get(key)
		if value == nil {
			return errors.New("unable to find object inside the database")
		}

		data = make([]byte, len(value))
		copy(data, value)

		return nil
	})
	if err != nil {
		return err
	}

	return UnmarshalObject(data, object)
}

// UpdateObject is a generic function used to update an object inside a bolt database.
func UpdateObject(db *bolt.DB, bucketName string, key []byte, object interface{}) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		data, err := MarshalObject(object)
		if err != nil {
			return err
		}

		err = bucket.Put(key, data)
		if err != nil {
			return err
		}

		return nil
	})
}

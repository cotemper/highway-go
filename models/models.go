package models

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

// ErrUsernameTaken is thrown when a user attempts to register a username that is taken.
var ErrUsernameTaken = errors.New("username already taken")

// Copy of auth.GenerateSecureKey to prevent cyclic import with auth library
func generateSecureKey() string {
	k := make([]byte, 32)
	io.ReadFull(rand.Reader, k)
	return fmt.Sprintf("%x", k)
}

// BytesToID converts a byte slice to a uint. This is needed because the
// WebAuthn specification deals with byte buffers, while the primary keys in
// our database are uints.
func BytesToID(buf []byte) uint {
	// TODO: Probably want to catch the number of bytes converted in production
	id, _ := binary.Uvarint(buf)
	return uint(id)
}

//TODO change stripe model
type SnrItem struct {
	ID string `json:"id"`
}

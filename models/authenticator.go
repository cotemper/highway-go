package models

import (
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
)

// Authenticator is a struct representing a WebAuthn authenticator, which is
// responsible for generating Credentials. For this demo, we map a single
// credential to a single authenticator.
type Authenticator struct {
	gorm.Model
	webauthn.Authenticator
}

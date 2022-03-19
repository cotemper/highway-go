package models

import (
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
)

// PlaceholderUsername is the default username to use if one isn't provided by
// the user (in the case of a placeholder)
const PlaceholderUsername = "testuser"

// PlaceholderUserIcon is the default user icon used when creating a new user
const PlaceholderUserIcon = "example.icon.duo.com/123/avatar.png"

type User struct {
	gorm.Model
	Did         string
	Jwt         Jwt
	Names       []string
	Username    string       `json:"name" sql:"not null;"`
	DisplayName string       `json:"display_name"`
	Icon        string       `json:"icon,omitempty"`
	Credentials []Credential `json:"credentials,omitempty"`
	Paid        bool         `json:"paid"`
	PiID        string       `json:"piid"`
	Created     time.Time    `json:"created"`
}

// WebAuthnID returns the user ID as a byte slice
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	return buf
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return u.Username
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnIcon returns the user's icon
func (u User) WebAuthnIcon() string {
	return u.Icon
}

// WebAuthnCredentials helps implement the webauthn.User interface by loading
// the user's credentials from the underlying database.
func (u User) WebAuthnCredentials() []webauthn.Credential {
	credentials := u.Credentials
	wcs := make([]webauthn.Credential, len(credentials))
	for i, cred := range credentials {
		credentialID, _ := base64.URLEncoding.DecodeString(cred.CredentialID)
		wcs[i] = webauthn.Credential{
			ID:            credentialID,
			PublicKey:     cred.PublicKey,
			Authenticator: cred.WebauthnAuthenticator(),
		}
	}
	return wcs
}

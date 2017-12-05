package userdb

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/yurinandayona-com/kuma/server"
)

// JWTManager is token manager of user DB.
//
// It implements server.UserVerifier.
type JWTManager struct {
	UserDB  *UserDB
	HMACKey []byte
}

var _ server.UserVerifier = &JWTManager{}

// JWTUserClaims is JWT claims containing user information.
type JWTUserClaims struct {
	jwt.StandardClaims

	ID   string `json:"https://github.com/yurinandayona-com/kuma/claim-types/user-id"`
	Name string `json:"https://github.com/yurinandayona-com/kuma/claim-types/user-name"`
}

// Verify verifies t as JWT and then returns a user bound this JWT or error.
func (jm *JWTManager) Verify(t string) (server.User, error) {
	claims, valid, err := jm.Parse(t)
	if err != nil {
		return nil, err
	}

	if valid {
		return jm.UserDB.Verify(claims.ID, claims.Name)
	}

	return nil, errors.New("invalid JWT")
}

// Parse parses t as JWT and then returns JWTUserClaims bound this JWT and
// flag which is token validation status and error.
func (jm *JWTManager) Parse(t string) (*JWTUserClaims, bool, error) {
	token, err := jwt.ParseWithClaims(t, &JWTUserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid JWT algorithm")
		}
		return jm.HMACKey, nil
	})

	if token != nil {
		if claims, ok := token.Claims.(*JWTUserClaims); ok {
			return claims, token.Valid, errors.Wrap(err, "invalid JWT token")
		}
	}

	return nil, false, errors.New("invalid JWT")
}

// Sign returns signed JWT with given user information and expiration.
func (jm *JWTManager) Sign(u *User, expiration time.Time) (string, error) {
	u, err := jm.UserDB.Verify(u.ID, u.Name)
	if err != nil {
		return "", err
	}

	claims := &JWTUserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
		ID:   u.ID,
		Name: u.Name,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jm.HMACKey)
}

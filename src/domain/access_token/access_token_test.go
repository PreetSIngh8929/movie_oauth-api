package access_token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenConstants(t *testing.T) {

	assert.EqualValues(t, 24, expirationTime, "Expiration time shopuld be 24 hrs")
}

func TestGetNewAccessToken(t *testing.T) {

	at := GetNewAccessToken(17)
	assert.False(t, at.IsExpired(), "brand new access token should not be expired")
	assert.EqualValues(t, "", at.AccessToken, "new access token should not have defined access token id")
	assert.True(t, at.ClientId == 0, "new access token should not have assosciated client id")
}

func TestAccessTokenIsExpired(t *testing.T) {
	at := AccessToken{}
	assert.True(t, at.IsExpired(), "empty access token should be expired by default")
	at.Expires = time.Now().UTC().Add(3 * time.Hour).Unix()
	assert.False(t, at.IsExpired(), "access token expiring three hours from now should not be expired")
}

func TestValidateEmptyAccessToken(t *testing.T) {
	at := AccessToken{
		AccessToken: "",
	}
	assert.EqualValues(t, "invalid access token id", at.Validate().Message())
}
func TestValidateInvalidUserId(t *testing.T) {
	at := AccessToken{
		AccessToken: "bcxhdYhj99",
		UserId:      -1,
	}
	assert.EqualValues(t, "invalid user id", at.Validate().Message())
}
func TestValidateInvalidClientId(t *testing.T) {
	at := AccessToken{
		AccessToken: "bcxhdYhj99",
		UserId:      1,
		ClientId:    -1,
	}
	assert.EqualValues(t, "invalid client id", at.Validate().Message())
}
func TestValidateInvalidExpirationTime(t *testing.T) {
	at := AccessToken{
		AccessToken: "bcxhdYhj99",
		UserId:      1,
		ClientId:    1,
		Expires:     -1,
	}
	assert.EqualValues(t, "invalid expiration time", at.Validate().Message())
}
func TestValidateNoError(t *testing.T) {
	at := AccessToken{
		AccessToken: "bcxhdYhj99",
		UserId:      1,
		ClientId:    1,
		Expires:     2,
	}
	assert.Nil(t, at.Validate())
}

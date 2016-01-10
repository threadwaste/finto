// +build integration

package finto

import (
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
)

func setupIntegrationTest() (rc *RoleSet) {
	client := sts.New(
		session.New(),
		&aws.Config{Credentials: credentials.NewEnvCredentials()},
	)

	rc = NewRoleSet(client)
	rc.SetRole("valid", os.Getenv("FINTO_VALID_ARN"))
	rc.SetRole("invalid", os.Getenv("FINTO_INVALID_ARN"))

	return
}

func TestSTSClientAssumeRole(t *testing.T) {
	rc := setupIntegrationTest()
	role, err := rc.Role("valid")

	if assert.NoError(t, err) {
		creds, _ := role.Credentials()

		assert.NotEmpty(t, creds.AccessKeyId)
		assert.NotEmpty(t, creds.Expiration)
		assert.NotEmpty(t, creds.SecretAccessKey)
		assert.NotEmpty(t, creds.SessionToken)
		assert.False(t, creds.IsExpired(), "Fresh credentials should not be expired")
	}
}

func TestSTSClientAssumeRoleFailure(t *testing.T) {
	rc := setupIntegrationTest()
	role, err := rc.Role("invalid")
	creds, err := role.Credentials()

	if assert.Error(t, err) {
		assert.Empty(t, creds.AccessKeyId)
		assert.Equal(t, time.Time{}, creds.Expiration)
		assert.Empty(t, creds.SecretAccessKey)
		assert.Empty(t, creds.SessionToken)
	}
}

package finto

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/assert"
)

var MockExpiry time.Time = time.Unix(11833862400, 0)

// A mock client that satisfies the AssumeRoleClient interface. For testing
// purposes.
type MockAssumeRoleClient struct{}

// Return a canned sts.AssumeRoleOutput.
func (c *MockAssumeRoleClient) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	mockId := *input.RoleArn + "-" + *input.RoleSessionName

	return &sts.AssumeRoleOutput{
		Credentials: &sts.Credentials{
			AccessKeyId:     aws.String(mockId),
			Expiration:      &MockExpiry,
			SecretAccessKey: aws.String("mock-key"),
			SessionToken:    aws.String("mock-token"),
		},
	}, nil
}

func TestCredentials(t *testing.T) {
	var (
		uxt = time.Now().Add(10 * time.Minute)
		xt  = time.Unix(0, 0)
	)

	cases := []struct {
		id, key, token string
		expiration     time.Time
		expired        bool
		result         *Credentials
	}{
		{
			"test-id",
			"test-key",
			"test-token",
			uxt,
			false,
			&Credentials{
				AccessKeyId:     "test-id",
				Expiration:      uxt,
				SecretAccessKey: "test-key",
				SessionToken:    "test-token",
			},
		},
		{
			"expired-id",
			"expired-key",
			"expired-token",
			xt,
			true,
			&Credentials{
				AccessKeyId:     "expired-id",
				Expiration:      xt,
				SecretAccessKey: "expired-key",
				SessionToken:    "expired-token",
			},
		},
	}

	for _, c := range cases {
		creds := &Credentials{}
		creds.SetCredentials(c.id, c.key, c.token)
		creds.SetExpiration(c.expiration, 0)

		assert.Equal(t, c.result, creds)
		assert.Equal(t, c.expired, creds.IsExpired())
	}
}

func TestRole(t *testing.T) {
	cases := []struct {
		arn, session string
		return_id    string
		expired      bool
	}{
		{
			"test-arn",
			"test-session",
			"test-arn-test-session",
			false,
		},
	}

	for _, c := range cases {
		r := NewRole(c.arn, c.session, &MockAssumeRoleClient{})
		creds, _ := r.Credentials()

		assert.Equal(t, c.expired, r.IsExpired())
		assert.Equal(t, c.return_id, creds.AccessKeyId)
	}
}

func TestRoleSet(t *testing.T) {
	rs := NewRoleSet(&MockAssumeRoleClient{})
	rs.SetRole("test-alias", "test-arn")
	rs.SetRole("active-alias", "active-arn")

	assert.Equal(t, []string{"active-alias", "test-alias"}, rs.Roles())

	role, err := rs.Role("test-alias")
	if assert.NoError(t, err) {
		assert.Equal(t, "test-arn", role.Arn())
		assert.Equal(t, "finto-test-alias", role.SessionName())
	}

	_, err = rs.Role("fake-role")
	assert.Error(t, err)
}

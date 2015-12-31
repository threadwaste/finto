package roles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

package roles

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

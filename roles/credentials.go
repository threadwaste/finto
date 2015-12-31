package roles

import (
	"time"
)

// Credentials represents a set of temporary credentials and their expiration.
type Credentials struct {
	AccessKeyId     string
	Expiration      time.Time
	SecretAccessKey string
	SessionToken    string
}

func (c *Credentials) IsExpired() bool {
	return c.Expiration.Before(time.Now())
}

func (c *Credentials) SetCredentials(id, key, token string) {
	c.AccessKeyId = id
	c.SecretAccessKey = key
	c.SessionToken = token
}

// Sets expiration. Accepts an offset to allow for early expiration. This helps
// avoid returning credentials that expire "in flight."
func (c *Credentials) SetExpiration(t time.Time, offset time.Duration) {
	c.Expiration = t.Add(-offset)
}

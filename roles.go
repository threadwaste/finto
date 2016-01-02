package finto

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
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

// AssumeRoleClient is a basic interface that wraps role assumption.
//
// AssumeRole takes the input of a role ARN and session name, and returns a set
// of credentials including: an access key ID, a secret access key, a session
// token, and their expiration.
//
// https://godoc.org/github.com/aws/aws-sdk-go/service/sts#AssumeRoleInput
// https://godoc.org/github.com/aws/aws-sdk-go/service/sts#AssumeRoleOutput
type AssumeRoleClient interface {
	AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error)
}

// Implements a role, the retrieval of its credentials, and management of their
// expiration.
type Role struct {
	arn         string      // The role's Amazon Resource Name
	creds       Credentials // The role's credentials
	sessionName string      // The session name recorded by assumption

	client AssumeRoleClient // An AssumeRoleClient for retrieving credentials
	m      sync.Mutex
}

func NewRole(a, s string, c AssumeRoleClient) *Role {
	return &Role{
		arn:         a,
		sessionName: s,
		client:      c,
	}
}

func (r *Role) Arn() string {
	return r.arn
}

func (r *Role) SessionName() string {
	return r.sessionName
}

// Returns whether the role's current credentials are expired.
func (r *Role) IsExpired() bool {
	r.m.Lock()
	defer r.m.Unlock()

	return r.isExpired()
}

func (r *Role) isExpired() bool {
	return r.creds.IsExpired()
}

// Returns the role's credentials. If expired, credentials are refreshed through
// the client.
func (r *Role) Credentials() (Credentials, error) {
	r.m.Lock()
	defer r.m.Unlock()

	if r.isExpired() {
		resp, err := r.client.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         aws.String(r.Arn()),
			RoleSessionName: aws.String(r.SessionName()),
		})

		if err != nil {
			return Credentials{}, err
		}

		creds := resp.Credentials
		r.creds.SetCredentials(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)
		r.creds.SetExpiration(*creds.Expiration, 300)
	}

	return r.creds, nil
}

// A collection of aliased roles.
type RoleSet struct {
	roles map[string]*Role

	client AssumeRoleClient
	m      sync.Mutex
}

func NewRoleSet(c AssumeRoleClient) *RoleSet {
	return &RoleSet{
		client: c,
		roles:  make(map[string]*Role),
	}
}

func (rs *RoleSet) Role(alias string) (*Role, error) {
	rs.m.Lock()
	defer rs.m.Unlock()

	if role, ok := rs.roles[alias]; ok {
		return role, nil
	}

	return &Role{}, fmt.Errorf("unknown role: %s", alias)
}

func (rs *RoleSet) Roles() (roles []string) {
	rs.m.Lock()
	defer rs.m.Unlock()

	roles = make([]string, len(rs.roles))

	i := 0
	for k := range rs.roles {
		roles[i] = k
		i += 1
	}

	sort.Strings(roles)
	return
}

// Set an alias's role configuration.
func (rs *RoleSet) SetRole(alias, arn string) {
	rs.m.Lock()
	defer rs.m.Unlock()

	rs.roles[alias] = NewRole(arn, fmt.Sprintf("finto-%s", alias), rs.client)
}

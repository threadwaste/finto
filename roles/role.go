package roles

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

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
	arn         string
	creds       Credentials
	sessionName string

	client AssumeRoleClient
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

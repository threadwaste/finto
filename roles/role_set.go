// A simple cache of AWS IAM roles and their temporary credentials.
package roles

import (
	"fmt"
	"sort"
	"sync"
)

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

func (rc *RoleSet) Role(alias string) (*Role, error) {
	rc.m.Lock()
	defer rc.m.Unlock()

	if role, ok := rc.roles[alias]; ok {
		return role, nil
	}

	return &Role{}, fmt.Errorf("unknown role: %s", alias)
}

func (rc *RoleSet) Roles() (roles []string) {
	rc.m.Lock()
	defer rc.m.Unlock()

	roles = make([]string, len(rc.roles))

	i := 0
	for k := range rc.roles {
		roles[i] = k
		i += 1
	}

	sort.Strings(roles)
	return
}

// Set an alias's role configuration.
func (rc *RoleSet) SetRole(alias, arn string) {
	rc.m.Lock()
	defer rc.m.Unlock()

	rc.roles[alias] = NewRole(arn, fmt.Sprintf("finto-%s", alias), rc.client)
}

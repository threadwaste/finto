package roles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

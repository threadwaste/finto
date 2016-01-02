package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const configExample = `
{

	"credentials": {
		"file": "a",
		"profile": "a"
	},
	"default_role": "1",
	"roles": {
		"1": "arn",
		"2": "arn"
	}
}
`

func setupConfigTests(t *testing.T) string {
	f, err := ioutil.TempFile("", "config-test")
	if err != nil {
		t.Fatal("Error creating file", err)
	}

	err = ioutil.WriteFile(f.Name(), []byte(configExample), 0644)
	if err != nil {
		t.Fatal("Error writing file", err)
	}

	return f.Name()
}

func teardownConfigTests(f string) {
	_ = os.Remove(f)
}

func ExampleStringRendering() {
	var c = Config{
		DefaultRole: "alias",
		Credentials: CredentialsConfig{
			File:    "File",
			Profile: "ID",
		},
		Roles: RolesConfig{
			"alias": "arn",
		},
	}

	fmt.Println(c.String())

	// Output:
	// {
	//   "default_role": "alias",
	//   "credentials": {
	//     "file": "File",
	//     "profile": "ID"
	//   },
	//   "roles": {
	//     "alias": "arn"
	//   }
	// }
}

func TestLoadsConfigFile(t *testing.T) {
	file := setupConfigTests(t)
	defer teardownConfigTests(file)

	expected := &Config{
		DefaultRole: "1",
		Credentials: CredentialsConfig{
			File:    "a",
			Profile: "a",
		},
		Roles: RolesConfig{
			"1": "arn",
			"2": "arn",
		},
	}

	c, err := LoadConfig(file)

	assert.Nil(t, err)
	assert.Equal(t, c, expected)
}

func TestLoadMissingConfigFile(t *testing.T) {
	_, err := LoadConfig("")
	assert.Error(t, err)
}

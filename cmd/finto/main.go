package main

import (
	"flag"
	"fmt"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/threadwaste/finto"
)

var (
	fintorc = flag.String("config", defaultRC(), "location of config file")

	addr = flag.String("addr", "169.254.169.254", "bind to addr")
	// TODO: logging for both main and handlers
	logfile = flag.String("log", "", "log to file")
	port    = flag.Uint("port", 16925, "listen on port")
)

func main() {
	flag.Parse()

	config, err := LoadConfig(*fintorc)
	if err != nil {
		panic(err)
	}

	// SharedCredentialsProvider defaults to file=~/.aws/credentials and
	// profile=default when provided zero-value strings
	stsClient := sts.New(session.New(), &aws.Config{
		Credentials: credentials.NewSharedCredentials(
			config.Credentials.File,
			config.Credentials.Profile,
		),
	})

	rs := finto.NewRoleSet(stsClient)
	for alias, arn := range config.Roles {
		rs.SetRole(alias, arn)
	}

	context, err := finto.InitFintoContext(rs, config.DefaultRole)
	if err != nil {
		fmt.Printf("initializing: %s\n", err)
	}

	router := finto.FintoRouter(&context)
	err = http.ListenAndServe(fmt.Sprint(*addr, ":", *port), router)
	if err != nil {
		panic(err)
	}
}

func homeDir() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get user: %s", err)
	}

	return currentUser.HomeDir, err
}

func defaultRC() string {
	dir, err := homeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(dir, ".fintorc")
}

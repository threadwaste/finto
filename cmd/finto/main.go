package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gorilla/handlers"
	"github.com/threadwaste/finto"
)

var (
	fintorc = flag.String("config", defaultRC(), "location of config file")

	addr    = flag.String("addr", "169.254.169.254", "bind to addr")
	logfile = flag.String("log", "", "log http to file")
	port    = flag.Uint("port", 16925, "listen on port")

	printver = flag.Bool("version", false, "print version")
)

func main() {
	flag.Parse()

	if *printver {
		fmt.Println("finto", finto.Version)
		os.Exit(0)
	}

	logdest, err := prepareLog(*logfile)
	if err != nil {
		panic(err)
	}
	defer logdest.Close()

	config, err := LoadConfig(*fintorc)
	if err != nil {
		panic(err)
	}

	// SharedCredentialsProvider defaults to file=~/.aws/credentials and
	// profile=default when provided zero-value strings
	rs := finto.NewRoleSet(sts.New(session.New(), &aws.Config{
		Credentials: credentials.NewSharedCredentials(
			config.Credentials.File,
			config.Credentials.Profile,
		),
	}))

	for alias, arn := range config.Roles {
		rs.SetRole(alias, arn)
	}

	context, err := finto.InitFintoContext(rs, config.DefaultRole)
	if err != nil {
		fmt.Println("warning: default role not set:", err)
	}

	router := finto.FintoRouter(&context)
	handler := handlers.LoggingHandler(logdest, router)
	err = http.ListenAndServe(fmt.Sprint(*addr, ":", *port), handler)
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

func prepareLog(filename string) (*os.File, error) {
	if filename != "" {
		return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	}

	return os.Stdout, nil
}

# finto

[![Build Status](https://travis-ci.org/threadwaste/finto.svg?branch=master)](https://travis-ci.org/threadwaste/finto)

**finto (-a)** /'finto (-a)/ *agg* **1** posticcio; artificiàle

## Overview

finto is a web server that emulates EC2 instance profile roles on a workstation
through STS's assume role function. It was born as an experiment to ease local
interaction with AWS services in a deeply-federated, role-based environment.
finto ships with a basic API for moving between roles, and handles credentials
caching and expiration.

## Installation

At its simplest:

    go get github.com/threadwaste/finto

## Usage

    Usage of finto:
      -addr="169.254.169.254": bind to addr
      -config="/home/demo/.fintorc": location of config file
      -log="": log http to file
      -port=16925: listen on port

While running, finto provides credentials to EC2 instance profile providers.
This provider is last in the default provider chain of each SDK. For more
information, refer to the official documentation on [EC2 instance profile
roles](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html) and the [standardized credentials interface](https://blogs.aws.amazon.com/security/post/Tx3D6U6WSFGOK2H/A-New-and-Standardized-Way-to-Manage-Credentials-in-the-AWS-SDKs).

Below is sample output of finto serving credentials to the AWS CLI:

    ~
    ❯ aws s3 ls --debug
    <truncated>
    2016-01-03 11:52:01,895 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: env
    2016-01-03 11:52:01,895 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: assume-role
    2016-01-03 11:52:01,895 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: shared-credentials-file
    2016-01-03 11:52:01,896 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: config-file
    2016-01-03 11:52:01,896 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: ec2-credentials-file
    2016-01-03 11:52:01,896 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: boto-config
    2016-01-03 11:52:01,897 - MainThread - botocore.credentials - DEBUG - Looking for credentials via: iam-role
    2016-01-03 11:52:01,902 - MainThread - botocore.vendored.requests.packages.urllib3.connectionpool - INFO - Starting new HTTP connection (1): 169.254.169.254
    2016-01-03 11:52:01,904 - MainThread - botocore.vendored.requests.packages.urllib3.connectionpool - DEBUG - "GET /latest/meta-data/iam/security-credentials/ HTTP/1.1" 200 5
    2016-01-03 11:52:02,259 - MainThread - botocore.vendored.requests.packages.urllib3.connectionpool - DEBUG - "GET /latest/meta-data/iam/security-credentials/example HTTP/1.1" 200 635
    2016-01-03 11:52:02,261 - MainThread - botocore.credentials - INFO - Found credentials from IAM Role: example
    <truncated>
    2016-01-03 11:52:03,282 - MainThread - botocore.hooks - DEBUG - Event after-call.s3.ListBuckets: calling handler <awscli.errorhandler.ErrorHandler object at 0x10483fc90>
    2016-01-03 11:52:03,282 - MainThread - awscli.errorhandler - DEBUG - HTTP Response Code: 200

finto also includes an API for bouncing between available roles. Helper
functions for bash and fish shells are available.

    ~
    ❯ curl 169.254.169.254/roles
    {"roles":["example","example2"]}
    ~
    ❯ curl 169.254.169.254/roles/example
    {"arn":"arn:aws:iam::123456789012:role/example","session_name":"finto-example"}
    ~
    ❯ curl 169.254.169.254/roles/example/credentials
    {
      "AccessKeyId": "<redacted>",
      "Code": "Success",
      "Expiration": "2016-01-03T19:40:30Z",
      "LastUpdated": "2015-07-07T23:06:33Z",
      "SecretAccessKey": "<redacted>",
      "Token": "<redacted>",
      "Type": "AWS-HMAC"
    }
    ~
    ❯ curl 169.254.169.254/latest/meta-data/iam/security-credentials/
    example
    ~
    ❯ curl -XPUT -d'{"alias":"example2"}' 169.254.169.254/roles
    {"active_role":"example2"}
    ~
    ❯ curl 169.254.169.254/latest/meta-data/iam/security-credentials/
    example2

## Configuration

finto uses a JSON configuration file to setup its credentials and the roles it
will serve. It currently uses a shared credentials provider only. Exluding the
credentials file or profile will use the defaults "~/.aws/credentials" and
"default," respectively.

    {
      "credentials": {
        "file": "/home/demo/.finto/credentials",
        "profile": "identity"
      },
      "roles": {
        "example": "arn:aws:iam::123456789012:role/example",
        "example2": "arn:aws:iam::123456789012:role/example2"
      }
      "default_role": "example",
    }

## Running

There are essentially two basic requirements for running finto:

  1. Routing the EC2 meta-data endpoint
  2. Using (or chaining to) the EC2 instance profile provider

The first can be achieved in several ways: interface aliasing, network
redirection, virtual machines, and so on. The wiki contains a couple of basic
examples.

The second is client-dependent. In the case of clients like the AWS CLI, the
user must clear a path to the EC2 instance profile provider. Multiple shared
credentials profiles can still be configured, and accessed with e.g. the
--profile option or AWS_DEFAULT_PROFILE environment variable.

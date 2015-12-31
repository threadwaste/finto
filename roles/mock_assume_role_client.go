package roles

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
)

// A mock client that satisfies the AssumeRoleClient interface. For testing
// purposes.
type MockAssumeRoleClient struct{}

// Return a canned sts.AssumeRoleOutput.
func (c *MockAssumeRoleClient) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	mockId := *input.RoleArn + "-" + *input.RoleSessionName
	expiry := time.Now().Add(15 * time.Minute)
	return &sts.AssumeRoleOutput{
		Credentials: &sts.Credentials{
			AccessKeyId:     aws.String(mockId),
			Expiration:      &expiry,
			SecretAccessKey: aws.String("mock-key"),
			SessionToken:    aws.String("mock-token"),
		},
	}, nil
}

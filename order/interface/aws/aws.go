package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	return region
}
func Dynamodb() *dynamodb.DynamoDB {
	return dynamodb.New(session.Must(
		session.NewSession(&aws.Config{
			Region: aws.String(GetRegion()),
		}),
	))
}

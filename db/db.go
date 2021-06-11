package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db *dynamodb.DynamoDB

func Init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	db = dynamodb.New(sess, &aws.Config{Endpoint: aws.String("http://localhost:8000")})
}

func GetDB() *dynamodb.DynamoDB {
	return db
}

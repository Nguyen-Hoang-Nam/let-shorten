package models

import (
	"fmt"

	"github.com/Nguyen-Hoang-Nam/let-shorten/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Url struct {
	Hash string `json:"Hash"`
	URL  string `json:"URL"`
	TTL  string `json:"TTL"`
}

func (u Url) RemoveUrlByHash(hash string) error {
	db := db.GetDB()

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Hash": {
				S: aws.String(hash),
			},
		},
		TableName: aws.String("Url"),
	}

	_, err := db.DeleteItem(params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (u Url) GetURLByHash(hash string) (*Url, error) {
	db := db.GetDB()

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Hash": {
				S: aws.String(hash),
			},
		},
		TableName: aws.String("Url"),
	}

	res, err := db.GetItem(params)
	if err != nil {
		return nil, err
	}

	var url *Url
	if err := dynamodbattribute.UnmarshalMap(res.Item, &url); err != nil {
		return nil, err
	}

	return url, err
}

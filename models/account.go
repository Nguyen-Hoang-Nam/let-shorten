package models

import (
	"fmt"

	"github.com/Nguyen-Hoang-Nam/let-shorten/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Account struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	ID       string `json:"Id"`
}

func (a Account) CreateAccount(account Account) error {
	db := db.GetDB()

	accountParams := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(account.Email),
			},
			"Password": {
				S: aws.String(account.Password),
			},
			"Id": {
				S: aws.String(account.ID),
			},
		},
		TableName: aws.String("Account"),
	}

	_, err := db.PutItem(accountParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	userParams := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(account.ID),
			},
			"Hashes": {
				L: []*dynamodb.AttributeValue{},
			},
		},
		TableName: aws.String("User"),
	}

	_, err = db.PutItem(userParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (a Account) GetAccount(email string) (*Account, error) {
	db := db.GetDB()

	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName:      aws.String("Account"),
		ConsistentRead: aws.Bool(true),
	}

	res, err := db.GetItem(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(res.Item) == 0 {
		return nil, nil
	}

	var account *Account
	if err = dynamodbattribute.UnmarshalMap(res.Item, &account); err != nil {
		return nil, err
	}

	return account, nil
}

func (a Account) UpdateAccount(email string, password string) error {
	db := db.GetDB()

	params := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName:        aws.String("Account"),
		UpdateExpression: aws.String("SET Password = :newval"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":newval": {
				S: aws.String(password),
			},
		},
	}

	_, err := db.UpdateItem(params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (a Account) DeleteAccount(email string) error {
	db := db.GetDB()

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String("Account"),
	}

	_, err := db.DeleteItem(params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

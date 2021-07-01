package models

import (
	"fmt"

	"github.com/Nguyen-Hoang-Nam/let-shorten/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Hashes []string

// func (cs *Hashes) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
// 	av.SS = make([]*string, 0, len(cs))
// 	for _, v := range cs {
// 		av.SS = append(av.SS, &v)
// 	}
// 	return nil
// }

type User struct {
	ID     string   `json:"Id,omitempty"`
	Hashes []string `json:"Hashes" dynamodbav:"Hashes"`
}

func (u User) GetByID(id string) (*User, error) {
	db := db.GetDB()
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName:      aws.String("User"),
		ConsistentRead: aws.Bool(true),
	}

	res, err := db.GetItem(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var user *User
	if err := dynamodbattribute.UnmarshalMap(res.Item, &user); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (u User) AddURLByID(url Url, id string) error {
	db := db.GetDB()

	params := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName:        aws.String("User"),
		UpdateExpression: aws.String("SET #hs = list_append(#hs, :vals)"),
		ExpressionAttributeNames: map[string]*string{
			"#hs": aws.String("Hashes"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":vals": {
				L: []*dynamodb.AttributeValue{
					{
						S: aws.String(url.Hash),
					},
				},
			},
		},
	}

	_, err := db.UpdateItem(params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	urlParams := &dynamodb.PutItemInput{
		TableName: aws.String("Url"),
		Item: map[string]*dynamodb.AttributeValue{
			"URL": {
				S: aws.String(url.URL),
			},
			"Hash": {
				S: aws.String(url.Hash),
			},
			"TTL": {
				N: aws.String(url.TTL),
			},
		},
	}

	_, err = db.PutItem(urlParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (u User) RemoveURLByID(id, position string) error {
	db := db.GetDB()

	query := fmt.Sprintf("REMOVE Hashes[%s]", position)
	params := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName:        aws.String("User"),
		UpdateExpression: aws.String(query),
	}

	_, err := db.UpdateItem(params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (u User) DeleteUser(id string) error {
	db := db.GetDB()

	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String("User"),
	}
	_, err := db.DeleteItem(params)
	if err != nil {
		fmt.Println(err)
		return err
	}

	findEmail := &dynamodb.QueryInput{
		TableName:              aws.String("Account"),
		KeyConditionExpression: aws.String("Id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {
				S: aws.String(id),
			},
		},
	}
	email, err1 := db.Query(findEmail)
	if err1 != nil {
		return err1
	}

	accountParams := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": email.Items[0]["Email"],
		},
		TableName: aws.String("Account"),
	}
	_, err = db.DeleteItem(accountParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

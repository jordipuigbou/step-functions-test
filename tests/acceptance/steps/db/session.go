package db

import (
	"context"
	"fmt"

	"github.com/TelefonicaTC2Tech/golium"
	"github.com/aws/aws-sdk-go/aws"
	aws_s "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Session struct {
	awsSession     *aws_s.Session
	DynamoDBClient *dynamodb.DynamoDB
	TestTableName  string
}

// const (
// 	stepFunctionPath = "/Users/zc01167/sam-app"
// )

// type executionOut struct {
// 	StatusCode int    `json:"statusCode"`
// 	Body       string `json:"body"`
// }

func (s *Session) SetAwsDynamoClient(ctx context.Context) error {
	awsConfig := &aws.Config{
		Endpoint:   aws.String(golium.Value(ctx, "[CONF:dynamoDBEndpoint]").(string)),
		DisableSSL: aws.Bool(true),
		Region:     aws.String(golium.Value(ctx, "[CONF:awsRegion]").(string)),
	}
	var err error
	if s.awsSession, err = aws_s.NewSession(awsConfig); err != nil {
		return fmt.Errorf("error creating aws session. %v", err)
	}
	s.DynamoDBClient = dynamodb.New(s.awsSession)
	return nil
}

func (s *Session) CreateDynamoDBTable(tableName, indexName string) error {
	logger := GetLogger()
	out, err := s.DynamoDBClient.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(indexName),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(indexName),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
	})
	if err != nil {
		return fmt.Errorf("error creating dynamodb table: %v", err)
	}
	s.TestTableName = tableName
	logger.LogMessage("Created " + tableName + " table with " + indexName + " index")
	logger.LogMessage(fmt.Sprint(*out.TableDescription.ItemCount))
	logger.LogMessage(*out.TableDescription.TableName)
	logger.LogMessage(*out.TableDescription.TableArn)
	logger.LogMessage(
		*out.TableDescription.KeySchema[len(out.TableDescription.KeySchema)-1].AttributeName)
	logger.LogMessage(*out.TableDescription.TableStatus)
	return nil
}

func (s *Session) ListDynamoDBTables() error {
	logger := GetLogger()
	pageNum := 0
	err := s.DynamoDBClient.ListTablesPages(&dynamodb.ListTablesInput{
		Limit: aws.Int64(10),
	},
		func(page *dynamodb.ListTablesOutput, lastPage bool) bool {
			logger.LogMessage(*page.TableNames[pageNum])
			pageNum++
			return pageNum <= 3
		})
	if err != nil {
		return fmt.Errorf("error listing dynamodb tables: %v", err)
	}
	return nil
}

func (s *Session) DeleteDynamoDBTable(tableName string) error {
	_, err := s.DynamoDBClient.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return fmt.Errorf("error deleting "+tableName+" table: %v", err)
	}
	return nil
}

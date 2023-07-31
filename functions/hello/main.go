package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return Response{}, err
	}

	dynamoClient := dynamodb.NewFromConfig(cfg)

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{
				Value: *aws.String("Test"),
			},
		},
	}

	result, err := dynamoClient.GetItem(ctx, getItemInput)
	if err != nil {
		return Response{}, err
	}

	if len(result.Item) == 0 {
		return Response{
			StatusCode: http.StatusNotFound,
			Body:       "Item not found",
		}, nil
	}

	return Response{StatusCode: 200, Body: "Not Implemented"}, nil
}

func main() {
	lambda.Start(Handler)
}

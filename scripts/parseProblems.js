const { DynamoDBClient, ListTablesCommand } = require("@aws-sdk/client-dynamodb");
const AWS = require('aws-sdk');

const dynamodb = AWS.DynamoDB();
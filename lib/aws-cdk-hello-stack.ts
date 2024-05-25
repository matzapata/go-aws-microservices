import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as goLambda from "@aws-cdk/aws-lambda-go-alpha"

export class AwsCdkHelloStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // hello lambda
    const helloLambda = new goLambda.GoFunction(this, 'GoLambda', {
      entry: 'src/hello/main.go', // lambdas works by in theory this should also
    });

    // Define the Lambda function
    const table = new dynamodb.Table(this, 'MyTable', {
      partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
    });
    const dynamoLambda = new goLambda.GoFunction(this, 'DynamoLambda', {
      entry: 'src/dynamo/main.go', // lambdas works by in theory this should also
      environment: {
        TABLE_NAME: table.tableName,
      },
    });
    table.grantReadWriteData(dynamoLambda);


    // Create a new API Gateway REST API
    const api = new apigateway.RestApi(this, 'GoLambdaAPI', {
      restApiName: 'GoLambda API',
      description: 'API for GoLambda',
    });

    // Add hello lambda to the gateway
    const helloLambdaIntegration = new apigateway.LambdaIntegration(helloLambda);
    const helloLambdaResource = api.root.addResource('hello');
    helloLambdaResource.addMethod('GET', helloLambdaIntegration);

    // Add dynamo lambda to the db
    const dynamoLambdaIntegration = new apigateway.LambdaIntegration(dynamoLambda);
    const dynamoLambdaResource = api.root.addResource('dynamo');
    dynamoLambdaResource.addMethod('GET', dynamoLambdaIntegration);
  }
}

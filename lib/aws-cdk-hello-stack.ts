import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as sns from 'aws-cdk-lib/aws-sns';
import * as sns_subscriptions from 'aws-cdk-lib/aws-sns-subscriptions';
import * as goLambda from "@aws-cdk/aws-lambda-go-alpha"

export class AwsCdkHelloStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ========================================
    // ============= api gateway ============== 
    // ========================================

    const api = new apigateway.RestApi(this, 'GoLambdaAPI', {
      restApiName: 'GoLambda API',
      description: 'API for GoLambda',
    });

    // ========================================
    // ============= sns topics =============== 
    // ========================================

    const topic = new sns.Topic(this, 'MyTopic');

    // ========================================
    // ============= hello lambda ============= 
    // ========================================

    const helloLambda = new goLambda.GoFunction(this, 'GoLambda', {
      entry: 'src/hello/main.go', // lambdas works by in theory this should also
    });

    const helloLambdaIntegration = new apigateway.LambdaIntegration(helloLambda);
    const helloLambdaResource = api.root.addResource('hello');
    helloLambdaResource.addMethod('GET', helloLambdaIntegration);

    // =========================================
    // ============= dynamo lambda =============
    // ========================================= 

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

    // grant permissions for lambda to write to table
    table.grantReadWriteData(dynamoLambda);

    const dynamoLambdaIntegration = new apigateway.LambdaIntegration(dynamoLambda);
    const dynamoLambdaResource = api.root.addResource('dynamo');
    dynamoLambdaResource.addMethod('GET', dynamoLambdaIntegration);


    // ============================================
    // ============= publisher lambda =============
    // ============================================

    // Define the first Lambda function (Publisher)
    const publisherLambda = new goLambda.GoFunction(this, 'MyPublisherLambda', {
      entry: 'src/publisher/main.go', // lambdas works by in theory this should also
      environment: {
        TOPIC_ARN: topic.topicArn,
      },
    });

    // Grant the publisher Lambda function publish permissions to the SNS topic
    topic.grantPublish(publisherLambda);

    const publisherLambdaIntegration = new apigateway.LambdaIntegration(publisherLambda);
    const publisherLambdaResource = api.root.addResource('publisher');
    publisherLambdaResource.addMethod('GET', publisherLambdaIntegration);


    // =============================================
    // ============= subscriber lambda =============
    // =============================================

    // Define the first Lambda function (Subscriber)
    const subscriberLambda = new goLambda.GoFunction(this, 'MySubscriberLambda', {
      entry: 'src/subscriber/main.go',
      environment: {
        TABLE_NAME: table.tableName,
      },
    });

    // Subscribe the subscriber Lambda function to the SNS topic
    topic.addSubscription(new sns_subscriptions.LambdaSubscription(subscriberLambda));

    // Grant the subscriber Lambda function read/write permissions to the DynamoDB table
    table.grantReadWriteData(subscriberLambda);

  }
}

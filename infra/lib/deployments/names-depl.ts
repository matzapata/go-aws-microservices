import * as goLambda from "@aws-cdk/aws-lambda-go-alpha"
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from "aws-cdk-lib/aws-dynamodb"
import * as sqs from 'aws-cdk-lib/aws-sqs';
import { DeploymentProps, SERVICES_BASE_PATH } from "./deployment";
import * as cdk from 'aws-cdk-lib';
import * as subscriptions from 'aws-cdk-lib/aws-sns-subscriptions';
import * as eventSources from 'aws-cdk-lib/aws-lambda-event-sources';


export class NamesDeployment extends Construct {
    private readonly name: string = "names"
    private readonly props: DeploymentProps

    constructor(scope: Construct, id: string, props: DeploymentProps) {
        super(scope, id)
        this.props = props

        // create resources
        const table = this.createDynamoTable()
        const queue = this.createSqsQueue()

        // subscribe queue to topics
        props.topics["CREATE_NAME_TOPIC"].addSubscription(new subscriptions.SqsSubscription(queue))

        // create rest lambda ========================================
        const restLambda = this.createRestLambda({
            environment: {
                TABLE_NAME: table.tableName,
                QUEUE_URL: queue.queueUrl
            }
        })
        table.grantReadWriteData(restLambda)
        this.addGatewayEventSource(restLambda, props.apigateway)


        // create events lambda ========================================
        const eventsLambda = this.createEventsLambda({
            environment: {
                TABLE_NAME: table.tableName,
                QUEUE_URL: queue.queueUrl
            }
        })
        table.grantReadWriteData(eventsLambda)
        queue.grantConsumeMessages(eventsLambda);
        this.addSqsEventSource(eventsLambda, queue)
    }

    private createDynamoTable() {
        return new dynamodb.Table(this, `${this.name}-table`, {
            partitionKey: { name: 'id', type: dynamodb.AttributeType.STRING },
            billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
        });
    }

    private createSqsQueue() {
        return new sqs.Queue(this, `${this.name}-queue`, {
            visibilityTimeout: cdk.Duration.seconds(30),
        });
    }

    private createRestLambda(props: { environment: { [key: string]: string } }): lambda.IFunction {
        const lambda = new goLambda.GoFunction(this, `${this.name}-rest-lambda`, {
            entry: path.join(SERVICES_BASE_PATH, "names", "cmd", "api"),
            environment: props.environment,
        });

        return lambda
    }

    private createEventsLambda(props: { environment: { [key: string]: string } }): lambda.IFunction {
        const lambda = new goLambda.GoFunction(this, `${this.name}-events-lambda`, {
            entry: path.join(SERVICES_BASE_PATH, "names", "cmd", "events"),
            environment: props.environment,
        });

        return lambda
    }

    private addGatewayEventSource(lambda: lambda.IFunction, gateway: apigateway.IRestApi) {
        const dynamoLambdaIntegration = new apigateway.LambdaIntegration(lambda);

        const dynamoLambdaResource = gateway.root.addResource(this.name);
        dynamoLambdaResource.addMethod('GET', dynamoLambdaIntegration);
    }

    private addSqsEventSource(lambda: lambda.IFunction, queue: sqs.IQueue) {
        lambda.addEventSource(new eventSources.SqsEventSource(queue, {
            batchSize: 1 // Process one message per Lambda invocation
        }));
    }
}
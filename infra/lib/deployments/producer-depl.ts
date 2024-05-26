import * as goLambda from "@aws-cdk/aws-lambda-go-alpha"
import { IFunction } from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as cdk from 'aws-cdk-lib';
import { DeploymentProps, SERVICES_BASE_PATH } from "./deployment";



export class ProducerDeployment extends Construct {
    name = "producer";
    props: DeploymentProps

    constructor(scope: Construct, id: string, props: DeploymentProps) {
        super(scope, id)
        this.props = props

        const lambda = this.createLambda({
            environment: {
                TOPIC_ARN: props.topics.CREATE_NAME_TOPIC.topicArn
            }
        })

        // grant permissions
        props.topics.CREATE_NAME_TOPIC.grantPublish(lambda);

        // add event sources
        this.addGatewayEventSource(lambda)
    }

    private createLambda(props: { environment: { [key: string]: string } }): IFunction {
        return new goLambda.GoFunction(this, `${this.name}-lambda`, {
            entry: path.join(SERVICES_BASE_PATH, "producer", "cmd", "api"),
            environment: props.environment
        });
    }

    private addGatewayEventSource(lambda: IFunction) {
        const publisherLambdaIntegration = new apigateway.LambdaIntegration(lambda);
        const publisherLambdaResource = this.props.apigateway.root.addResource(this.name);
        publisherLambdaResource.addMethod('GET', publisherLambdaIntegration);
    }

}
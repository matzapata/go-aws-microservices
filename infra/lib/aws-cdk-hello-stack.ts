import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as sns from 'aws-cdk-lib/aws-sns';
import * as eventtargets from 'aws-cdk-lib/aws-events-targets';
import * as eventbridge from 'aws-cdk-lib/aws-events';
import * as lambdatargets from 'aws-cdk-lib/aws-lambda-event-sources';
import * as goLambda from "@aws-cdk/aws-lambda-go-alpha"
import * as path from "path"
import { HelloDeployment } from './deployments/hello-depl';
import { NamesDeployment } from './deployments/names-depl';
import { ProducerDeployment } from './deployments/producer-depl';
import { DeploymentProps } from './deployments/deployment';

const SERVICES_BASE_PATH = path.join(__dirname, "..", "..", "services")
console.log("Services path", SERVICES_BASE_PATH)

export class AwsCdkHelloStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // ============================================
    // ============= Shared Resources =============
    // ============================================

    // create gateway
    const api = new apigateway.RestApi(this, 'GoLambdaAPI', {
      restApiName: 'GoLambda API',
      description: 'Golang serverless microservices',
    });

    // create topics for microservices
    const CREATE_NAME_TOPIC = new sns.Topic(this, 'ImportantEventsTopic', {
      displayName: 'Important Events Topic'
    });

    // shared resources
    const shared: DeploymentProps = {
      apigateway: api, 
      topics: {
        CREATE_NAME_TOPIC
      }
    }

    // ================================================
    // ============= Services Deployments =============
    // ================================================

    // Simple hello endpoint
    new HelloDeployment(this, "HelloDeployment", shared)

    // Names microservice
    new NamesDeployment(this, "NamesDeployment", shared)

    // Example events producer
    new ProducerDeployment(this, "ProducerDeployment", shared)


  }
}

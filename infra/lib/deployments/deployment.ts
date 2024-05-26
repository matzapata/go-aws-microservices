import { RestApi } from "aws-cdk-lib/aws-apigateway";
import { EventBus } from "aws-cdk-lib/aws-events";
import { Topic } from "aws-cdk-lib/aws-sns";
import { Queue } from "aws-cdk-lib/aws-sqs";
import * as path from "path"


export interface DeploymentProps {
    apigateway: RestApi,
    topics: { [key: string]: Topic }
}


export const SERVICES_BASE_PATH = path.join(__dirname, "..", "..", "..", "services")
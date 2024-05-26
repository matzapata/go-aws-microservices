import * as goLambda from "@aws-cdk/aws-lambda-go-alpha"
import { IFunction } from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import * as apigateway from 'aws-cdk-lib/aws-apigateway';
import { DeploymentProps, SERVICES_BASE_PATH } from "./deployment";

export class HelloDeployment extends Construct {
    private readonly name: string = "hello"
    private readonly props: DeploymentProps

    constructor(scope: Construct, id: string, props: DeploymentProps) {
        super(scope, id)
        this.props = props

        const lambda = this.createLambda()
        this.createGatewayRules(lambda)

    }

    private createLambda(): IFunction {
        return new goLambda.GoFunction(this, `${this.name}-lambda`, {
            entry: path.join(SERVICES_BASE_PATH, "hello", "cmd", "api"),
        });
    }

    private createGatewayRules(lambda: IFunction) {
        const helloLambdaIntegration = new apigateway.LambdaIntegration(lambda);

        const helloLambdaResource = this.props.apigateway.root.addResource(this.name);
        helloLambdaResource.addMethod('GET', helloLambdaIntegration);
    }

}
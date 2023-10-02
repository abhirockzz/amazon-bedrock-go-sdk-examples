package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const functionDir = "../function"

type BedrockLambdaGolangStackProps struct {
	awscdk.StackProps
}

func NewBedrockLambdaGolangStack(scope constructs.Construct, id string, props *BedrockLambdaGolangStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("bedrock-basic-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime: awslambda.Runtime_GO_1_X(),
			Entry:   jsii.String(functionDir),
			Timeout: awscdk.Duration_Seconds(jsii.Number(15)),
		})

	function.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("bedrock:*"),
		Effect:    awsiam.Effect_ALLOW,
		Resources: jsii.Strings("*"),
	}))

	awscdk.NewCfnOutput(stack, jsii.String("function-name"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("function-name"),
			Value:      function.FunctionName()})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewBedrockLambdaGolangStack(app, "BedrockLambdaGolangStack", &BedrockLambdaGolangStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}

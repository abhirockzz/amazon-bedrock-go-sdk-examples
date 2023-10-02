package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const functionDir = "../function"

type BedrockStableDiffusionLambdaGolangStackProps struct {
	awscdk.StackProps
}

func NewBedrockStableDiffusionLambdaGolangStack(scope constructs.Construct, id string, props *BedrockStableDiffusionLambdaGolangStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	targetBucket := awss3.NewBucket(stack, jsii.String("target-s3-bucket"), &awss3.BucketProps{
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("bedrock-imagegen-s3"),
		&awscdklambdagoalpha.GoFunctionProps{
			Runtime:     awslambda.Runtime_GO_1_X(),
			Entry:       jsii.String(functionDir),
			Timeout:     awscdk.Duration_Seconds(jsii.Number(15)),
			Environment: &map[string]*string{"BUCKET_NAME": targetBucket.BucketName()},
		})

	function.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("bedrock:*"),
		Effect:    awsiam.Effect_ALLOW,
		Resources: jsii.Strings("*"),
	}))

	// Grant write permissions to the Lambda function for the S3 bucket via the function
	function.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("s3:PutObject"),
		Resources: jsii.Strings(*targetBucket.BucketArn() + "/*"),
	}))

	awscdk.NewCfnOutput(stack, jsii.String("function-name"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("function-name"),
			Value:      function.FunctionName()})

	awscdk.NewCfnOutput(stack, jsii.String("target-s3-bucket-name"),
		&awscdk.CfnOutputProps{
			ExportName: jsii.String("target-s3-bucket-name"),
			Value:      targetBucket.BucketName()})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewBedrockStableDiffusionLambdaGolangStack(app, "BedrockLambdaGolangStack", &BedrockStableDiffusionLambdaGolangStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}

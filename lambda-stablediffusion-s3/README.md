# AWS Lambda function to invoke Bedrock stable diffusion for image generation

To deploy function:

```shell
cd cdk

export DOCKER_DEFAULT_PLATFORM=linux/amd64
cdk deploy
```

Invoke function:

```shell
export FUNCTION_NAME=<check CDK output>

aws lambda invoke --function-name $FUNCTION_NAME --payload '{"input": "Sri lanka tea plantation."}' --cli-binary-format raw-in-base64-out /dev/stdout
```

Check S3:

```shell
aws s3 ls

aws s3 cp <file name in output> .
```

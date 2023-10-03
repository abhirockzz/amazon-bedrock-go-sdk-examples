# Basic AWS Lambda function to invoke Bedrock

To deploy function:

```shell
cd cdk

export DOCKER_DEFAULT_PLATFORM=linux/amd64
cdk deploy
```

Invoke function:

```shell
export FUNCTION_NAME=<check CDK output>

aws lambda invoke --function-name $FUNCTION_NAME --payload '{"input": "Convert this into SQL statement: Pick top three rows from the customer table based on their total expenditure"}' --cli-binary-format raw-in-base64-out output.txt

aws lambda invoke --function-name $FUNCTION_NAME --payload '{"input": "Suggest an outline for a blog post based on a title.\nTitle: How I put the pro in prompt engineering"}' --cli-binary-format raw-in-base64-out output.txt
```

Verify:

```shell
cat output.txt
```
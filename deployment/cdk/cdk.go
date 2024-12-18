package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	flybuysfunction "gitlab.com/flybuys/app/golang-cdk-shared/middleware/function"
	"gitlab.com/flybuys/app/golang-cdk-shared/middleware/tagging"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	_, _ = flybuysfunction.CreateApigatewayAndFunctions(stack)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	// tags
	jsonStr := tagging.GetTagsJsonString()
	tags, _ := tagging.ParseTagsFromJson(jsonStr)

	NewCdkStack(app, "ApiGatewayCdkStack", &CdkStackProps{
		awscdk.StackProps{
			Env:       env(),
			StackName: jsii.String(os.Getenv("STACK_NAME")),
			Tags:      &tags,
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT_ID")),
		Region:  jsii.String(os.Getenv("AWS_REGION")),
	}
}

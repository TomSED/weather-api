package function

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

const ParameterStoreLayerArn = "arn:aws:lambda:ap-southeast-2:665172237481:layer:AWS-Parameters-and-Secrets-Lambda-Extension:12"

func GetDatadogLayer(stack awscdk.Stack) awslambda.ILayerVersion {
	datadogExtensionLayerArn := fmt.Sprintf("arn:aws:lambda:%s:464622532012:layer:Datadog-Extension:63", *stack.Region())
	return awslambda.LayerVersion_FromLayerVersionArn(stack, jsii.String("DatadogLayer"), jsii.String(datadogExtensionLayerArn))
}

func GetParameterStoreLayer(stack awscdk.Stack) awslambda.ILayerVersion {
	return awslambda.LayerVersion_FromLayerVersionArn(stack, jsii.String("ParameterStoreLayer"), jsii.String(ParameterStoreLayerArn))
}

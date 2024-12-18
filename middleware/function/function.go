package function

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

type FunctionRestApiProps struct {
	RequestMethod string
	Api           awsapigateway.RestApi
	ApiPath       string
}

type FunctionProps struct {
	StackName        string
	FunctionDomain   string
	FunctionID       string
	FunctionCodePath string
	// OverrideLogicalID replaces the default logical ID with FunctionID. This is required for openapi apigateway integration
	OverrideLogicalID bool
	RestApiProps      *FunctionRestApiProps
	VpcId             string
	FunctionTimeOut   awscdk.Duration
	Layers            []awslambda.ILayerVersion
	DatadogConfig     DatadogConfig
	Environment       map[string]string
	privateSubnets    []awsec2.ISubnet
}

type DatadogConfig struct {
	ApiKeySecretArn string
	ServiceTag      string
	EnvTag          string
}

func setDatadogEnv(env map[string]string, config DatadogConfig) map[string]string {
	env["DD_API_KEY_SECRET_ARN"] = config.ApiKeySecretArn
	env["DD_SITE"] = "datadoghq.com"
	env["DD_SERVERLESS_LOGS_ENABLED"] = "true"
	env["DD_SERVICE"] = config.ServiceTag
	env["DD_ENV"] = config.EnvTag

	commitTag := os.Getenv("CI_COMMIT_TAG")
	commitShortSHA := os.Getenv("CI_COMMIT_SHORT_SHA")

	if commitTag != "" {
		env["DD_TAGS"] = fmt.Sprintf("version:%s", commitTag)
	} else {
		env["DD_TAGS"] = fmt.Sprintf("version:%s", commitShortSHA)
	}

	return env
}

func transformEnvironment(env map[string]string) *map[string]*string {
	out := map[string]*string{}
	for key, value := range env {
		out[key] = jsii.String(value)
	}
	return &out
}

func NewFunction(stack awscdk.Stack, props *FunctionProps) awslambda.Function {
	vpc := awsec2.Vpc_FromLookup(stack, jsii.String(fmt.Sprintf("ExistingVPC-%v", props.FunctionID)), &awsec2.VpcLookupOptions{
		VpcId: jsii.String(props.VpcId),
	})

	functionName := fmt.Sprintf("%s-%s", props.StackName, props.FunctionID)

	if len(functionName) > 64 {
		functionName = fmt.Sprintf("%s-%s-%s", props.DatadogConfig.EnvTag, props.FunctionDomain, props.FunctionID)
	}

	if len(functionName) > 64 {
		log.Fatalf("%s name too long", functionName)
	}

	lambdaRole := awsiam.NewRole(stack, jsii.String(fmt.Sprintf("%v-LambdaRole", functionName)), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaVPCAccessExecutionRole")),
		},
	})

	// Create the policy statement for Secrets Manager access for datadog api key
	policyStatementSecret := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("secretsmanager:GetSecretValue"),
		Resources: jsii.Strings(props.DatadogConfig.ApiKeySecretArn),
	})
	lambdaRole.AddToPolicy(policyStatementSecret)

	// Define the policy statement
	policyStatementParameterStore := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("ssm:GetParameter", "ssm:GetParameters", "ssm:GetParametersByPath"),
		Resources: jsii.Strings("arn:aws:ssm:*:*:parameter/mobile/bff/*", "arn:aws:ssm:*:*:parameter/aws/reference/secretsmanager/*"),
	})
	lambdaRole.AddToPolicy(policyStatementParameterStore)

	// Define the policy statement
	policyStatementDenyCloudwatch := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("logs:CreateLogGroup", "logs:CreateLogStream", "logs:PutLogEvents"),
		Resources: jsii.Strings("arn:aws:logs:*:*:*"),
		Effect:    awsiam.Effect_DENY,
	})
	lambdaRole.AddToPolicy(policyStatementDenyCloudwatch)

	lambdaFunction := awslambda.NewFunction(stack, jsii.String(props.FunctionID), &awslambda.FunctionProps{
		FunctionName: jsii.String(functionName),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String(props.FunctionCodePath), nil),
		Vpc:          vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			Subnets: &props.privateSubnets,
		},
		Timeout:     props.FunctionTimeOut,
		Layers:      &props.Layers,
		Environment: transformEnvironment(setDatadogEnv(props.Environment, props.DatadogConfig)),
		Role:        lambdaRole,
	})

	// For attaching RestApi from LambdaFunction directly
	if props.RestApiProps != nil {
		newResource := props.RestApiProps.Api.Root().ResourceForPath(jsii.String(props.RestApiProps.ApiPath))
		newResource.AddMethod(jsii.String(props.RestApiProps.RequestMethod), awsapigateway.NewLambdaIntegration(lambdaFunction, nil), &awsapigateway.MethodOptions{
			ApiKeyRequired: jsii.Bool(true),
		})
	}

	if props.OverrideLogicalID {
		lambdaFunction.Node().DefaultChild().(awscdk.CfnResource).OverrideLogicalId(jsii.String(props.FunctionID))
	}
	return lambdaFunction
}

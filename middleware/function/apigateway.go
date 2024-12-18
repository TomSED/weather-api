package function

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awswafv2"
	"github.com/aws/jsii-runtime-go"
)

type SpecRestApiGatewayProps struct {
	Env                         string
	StackName                   string
	RestApiID                   string
	StageName                   string
	FunctionDomain              string
	CustomDomainName            string
	DomainNameAliasHostedZoneId string
	DomainNameAliasTarget       string
	DefinitionFilePath          string
}

// NewSpecRestApiGateway lets you define a RestApi DefinitionBody using an openapi.yaml document
func NewSpecRestApiGateway(stack awscdk.Stack, props *SpecRestApiGatewayProps) awsapigateway.SpecRestApi {
	var assetLocation *string
	if len(props.Env) == 0 || props.Env == "local" {
		assetLocation = jsii.String(props.DefinitionFilePath)
	} else {
		assetLocation = awss3assets.NewAsset(stack, jsii.String("OpenAPIDefinition"), &awss3assets.AssetProps{
			Path: jsii.String(props.DefinitionFilePath),
		}).S3ObjectUrl()
	}

	data := awscdk.Fn_Transform(jsii.String("AWS::Include"), &map[string]interface{}{
		"Location": assetLocation,
	})

	api := awsapigateway.NewSpecRestApi(stack, jsii.String(props.RestApiID), &awsapigateway.SpecRestApiProps{
		ApiDefinition: awsapigateway.AssetApiDefinition_FromInline(data),
		RestApiName:   jsii.String(fmt.Sprintf("%s-%s", props.StackName, props.RestApiID)),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(props.StageName),
		},
	})
	domainName := awsapigateway.DomainName_FromDomainNameAttributes(stack, jsii.String("customDomainName"), &awsapigateway.DomainNameAttributes{
		DomainName:                  jsii.String(props.CustomDomainName),
		DomainNameAliasHostedZoneId: jsii.String(props.DomainNameAliasHostedZoneId),
		DomainNameAliasTarget:       jsii.String(props.DomainNameAliasTarget),
	})
	awsapigateway.NewBasePathMapping(stack, jsii.String("BasePathMapping"), &awsapigateway.BasePathMappingProps{
		DomainName: domainName,
		RestApi:    api,
		BasePath:   jsii.String(props.FunctionDomain),
	})

	return api
}

type LambdaExecutionRoleProps struct {
	RoleID string
	// OverrideLogicalID replaces the default logical ID with RoleID. This is required for openapi apigateway integration
	OverrideLogicalID bool
	Lambdas           []awslambda.Function
}

func NewSpecRestApiLambdaExecutionRole(stack awscdk.Stack, props *LambdaExecutionRoleProps) awsiam.Role {
	lambdaResources := []*string{}
	for _, lambda := range props.Lambdas {
		lambdaResources = append(lambdaResources, jsii.String(*lambda.FunctionArn()))
	}

	lambdaExecutionRole := awsiam.NewRole(stack, jsii.String(props.RoleID), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), nil),
		InlinePolicies: &map[string]awsiam.PolicyDocument{
			"ApiGatewayExecutionPolicy": awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
				Statements: &[]awsiam.PolicyStatement{
					awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
						Actions: &[]*string{
							jsii.String("lambda:*"),
						},
						Effect:    awsiam.Effect_ALLOW,
						Resources: &lambdaResources,
					}),
				},
			}),
		},
	})

	if props.OverrideLogicalID {
		lambdaExecutionRole.Node().DefaultChild().(awscdk.CfnResource).OverrideLogicalId(jsii.String(props.RoleID))
	}

	return lambdaExecutionRole
}

type RestApiGatewayProps struct {
	StackName                   string
	RestApiID                   string
	StageName                   string
	FunctionDomain              string
	CustomDomainName            string
	DomainNameAliasHostedZoneId string
	DomainNameAliasTarget       string
	ApiKeyID                    string
	WafArn                      string
}

func NewRestApiGateway(stack awscdk.Stack, props *RestApiGatewayProps) awsapigateway.RestApi {
	api := awsapigateway.NewRestApi(stack, jsii.String(props.RestApiID), &awsapigateway.RestApiProps{
		RestApiName: jsii.String(fmt.Sprintf("%s-%s", props.StackName, props.RestApiID)),
		DeployOptions: &awsapigateway.StageOptions{
			StageName: jsii.String(props.StageName),
		},
	})

	// associate existing WAF by ARN
	_ = awswafv2.NewCfnWebACLAssociation(stack, jsii.String(fmt.Sprintf(("%s-waf-association"), props.RestApiID)), &awswafv2.CfnWebACLAssociationProps{
		ResourceArn: api.DeploymentStage().StageArn(),
		WebAclArn:   jsii.String(props.WafArn),
	})

	domainName := awsapigateway.DomainName_FromDomainNameAttributes(stack, jsii.String("customDomainName"), &awsapigateway.DomainNameAttributes{
		DomainName:                  jsii.String(props.CustomDomainName),
		DomainNameAliasHostedZoneId: jsii.String(props.DomainNameAliasHostedZoneId),
		DomainNameAliasTarget:       jsii.String(props.DomainNameAliasTarget),
	})
	awsapigateway.NewBasePathMapping(stack, jsii.String("BasePathMapping"), &awsapigateway.BasePathMappingProps{
		DomainName: domainName,
		RestApi:    api,
		BasePath:   jsii.String(props.FunctionDomain),
	})

	// Import existing API Key
	apiKey := awsapigateway.ApiKey_FromApiKeyId(stack, jsii.String("ImportedApiKey"), jsii.String(props.ApiKeyID))

	// Attach the API Key to a Usage Plan
	usagePlan := awsapigateway.NewUsagePlan(stack, jsii.String("UsagePlan"), &awsapigateway.UsagePlanProps{
		Name: jsii.String(fmt.Sprintf("%s-%s-usage-plan", props.StackName, props.RestApiID)),
	})
	usagePlan.AddApiKey(apiKey, nil)
	usagePlan.AddApiStage(&awsapigateway.UsagePlanPerApiStage{
		Api:   api,
		Stage: api.DeploymentStage(),
	})

	return api
}

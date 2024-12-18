package function

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

func getPrivateSubnets(vpcId string, region string) []string {
	// Initialize a session in the desired region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	// Create an EC2 service client
	svc := ec2.New(sess)

	// Describe subnets with filters
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcId)},
			},
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String("Private*")},
			},
		},
	}

	result, err := svc.DescribeSubnets(input)
	if err != nil {
		log.Fatalf("failed to describe subnets: %v", err)
	}

	var ret []string
	for _, subnet := range result.Subnets {
		ret = append(ret, *subnet.SubnetId)
	}
	return ret
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func checkDotEnv(dotEnvFilePath string) {
	if _, err := os.Stat(dotEnvFilePath); os.IsNotExist(err) {
		// If .env does not exist, copy from .env.template
		err := copyFile("../../.env.template", dotEnvFilePath)
		if err != nil {
			log.Fatalf("Failed to copy env.template file: %v ",err)
		}
		fmt.Println(".env file created from .env.template")
	} else {
		fmt.Println(".env file already exists")
	}
}

func getParameters(paramName, region string) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}

	// Create an SSM client
	svc := ssm.New(sess)

	// Get the parameter value
	result, err := svc.GetParameter(&ssm.GetParameterInput{
		Name: aws.String(paramName),
	})
	if err != nil {
		log.Fatalf("unable to get parameter, %v", err)
	}

	return *result.Parameter.Value
}

func GetEnvFromDotEnv(dotEnvFilePath string, env string, region string) map[string]string {
	checkDotEnv(dotEnvFilePath)
	file, err := os.Open(dotEnvFilePath)
	if err != nil {
		log.Fatalf("Failed to open dotEnvFile: %v ",err)
	}
	defer file.Close()

	envMap := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 && env == "local" {
				envMap[parts[0]] = parts[1]
			} else if env != "local" {
				envMap[parts[0]] = getParameters(fmt.Sprintf("/mobile/bff/%s/%s", env, parts[0]), region)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to scan ENV dot file: %v ",err)
	}
	envMap["ENVIRONMENT"] = env
	return envMap
}

func CreateApigatewayAndFunctions(stack awscdk.Stack) (awsapigateway.RestApi, map[string]awslambda.Function) {
	vpcId := os.Getenv("VPC_ID")
	functionDomain := os.Getenv("FUNCTION_DOMAIN")
	stageName := os.Getenv("STAGE_DOMAIN")
	apiName := os.Getenv("API_NAME")
	stackName := os.Getenv("STACK_NAME")
	env := os.Getenv("ENVIRONMENT")

	customDomainName := awscdk.Fn_ImportValue(jsii.String(fmt.Sprintf("CustomDomainNameOutput-%s", env)))
	domainNameAliasHostedZoneId := awscdk.Fn_ImportValue(jsii.String(fmt.Sprintf("CustomDomainNameAliasHostedZoneIdOutput-%s", env)))
	domainNameAliasTarget := awscdk.Fn_ImportValue(jsii.String(fmt.Sprintf("CustomDomainNameAliasDomainNameOutput-%s", env)))
	datadogApiKeySecretArn := awscdk.Fn_ImportValue(jsii.String(fmt.Sprintf("DatadogApiKeySecretArnOutput-%s", env)))
	apiKeyID := awscdk.Fn_ImportValue(jsii.String(fmt.Sprintf("MobileBffApiKeyIdOutput-%s", env)))
	webAclArn := awscdk.Fn_ImportValue(jsii.String(fmt.Sprintf("MobileBffWafWebAclArn-%s", env)))

	api := NewRestApiGateway(stack, &RestApiGatewayProps{
		StackName:                   stackName,
		RestApiID:                   apiName,
		StageName:                   stageName,
		FunctionDomain:              functionDomain,
		CustomDomainName:            *customDomainName,
		DomainNameAliasHostedZoneId: *domainNameAliasHostedZoneId,
		DomainNameAliasTarget:       *domainNameAliasTarget,
		ApiKeyID:                    *apiKeyID,
		WafArn:                      *webAclArn,
	})

	lambdaEnvs := GetEnvFromDotEnv("../../.env", env, *stack.Region())

	functions := GenerateLambdasByOpenAPITemplatePath(stack, &GenerateLambdasProps{
		OpenAPITemplatePath: "../../api/openapi.yaml",
		StackName:           stackName,
		FunctionDomain:      functionDomain,
		Api:                 api,
		VpcId:               vpcId,
		FunctionTimeOut:     awscdk.Duration_Seconds(jsii.Number(30)),
		DatadogConfig: DatadogConfig{
			ApiKeySecretArn: *datadogApiKeySecretArn,
			ServiceTag:      os.Getenv("DATADOG_SERVICE_TAG"),
			EnvTag:          env,
		},
		Environment: lambdaEnvs,
	})

	return api, functions
}

package function

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"

	"github.com/go-openapi/loads"
)

type GenerateLambdasProps struct {
	OpenAPITemplatePath string
	StackName           string
	FunctionDomain      string
	Api                 awsapigateway.RestApi
	VpcId               string
	FunctionTimeOut     awscdk.Duration
	DatadogConfig       DatadogConfig
	Environment         map[string]string
}

func convertEndpoint(endpoint string) string {
	// Replace slashes with hyphens
	endpoint = strings.ReplaceAll(endpoint, "/", "-")
	// Remove curly braces
	endpoint = strings.ReplaceAll(endpoint, "{", "")
	endpoint = strings.ReplaceAll(endpoint, "}", "")
	return endpoint
}

func GenerateLambdasByOpenAPITemplatePath(stack awscdk.Stack, props *GenerateLambdasProps) map[string]awslambda.Function {
	doc, err := loads.Spec(props.OpenAPITemplatePath)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI template: %v ", err)
	}

	// Access the spec
	spec := doc.Spec()

	returnFunctions := make(map[string]awslambda.Function)

	var privateSubnets []awsec2.ISubnet
	if os.Getenv("ENVIRONMENT") == "local" {
		privateSubnets = []awsec2.ISubnet{}
	} else {
		for _, subnetId := range getPrivateSubnets(props.VpcId, *stack.Region()) {
			subnet := awsec2.Subnet_FromSubnetId(stack, jsii.String(fmt.Sprintf("ImportedSubnet-%s", subnetId)), jsii.String(subnetId))
			privateSubnets = append(privateSubnets, subnet)
		}
	}

	layers := *new([]awslambda.ILayerVersion)
	if os.Getenv("ENVIRONMENT") == "local" {
		layers = nil
	} else {
		datadogLayer := GetDatadogLayer(stack)
		parameterStoreLayer := GetParameterStoreLayer(stack)
		layers = []awslambda.ILayerVersion{datadogLayer, parameterStoreLayer}
	}

	// Iterate over paths
	for path, pathItem := range spec.Paths.Paths {
		fmt.Printf("Path: %s\n", path)
		v := reflect.ValueOf(pathItem.PathItemProps)
		for i := 0; i < v.NumField(); i++ {
			method := v.Type().Field(i).Name
			if !v.Field(i).IsNil() {
				fmt.Printf("  Method: %s, Path: %s\n", method, path)
				var functionID string
				if path == "/" {
					functionID = fmt.Sprintf("%s-root", strings.ToLower(method))
				} else {
					functionID = fmt.Sprintf("%s%s", strings.ToLower(method), convertEndpoint(path))
				}

				this := NewFunction(stack, &FunctionProps{
					StackName:         props.StackName,
					FunctionDomain:    props.FunctionDomain,
					FunctionID:        functionID,
					FunctionCodePath:  fmt.Sprintf("../../dist/%s", functionID),
					OverrideLogicalID: false,
					RestApiProps:      &FunctionRestApiProps{RequestMethod: strings.ToUpper(method), Api: props.Api, ApiPath: path},
					VpcId:             props.VpcId,
					FunctionTimeOut:   props.FunctionTimeOut,
					Layers:            layers,
					DatadogConfig:     props.DatadogConfig,
					Environment:       props.Environment,
					privateSubnets:    privateSubnets,
				})

				returnFunctions[functionID] = this
			}
		}
	}

	return returnFunctions
}

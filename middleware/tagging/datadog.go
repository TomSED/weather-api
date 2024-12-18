// Create or update service definition using schema v2-2 returns "CREATED" response

package tagging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

func addRepoLinks(config Config, serviceDefinition *datadogV2.ServiceDefinitionV2Dot2) {
	for i := range serviceDefinition.Links {
		if serviceDefinition.Links[i].Name == config.RepositoryName {
			serviceDefinition.Links[i].Type = "repo"
			serviceDefinition.Links[i].Url = config.RepositoryUrl
			serviceDefinition.Links[i].Provider = datadog.PtrString("Gitlab")
			return
		}
	}

	serviceDefinition.Links = append(serviceDefinition.Links, datadogV2.ServiceDefinitionV2Dot2Link{
		Name:     config.RepositoryName,
		Type:     "repo",
		Provider: datadog.PtrString("Gitlab"),
		Url:      config.RepositoryUrl,
	})
}

func addTags(leanixData FactSheet, serviceDefinition *datadogV2.ServiceDefinitionV2Dot2) {
	serviceDefinition.Tags = append(serviceDefinition.Tags, fmt.Sprintf("asset-id:%s", leanixData.ID),
		fmt.Sprintf("partner-name:%s", leanixData.Partner),
		fmt.Sprintf("data-classification:%s", leanixData.DataClassification))
}

func createDatadogServiceDefinition(config Config, leanixData FactSheet) {
	os.Setenv("DD_API_KEY", config.DataDogApiKey)
	os.Setenv("DD_APP_KEY", config.DataDogAppKey)
	os.Setenv("DD_SITE", "datadoghq.com")

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()

	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewServiceDefinitionApi(apiClient)
	svcDefinition := datadogV2.NewServiceDefinitionV2Dot2(config.DatadogServiceName, datadogV2.SERVICEDEFINITIONV2DOT2VERSION_V2_2)
	resp, r, err := api.GetServiceDefinition(ctx, config.DatadogServiceName, *datadogV2.NewGetServiceDefinitionOptionalParameters().WithSchemaVersion(datadogV2.SERVICEDEFINITIONSCHEMAVERSIONS_V2_2))

	if err != nil {
		if r.StatusCode == 404 {
			fmt.Println("Service definition not found:", config.DatadogServiceName)
		} else {
			fmt.Printf("Error when calling `ServiceDefinitionApi.GetServiceDefinition`: %v\n", err)
			log.Fatalf("Full HTTP response: %v\n", r)
		}
	} else {
		svcDefinition = resp.Data.Attributes.Schema.ServiceDefinitionV2Dot2
	}

	ddType, err := datadogV2.NewServiceDefinitionV2Dot2TypeFromValue(config.DatadogServiceType)
	if err != nil {
		log.Fatalf("Failed parse datadog type: %v", err)
	}
	svcDefinition.SetType(*ddType)
	svcDefinition.SetTeam(leanixData.CT)
	svcDefinition.SetTier(leanixData.AST)
	svcDefinition.SetDescription(leanixData.Description)
	svcDefinition.SetApplication(leanixData.Name)

	addRepoLinks(config, svcDefinition)
	addTags(leanixData, svcDefinition)

	test, _ := svcDefinition.MarshalJSON()
	fmt.Print(string(test))
	body := datadogV2.ServiceDefinitionsCreateRequest{ServiceDefinitionV2Dot2: svcDefinition}

	_, r, err = api.CreateOrUpdateServiceDefinitions(ctx, body)

	if err != nil {
		fmt.Printf("Error when calling `ServiceDefinitionApi.CreateOrUpdateServiceDefinitions`: %v\n", err)
		log.Fatalf("Full HTTP response: %v\n", r)
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `ServiceDefinitionApi.CreateOrUpdateServiceDefinitions`:\n%s\n", responseContent)
}

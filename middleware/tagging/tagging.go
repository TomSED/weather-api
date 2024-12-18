package tagging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const MOBILE_BFF_LEANIX_ID = "9d6d7a9e-ee31-4cb6-a2c2-f20da19ab8d9"

type Config struct {
	DataDogApiKey      string `json:"data_dog_api_key"`
	DataDogAppKey      string `json:"data_dog_app_key"`
	DatadogServiceName string `json:"datadog_service_name"`
	DatadogServiceType string `json:"datadog_service_type"`
	LeanIXAPIKey       string `json:"leanix_api_key"`
	LeanIXDomain       string `json:"leanix_domain"`
	LeanIXRequestURL   string `json:"leanix_request_url"`
	LeanIXAuthURL      string `json:"leanix_auth_url"`
	ReleaseVersion     string `json:"release_version"`
	Lifecycle          string `json:"lifecycle"`
	RepositoryUrl      string `json:"repository_url"`
	RepositoryName     string `json:"repository_name"`
	GitlabJobToken     string `json:"gitlab_job_token"`
	RepoProjectID      string `json:"repo_project_id"`
	TagSchemaVersion   string `json:"tag_schema_version"`
}

func loadConfig() Config {
	config := Config{}

	config.DataDogApiKey = getEnv("SVC_CATALOG_DD_API_KEY", "None")
	config.DataDogAppKey = getEnv("SVC_CATALOG_DD_APP_KEY", "None")

	config.LeanIXAPIKey = getEnv("LEANIX_API_KEY", "None")
	config.LeanIXDomain = getEnv("LEANIX_DOMAIN", "flybuys.leanix.net")
	config.LeanIXRequestURL = "https://" + config.LeanIXDomain + "/services/pathfinder/v1/graphql"
	config.LeanIXAuthURL = "https://" + config.LeanIXDomain + "/services/mtm/v1/oauth2/token"

	config.ReleaseVersion = getEnv("CI_COMMIT_TAG", getEnv("CI_COMMIT_SHORT_SHA", "None"))
	config.Lifecycle = mapLifecycle(getEnv("ENVIRONMENT", "None"))

	config.RepositoryUrl = getEnv("CI_PROJECT_URL", "None")
	config.TagSchemaVersion = getEnv("TAG_SCHEMA_VERSION", "None")
	config.RepositoryName = getEnv("CI_PROJECT_NAME", "None")
	config.GitlabJobToken = getEnv("FLYBUYS_GITLAB_API_PAT", "None")
	config.RepoProjectID = getEnv("CI_PROJECT_ID", "None")

	config.DatadogServiceName = getEnv("DATADOG_SVC_NAME", config.RepositoryName)

	config.DatadogServiceType = getEnv("DATADOG_SVC_TYPE", "custom")

	return config
}

func mapLifecycle(lifecycle string) string {
	switch lifecycle {
	case "prod":
		return "production"
	case "dev":
		return "development"
	case "eng":
		return "engineering"
	default:
		return lifecycle
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func jsonStructToMapStringString(data interface{}) (map[string]string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]string
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ParseTagsFromJson(jsonStr string) (map[string]*string, error) {
	// Unmarshal JSON into a map of strings
	var tempMap map[string]string
	err := json.Unmarshal([]byte(jsonStr), &tempMap)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	// Create a new map with string pointers
	ptrMap := make(map[string]*string)
	for key, value := range tempMap {
		valueCopy := value // create a copy to avoid pointer to loop variable
		ptrMap[key] = &valueCopy
	}

	return ptrMap, nil
}

func GetTagsJsonString() string {
	if os.Getenv("ENVIRONMENT") == "local" {
		return `{"flybuys:asset:id":"9d6d7a9e-ee31-4cb6-a2c2-f20da19ab8d9"}`
	}

	config := loadConfig()

	if (config.LeanIXAPIKey == "None") || (config.LeanIXAPIKey == "") {
		log.Fatalf("LEANIX_API_KEY is not set")
	}
	factsheet, err := getFactsheetByID(config, MOBILE_BFF_LEANIX_ID)
	if err != nil {
		log.Fatal("Error getFactsheetByID:", err)
	}
	fmt.Printf("Leanix factsheet: %v \n", factsheet)

	// generate flybuys tags
	asset := mapLeanixToFlybuysAsset(config, factsheet)
	assetJsonData, err := json.Marshal(asset)
	if err != nil {
		log.Fatal("Error marshalling to JSON:", err)
	}
	fmt.Println("Asset tag JSON string:", string(assetJsonData))

	// label repo
	if config.GitlabJobToken == "None" || config.GitlabJobToken == "" {
		log.Fatal("FLYBUYS_GITLAB_API_PAT is not set")
	}
	labels, err := leanixToRepoLabels(factsheet)
	if err != nil {
		log.Fatal("Error leanixToRepoLabels:", err)
	}
	labelRepo(config.GitlabJobToken, config.RepoProjectID, labels)
	fmt.Printf("Label repo finished, Labels: %v \n", labels)

	// create datadog service definition
	if os.Getenv("DATADOG_SVC_OUTPUT") != "true" {
		fmt.Print("DATADOG_SVC_OUTPUT is not true, skip create datadog service definition")
		return string(assetJsonData)
	}
	fmt.Println("Creating datadog service definition: ")
	if config.DataDogApiKey == "None" || config.DataDogApiKey == "" {
		log.Fatal("SVC_CATALOG_DD_API_KEY is not set")
	}
	if config.DataDogAppKey == "None" || config.DataDogAppKey == "" {
		log.Fatal("SVC_CATALOG_DD_APP_KEY is not set")
	}
	createDatadogServiceDefinition(config, factsheet)

	return string(assetJsonData)
}

package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-openapi/loads"
)

func convertEndpoint(endpoint string) string {
	// Replace slashes with hyphens
	endpoint = strings.ReplaceAll(endpoint, "/", "-")
	// Remove curly braces
	endpoint = strings.ReplaceAll(endpoint, "{", "")
	endpoint = strings.ReplaceAll(endpoint, "}", "")
	return endpoint
}

func generateLambdaEntry(folderPath string) error {

	lambdaPath := filepath.Join(folderPath, "main.go")

	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return err
	}

	if _, err := os.Stat(lambdaPath); err == nil {
		fmt.Println("File exists: ", lambdaPath)
		return nil
	} else if os.IsNotExist(err) {
		// Create the file since it does not exist
		err := os.WriteFile(lambdaPath, []byte(LAMBDA_CONTENT), 0644)
		if err != nil {
			return err
		}
		fmt.Println("Created", lambdaPath)
	} else {
		fmt.Println("Error checking file:", err)
		return err
	}
	return nil
}

// this function will auto generator lambda entry files under cmd folder
// example for use this function
// then we can use command `go generate`

// package main

// //go:generate echo "Generating code..."
// //go:generate go run generate.go
// import (
// 	"gitlab.com/flybuys/app/golang-cdk-shared/middleware/generator"
// )

// func main() {
// 	generator.GenerateFilesByOpenApiTemplate("./openapi.yaml")
// }

func GenerateFilesByOpenApiTemplate(openApiTemplatePath string) error {
	doc, err := loads.Spec(openApiTemplatePath)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI template: %v ", err)
	}

	// Access the spec
	spec := doc.Spec()

	// Iterate over paths
	for path, pathItem := range spec.Paths.Paths {
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
				err := generateLambdaEntry(fmt.Sprintf("../cmd/%s", functionID))
				if err != nil {
					log.Fatalf("Failed to generate lambda Entry: %v ",err)
				}
			}
		}
	}
	return nil
}

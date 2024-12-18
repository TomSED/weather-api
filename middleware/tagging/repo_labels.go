package tagging

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

func splitLabelName(labelName string) (string, string) {
	parts := strings.SplitN(labelName, "::", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func updateLabel(gitlabClient *gitlab.Client, pid, lid string, opt *gitlab.UpdateLabelOptions) (*gitlab.Label, *gitlab.Response, error) {
	u := fmt.Sprintf("projects/%s/labels/%s", gitlab.PathEscape(pid), gitlab.PathEscape(lid))

	req, err := gitlabClient.NewRequest(http.MethodPut, u, opt, nil)
	if err != nil {
		return nil, nil, err
	}

	l := new(gitlab.Label)
	resp, err := gitlabClient.Do(req, l)
	if err != nil {
		return nil, resp, err
	}

	return l, resp, nil
}

func labelRepo(token, projectID string, labels RepoLabels) {
	labelStringMap, err := jsonStructToMapStringString(labels)
	if err != nil {
		log.Fatalf("Failed to parse to labels string mapping: %v", err)
	}

	gl, err := gitlab.NewClient(token, gitlab.WithBaseURL("https://gitlab.com"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	_, _, err = gl.Projects.GetProject(projectID, nil)
	if err != nil {
		log.Fatalf("Failed to get project: %v", err)
	}

	gitlabLabels, _, err := gl.Labels.ListLabels(projectID, nil)
	if err != nil {
		log.Fatalf("Failed to list labels: %v", err)
	}

	color := "#0033CC"

	for key, value := range labelStringMap {
		newLabel := fmt.Sprintf("%s::%s", key, value)
		createNewLabel := true

		for _, item := range gitlabLabels {
			flybuysLabelKey, _ := splitLabelName(item.Name)
			if item.Name == newLabel {
				createNewLabel = false
				break
			} else if flybuysLabelKey == key && item.Name != newLabel {
				createNewLabel = false
				item.Name = newLabel
				_, _, err := updateLabel(gl, projectID, strconv.Itoa(item.ID), &gitlab.UpdateLabelOptions{
					Name:    &item.Name,
					NewName: &newLabel,
				})
				if err != nil {
					log.Fatalf("Failed to update label: %v", err)
				}
				break
			}
		}
		if createNewLabel {
			_, _, err := gl.Labels.CreateLabel(projectID, &gitlab.CreateLabelOptions{
				Name:  &newLabel,
				Color: &color,
			})
			if err != nil {
				log.Fatalf("Failed to create label: %v", err)
			}
		}
	}
}

type RepoLabels struct {
	AssetID            string `json:"flybuys:asset:id"`
	AssetName          string `json:"flybuys:asset:name"`
	AssetTier          string `json:"flybuys:asset:tier"`
	DataClassification string `json:"flybuys:asset:data-classification"`
	AssetTeam          string `json:"flybuys:asset:team"`
	AssetPartner       string `json:"flybuys:asset:partner,omitempty"`
	AssetCustodian     string `json:"flybuys:asset:custodian,omitempty"`
}

func leanixToRepoLabels(leanixFactSheet FactSheet) (RepoLabels, error) {
	return RepoLabels{
		AssetID:            leanixFactSheet.ID,
		AssetName:          strings.ReplaceAll(strings.ToLower(leanixFactSheet.Name), " ", "-"),
		AssetTier:          leanixFactSheet.AST,
		DataClassification: leanixFactSheet.DataClassification,
		AssetTeam:          leanixFactSheet.CT,
		AssetPartner:       leanixFactSheet.Partner,
		AssetCustodian:     leanixFactSheet.CTO,
	}, nil
}

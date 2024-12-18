package tagging

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type FlybuysAsset struct {
	AssetID               string `json:"flybuys:asset:id"`
	AssetName             string `json:"flybuys:asset:name"`
	AssetTier             string `json:"flybuys:asset:tier"`
	DataClassification    string `json:"flybuys:asset:data-classification,omitempty"`
	AssetTeam             string `json:"flybuys:asset:team"`
	AssetPartner          string `json:"flybuys:asset:partner,omitempty"`
	Environment           string `json:"environment"`
	ApplicationVersion    string `json:"flybuys:application:version"`
	RepositoryURL         string `json:"flybuys:application:repository-url"`
	TagSchemaVersion      string `json:"flybuys:tag:schema-version"`
	AssetCustodian        string `json:"flybuys:asset:custodian,omitempty"`
	MonitorBoardURL       string `json:"flybuys:application:monitor-board-url,omitempty"`
	ApplicationURL        string `json:"flybuys:application:url,omitempty"`
	SlackChannelName      string `json:"flybuys:team:slack-channel-name,omitempty"`
	SlackChannelID        string `json:"flybuys:team:slack-channel-id,omitempty"`
	SlackChannelURL       string `json:"flybuys:team:slack-channel-url,omitempty"`
	JiraProject           string `json:"flybuys:team:jira-project,omitempty"`
	CDKVersion            string `json:"flybuys:cdk:version,omitempty"`
	AutomationSchedule    string `json:"flybuys:automation:schedule,omitempty"`
	AutomationBackup      string `json:"flybuys:automation:backup,omitempty"`
	ApplicationCostCentre string `json:"flybuys:application:cost-centre,omitempty"`
}

func mapLeanixToFlybuysAsset(config Config, leanixFactSheet FactSheet) FlybuysAsset {
	retFlybuysAsset := FlybuysAsset{
		// map leanix to FlybusAsseet
		AssetID:            leanixFactSheet.ID,
		AssetName:          leanixFactSheet.Name,
		AssetTier:          leanixFactSheet.AST,
		DataClassification: leanixFactSheet.DataClassification,
		AssetTeam:          leanixFactSheet.CT,
		AssetPartner:       leanixFactSheet.Partner,
		Environment:        config.Lifecycle,
		ApplicationVersion: config.ReleaseVersion,
		RepositoryURL:      config.RepositoryUrl,
		TagSchemaVersion:   config.TagSchemaVersion,

		// map optional tag to FlybusAsseet
		AssetCustodian:        leanixFactSheet.CTO,
		MonitorBoardURL:       getEnv("MONITOR_BOARD_URL", ""),
		ApplicationURL:        getEnv("APPLICATION_URL", ""),
		SlackChannelName:      getEnv("SLACK_CHANNEL_NAME", ""),
		SlackChannelID:        getEnv("SLACK_CHANNEL_ID", ""),
		SlackChannelURL:       getEnv("SLACK_CHANNEL_URL", ""),
		JiraProject:           getEnv("JIRA_PROJECT", ""),
		CDKVersion:            getEnv("CDK_VERSION", ""),
		AutomationSchedule:    getEnv("SCHEDULE", ""),
		AutomationBackup:      getEnv("BACKUP", ""),
		ApplicationCostCentre: getEnv("COST_CENTRE", ""),
	}
	// Convert name to lowercase and replace space with "-"
	if retFlybuysAsset.AssetName != "" {
		retFlybuysAsset.AssetName = strings.ReplaceAll(strings.ToLower(retFlybuysAsset.AssetName), " ", "-")
	}

	return retFlybuysAsset
}

func writeTagsToOutputFile(logger *log.Logger, tags map[string]interface{}, output string) bool {
	file, err := os.Create(output)
	if err != nil {
		logger.Println("File not found, name:", output)
		return false
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tags); err != nil {
		logger.Println("Error writing to file:", err)
		return false
	}

	logger.Println("Updated", output, "with latest contents")
	return true
}

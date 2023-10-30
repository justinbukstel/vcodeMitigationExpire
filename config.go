package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

type config struct {
	Auth struct {
		CredsFile string `json:"credsFile"`
	} `json:"auth"`

	Mode struct {
		LogOnly           bool `json:"logOnly"`
		RejectMitigations bool `json:"rejectMitigations"`
	} `json:"mode"`

	TargetMitigations struct {
		PotentialFalsePositive bool `json:"potentialFalsePositive"`
		MitigatedByDesign      bool `json:"mitigatedByDesign"`
		MitigationByOSEnv      bool `json:"mitigationByOSEnv"`
		MitigatedByNetworkEnv  bool `json:"mitigatedByNetworkEnv"`
		ReviewedNoActionTaken  bool `json:"reviewedNoActionTaken"`
		RemediatedByUser       bool `json:"remediatedByUser"`
		ReportedToLibraryMaintainer	bool `json:"reportedToLibraryMaintainer"`
		AcceptTheRisk      bool 'json:"acceptTheRisk"
	} `json:"targetMitigations"`

	CommentText struct {
		RequireCommentText bool   `json:"requireCommentText"`
		Text               string `json:"text"`
	} `json:"commentText"`

	AppScope struct {
		LimitAppList    bool   `json:"limitAppList"`
		AppListTextFile string `json:"appListTextFile"`
	} `json:"appScope"`

	ExpirationDetails struct {
		DateFlawFound            bool   `json:"DateFlawFound"`
		DateOfMitigationApproval bool   `json:"dateOfMitigationApproval"`
		SpecificDate             bool   `json:"specificDate"`
		Date                     string `json:"date"`
		DaysToExpire             int    `json:"daysToExpire"`
		RejectionComment         string `json:"rejectionComment"`
	} `json:"expirationDetails"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Veracode username")
}

func parseConfig() config {

	flag.Parse()

	//READ CONFIG FILE
	var config config

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	// VALID AT LEAST ONE MODE IS SET
	modeCounter := 0
	if config.Mode.LogOnly == true {
		modeCounter++
	}
	if config.Mode.RejectMitigations == true {
		modeCounter++
	}
	if modeCounter == 0 {
		log.Fatal("One mode must be set to be set to true")
	}
	if modeCounter > 1 {
		log.Fatal("Only one mode is allowed")
	}

	// VALIDATE AT LEAST ONE TARGET MITIGATION IS SET
	if config.TargetMitigations.PotentialFalsePositive == false &&
		config.TargetMitigations.MitigatedByDesign == false &&
		config.TargetMitigations.MitigationByOSEnv == false &&
		config.TargetMitigations.MitigatedByNetworkEnv == false &&
		config.TargetMitigations.ReviewedNoActionTaken == false &&
		config.TargetMitigations.RemediatedByUser == false &&
		config.TargetMitigations.ReportedToLibraryMaintainer == false &&
	        config.TargetMitigations.AcceptTheRisk == false {
		log.Fatal("at least one target mitigation must be set to true")
	}

	// VALIDATE EXPIRATION CONFIG
	expTypeCounter := 0
	if config.ExpirationDetails.DateFlawFound == true {
		expTypeCounter++
	}
	if config.ExpirationDetails.DateOfMitigationApproval == true {
		expTypeCounter++
	}
	if config.ExpirationDetails.SpecificDate == true {
		expTypeCounter++
	}
	if expTypeCounter == 0 {
		log.Fatal("One expiration trigger needs to be set to true")
	}
	if expTypeCounter > 1 {
		log.Fatal("Only one expiration trigger is allowed")
	}

	return config
}

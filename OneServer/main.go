package main

import (
	"OneServer/logs"
	"OneServer/server"
	"log"
)

func main() {

	var profilePath string
	profilePath = "./profile.json"

	ts := server.NewTeamServer()
	if profilePath != "" {
		err := ts.SetProfile(profilePath)
		if err != nil {
			log.Fatalf("Error loading profile: %v", err)
		}
	}

	ts.Profile.Validate()

	logs.Logger = logs.InitLogger(ts.Profile.ZapConfig)

	ts.Start()
}

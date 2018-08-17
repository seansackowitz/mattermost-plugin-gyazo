package main

import (
	"github.com/mattermost/mattermost-server/plugin"
)

func main() {
	plugin.ClientMain(&Plugin{})
}

type GyazoJson struct {
	Version       string  `json:"version"`
	Type          string  `json:"type"`
	Provider_name string  `json:"provider_name"`
	Provider_url  string  `json:"provider_url"`
	Url           string  `json:"url"`
	Html          string  `json:"html"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Scale         float32 `json:"scale"`
}

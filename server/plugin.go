package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// Plugin the main struct for everything
type Plugin struct {
	api   plugin.API
	links atomic.Value
}

// OnActivate is invoked when the plugin is activated.
func (p *Plugin) OnActivate(api plugin.API) error {
	p.api = api

	if err := p.OnConfigurationChange(); err != nil {
		return err
	}

	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	var c Configuration
	err := p.api.LoadPluginConfiguration(&c)
	if err != nil {
		return err
	}

	return nil
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get("https://api.gyazo.com/api/oembed?url=" + url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (p *Plugin) FilterPost(post *model.Post) (*model.Post, string) {
	message := post.Message
	r, _ := regexp.Compile("(https://gyazo.com/)[a-z0-9]*")
	link := r.FindString(message)
	if len(link) > 0 {
		return p.FormatPost(post, link), ""
	}

	return post, ""
}

func (p *Plugin) FormatPost(post *model.Post, gyazo string) *model.Post {
	gyazojson := new(GyazoJson)
	getJson(gyazo, gyazojson)
	url := gyazojson.Url
	if gyazojson.Type == "video" {
		url = strings.Replace(gyazo, "https://", "https://i.", -1) + ".gif"
	}
	post.Message = "[Gyazo](" + gyazo + ")\n" + url

	mlog.Error(post.Message)
	return post
}

// MessageWillBePosted is invoked when a message is posted by a user before it is commited
// to the database.
func (p *Plugin) MessageWillBePosted(post *model.Post) (*model.Post, string) {
	mlog.Error(post.Message)
	return p.FilterPost(post)
}

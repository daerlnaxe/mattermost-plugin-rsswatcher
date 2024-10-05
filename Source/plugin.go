package main

/*
Inspired by rssfeed to learn golang and make my first plugin, big thanks to William B Ernest.

Great job with Welcome Bot, helpful to understand how to create a bot.

*/

import(
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/pluginapi"

)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
// Named Test and used by ServeHTTP
type RssWatcherPlugin struct {
	plugin.MattermostPlugin

	// Need by the bot
	client *pluginapi.Client

	// Keep bot UserID after initialisation
	botUserID       string

	// Flag to stop subProcess
	corePosterFlag bool
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *RssWatcherPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

// This example demonstrates a plugin that handles HTTP requests which respond by greeting the
// world.
func main() {
	plugin.ClientMain(&RssWatcherPlugin{})
}
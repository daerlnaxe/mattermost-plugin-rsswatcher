package main

/*
Commands list used

Called by activate.go

*/

import (
	"context"
	"fmt"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"

	//"github.com/mattermost/mattermost/server/public/RssWatcherPlugin"
//	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"strings"
)


// The help you see when you type /rssw {*}
const COMMAND_HELP = `* |/rssw sub {url}| - Connect your Mattermost channel to a RSS feed.
* |/rssw ls| - List RSS Subscribtions for this Mattermost.
`

const (
	commandTriggerSub                 = "sub"
	commandTriggerList                 = "ls"
	commandTriggerHelp                 = "help"
)


// ---- API

// Commands description
// API
func getCommand() *model.Command {
	return &model.Command{
		// Slack command "/rssw"
		Trigger:          "rssw",
		DisplayName:      "RssWatcher",
		Description:      "Allows user to subscribe to an RSS feed.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: ls, sub, rem, help",
		AutoCompleteHint: "[command]",
		AutocompleteData: getAutocompleteData(),
		

	}
}

// Will interact with user when he is using Slack commands
// API
func (p *RssWatcherPlugin) ExecuteCommand(_ *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	// Split for command args
	split := strings.Fields(args.Command)
	
	command := split[0]
	parameters := []string{}
	action := ""

	// No param
	if len(split) > 1 {
		action = split[1]
	}

	// With params
	if len(split) > 2 {
		parameters = split[2:]
	}

	// Doesn't concern RssWatcher => pass
	if command != "/rssw" {
		return &model.CommandResponse{}, nil
	}

	// Managing actions
	switch action {
	// help
	case "help":
		text := "###### Mattermost RSSFeed Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return getCommandResponse("ephemral", text), nil

	// Subscribing to a Rss feed
	case "sub":
		if len(parameters) == 0 {
			return getCommandResponse("ephemeral", "Please specify a url."), nil
		} else if len(parameters) > 1 {
			return getCommandResponse("ephemeral", "Please specify a valid url."), nil
		}

		url := parameters[0]

		// Try to register the RSS feed
		if err := p.subscribe(context.Background(), args.ChannelId, url); err != nil {
			return getCommandResponse("ephemeral", err.Error()), nil
		}
		
		// Return this message if succesfull
		return getCommandResponse("ephemeral", fmt.Sprintf("Successfully subscribed to %s.", url)), nil
		
		//return getCommandResponse("ephemeral", "Ca marchera un jour, ça marchera..."), nil

	// list of subscribtions for the channel
	case "ls"	:
		
		subscriptions, err := p.getSubscriptions()

		if err != nil {
			return getCommandResponse("ephemeral", err.Error()), nil
		}

		// Check Subscribtions
		index := 0
		txt := "### RSS Subscriptions in this channel\n"

		//commits := map[Subscription]int{}
		

		/*for i := 0; i < 10; i++ {
			value := subscriptions.Subscriptions[i]

			if value.ChannelID == args.ChannelId {
				index ++
				txt += fmt.Sprintf("* %s) `%s`\n", index,value.URL)
			}			
			
			
		}*/
		
		for _, value := range subscriptions.Subscriptions {
			
			if value.ChannelID == args.ChannelId {
				txt += fmt.Sprintf("* `%d`) `%s`\n", index+1,value.URL)
				index++
			}			
			
		}

		// If no subscription, manage differently
		if(index==0){
			return getCommandResponse("ephemeral", "There is no subscribtion on this channel"), nil
		// There is entries
		}else{
			return getCommandResponse("ephemeral", txt), nil
		}
		
	

	// Default
	default:
		text := "###### Mattermost RssWatcher Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return getCommandResponse("ephemeral", text), nil
	}

}


// ---- END API



/* Le type CommandResponse a besoin maintenant que type soit une string.Add a RSS to watch
	- canal: in_channel
*/
func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     BOT_DISPLAY_NAME,
		//IconURL:      RSSFEED_ICON_URL,
		Type:         "ephemeral",
	}
}




// Check if user has admin role
func (p *RssWatcherPlugin) hasSysadminRole(userID string) (bool, error) {
	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		return false, appErr
	}
	if !strings.Contains(user.Roles, "system_admin") {
		return false, nil
	}
	return true, nil
}

/*
	Caution: cname of Autocomplete determinate the name of the slack command
	En plus de Help
*/
func getAutocompleteData() *model.AutocompleteData {
	rsswatcherbot := model.NewAutocompleteData("rssw", "[command]",
		"Available commands: help, sub, ls")
		
		sub := model.NewAutocompleteData("sub", "", "Add a RSS to watch")
		rsswatcherbot.AddCommand(sub)

		ls := model.NewAutocompleteData("ls", "", "List Rss for the channel")
		rsswatcherbot.AddCommand(ls)
		/*
	preview := model.NewAutocompleteData("preview", "[team-name]", "Preview the welcome message for the given team name")
	preview.AddTextArgument("Team name to preview welcome message", "[team-name]", "")
	welcomebot.AddCommand(preview)



	setChannelWelcome := model.NewAutocompleteData("set_channel_welcome", "[welcome-message]", "Set the welcome message for the channel")
	setChannelWelcome.AddTextArgument("Welcome message for the channel", "[welcome-message]", "")
	welcomebot.AddCommand(setChannelWelcome)

	getChannelWelcome := model.NewAutocompleteData("get_channel_welcome", "", "Print the welcome message set for the channel")
	welcomebot.AddCommand(getChannelWelcome)

	deleteChannelWelcome := model.NewAutocompleteData("delete_channel_welcome", "", "Delete the welcome message for the channel")
	welcomebot.AddCommand(deleteChannelWelcome)
*/
	return rsswatcherbot
}
//
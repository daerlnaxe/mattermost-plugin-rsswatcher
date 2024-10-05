package main

import (

	//"io/ioutil"
	//"path/filepath"


	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/pluginapi"
	"github.com/pkg/errors"

)

const (
	BOTNAME="rsswatcherbot"
	BOT_DISPLAY_NAME = "RssWatcher Bot"
	BOT_DESCRIPTION = "A bot account created by the RssWatcher plugin."

)

// First when plugin start
// Remember to use the plugin name
func (p *RssWatcherPlugin) OnActivate() error{
	/*
	// Verifyng Bot exists
	_, err := p.ensureBotExists()

	if err != nil {
		p.API.LogError("RssWatcher >> OnActivate >> Failed to find "+BOT_DISPLAY_NAME+" user", "err", err)
		return err
	}

	*/

	// New client
	p.client = pluginapi.NewClient(p.API, p.Driver)

	// Bot informations
	bot := &model.Bot{
		Username:    BOTNAME,
		DisplayName: BOT_DISPLAY_NAME,
		Description: BOT_DESCRIPTION,
	}

	// Bot creation
	botUserID, appErr := p.client.Bot.EnsureBot(bot)
	
	if appErr != nil {
		return errors.Wrap(appErr, "failed to ensure bot user")
	}

	// Registering commands for the bot
	err := p.API.RegisterCommand(getCommand())
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}


	p.botUserID = botUserID

	// Launch core to automatizing posts
	p.corePosterFlag= true
	go p.initCorePoster()

	//
	return nil
}


// Verifying bot is ready
// Remember to use the plugin name
func (p *RssWatcherPlugin) ensureBotExists() (string, *model.AppError) {
	p.API.LogDebug("Ensuring " + BOT_DISPLAY_NAME + " exists")

	// Trying to create the bot
	bot, createErr := p.API.CreateBot(&model.Bot{
			Username:    BOTNAME,
			DisplayName: BOT_DISPLAY_NAME,
			Description: "Allows users to subscribe to RSS feeds.",})

	// Failed to create bot
	if createErr != nil {
		p.API.LogDebug("RssWatcher >> "+ BOT_DISPLAY_NAME + " not created. Attempting to find existing one...")

		p.API.LogError("RssWatcher >> "+createErr.Message);

		// Verifying user exists with the same name
		userBot, err := p.API.GetUserByUsername(BOTNAME)

		// No user with this name <-- On devrait mettre une erreur
		if err != nil || userBot == nil {
			p.API.LogError("RssWatcher >> No user with the name "+BOT_DISPLAY_NAME+": userBot", "err", err)
			return "", err
		}

		// Verifying bot exist
		bot, err = p.API.GetBot(userBot.Id, true)


		if err != nil {
			p.API.LogError("RssWatcher >> Failed to find "+BOT_DISPLAY_NAME, "err", err)
			return "", err
		}

		p.API.LogDebug("RssWatcher >> Found " + BOT_DISPLAY_NAME)

	}else{
		/*
		if err := p.setBotProfileImage(bot.UserId); err != nil {
			p.API.LogError("Failed to set profile image for bot", "err", err)
		}*/

		p.API.LogDebug("RssWatcher >> "+ BOT_DISPLAY_NAME + " created")
	}

	p.botUserID = bot.UserId

	return bot.UserId, nil
}
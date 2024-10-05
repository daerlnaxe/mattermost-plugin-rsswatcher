package main

import(
	"fmt"
	
	"time"
	"github.com/mattermost/mattermost/server/public/model"
	// RSS/ATOM Parser
	"github.com/mmcdole/gofeed"
	// Convert HTML To Markdown
	"github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)


const (
	DEBUG = true
)

func (p *RssWatcherPlugin) initCorePoster() {
	wakeUpTime, err := p.getWakeUpTime()
	
	if err != nil {
		p.API.LogError(err.Error())
	}


	// golang seems to not implement a while loop, for works.
	for p.corePosterFlag {
		
		// Run poster manager
		err := p.subscribtionManager()
		

		if err != nil {
			p.API.LogError(err.Error())

		}

		// Impossible to fall under a minute (for the)
		// Mandatory to use time.Duration() don't accept int * minutes
		time.Sleep(time.Duration( wakeUpTime) * time.Minute)
	}
}



func (p *RssWatcherPlugin) subscribtionManager()error{
	// get the stored 'Subscription" map (or new)
	currentSubscriptions, err := p.getSubscriptions()
	if err != nil {
		p.API.LogError(err.Error())
		return err
	}


	
	for _, value := range currentSubscriptions.Subscriptions {
		//err := p.createBotPost(value.ChannelID, "I'm still standing yeah yeah yeah", "")
		//err := p.parseContent(value)
	
		fp := gofeed.NewParser()
		//feed, err := fp.ParseURL("http://feeds.twit.tv/twit.xml")
		feed, err := fp.ParseURL("https://www.jeuxvideo.com/rss/rss.xml")

		if err != nil {
			p.API.LogError(">>>> Parsing Error"+ err.Error())
		}

		if DEBUG {
			fmt.Println(feed.Title)			
		}

		// Fonctionne sans image si j'ai bien compris
		p.sendItems(value.ChannelID, feed.Items);
		//p.sendItems(value.ChannelID, feed.Images;
	}
	
	
	return nil
}

/*
// Will Parse RSS Content thanksfull
func(p* RssWatcherPlugin) parseContent(subscribtion *Subscription) (error, FEED) {

	if err != nil {
		p.API.LogError(err.Error())
		return err
	}

	return err, fp

}
*/

// Loop on all Items and send them
func (p *RssWatcherPlugin) sendItems(channelID string, items []*gofeed.Item) {
	
	// !!!! Rajouter un comparateur pour ne pas reposter les mêmes messages
	converter := md.NewConverter("", true, nil)

		// Use the `GitHubFlavored` plugin from the `plugin` package.
		converter.Use(plugin.GitHubFlavored())

	for _, value := range items {
		if DEBUG{
			/*
			fmt.Println("--------------------------------------------------------------------------------------------------------------")

			fmt.Println(value.Title)
			//fmt.Println(value.Description)
			//fmt.Println(value.Content)
			/*fmt.Println("Enclosures")
			fmt.Println(value.Enclosures)
			fmt.Println("Links")
			fmt.Println(value.Links)
			//fmt.Println("Image")
			//fmt.Println(value.Image)*/
			fmt.Println("--------------------------------------------------------------------------------------------------------------")		
		}

		
		message := "##### "

		/*
		if config.FormatTitle {
			message ="##### "
		}*/
		
		message += value.Title + "\n"
		//message += fmt.Sprintf("[%s](%s)\n",value.Link,value.Link)
		message += fmt.Sprintf("%s \n",value.Link)

		message += "\n**Summary**:\n"
		/*

			if config.FormatTitle {
				post = post + "###### "
			}
			post = post + item.Title + "\n"
		}*/
	
			/*
		if config.ShowRSSLink {
			post = post + strings.TrimSpace(item.Link) + "\n"
		}
			*/
		
		/*
			if config.ShowDescription {
			post = post + html2md.Convert(item.Description) + "\n"
		}
		*/
		
		p.API.LogError("Try Convert string to md")

		markdown, err := converter.ConvertString(value.Description)
			

		if err != nil {
			p.API.LogError("%s",err)
		  }

		
		if(len(markdown) > 200){
			p.API.LogError("cut markdown")
			message += markdown[:200]+ "[...]\n\n"
		}

		p.API.LogError("Published: ")
		message+= fmt.Sprintf("*%s", value.Published)

		p.API.LogError("Author:")
		if(value.Author != nil){
			message +=  fmt.Sprintf(" by %s*\n", value.Author.Name)
		}
		

		//message += value.Link


		

		//message += fmt.Sprintf( "![%s](%s \"Mattermost Icon\")",value.Image.Title, value.Image.URL)
		//message += fmt.Sprintf( "![%s](%s =200 \"Mattermost Icon\")",value.Image.Title, value.Image.URL)
//		message += "----mwahahhaha\n----"
/*
		message += fmt.Sprintf( "![%s](%s =200 \"Hover text\")",value.Image.Title, value.Image.URL)
		message += "----mwahahhaha\n----"
*/
		
		  
  
		
///		message += "-----\r\n"
		/*markdown2, err := converter.ConvertString(value.Content)
		message+=markdown2*/




		p.createBotPost(channelID, message, "")
	}



}


func (p *RssWatcherPlugin) createBotPost(channelID string, message string, postType string) error {
/*
	mee :=`json:"attachments": [
        {
            "fallback": "test",
            "color": "#FF8000",
            "pretext": "This is optional pretext that shows above the attachment.",
            "text": "This is the text of the attachment. It should appear just above an image of the Mattermost logo. The left border of the attachment should be colored orange, and below the image it should include additional fields that are formatted in columns. At the top of the attachment, there should be an author name followed by a bolded title. Both the author name and the title should be hyperlinks.",
            "author_name": "Mattermost",
            "author_icon": "https://mattermost.com/wp-content/uploads/2022/02/icon_WS.png",
            "author_link": "https://mattermost.org/",
            "title": "Example Attachment",
            "title_link": "https://developers.mattermost.com/integrate/reference/message-attachments/",
            "fields": [
                {
                    "short":false,
                    "title":"Long Field",
                    "value":"Testing with a very long piece of text that will take up the whole width of the table. And then some more text to make it extra long."
                },
                {
                    "short":true,
                    "title":"Column One",
                    "value":"Testing"
                },
                {
                    "short":true,
                    "title":"Column Two",
                    "value":"Testing"
                },
                {
                    "short":false,
                    "title":"Another Field",
                    "value":"Testing"
                }
            ],
            "image_url": "https://mattermost.com/wp-content/uploads/2022/02/icon_WS.png"
        }
    ]
	}
	`
*/
		
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: channelID,
		Message:   message,
		//Type:      "custom_git_pr",
		/*Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": botDisplayName,
		},*/
		
	}
	



	if _, err := p.API.CreatePost(post); err != nil {
		p.API.LogError(
			"We could not create the RSS post",
			"user_id", post.UserId,
			"err", err.Error(),
		)
	}

	return nil
}

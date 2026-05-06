package options

import "github.com/fgrzl/json/polymorphic"

func init() {
	polymorphic.RegisterType[SlackUsersCollectorOptions]()
	polymorphic.RegisterType[SlackChannelsCollectorOptions]()
	polymorphic.RegisterType[SlackChannelMessagePostActionOptions]()
	polymorphic.RegisterType[SlackUserMessagePostActionOptions]()
}

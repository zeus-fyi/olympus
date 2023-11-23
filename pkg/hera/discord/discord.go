package hera_discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

type Discord struct {
	*discordgo.Session
}

var DiscordClient Discord

func InitDiscordClient(ctx context.Context, token string) Discord {
	dg, err := discordgo.New(token)
	if err != nil {
		panic(err)
	}
	DiscordClient.Session = dg
	return DiscordClient
}

func (d *Discord) ListAllChannels(ctx context.Context) ([]*discordgo.Channel, error) {
	var allChannels []*discordgo.Channel
	guilds, err := d.UserGuilds(0, "", "")
	if err != nil {
		log.Err(err).Msg("Error getting guilds")
		return nil, err
	}

	for _, guild := range guilds {
		channels, cerr := d.GuildChannels(guild.ID)
		if cerr != nil {
			fmt.Printf("Error getting channels for guild %s: %s\n", guild.Name, err)
			continue
		}
		allChannels = append(allChannels, channels...)
	}
	return allChannels, nil
}

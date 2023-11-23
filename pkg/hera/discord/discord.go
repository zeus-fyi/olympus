package hera_discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DiscordWrapper struct {
	DC *discordgo.Session
}

var DiscordClient *DiscordWrapper

func InitDiscordClient(ctx context.Context, token string) {
	dg, err := discordgo.New(token)
	if err != nil {
		panic(err)
	}
	DiscordClient = &DiscordWrapper{
		DC: dg,
	}
	return
}

func (d *DiscordWrapper) ListAllChannels(ctx context.Context) ([]*discordgo.Channel, error) {
	var allChannels []*discordgo.Channel
	guilds, err := d.DC.UserGuilds(0, "", "")
	if err != nil {
		return nil, err
	}

	for _, guild := range guilds {
		channels, cerr := d.DC.GuildChannels(guild.ID)
		if cerr != nil {
			fmt.Printf("Error getting channels for guild %s: %s\n", guild.Name, cerr)
			continue
		}
		allChannels = append(allChannels, channels...)
	}
	return allChannels, nil
}

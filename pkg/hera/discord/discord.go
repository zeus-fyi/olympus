package hera_discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	QuickNodeChannelID = "685243210829922350"
	ZeusfyiChannelID   = "1018610566572544080"
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

func (d *DiscordWrapper) FetchChatMessages(ctx context.Context, chID, afterID string, limit int) ([]*discordgo.Message, error) {
	messages, err := d.DC.ChannelMessages(chID, limit, "", afterID, "")
	if err != nil {
		fmt.Println("Error retrieving messages:", err)
		return nil, err
	}
	return messages, nil
}

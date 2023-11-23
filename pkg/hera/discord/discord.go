package hera_discord

import (
	"context"
)

type Discord struct {
}

var DiscordClient Discord

func InitDiscordClient(ctx context.Context, id, secret, u, pw string) (Discord, error) {
	//token := "YOUR_BOT_TOKEN" // Replace with your Discord bot token
	//
	//dg, err := discordgo.New("Bot " + token)
	//if err != nil {
	//	fmt.Println("error creating Discord session,", err)
	//	return
	//}

	return Discord{}, nil
}

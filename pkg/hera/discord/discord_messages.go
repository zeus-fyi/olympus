package hera_discord

import "time"

type ChannelMessages struct {
	Guild struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"guild"`
	Channel struct {
		Id         string `json:"id"`
		CategoryId string `json:"categoryId"`
		Category   string `json:"category"`
		Name       string `json:"name"`
		Topic      string `json:"topic"`
	} `json:"channel"`
	Messages     []Message `json:"messages"`
	MessageCount int       `json:"messageCount"`
}

type Message struct {
	Author struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Roles    []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
	} `json:"author"`
	Content  string `json:"content"`
	Id       string `json:"id"`
	Mentions []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Roles    []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"roles"`
	} `json:"mentions"`
	Reactions []struct {
		Count int `json:"count"`
		Emoji struct {
			Code string `json:"code"`
		} `json:"emoji"`
		Users []struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Nickname string `json:"nickname"`
		} `json:"users"`
	} `json:"reactions"`
	TimestampEdited time.Time `json:"timestampEdited"`
	Type            string    `json:"type"`
	Reference       struct {
		ChannelId string `json:"channelId"`
		GuildId   string `json:"guildId"`
		MessageId string `json:"messageId"`
	} `json:"reference,omitempty"`
}

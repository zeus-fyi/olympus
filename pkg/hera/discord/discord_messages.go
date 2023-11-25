package hera_discord

import "time"

type Guild struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Channel struct {
	Id         string `json:"id"`
	CategoryId string `json:"categoryId"`
	Category   string `json:"category"`
	Name       string `json:"name"`
	Topic      string `json:"topic"`
}

type ChannelMessages struct {
	Guild        Guild     `json:"guild"`
	Channel      Channel   `json:"channel"`
	Messages     []Message `json:"messages"`
	MessageCount int       `json:"messageCount"`
}

type Author struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Roles    []Role `json:"roles"`
}

type Mention struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Roles    []Role `json:"roles"`
}

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Reaction struct {
	Count int    `json:"count"`
	Emoji Emoji  `json:"emoji"`
	Users []User `json:"users"`
}

type Emoji struct {
	Code string `json:"code"`
}
type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

type Reference struct {
	ChannelId string `json:"channelId"`
	GuildId   string `json:"guildId"`
	MessageId string `json:"messageId"`
}

func (r Reference) IsEmpty() bool {
	return r.ChannelId == "" && r.GuildId == "" && r.MessageId == ""
}

type Message struct {
	Author          Author     `json:"author"`
	Content         string     `json:"content"`
	Id              string     `json:"id"`
	Mentions        []Mention  `json:"mentions"`
	Reactions       []Reaction `json:"reactions"`
	TimestampEdited time.Time  `json:"timestampEdited"`
	Type            string     `json:"type"`
	Reference       Reference  `json:"reference,omitempty"`
}

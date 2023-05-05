package server

import (
	"fmt"
	"time"

	"github.com/wakscord/new-wakscord-node/env"
)

var (
	status = nodeStatus{
		Info: nodeInfo{
			NodeID: env.GetInt("ID", 0),
			Owner:  env.GetString("OWNER", "Unknown"),
		},
		Pending: nodePending{
			Total:    0,
			Messages: 0,
			Tasks:    0,
		},
		Processed: 0,
		Uptime:    0,
	}

	deletedWebhooks = map[string]struct{}{}
	tasks           = make(chan task, 100)

	startTime = time.Now()

	serverKey = fmt.Sprintf("Bearer %s", env.GetString("KEY", "wakscord"))
)

type nodeStatus struct {
	Info       nodeInfo    `json:"info"`
	Pending    nodePending `json:"pending"`
	Processed  int32       `json:"processed"`
	Deleted    int         `json:"deleted"`
	Uptime     int         `json:"uptime"`
	Goroutines int         `json:"goroutines"`
}

type nodeInfo struct {
	NodeID int    `json:"node_id"`
	Owner  string `json:"owner"`
}

type nodePending struct {
	Total    int32 `json:"total"`
	Messages int32 `json:"messages"`
	Tasks    int32 `json:"tasks"`
}

type requestPayload struct {
	Keys []string      `json:"keys"`
	Data WebhookParams `json:"data"`
}

type task struct {
	chunks [][]string
	data   WebhookParams
}

// whole discord structures are copied from bwmarrin/discordgo
type WebhookParams struct {
	Content         string                  `json:"content,omitempty"`
	Username        string                  `json:"username,omitempty"`
	AvatarURL       string                  `json:"avatar_url,omitempty"`
	TTS             bool                    `json:"tts,omitempty"`
	Embeds          []*MessageEmbed         `json:"embeds,omitempty"`
	AllowedMentions *MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           int                     `json:"flags,omitempty"`
}

type MessageEmbed struct {
	URL         string                 `json:"url,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	Color       int                    `json:"color,omitempty"`
	Footer      *MessageEmbedFooter    `json:"footer,omitempty"`
	Image       *MessageEmbedImage     `json:"image,omitempty"`
	Thumbnail   *MessageEmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *MessageEmbedVideo     `json:"video,omitempty"`
	Provider    *MessageEmbedProvider  `json:"provider,omitempty"`
	Author      *MessageEmbedAuthor    `json:"author,omitempty"`
	Fields      []*MessageEmbedField   `json:"fields,omitempty"`
}

type MessageEmbedFooter struct {
	Text         string `json:"text,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedImage is a part of a MessageEmbed struct.
type MessageEmbedImage struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedThumbnail is a part of a MessageEmbed struct.
type MessageEmbedThumbnail struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

// MessageEmbedVideo is a part of a MessageEmbed struct.
type MessageEmbedVideo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// MessageEmbedProvider is a part of a MessageEmbed struct.
type MessageEmbedProvider struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// MessageEmbedAuthor is a part of a MessageEmbed struct.
type MessageEmbedAuthor struct {
	URL          string `json:"url,omitempty"`
	Name         string `json:"name"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

// MessageEmbedField is a part of a MessageEmbed struct.
type MessageEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type MessageAllowedMentions struct {
	Parse       []string `json:"parse"`
	Roles       []string `json:"roles,omitempty"`
	Users       []string `json:"users,omitempty"`
	RepliedUser bool     `json:"replied_user"`
}

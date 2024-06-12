package controllers

import (
	"os"

	"github.com/bartosian/notibot/internal/core/ports"
	"github.com/bartosian/notibot/pkg/l0g"
)

type NotifierController struct {
	notifierGateway ports.NotifierGateway
	discordGateway  ports.DiscordGateway
	cherryGateway   ports.CherryGateway
	logger          l0g.Logger
	targetChannel   string
	serverID        string
}

func NewNotifierController(
	notifierGateway ports.NotifierGateway,
	discordGateway ports.DiscordGateway,
	cherryGateway ports.CherryGateway,
	logger l0g.Logger,
) ports.NotifierController {
	return &NotifierController{
		notifierGateway: notifierGateway,
		discordGateway:  discordGateway,
		cherryGateway:   cherryGateway,
		logger:          logger,
		targetChannel:   os.Getenv("DISCORD_CHANNEL"),
		serverID:        os.Getenv("CHERRY_SERVER_ID"),
	}
}

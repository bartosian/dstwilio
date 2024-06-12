package discordgw

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/bartosian/notibot/internal/core/ports"
	"github.com/bartosian/notibot/pkg/l0g"
)

type Gateway struct {
	session *discordgo.Session
	logger  l0g.Logger
}

func NewDiscordGateway(logger l0g.Logger) (ports.DiscordGateway, error) {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		err = fmt.Errorf("error creating discord session: %w", err)

		logger.Error("error creating discord session:", err, nil)

		return nil, err
	}

	return &Gateway{
		session: session,
		logger:  logger,
	}, nil
}

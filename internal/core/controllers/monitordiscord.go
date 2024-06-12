package controllers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/sync/errgroup"

	"github.com/bartosian/notibot/internal/core/ports"
)

const (
	callVoiceTemplate   = "RECEIVED MESSAGE FROM %s IN %s CHANNEL"
	messageTextTemplate = "📢 [ RECEIVED MESSAGE FROM %s IN %s CHANNEL ]"
	userName            = "UptimeRobot"
)

// MonitorDiscord starts monitoring the Discord channel for new messages by adding the message handler.
// It returns an error if there is an issue adding the message handler.
func (c *NotifierController) MonitorDiscord() error {
	var errGroup errgroup.Group

	errGroup.Go(func() error {
		messageHandler := c.newMessageHandler()

		if err := c.discordGateway.AddHandler(messageHandler); err != nil {
			c.logger.Error("error adding message handler", err, nil)
			return err
		}

		return nil
	})

	if err := errGroup.Wait(); err != nil {
		c.logger.Error("unexpected error occurred while monitoring Discord", err, nil)
		return err
	}

	return nil
}

func (c *NotifierController) newMessageHandler() ports.MessageHandler {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		messageChannel, err := s.Channel(m.ChannelID)
		if err != nil {
			c.logger.Error("error getting channel by ID", err, map[string]interface{}{"channelID": m.ChannelID})
			return
		}

		c.logger.Info("received new message", map[string]interface{}{
			"author":  m.Author.Username,
			"channel": messageChannel.Name,
			"content": m.Content,
		})

		if messageChannel.Name == c.targetChannel {
			if err := c.handleNewMessage(m.Author.Username); err != nil {
				c.logger.Error("error handling new message", err, map[string]interface{}{
					"author": m.Author.Username,
				})
			}
		}
	}
}

func (c *NotifierController) handleNewMessage(username string) error {
	if username == userName {
		if err := c.cherryGateway.RestartServerByID(c.serverID); err != nil {
			return fmt.Errorf("error restarting server: %w", err)
		}
	}

	if err := c.notifierGateway.CreateCall(buildCallVoice(username, c.targetChannel)); err != nil {
		return fmt.Errorf("error creating call: %w", err)
	}

	if err := c.notifierGateway.SendMessage(buildMessageText(username, c.targetChannel)); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func buildCallVoice(username, channel string) string {
	return fmt.Sprintf(callVoiceTemplate, username, channel)
}

func buildMessageText(username, channel string) string {
	return fmt.Sprintf(messageTextTemplate, username, channel)
}

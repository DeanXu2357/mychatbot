package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	session *discordgo.Session
}

func New(dcToken string) (*Handler, error) {
	ds, err := discordgo.New("Bot " + dcToken)
	if err != nil {
		return nil, err
	}

	ds.AddHandler(messageCreate)

	return &Handler{
		session: ds,
	}, nil
}

func (h *Handler) Handle() error {
	go func() {
		if err := h.session.Open(); err != nil {
			log.Panic(err)
		}
	}()

	return nil
}

func (h *Handler) Close() error {
	log.Println("discord handler close")
	return h.session.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

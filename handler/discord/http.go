package discord

import (
	"context"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/DeanXu2357/mychatbot/llm"
)

type Handler struct {
	session *discordgo.Session
	ollama  *llm.Ollama
}

func New(dcToken string, ollama *llm.Ollama) (*Handler, error) {
	ds, err := discordgo.New("Bot " + dcToken)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		session: ds,
		ollama:  ollama,
	}

	ds.AddHandler(h.interaction(context.Background()))

	return h, nil
}

func (h *Handler) AddHandler(handler any) {
	h.session.AddHandler(handler)
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

func (h *Handler) interaction(ctx context.Context) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// If the message is "ping" reply with "Pong!"
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}

		if len(m.Content) > 4 && m.Content[:4] == "/llm" {
			gctx, ok := ctx.Value("generateContext").([]int)
			if !ok {
				gctx = []int{}
			}

			resp, generateCTX, errG := h.ollama.Generate(ctx, gctx, m.Content[3:])
			if errG != nil {
				log.Print(errG)
				s.ChannelMessageSend(m.ChannelID, "Someone tell poyu, there is a problem with my AI.")
			}

			ctx = context.WithValue(ctx, "generateContext", generateCTX)

			s.ChannelMessageSend(m.ChannelID, resp)
		}
	}
}

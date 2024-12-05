package discord

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/DeanXu2357/mychatbot/llm"
	"github.com/DeanXu2357/mychatbot/service/probe"
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

	ds.ChannelMessageSend("1309066554201079809", "Hello, here is your chatbot")

	return h, nil
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

func (h *Handler) MonitorTailscaleService(ctx context.Context, name string) {
	notifier := probe.NewTailscaleNotifier(name) //"sony-xq-dq72")
	okCh := notifier.OK()
	notOkCh := notifier.NotOK()
	ok := true

	go func() {
		for {
			select {
			case <-okCh:
				if !ok {
					h.session.ChannelMessageSend("1309066554201079809", fmt.Sprintf("TailScale(%s) is OK at %s", name, time.Now().Local().String()))
					ok = true
				}
			case <-notOkCh:
				h.session.ChannelMessageSend("1309066554201079809", fmt.Sprintf("TailScale(%s) is not OK at %s", name, time.Now().Local().String()))
				ok = false
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (h *Handler) Shutdown() {
	h.session.ChannelMessageSend("1309066554201079809", "Bye bye ~")
}

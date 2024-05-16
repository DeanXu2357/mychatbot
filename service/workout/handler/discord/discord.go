package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/DeanXu2357/mychatbot/service/workout"
)

const (
	CommandPing      = "ping"
	CommandHelp      = "help"
	CommandListEvent = "listevent"
	CommandAddEvent  = "addevent"
	CommandDelEvent  = "rmevent"
	CommandListRec   = "listrecord"
	CommandAddRec    = "addrecord"
	CommandDelRec    = "rmrecord"
)

type Handler struct {
	Record workout.RecordEditor
	Event  workout.EventEditor
}

func (h *Handler) HandleWorkoutRecord(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// fixme: use context from session
	ctx := context.Background()

	if len(m.Content) > 7 && m.Content[:7] == "workout" {
		content := m.Content[8:]

		split := strings.Split(content, " ")
		command := split[0]

		switch command {
		case CommandPing:
			s.ChannelMessageSend(m.ChannelID, "Pong! Workout module is working!")
		case CommandListEvent:
			events, errE := h.Event.Events(ctx, m.Author.ID)
			if errE != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to list events: %v", errE))
			}

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Events: %v", events))
		case CommandAddEvent:
			userID := m.Author.ID

			eventName := split[1]

			tagsString := split[2]
			tags := strings.Split(tagsString, ",")

			if _, err := h.Event.Create(ctx, userID, eventName, tags); err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to create event: %v", err))
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Event %s created", eventName))
			}
		case CommandDelEvent:
		case CommandListRec:
		case CommandAddRec:
		case CommandDelRec:
		case CommandHelp:
		default:
		}
	}
}

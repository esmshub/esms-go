package commentary

import (
	"github.com/esmshub/esms-go/engine/events"
)

type CommentaryEventHandler struct {
	provider CommentaryProvider
}

func NewEventHandler(provider CommentaryProvider) events.EventHandler {
	return &CommentaryEventHandler{
		provider: provider,
	}
}

func (h *CommentaryEventHandler) Handle(event events.Event) ([]events.Event, error) {
	err := h.provider.WriteCommentary(event)
	return []events.Event{}, err
}

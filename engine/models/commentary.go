package models

type CommentaryProvider interface {
	GetEventText(event string, args ...any) string
}

package common

const (
	MATCHTYPE_LEAGUE = iota
	MATCHTYPE_CUP
	MATCHTYPE_FRIENDLY
)

type MATCHTYPE int

type CommentaryProvider interface {
	GetEventText(event string, args ...any) string
}

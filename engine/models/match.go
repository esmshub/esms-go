package models

import "github.com/esmshub/esms-go/engine/common"

type Match struct {
	HomeTeam   *TeamConfig
	AwayTeam   *TeamConfig
	Commentary common.CommentaryProvider
}

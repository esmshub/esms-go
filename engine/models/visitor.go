package models

type Visitor interface {
	VisitTeam(*TeamConfig)
}

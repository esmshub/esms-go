package models

type Observer interface {
	Update(string, Subject)
}

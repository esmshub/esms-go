package store

type Store interface {
	Save(obj, options any) error
}

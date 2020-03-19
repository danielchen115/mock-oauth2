package db

type Collection interface {
	Insert(data map[string]interface{}) uint64
}

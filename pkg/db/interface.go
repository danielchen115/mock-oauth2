package db

type Collection interface {
	Connect(URI string, database string) Collection
	Insert(data map[string]interface{}) uint64
}
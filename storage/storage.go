package storage

// Database is an interface that describes the DB operations
type Database interface {
	ListSubscribers() (*model.Subscriber, error)
	AddSubscriber(*model.Subscriber) (*model.Subscriber, error)
}
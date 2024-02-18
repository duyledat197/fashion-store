package pubsub

type Subscriber interface {
	Subscribe(topic string, subscribeFn func(key, value []byte))
}

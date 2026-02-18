//go:build !kafka

package pubsub

type Publisher interface {
	PublishBlock([]byte) error
	PublishTx([]byte) error
	StartEventHandler()
	Close()
}

type StubPublisher struct{}

func NewPublisher(url string) (Publisher, error) {
	return &StubPublisher{}, nil
}

func (p *StubPublisher) PublishBlock(data []byte) error {
	return nil
}

func (p *StubPublisher) PublishTx(data []byte) error {
	return nil
}

func (p *StubPublisher) StartEventHandler() {}

func (p *StubPublisher) Close() {}
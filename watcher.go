package watcher

import (
	"context"
	"fmt"
	"github.com/tryfix/kstream/consumer"
	"github.com/tryfix/kstream/data"
	"github.com/tryfix/log"
	"runtime"
	"sync"
)

// Watcher implements casbin policy update watcher
type Watcher struct {
	ctx      context.Context
	connMu   *sync.RWMutex
	consumer consumer.PartitionConsumer
	//consumer consumer.Consumer
	*config
}

type config struct {
	topic        string
	synced       chan bool
	brokers      []string
	logger       log.Logger
	callbackFunc ExecuteFunc
}

// ExecuteFunc executed on each kafka message receive
type ExecuteFunc func(record *data.Record) error

// SetLogger set the logger for watcher
func (c *config) SetLogger(l log.Logger) {
	c.logger = l
}

func (c *config) apply() error {
	if c.logger == nil {
		c.logger = log.NewNoopLogger()
	}
	if c.callbackFunc == nil {
		return fmt.Errorf(`callback function not found`)
	}
	if len(c.brokers) == 0 {
		return fmt.Errorf(`kafka brokers not found`)
	}
	if c.topic == "" {
		return fmt.Errorf(`kafka topic not found`)
	}
	return nil
}

// NewConfig return a config object to use when creating the KafkaWatcher instance
// this will not take logger as config since its optional. If logging is required by the
// they can set it calling SetLogger function.
func NewConfig(topic string, brokers []string, callbackFunc ExecuteFunc) *config {
	c := new(config)
	c.topic = topic
	c.brokers = brokers
	c.callbackFunc = callbackFunc
	c.synced = make(chan bool, 1)
	return c
}

// New creates a new kafka watcher instance to use for casbin watcher
func New(ctx context.Context, cfg *config) (*Watcher, error) {
	if err := cfg.apply(); err != nil {
		return nil, err
	}

	w := &Watcher{
		ctx:    ctx,
		config: cfg,
		connMu: &sync.RWMutex{},
	}

	runtime.SetFinalizer(w, finalizer)

	//w.producer = NewProducer(cfg.brokers, cfg.logger)
	w.consumer = newPartitionConsumer(cfg.brokers, cfg.logger)
	_ = w.setUpdateCallback(cfg.callbackFunc)

	go w.Subscribe(w.config.synced)

	<-w.config.synced

	return w, nil
}

func (w *Watcher) setUpdateCallback(callbackFunc ExecuteFunc) error {
	w.callbackFunc = func(record *data.Record) error {
		w.connMu.RLock()
		defer w.connMu.RUnlock()
		return callbackFunc(record)
	}
	return nil
}

// Close stops and releases the watcher, the callback function will not be called any more.
func (w *Watcher) Close() {
	finalizer(w)
}

func finalizer(w *Watcher) {
	w.connMu.Lock()
	defer w.connMu.Unlock()

	err := w.consumer.Close()
	if err != nil {
		w.logger.Error(fmt.Errorf("unable to close the consumer, err:%v", err))
	}

	w.callbackFunc = nil
}

// Subscribe function will listen to all incoming messages from provided kafka topic
// and will call callbackFunc on each message.
func (w *Watcher) Subscribe(synced chan<- bool) {
	// If needed to use multiple partitions
	// client, err := sarama.NewClient(w.brokers, sarama.NewConfig())
	// p, err := client.Partitions(w.topic)

	partition, err := w.consumer.Consume(w.topic, 0, consumer.Earliest)
	if err != nil {
		w.logger.Fatal(err)
	}
	for event := range partition {
		switch record := event.(type) {
		case *data.Record:
			err := w.callbackFunc(record)
			if err != nil {
				w.logger.Error(fmt.Errorf(`error executing the message processor func, err:%v`, err))
				continue
			}
		case *consumer.PartitionEnd:
			synced <- true
		}

	}
}

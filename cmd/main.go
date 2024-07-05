package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/CaelRowley/geth-indexer-service/cmd/server"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
)

var topic = "testTopic"

func produce() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		slog.Error("failed to create kafka producer", "err", err)
		os.Exit(1)
	}
	defer p.Close()

	for i := 0; i < 1000; i++ {
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("foo: " + strconv.Itoa(i)),
			Headers:        []kafka.Header{{Key: "testHeader", Value: []byte("bar")}},
		}, nil)
		if err != nil {
			slog.Error("failed to produce kafka message", "err", err)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func consume() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        "localhost:9092",
		"group.id":                 "testGroup",
		"session.timeout.ms":       6000,
		"auto.offset.reset":        "earliest",
		"enable.auto.offset.store": false,
	})
	if err != nil {
		slog.Error("failed to create kafka consumer", "err", err)
		os.Exit(1)
	}
	defer c.Close()

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		slog.Error("failed to subscribe to kafka topics", "err", err)
		os.Exit(1)
	}

	for {
		ev := c.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("Partition: %s Message: %s\n", e.TopicPartition, string(e.Value))
			if e.Headers != nil {
				fmt.Printf("Headers: %v\n", e.Headers)
			}
			_, err := c.StoreMessage(e)
			if err != nil {
				slog.Error("failed to store kafka offset after message", "err", err)
			}
		case kafka.Error:
			slog.Error("failed to store kafka offset after message", "err", e, "code", e.Code())
			if e.Code() == kafka.ErrAllBrokersDown {
				break
			}
		}
	}
}

func main() {
	go produce()
	go consume()

	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file", "err", err)
		os.Exit(1)
	}

	var serverCfg server.ServerConfig
	flag.StringVar(&serverCfg.Port, "port", "8080", "Port where the service will run")
	flag.BoolVar(&serverCfg.Sync, "sync", false, "Sync blocks on node with db")
	flag.Parse()

	slog.Info("flags set", "Port", serverCfg.Port, "Sync", serverCfg.Sync)

	s, err := server.New(serverCfg)
	if err != nil {
		slog.Error("failed to create server", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := s.Start(ctx); err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}

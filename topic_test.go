package main

import (
	"testing"
	"time"
)

func TestNewTopic(t *testing.T) {
	partitionCount := 3
	topic := NewTopic(partitionCount)

	if len(topic.partitions) != partitionCount {
		t.Errorf("Expected %d partitions, but got %d", partitionCount, len(topic.partitions))
	}

	for i := 0; i < partitionCount; i++ {
		if len(topic.partitions[i]) != 0 {
			t.Errorf("Expected partition %d to be empty, but it has %d messages", i, len(topic.partitions[i]))
		}
	}
}

func TestProduce(t *testing.T) {
	topic := NewTopic(3)
	msg := "Test message"
	partition := 1

	topic.produce(msg, partition)

	if len(topic.partitions[partition]) != 1 {
		t.Errorf("Expected 1 message in partition %d, but got %d", partition, len(topic.partitions[partition]))
	}

	if topic.partitions[partition][0] != msg {
		t.Errorf("Expected message '%s' in partition %d, but got '%s'", msg, partition, topic.partitions[partition][0])
	}
}

func TestConsume(t *testing.T) {
	topic := NewTopic(3)
	msg := "Test message"
	partition := 1

	topic.produce(msg, partition)

	done := make(chan struct{})
	go func() {
		topic.consume(partition, done)
	}()

	// Даём время консьюмеру обработать сообщение
	time.Sleep(2 * time.Second)

	// Завершаем работу консьюмера
	close(done)

	if len(topic.partitions[partition]) != 0 {
		t.Errorf("Expected 0 messages in partition %d, but got %d", partition, len(topic.partitions[partition]))
	}
}

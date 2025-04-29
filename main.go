package main

import (
	"fmt"
	"sync"
	"time"
)

// Имитация топика Kafka с партициями
type Topic struct {
	partitions map[int][]string // Сообщения распределяются по партициям
	mu         sync.Mutex
}

// Создаём новый "топик"
func NewTopic(partitionCount int) *Topic {
	t := &Topic{partitions: make(map[int][]string)}
	for i := 0; i < partitionCount; i++ {
		t.partitions[i] = []string{} // Создаём пустые партиции
	}
	return t
}

// Функция продюсера: отправляем сообщения в партиции
func (t *Topic) produce(msg string, partition int) {
	t.mu.Lock()
	t.partitions[partition] = append(t.partitions[partition], msg)
	t.mu.Unlock()
	fmt.Printf("Продюсер отправил в партицию %d: %s\n", partition, msg)
}

// Функция консьюмера: читаем сообщения из партиций
func (t *Topic) consume(partition int) {
	for {
		t.mu.Lock()
		if len(t.partitions[partition]) > 0 {
			msg := t.partitions[partition][0]
			t.partitions[partition] = t.partitions[partition][1:] // Удаляем прочитанное сообщение
			fmt.Printf("Консьюмер обработал из партиции %d: %s\n", partition, msg)
		}
		t.mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// Создаём топик с 3 партициями
	topic := NewTopic(3)

	// Запускаем консьюмеров для каждой партиции
	for i := 0; i < 3; i++ {
		go topic.consume(i)
	}

	// Отправляем сообщения в разные партиции
	for i := 1; i <= 5; i++ {
		topic.produce(fmt.Sprintf("Заказ %d", i), i%3) // Разделяем по 3 партициям
		time.Sleep(500 * time.Millisecond)
	}
}

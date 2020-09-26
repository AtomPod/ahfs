package setting

import (
	"time"

	"github.com/czhj/ahfs/modules/log"
	"github.com/spf13/viper"
)

type QueueSetting struct {
	Type         string
	Name         string
	Length       int
	BatchLength  int
	QueueName    string
	BoostWorker  int
	MaxWorkers   int
	Workers      int
	BlockTimeout time.Duration
	BoostTimeout time.Duration
}

var (
	Queue *QueueSetting
)

func newQueueService() {
	prefix := "queue.default"
	viper.SetDefault("queue.default", map[string]interface{}{
		"type":          "channel",
		"name":          "default",
		"batch_length":  20,
		"queue_name":    "default",
		"block_timeout": time.Duration(1) * time.Second,
		"boost_timeout": time.Duration(3) * time.Minute,
		"length":        20,
		"boost_worker":  5,
		"max_workers":   10,
		"workers":       1,
	})
	queueConfig := viper.Sub(prefix)
	Queue = parseQueueSetting("", queueConfig)
	log.Info("Queue Service Enabled")
}

func parseQueueSetting(name string, q *viper.Viper) *QueueSetting {
	queue := &QueueSetting{}
	queue.Type = q.GetString("type")
	queue.Name = q.GetString("name")
	queue.BatchLength = q.GetInt("batch_length")
	queue.QueueName = q.GetString("queue_name")
	queue.BlockTimeout = q.GetDuration("block_timeout")
	queue.BoostTimeout = q.GetDuration("boost_timeout")
	queue.Length = q.GetInt("length")
	queue.BoostWorker = q.GetInt("boost_worker")
	queue.MaxWorkers = q.GetInt("max_workers")
	queue.Workers = q.GetInt("workers")

	return queue
}

func GetQueueSetting(name string) *QueueSetting {
	prefix := "queue." + name

	viper.SetDefault(prefix, map[string]interface{}{
		"type":          Queue.Type,
		"batch_length":  Queue.BatchLength,
		"block_timeout": Queue.BlockTimeout,
		"boost_timeout": Queue.BoostTimeout,
		"length":        Queue.Length,
		"boost_worker":  Queue.BoostWorker,
		"max_workers":   Queue.MaxWorkers,
		"workers":       Queue.Workers,
		"name":          name,
	})

	queueConfig := viper.Sub(prefix)
	return parseQueueSetting(name, queueConfig)
}

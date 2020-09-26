package queue

import (
	"encoding/json"
	"fmt"

	"github.com/czhj/ahfs/modules/log"
	"github.com/czhj/ahfs/modules/setting"
	"go.uber.org/zap"
)

func validType(t string) (Type, error) {
	if len(t) == 0 {
		return ChannelQueueType, nil
	}

	for _, typ := range ResgiteredTypes() {
		if t == string(typ) {
			return typ, nil
		}
	}
	return ChannelQueueType, fmt.Errorf("Unknow queue type: %s defaulting to: %s", t, string(ChannelQueueType))
}

func getQueueSetting(name string) (*setting.QueueSetting, []byte) {
	q := setting.GetQueueSetting(name)

	queueParam := make(map[string]interface{})
	queueParam["name"] = q.Name
	queueParam["batchLength"] = q.BatchLength
	queueParam["boostTimeout"] = q.BoostTimeout
	queueParam["blockTimeout"] = q.BlockTimeout
	queueParam["boostWorker"] = q.BoostWorker
	queueParam["queueLength"] = q.Length
	queueParam["queueName"] = q.QueueName
	queueParam["maxWorkers"] = q.MaxWorkers
	queueParam["workers"] = q.Workers

	params, _ := json.Marshal(queueParam)

	return q, params
}

func CreateQueue(name string, handle HandleFunc, exemplar interface{}) Queue {
	q, cfg := getQueueSetting(name)

	typ, err := validType(q.Type)
	if err != nil {
		log.Error("Invalid type provided for queue, use default", zap.String("type", q.Type), zap.String("name", name), zap.String("default", string(typ)), zap.Error(err))
	}

	queue, err := NewQueue(typ, handle, cfg, exemplar)
	if err != nil {
		log.Error("Unable to create queue", zap.String("name", name), zap.Error(err))
		return nil
	}
	return queue
}

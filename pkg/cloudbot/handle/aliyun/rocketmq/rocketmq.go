package rocketmq

import (
	"sync"

	"github.com/atompi/cloudbot/pkg/aliyun/rocketmq"
	"github.com/atompi/cloudbot/pkg/cloudbot/handle/options"
	"github.com/atompi/cloudbot/pkg/dataio"
	"github.com/atompi/cloudbot/pkg/utils"
	"go.uber.org/zap"
)

func CreateTopicHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	data, err := utils.DataToMap(&res)
	if err != nil {
		zap.S().Errorf("data convert error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	for _, row := range *data {
		wg.Add(1)
		ch <- 1
		go rocketmq.CreateTopic(ch, &wg, t, row["instanceId"], row["topicName"], row["messageType"], row["remark"])
	}

	wg.Wait()
	return nil
}

func CreateConsumerGroupHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	data, err := utils.DataToMap(&res)
	if err != nil {
		zap.S().Errorf("data convert error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	for _, row := range *data {
		wg.Add(1)
		ch <- 1
		go rocketmq.CreateConsumerGroup(ch, &wg, t, row["instanceId"], row["consumerGroupId"], row["deliveryOrderType"], row["consumeRetryPolicy"], row["maxRetryTimes"], row["deadLetterTargetTopic"], row["remark"])
	}

	wg.Wait()
	return nil
}

func UpdateConsumerGroupHandler(t options.TaskOptions) error {
	res, err := dataio.InputCSV(t.Input)
	if err != nil {
		zap.S().Errorf("input error: %v", err)
		return err
	}

	data, err := utils.DataToMap(&res)
	if err != nil {
		zap.S().Errorf("data convert error: %v", err)
		return err
	}

	wg := sync.WaitGroup{}
	ch := make(chan int, t.Threads)

	for _, row := range *data {
		wg.Add(1)
		ch <- 1
		go rocketmq.UpdateConsumerGroup(ch, &wg, t, row["instanceId"], row["consumerGroupId"], row["deliveryOrderType"], row["consumeRetryPolicy"], row["maxRetryTimes"], row["deadLetterTargetTopic"], row["remark"])
	}

	wg.Wait()
	return nil
}

package core

import (
	"sort"
)

/*******************************************************************************
  Metric统计功能：从一些用户ID得到这些人身上的统计指标
  *******************************************************************************/

type Metric struct {
	TagId   uint32  `json:"tag_id"`
	Options Options `json:"options"`
}

type Option struct {
	Id      int32              `json:"id"`
	Count   uint32             `json:"-"`
	Percent float32            `json:"percent"`
	TGI     float32            `json:"tgi"`
	Info    *DisplayOptionInfo `json:"info,omitempty"`
}

type DisplayOptionInfo struct {
	Category   string `json:"category"`
	TagName    string `json:"tag_name"`
	OptionName string `json:"option_name"`
}

// 下面这些结构体和函数为了方便排序
type Options []Option

func (options Options) Len() int {
	return len(options)
}
func (options Options) Swap(i, j int) {
	options[i], options[j] = options[j], options[i]
}
func (options Options) Less(i, j int) bool {
	// 为了从大到小排序，这实际上实现的是More的功能
	return options[i].Percent > options[j].Percent
}

type OptionsByTGI []Option

func (options OptionsByTGI) Len() int {
	return len(options)
}
func (options OptionsByTGI) Swap(i, j int) {
	options[i], options[j] = options[j], options[i]
}
func (options OptionsByTGI) Less(i, j int) bool {
	// 为了从大到小排序，这实际上实现的是More的功能
	return options[i].TGI > options[j].TGI
}

// 从满足query的用户里得到指标（metricsArray中指定）的统计结果
// 包括满足条件的用户总数和各个metric的指标
func (service *LookupService) GetMetricStats(query []string, metricsArray []uint32, sortByTGI bool) (int, []Metric, error) {
	totalUser := service.GetUserCount(query)

	var output []Metric
	for _, metric := range metricsArray {
		if _, ok := service.tagInfos[metric]; !ok {
			continue
		}

		node := service.tagInfos[metric]
		ch := make(chan Option, len(node.Options))

		var options Options
		for _, optionId := range node.Options {
			service.metricWorkerChannel <- MetricWorkerOption{
				totalUser:     totalUser,
				query:         query,
				metric:        metric,
				optionId:      optionId,
				optionChannel: ch,
			}
		}
		for _ = range node.Options {
			options = append(options, <-ch)
		}

		// 排序后添加
		if sortByTGI {
			sort.Sort(OptionsByTGI(options))
		} else {
			sort.Sort(options)
		}
		output = append(output, Metric{
			TagId:   metric,
			Options: options,
		})
	}

	return totalUser, output, nil
}

type MetricWorkerOption struct {
	totalUser     int
	query         []string
	metric        uint32
	optionId      int32
	optionChannel chan Option
}

func (service *LookupService) InitMetricWorker() {
	numWorker := service.initOptions.MetricWorkerChannelBufferLength
	service.metricWorkerChannel = make(chan MetricWorkerOption, numWorker)
	for i := 0; i < numWorker; i++ {
		go service.metricWorker()
	}
}

func (service *LookupService) metricWorker() {
	for {
		workerOption := <-service.metricWorkerChannel

		optionKey := GetOptionKey(workerOption.metric, workerOption.optionId)
		query := append(workerOption.query, optionKey)
		optionUserCount := service.GetUserCount(query)

		percent := float32(optionUserCount) / float32(workerOption.totalUser)
		tgi := float32(0.0)
		if v, ok := service.optionCounts[optionKey]; ok {
			if v != 0 {
				tgi = percent * float32(len(service.userIds)) / float32(v) * 100.0
			}
		}

		workerOption.optionChannel <- Option{
			Count:   uint32(optionUserCount),
			Percent: percent,
			TGI:     tgi,
			Id:      workerOption.optionId,
		}
	}
}

package core

import (
	"flag"
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
	"runtime"
)

var (
	dataFiles     = flag.String("data_files", "", "包含用户属性的数据文件，支持wildcard")
	emptyFile     = flag.String("empty_file", "../data/empty.txt", "空文件，用于启动tag搜索引擎")
	tagOptionFile = flag.String("tag_option_file", "../data/tag_option.csv", "从哪个日期的数据中搜索")
)

type LookupService struct {
	searcher engine.Engine // 用于搜索人群
	userIds  []uint64      // 全部人群的ID

	optionInfos      map[string]OptionInfo // 快速查找option节点信息，key为<tag id>:<option id>格式
	optionInfosIndex []string              // 用于生成倒排索引中的docId
	optionSearcher   engine.Engine         // 用于搜索option

	tagInfos      map[uint32]TagInfo // 快速查找tag节点信息，key为<tag id>格式
	tagInfosIndex []uint32           // 用于生成倒排索引中的docId
	tagSearcher   engine.Engine      // 用于搜索tag

	optionCounts map[string]uint32 // 全部人群中各个option的人数
	tagCounts    map[uint32]uint32 // 全部人群中各个tag的人数

	optionUserPair map[uint32]int // 用于统计option-user的对数，基本正比于该tag需要的内存量

	tagWhitelist map[uint32]bool

	metricWorkerChannel chan MetricWorkerOption

	initOptions LookupServiceInitOptions // 初始化参数
}

type LookupServiceInitOptions struct {
	UserSearcherInitOptions         types.EngineInitOptions
	DataFiles                       string
	MetricWorkerChannelBufferLength int
	TagOptionFile                   string
	TagSearcherInitOptions          types.EngineInitOptions
}

// 使用前必须初始化
func (service *LookupService) Init() {
	service.initOptions = LookupServiceInitOptions{
		UserSearcherInitOptions: types.EngineInitOptions{
			NotUsingSegmenter: true,
			IndexerInitOptions: &types.IndexerInitOptions{
				IndexType: types.DocIdsIndex,
			},
			UsePersistentStorage: false,
			NumShards:            runtime.NumCPU() * 8,
			RankerBufferLength:   runtime.NumCPU() * 2,
			IndexerBufferLength:  runtime.NumCPU() * 2,
			NumSegmenterThreads:  runtime.NumCPU() * 2,
		},
		DataFiles:                       *dataFiles,
		MetricWorkerChannelBufferLength: runtime.NumCPU() * 2,
		TagOptionFile:                   *tagOptionFile,
		TagSearcherInitOptions: types.EngineInitOptions{
			SegmenterDictionaries: *emptyFile,
			StopTokenFile:         *emptyFile,
			NumShards:             8,
			IndexerInitOptions: &types.IndexerInitOptions{
				IndexType: types.LocationsIndex,
			},
		},
	}
	service.loadTagWhitelist()
	service.InitMetricWorker()
	service.loadTagOptionSearcher()
	service.loadUserSearcher()
}

// 关闭资源
func (service *LookupService) Close() {
	service.searcher.Close()
}

func (service *LookupService) GetTotalNumUsers() int {
	return len(service.userIds)
}

package json_service

import (
	core "../core"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type MetricParas struct {
	// 仅统计满足条件的人数，不计算metric
	CountUserOnly bool `schema:"count_users_only"`
	// query: 搜索人群使用的条件，格式 <tag_id><tag_option>,<tag_id><tag_option>,...
	Query string `schema:"query"`
	// metrics: 显示的指标，格式 <tag_id>,<tag_id>,...
	Metrics string `schema:"metrics"`
	// show_info: 不为空时打印option名等信息
	ShowInfo bool `schema:"show_info"`
	// short_by_tgi: 不为空时按照TGI倒排序，为空时按照百分比排序
	SortByTGI bool `schema:"sort_by_tgi"`
	// min_percent: 过滤的最低百分比
	MinPercent float32 `schema:"min_percent"`
	// min_count: 规律的最低count
	MinCount uint32 `schema:"min_count"`
	// min_tgi: 过滤的最低TGI
	MinTGI float32 `schema:"min_tgi"`
}

type MetricJsonResponse struct {
	TotalNumUsers int           `json:"total_num_users"`
	Metrics       []core.Metric `json:"metrics"`
}

// 生成metric
// JSON参数见MetricParas结构体
func (service *JsonService) MetricJsonRpcService(w http.ResponseWriter, req *http.Request) {
	var paras MetricParas
	if err := service.decoder.Decode(&paras, req.URL.Query()); err != nil {
		WriteErrResponse(w, err)
		return
	}

	// 搜索全部满足条件的用户ID
	q, err := BuildQuery(paras.Query)
	if err != nil {
		WriteErrResponse(w, err)
		return
	}
	if paras.CountUserOnly {
		count := service.lookupService.GetUserCount(q)
		totalNumUsers := count
		response, _ := json.Marshal(&MetricJsonResponse{
			TotalNumUsers: totalNumUsers,
		})
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(response))
		return
	}

	// 统计metric指标
	var m []uint32
	if m, err = BuildMetrics(paras.Metrics); err != nil {
		WriteErrResponse(w, err)
		return
	}
	total, result, err := service.lookupService.GetMetricStats(q, m, paras.SortByTGI)
	if err != nil {
		WriteErrResponse(w, err)
		return
	}
	log.Printf("满足条件的id数: %d", total)

	var filteredResult []core.Metric
	for i, _ := range result {
		var filteredMetric core.Metric
		filteredMetric.TagId = result[i].TagId
		for _, option := range result[i].Options {
			if option.Percent >= paras.MinPercent &&
				option.Count >= paras.MinCount &&
				option.TGI >= paras.MinTGI {
				filteredMetric.Options = append(filteredMetric.Options, option)
			}
		}
		filteredResult = append(filteredResult, filteredMetric)
	}

	// 整理为输出格式
	if paras.ShowInfo {
		for i, _ := range filteredResult {
			for j, _ := range filteredResult[i].Options {
				info := service.lookupService.GetOptionInfo(
					filteredResult[i].TagId, filteredResult[i].Options[j].Id)
				filteredResult[i].Options[j].Info = &core.DisplayOptionInfo{
					Category:   info.Category,
					TagName:    info.TagName,
					OptionName: info.OptionName,
				}
			}
		}
	}
	response, _ := json.Marshal(&MetricJsonResponse{
		TotalNumUsers: total,
		Metrics:       filteredResult,
	})
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(response))
}

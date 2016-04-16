package json_service

import (
	core "../core"
	"encoding/json"
	"io"
	"net/http"
)

type TagParas struct {
	// query: 搜索字符串
	Query string `schema:"query"`

	TagId uint32 `schema:"tag_id"`
}

type TagJsonResponse struct {
	Tags []core.TagInfo `json:"tags"`
}

// 搜索满足query条件的tag
// JSON参数见TagParas结构体
func (service *JsonService) TagJsonRpcService(w http.ResponseWriter, req *http.Request) {
	var paras TagParas
	if err := service.decoder.Decode(&paras, req.URL.Query()); err != nil {
		WriteErrResponse(w, err)
		return
	}

	if paras.TagId != 0 {
		response, _ := json.Marshal(&TagJsonResponse{
			Tags: []core.TagInfo{*service.lookupService.GetTagInfo(paras.TagId)},
		})
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(response))
		return
	}
	tags := service.lookupService.SearchTag(paras.Query)

	// 整理为输出格式
	response, _ := json.Marshal(&TagJsonResponse{
		Tags: tags,
	})
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(response))
}

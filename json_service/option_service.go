package json_service

import (
	core "../core"
	"encoding/json"
	"io"
	"net/http"
)

type OptionParas struct {
	// query: 搜索字符串
	Query string `schema:"query"`

	OptionKey string `schema:"option_key"`
}

type OptionJsonResponse struct {
	Options []core.OptionInfo `json:"options"`
}

// 搜索满足query条件的option
// JSON参数见OptionParas结构体
func (service *JsonService) OptionJsonRpcService(w http.ResponseWriter, req *http.Request) {
	var paras OptionParas
	if err := service.decoder.Decode(&paras, req.URL.Query()); err != nil {
		WriteErrResponse(w, err)
		return
	}

	if paras.OptionKey != "" {
		response, _ := json.Marshal(&OptionJsonResponse{
			Options: []core.OptionInfo{*service.lookupService.GetOptionInfoWithKey(paras.OptionKey)},
		})
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(response))
		return
	}

	tags := service.lookupService.SearchOption(paras.Query)

	// 整理为输出格式
	response, _ := json.Marshal(&OptionJsonResponse{
		Options: tags,
	})
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(response))
}

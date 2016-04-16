package json_service

import (
	"fmt"
	"io"
	"net/http"
)

// 返回服务器信息
func (service *JsonService) StatsService(w http.ResponseWriter, req *http.Request) {
	content := fmt.Sprintf("总人数 %d<br><br>", service.lookupService.GetTotalNumUsers())

	tagStats := service.lookupService.GetTagStats()
	content += fmt.Sprintf("标签数 %d<br>", len(tagStats))
	content += "<table border=1 style=\"BORDER-COLLAPSE: collapse\">"
	content += "<tr><td><b>标签名</b></td><td><b>标签ID</b></td><td><b>Option数</b></td><td><b>人数</b></td><td><b>option-人</b></td></tr>"
	for _, s := range tagStats {
		content += fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%d</td><td>%d</td><td>%d</td></tr>", s.TagName, s.TagId, s.OptionCount, s.UserCount, s.OptionUserPair)
	}
	content += "</table>"

	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, content)
}

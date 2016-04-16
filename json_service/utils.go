package json_service

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ErrorJsonResponse struct {
	Error string `json:"error"`
}

func BuildMetrics(metricsString string) ([]uint32, error) {
	if metricsString == "" {
		return nil, errors.New("metrics不能为空")
	}
	metricsData := strings.Split(metricsString, ",")
	var metrics []uint32
	for _, v := range metricsData {
		if len(v) == 0 {
			log.Printf("非法metrics")
			continue
		}
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			log.Printf("非法metric: %v", v)
			continue
		}
		metrics = append(metrics, uint32(id))
	}
	return metrics, nil
}

func BuildQuery(query string) ([]string, error) {
	if query == "" {
		return nil, errors.New("query不能为空")
	}
	labels := strings.Split(query, ",")
	return labels, nil
}

func WriteErrResponse(w http.ResponseWriter, err error) {
	response, _ := json.Marshal(&ErrorJsonResponse{err.Error()})
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(response))
	log.Printf("error : %s", err)
}

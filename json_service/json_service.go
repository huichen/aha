package json_service

import (
	core "../core"
	"github.com/gorilla/schema"
)

// 各个service在其余文件中实现
type JsonService struct {
	lookupService core.LookupService
	decoder       *schema.Decoder
}

func (service *JsonService) Init() {
	service.decoder = schema.NewDecoder()
	service.lookupService.Init()
}

func (service *JsonService) Close() {
	service.lookupService.Close()
}

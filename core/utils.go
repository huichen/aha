package core

import (
	"fmt"
)

// 从tag id和option得到option节点的信息
func (service *LookupService) GetOptionInfo(tagId uint32, tagOption int32) *OptionInfo {
	optionKey := GetOptionKey(tagId, tagOption)
	if n, ok := service.optionInfos[optionKey]; ok {
		return &n
	}
	return nil
}

// 从tagId和tagOption得到<tag id>:<tag option>字符串，用于标示一个option
func GetOptionKey(tagId uint32, tagOption int32) string {
	return fmt.Sprintf("%d:%d", tagId, tagOption)
}

func (service *LookupService) GetTagInfo(tagId uint32) *TagInfo {
	if n, ok := service.tagInfos[tagId]; ok {
		return &n
	}
	return nil
}

func (service *LookupService) GetOptionInfoWithKey(optionKey string) *OptionInfo {
	if n, ok := service.optionInfos[optionKey]; ok {
		return &n
	}
	return nil
}

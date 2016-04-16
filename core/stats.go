package core

import (
	"sort"
)

type TagStats struct {
	TagName        string
	TagId          uint32
	UserCount      uint32
	OptionCount    int
	OptionUserPair int
}

func (service *LookupService) GetTagStats() (output []TagStats) {
	for k, v := range service.tagCounts {
		stat := TagStats{}
		stat.TagName = service.tagInfos[k].TagName
		stat.TagId = service.tagInfos[k].TagId
		stat.UserCount = v
		stat.OptionCount = len(service.tagInfos[k].Options)
		stat.OptionUserPair = service.optionUserPair[k]
		output = append(output, stat)
	}
	sort.Sort(TagStatsArray(output))
	return
}

type TagStatsArray []TagStats

func (array TagStatsArray) Len() int {
	return len(array)
}
func (array TagStatsArray) Swap(i, j int) {
	array[i], array[j] = array[j], array[i]
}
func (array TagStatsArray) Less(i, j int) bool {
	return array[i].TagId < array[j].TagId
}

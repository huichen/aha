package core

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/huichen/wukong/types"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*******************************************************************************
  搜索功能：从query搜索得到满足条件tag和option
  *******************************************************************************/

type OptionInfo struct {
	Category   string `json:"category"`
	TagName    string `json:"tag_name"`
	OptionName string `json:"option_name"`

	TagId    uint32 `json:"tag_id"`
	OptionId int32  `json:"option_id"`
}

type TagInfo struct {
	Category string `json:"category"`
	TagName  string `json:"tag_name"`

	TagId   uint32  `json:"tag_id"`
	Options []int32 `json:"-"`
}

func (service *LookupService) loadTagOptionSearcher() error {
	gob.Register(TagScoringFields{})

	service.tagSearcher.Init(service.initOptions.TagSearcherInitOptions)
	service.optionSearcher.Init(service.initOptions.TagSearcherInitOptions)

	service.tagInfos = make(map[uint32]TagInfo)
	service.optionInfos = make(map[string]OptionInfo)

	file, err := os.Open(service.initOptions.TagOptionFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		if len(data) != 5 {
			continue
		}
		node := OptionInfo{}
		node.Category = data[4]
		node.TagName = data[1]
		node.OptionName = data[3]

		var err error
		var tagId uint64
		tagId, err = strconv.ParseUint(data[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		node.TagId = uint32(tagId)
		var optionId int64
		optionId, err = strconv.ParseInt(data[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		node.OptionId = int32(optionId)

		optionKey := GetOptionKey(node.TagId, node.OptionId)
		service.optionInfos[optionKey] = node

		if _, ok := service.tagInfos[node.TagId]; !ok {
			tagNode := TagInfo{
				Category: node.Category,
				TagName:  node.TagName,
				TagId:    node.TagId,
			}
			service.tagInfos[node.TagId] = tagNode
		}
	}

	// 将option加入tag下
	for _, v := range service.optionInfos {
		tagId := v.TagId
		optionId := v.OptionId
		if node, ok := service.tagInfos[tagId]; !ok {
			log.Fatal("无法载入tagid")
		} else {
			node.Options = append(node.Options, optionId)
			service.tagInfos[tagId] = node
		}
	}

	log.Print("添加索引")

	count := 0
	for tagId, node := range service.tagInfos {
		service.tagInfosIndex = append(service.tagInfosIndex, tagId)
		contentOption := node.TagName
		service.tagSearcher.IndexDocument(uint64(count), types.DocumentIndexData{
			Content: contentOption,
			Fields: TagScoringFields{
				TagId: node.TagId,
			},
		})
		count++
	}
	service.tagSearcher.FlushIndex()
	log.Printf("索引了%d个tag\n", len(service.tagInfos))

	count = 0
	for optionKey, node := range service.optionInfos {
		service.optionInfosIndex = append(service.optionInfosIndex, optionKey)
		contentOption := fmt.Sprintf("%s %s", node.TagName, node.OptionName)
		service.optionSearcher.IndexDocument(uint64(count), types.DocumentIndexData{
			Content: contentOption,
			Fields: TagScoringFields{
				TagId:    node.TagId,
				OptionId: node.OptionId,
			},
		})
		count++
	}
	service.optionSearcher.FlushIndex()
	log.Printf("索引了%d个option\n", len(service.optionInfos))

	return nil
}

type TagScoringFields struct {
	TagId    uint32
	OptionId int32
}

type TagScoringCriteria struct {
}

func (criteria TagScoringCriteria) Score(
	doc types.IndexedDocument, fields interface{}) []float32 {
	if reflect.TypeOf(fields) != reflect.TypeOf(TagScoringFields{}) {
		return []float32{}
	}
	tsf := fields.(TagScoringFields)
	output := make([]float32, 3)
	output[0] = -float32(doc.TokenProximity)
	output[1] = -float32(tsf.TagId)
	output[2] = -float32(tsf.OptionId)
	return output
}

func (service *LookupService) SearchTag(query string) []TagInfo {
	output := service.tagSearcher.Search(types.SearchRequest{
		Text: query,
		RankOptions: &types.RankOptions{
			ScoringCriteria: &TagScoringCriteria{},
			OutputOffset:    0,
			MaxOutputs:      1000,
		},
	})

	var tags []TagInfo
	for _, doc := range output.Docs {
		if _, ok := service.tagCounts[service.tagInfosIndex[doc.DocId]]; ok {
			tags = append(tags, service.tagInfos[service.tagInfosIndex[doc.DocId]])
		}
	}

	return tags
}

func (service *LookupService) SearchOption(query string) []OptionInfo {
	output := service.optionSearcher.Search(types.SearchRequest{
		Text: query,
		RankOptions: &types.RankOptions{
			ScoringCriteria: &TagScoringCriteria{},
			OutputOffset:    0,
			MaxOutputs:      1000,
		},
	})

	var options []OptionInfo
	for _, doc := range output.Docs {
		if _, ok := service.optionCounts[service.optionInfosIndex[doc.DocId]]; ok {
			options = append(options, service.optionInfos[service.optionInfosIndex[doc.DocId]])
		}
	}

	return options
}

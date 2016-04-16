package core

import (
	"bufio"
	"github.com/huichen/wukong/types"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

/*******************************************************************************
  搜索功能：从逻辑规则搜索得到满足条件的用户ID
  *******************************************************************************/

// 初始化搜索引擎
func (service *LookupService) loadUserSearcher() error {
	service.searcher.Init(service.initOptions.UserSearcherInitOptions)
	service.optionCounts = make(map[string]uint32)
	service.tagCounts = make(map[uint32]uint32)
	service.optionUserPair = make(map[uint32]int)
	t1 := time.Now()

	files, _ := filepath.Glob(service.initOptions.DataFiles)
	count := 0
	for _, f := range files {
		// 打开数据文件
		log.Printf("打开数据文件 %s", f)
		file, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// 向数据库中添加消费者信息
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()
			size := len(text)

			// 跳过日期
			index := 0
			for ; index < size && text[index] != ','; index++ {
			}
			if index == size {
				log.Fatal("错误的用户属性格式")
			}

			// 找到userId的起始位置
			index++
			oldIndex := index
			for ; index < size && text[index] != ','; index++ {
			}
			if index == size {
				log.Fatal("错误的用户属性格式")
			}

			// 解析userId
			userId, err := strconv.ParseUint(text[oldIndex:index], 10, 64)
			if err != nil {
				log.Fatal(err)
			}

			tagFlag := make(map[uint32]bool)

			// 提取用户属性
			var labels []string
			for {
				// 找到下一个逗号
				index++
				oldIndex := index
				for ; index < size && text[index] != ','; index++ {
				}

				// 找到分号
				cindex := oldIndex
				for ; cindex < index && text[cindex] != ':'; cindex++ {
				}

				if cindex == oldIndex || cindex == index {
					log.Fatal("无法解析<tag id>:<option id>")
				}

				var tagId uint64
				if tagId, err = strconv.ParseUint(text[oldIndex:cindex], 10, 32); err != nil {
					log.Fatal(err)
				}
				if _, ok := service.tagWhitelist[uint32(tagId)]; ok {
					if _, err = strconv.ParseInt(text[cindex+1:index], 10, 32); err != nil {
						log.Fatal(err)
					}

					labels = append(labels, text[oldIndex:index])
					service.optionCounts[text[oldIndex:index]] += 1

					tagFlag[uint32(tagId)] = true
					service.optionUserPair[uint32(tagId)]++

				}
				if index == size {
					break
				}
			}
			if len(labels) == 0 {
				continue
			}
			for k, _ := range tagFlag {
				service.tagCounts[k] += 1
			}

			// 索引
			service.userIds = append(service.userIds, userId)
			go service.searcher.IndexDocument(userId,
				types.DocumentIndexData{
					Labels: labels,
				})

			count++
			if count%10000 == 0 {
				log.Printf("已经索引了%d万个消费者", count/10000)
			}
		}
	}

	service.searcher.FlushIndex()

	t2 := time.Now()
	t := t2.Sub(t1).Seconds()
	log.Printf("载入用户索引耗时%f秒", float64(t))
	return nil
}

func (service *LookupService) SearchUsers(query []string) []uint64 {
	request := types.SearchRequest{
		Labels:    service.sortQuery(query),
		Orderless: true,
	}

	response := service.searcher.Search(request)
	var ids []uint64

	for _, doc := range response.Docs {
		ids = append(ids, doc.DocId)
	}

	return ids
}

func (service *LookupService) GetUserCount(query []string) int {
	request := types.SearchRequest{
		Labels:        service.sortQuery(query),
		CountDocsOnly: true,
	}

	response := service.searcher.Search(request)
	return response.NumDocs
}

// 下面这些结构体和函数为了方便排序
type Query struct {
	key   string
	count uint32
}

type Queries []Query

func (queries Queries) Len() int {
	return len(queries)
}
func (queries Queries) Swap(i, j int) {
	queries[i], queries[j] = queries[j], queries[i]
}
func (queries Queries) Less(i, j int) bool {
	return queries[i].count < queries[j].count
}

func (service *LookupService) sortQuery(query []string) (output []string) {
	var qs Queries
	for _, q := range query {
		if v, ok := service.optionCounts[q]; ok {
			qs = append(qs, Query{key: q, count: v})
		} else {
			qs = append(qs, Query{key: q, count: 0})
		}
	}
	sort.Sort(qs)
	for _, q := range qs {
		output = append(output, q.key)
	}
	return
}

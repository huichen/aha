package core

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
)

var (
	tagWhitelist = flag.String("tag_whitelist", "../data/whitelist.csv", "tag白名单")
)

func (service *LookupService) loadTagWhitelist() error {
	service.tagWhitelist = make(map[uint32]bool)

	file, err := os.Open(*tagWhitelist)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		tagId, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		service.tagWhitelist[uint32(tagId)] = true
	}
	return nil
}

package util

import (
	"github.com/importcjj/sensitive"
	"log"
)

var Filter *sensitive.Filter

const path = "./config/sensitiveDict.txt"

func InitWordsFilter() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict(path)
	if err != nil {
		log.Println("filter init failed:", err.Error())
	}
	log.Println("filter init successfully!")
}

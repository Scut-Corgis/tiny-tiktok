package util

import (
	"github.com/importcjj/sensitive"
	"log"
)

var Filter *sensitive.Filter

var path = "./"

func InitwordsFilter() {
	filter := sensitive.New()
	err := filter.LoadWordDict(path)
	if err != nil {
		log.Println("InitFilter failed:", err.Error())
	}
}

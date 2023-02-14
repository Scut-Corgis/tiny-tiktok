package util

import (
	"github.com/linvon/cuckoo-filter"
	"log"
)

var CuckooFilterUserName *cuckoo.Filter // 新建一个过滤器过滤用户名

func InitCuckooFilter() {
	CuckooFilterUserName = cuckoo.NewFilter(4, 8, 100000, cuckoo.TableTypePacked)
	log.Println("CuckooFilter init successfully!")
}

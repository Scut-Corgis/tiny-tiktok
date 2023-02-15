package util

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/linvon/cuckoo-filter"
	"log"
)

var CuckooFilterUserName *cuckoo.Filter // 新建一个过滤器过滤用户名

func InitCuckooFilter() {
	CuckooFilterUserName = cuckoo.NewFilter(4, 8, 100000, cuckoo.TableTypePacked)
	names := dao.QueryAllNames()
	for _, name := range names {
		CuckooFilterUserName.Add([]byte(name))
	}
	log.Println("CuckooFilter init successfully!")
}

package redis

import (
	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"github.com/linvon/cuckoo-filter"
	"log"
	"strconv"
)

var CuckooFilterUserName *cuckoo.Filter // 过滤器过滤用户名
var CuckooFilterVideoId *cuckoo.Filter  // 过滤器过滤视频id

func InitCuckooFilter() {
	CuckooFilterUserName = cuckoo.NewFilter(4, 8, 100000, cuckoo.TableTypePacked)
	names := dao.QueryAllNames()
	for _, name := range names {
		CuckooFilterUserName.Add([]byte(name))
	}

	CuckooFilterVideoId = cuckoo.NewFilter(4, 8, 100000, cuckoo.TableTypePacked)
	videos := dao.QueryAllVideoIds()
	for _, id := range videos {
		CuckooFilterVideoId.Add([]byte(strconv.FormatInt(id, 10)))
	}
	log.Println("CuckooFilter init successfully!")
}

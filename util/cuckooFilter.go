package util

import (
	"github.com/linvon/cuckoo-filter"
)

var CuckooFilter *cuckoo.Filter

func InitCuckooFilter() {
	CuckooFilter = cuckoo.NewFilter(4, 8, 100000, cuckoo.TableTypePacked)
}

//func main() {
//	cf := cuckoo.NewFilter(4, 8, 100000, cuckoo.TableTypePacked)
//	fmt.Println(cf.Info())
//	fmt.Println(cf.FalsePositiveRate())
//
//	a := []byte("A")
//	cf.Add(a)
//	fmt.Println(cf.Contain(a))
//	fmt.Println(cf.Size())
//
//	b, _ := cf.Encode()
//	ncf, _ := cuckoo.Decode(b)
//	fmt.Println(ncf.Contain(a))
//
//	cf.Delete(a)
//	fmt.Println(cf.Size())
//}

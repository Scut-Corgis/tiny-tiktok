package dao

import (
	"github.com/Scut-Corgis/tiny-tiktok/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func Init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // 彩色打印
		},
	)
	var err error
	dsn := config.MySQL_username + ":" + config.MySQL_password + "@tcp(" + config.MySQL_IP + ":" + config.MySQL_host + ")/" + config.MySQL_database + "?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true"
	// 为MySQL开启interpolateparams以减少roundtrip
	// 想要正确的处理time.Time,需要带上 parseTime 参数，
	// 要支持完整的UTF-8编码，需要将 charset=utf8 更改为 charset=utf8mb4
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true, // 禁用GORM将写操作封装在事务内运行，提升性能
		PrepareStmt:            true, // 缓存预编译语句，提高后续的调用速度
	})

	if err != nil {
		log.Panicln("err:", err.Error())
	}
	log.Println("mysql has connected!")

	// //重建数据库
	// RebuildTable()
	// log.Println("Rebuld database successfully!")

	// 生成虚拟数据
	// FakeUsers(10)
	// log.Println("fake users generate successfully!")

	// FakeFollows(10)
	// log.Println("fake follows generate successfully!")

	//生成虚拟数据
	// FakeUsers(10)
	// log.Println("fake users generate successfully!")
	// FakeFollows(10)
	// log.Println("fake follows generate successfully!")
	// FakeVideos(10)
	// log.Println("fake videos generate successfully!")
	// FakeComments(10)
	// log.Println("fake comments generate successfully!")
	// FakeLikes(10)
	// log.Println("fake likes generate successfully!")
}

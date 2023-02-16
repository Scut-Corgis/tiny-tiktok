package service

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/Scut-Corgis/tiny-tiktok/middleware/redis"
	"github.com/Scut-Corgis/tiny-tiktok/util"

	"github.com/Scut-Corgis/tiny-tiktok/dao"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	RelationServiceImpl
	LikeServiceImpl
}

// QueryUserByName 根据name获取User对象
func (UserServiceImpl) QueryUserByName(name string) dao.User {
	user, err := dao.QueryUserByName(name)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("User not found!")
		return user
	}
	log.Println("Query user successfully!")
	return user
}

// QueryUserById 根据id获取User对象
func (UserServiceImpl) QueryUserById(id int64) dao.User {
	user, err := dao.QueryUserById(id)
	if err != nil {
		log.Println("error:", err.Error())
		log.Println("User not found!")
		return user
	}
	log.Println("Query user successfully!")
	return user
}

// QueryUserRespById 根据id获取UserResp对象
func (UserServiceImpl) QueryUserRespById(id int64) (dao.UserResp, error) {
	rsi := RelationServiceImpl{}
	userInfo := dao.UserResp{}
	user, err := dao.QueryUserById(id)
	if err != nil {
		log.Println(err.Error())
		return userInfo, err
	}
	userInfo.FollowerCount = rsi.CountFollowers(id) // 统计粉丝数量
	userInfo.FollowCount = rsi.CountFollowings(id)  // 统计关注博主的数量
	userInfo.Id = user.Id
	userInfo.Name = user.Name
	return userInfo, err
}

// Register 用户注册，返回状态码和状态信息
func (UserServiceImpl) Register(username string, password string) (int64, int32, string) {
	rand.Seed(time.Now().UnixNano())
	value := strconv.Itoa(rand.Int())
	lock := redis.Lock(username, value) // 加锁
	if lock {
		log.Println("Add lock successfully!")
		// 布谷鸟过滤器过滤
		if redis.CuckooFilterUserName.Contain([]byte(username)) {
			return -1, 1, "User already exist!"
		}
		user, _ := dao.QueryUserByName(username)
		if username == user.Name {
			return -1, 1, "User already exist!"
		} else {
			encoderPassword, err := HashEncode(password)
			if err != nil {
				return -1, 1, "Incorrect password format!"
			}
			newUser := dao.User{
				Name:     username,
				Password: encoderPassword,
			}
			usr, err := dao.InsertUser(newUser)
			if err != nil {
				return -1, 1, "Insert user failed！"
			}
			unlock := redis.Unlock(username) // 解锁
			if !unlock {
				return -1, 1, "Register failed!"
			}
			log.Println("Unlock successfully!")
			// 添加布谷鸟过滤器
			redis.CuckooFilterUserName.Add([]byte(username))
			// 添加redis缓存
			UserInsertRedis(usr.Id, user.Name)
			return usr.Id, 0, "Register successfully!"
		}
	} else {
		return -1, 1, "Wait for register!"
	}
}

// Login 用户登录，返回状态码和状态信息
func (UserServiceImpl) Login(username string, password string) (int32, string) {
	// 布谷鸟过滤器过滤
	if !redis.CuckooFilterUserName.Contain([]byte(username)) {
		return 1, "User doesn't exist!"
	}
	user, _ := dao.QueryUserByName(username)
	if ComparePasswords(user.Password, password) {
		return 0, "Login success"
	} else {
		return 1, "Username or Password error"
	}
}

// HashEncode 加密密码
func HashEncode(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePasswords 验证密码，password1为加密的密码，password2为待验证的密码
func ComparePasswords(password1 string, password2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
	if err != nil {
		return false
	}
	return true
}

func UserInsertRedis(id int64, name string) bool {
	// 插入键值对 key:user_id value:username
	redisIdKey := util.User_Id_Key + strconv.FormatInt(id, 10)
	if err := redis.RedisDb.Set(redis.Ctx, redisIdKey, name, redis.RandomTime()).Err(); err != nil {
		log.Println("Insert key:comment_id value:video_id into redis failed!")
		return false
	}
	return true
}

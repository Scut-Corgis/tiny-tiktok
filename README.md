# tiny-tiktok

## version 1.0

无任何优化

# 目录：
config为配置文件目录, 表结构在里面

controller ：控制器层，只写客户端调用接口回调函数的基本逻辑，核心逻辑实现都在service实现

dao： 为数据库操作的封装，与数据库的底层操作的封装都在里面实现

service层 ：业务核心逻辑
# 数据库：
* 目前表结构只有一级索引，无字段索引
* 自增初始值都为1000，随便取的

> users表上用户名上建立了索引， 因为这个用户名查询操作很频繁，但注意不是唯一索引，所以注册用户时记得检查用户名唯一性

> 实现自己模块过程中如果发现频繁需要某字段作为索引以提升性能，可以修改表结构，并在文档中说明原因

数据库名请取为`tiktok`，正如`dao/initDb.go`规定的那样，请修改对应的用户名和密码为自己的

`dsn := "root:123456@tcp(127.0.0.1:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"`

如以上用户名为 root, 密码为123456, 数据库名为tiktok

CRUD接口说明 : https://gorm.cn/zh_CN/docs/connecting_to_the_database.html

## jwt-auth:

> 位于`middleware/jwt`路径下

token生成和确认， 目前token中只放置了username

鉴权已注射于gin路由中，会最先执行，若通过鉴权，则会将解析出的username放于 gin.context的键值对中，可通过调用`username := context.GetString("username")`提取；若没通过鉴权，则不会运行controller代码

未测试，需等待用户注册接口完成，用户注册产生token清调用`token.go - GenerateToken(name string)` 

**其他** 

* 没有采用 `jti`, 因此有重放攻击危险，不打算考虑此问题

* jwt可选字段中，只使用了过期时间为24h，其他如发行方、接收方字段均未使用 

## videoController

### pubish - 发布视频

文件服务器 ： 
1. 配置ftp服务器，用于service服务器发送视频文件
2. 安装ffmpeg命令 并通过 ssh 连接
3. Nginx对外提供获取视频和封面的服务 (均为静态资源)

远端服务器视频文件路径/home/ftpdata/video/，图片文件路径/home/ftpdata/images/


> ffmpeg命令于文件服务器中执行，因此nginx ftp ffmpeg均在一台服务器上

**ssh服务器**

搭建：
https://www.bilibili.com/video/BV1rz4y1R7DA/?spm_id_from=333.337.search-card.all.click&vd_source=1e3090bc7a88f02cda5247bc11cd548d

ffmpeg : `sudo snap install ffmpeg`

openssl采用口令方式登陆，服务器提供公钥给客户端，客户端用公钥加密自己的密码后发回，服务器用私钥解密验证

当前实现客户端方未验证服务器公钥是否正确，若实际生产环境需得到服务器公钥到配置文件中，防止中间人截获并伪装服务器


**linux搭建ftp服务器**

**linux配置Nginx - http服务**



目前各服务器地址均为127.0.0.1

用户调动`publish` -> service服务器读取data数据 -> 将视频文件发往nginx -> 通过ssh调用ffmpeg服务得到视频起始帧图片并存于nginx服务器 -> 图片存于本地


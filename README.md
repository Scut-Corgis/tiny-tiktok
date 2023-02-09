# tiny-tiktok

## version 1.0

无任何优化

### 项目架构：
config为配置文件目录, 表结构在里面

controller ：控制器层，只写客户端调用接口回调函数的基本逻辑，核心逻辑实现都在service实现

dao： 为数据库操作的封装，与数据库的底层操作的封装都在里面实现

service层 ：业务核心逻辑
### 数据库：

#### 数据库简介

* 目前表结构只有一级索引，无字段索引
* 自增初始值都为1000，随便取的

> users表上用户名上建立了索引， 因为这个用户名查询操作很频繁，但注意不是唯一索引，所以注册用户时记得检查用户名唯一性

> 实现自己模块过程中如果发现频繁需要某字段作为索引以提升性能，可以修改表结构，并在文档中说明原因

数据库名请取为`tiktok`，正如`dao/initDb.go`规定的那样，请修改对应的用户名和密码为自己的

`dsn := "root:123456@tcp(127.0.0.1:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local"`

如以上用户名为 root, 密码为123456, 数据库名为tiktok

CRUD接口说明 : https://gorm.cn/zh_CN/docs/connecting_to_the_database.html

#### 虚拟数据生成

```txt
└─ fakeDataGenerator.go
    ├─ RebuildTable  // 重建数据库
    ├─ FakeUsers  // 生成 user 数据
    ├─ FakeFollows  // 生成 follow 数据
    ├─ FakeVideos  // 生成 video 数据
    ├─ FakeComments  // 生成 comment 数据
————————————————
```

可以在 initDao 里调用 `RebuildTable` 函数重建数据库

```go
// 需要修改一下代码里的 tableStruct.sql 路径
cmd := exec.Command("sh", "绝对路径")
```

### jwt-auth:

> 位于`middleware/jwt`路径下

token生成和确认， 目前token中只放置了username

鉴权已注射于gin路由中，会最先执行，若通过鉴权，则会将解析出的username放于 gin.context的键值对中，可通过调用`username := context.GetString("username")`提取；若没通过鉴权，则不会运行controller代码

未测试，需等待用户注册接口完成，用户注册产生token清调用`token.go - GenerateToken(name string)` 

**其他** 

* 没有采用 `jti`, 因此有重放攻击危险，不打算考虑此问题

* jwt可选字段中，只使用了过期时间为24h，其他如发行方、接收方字段均未使用 

### videoController

#### pubish - 发布视频

文件服务器 ： 
1. 配置ftp服务器，用于service服务器发送视频文件
2. 安装ffmpeg命令 并通过 ssh 连接
3. Nginx对外提供获取视频和封面的服务 (均为静态资源)

**核心逻辑：**


用户调动`publish` -> service服务器读取data数据 -> 将视频文件发往nginx -> 通过ssh调用ffmpeg服务得到视频起始帧图片并存于nginx服务器 -> 图片存于本地


> ffmpeg命令于文件服务器中执行，因此nginx ftp ffmpeg均在一台服务器上

#### ssh服务器

搭建：
https://www.bilibili.com/video/BV1rz4y1R7DA/?spm_id_from=333.337.search-card.all.click&vd_source=1e3090bc7a88f02cda5247bc11cd548d

ffmpeg : `sudo snap install ffmpeg`

wpy云服务器上：`sudo apt-get install ffmpeg`

> ssh和ffmpeg 配置好之后检查是否成功
>
> 搭建环境是否成功测试 : 
>
> 1.修改config.go文件下变量Ssh_addr_port、Ssh_username、Ssh_password为自己服务器对应的参数
>
> 2.创建config.go文件下变量Ftp_video_path、Ftp_image_path对应的路径（可自己定义）
>
> 3.将data目录下的bear.mp4文件存放在config.go文件的Ftp_video_path变量记录的路径下
>
> 4.cd到ffmpeg文件夹下 `go test`
>
> 5.Ftp_image_path路径下出现bear.jpg和bear2.jpg文件 则配置成功

openssl采用口令方式登陆，服务器提供公钥给客户端，客户端用公钥加密自己的密码后发回，服务器用私钥解密验证

openssl默认60s断开连接，因此需要加入应用层客户端心跳。

当前实现客户端方未验证服务器公钥是否正确，若实际生产环境需得到服务器公钥到配置文件中，防止中间人截获并伪装服务器


#### ftp服务器

搭建：

``` sh
sudo apt install vsftpd

systemctl enable vsftpd.service 

systemctl start vsftpd.service

systemctl status vsftpd.service

vim /etc/vsftpd/vsftpd.conf
```

增加或修改以下配置项

``` sh
# Example config file /etc/vsftpd.conf

listen=YES
listen_port=21
#
listen_ipv6=NO
# ftp登陆目录  -    改成自己的，取config.Ftp_video_path中去掉/video的路径
local_root=/home/hjg/ftpdata 
# 必须打开写权限
write_enable=YES
# 永不杀死空闲连接
idle_session_timeout=0
```

重启服务

`systemctl restart vsftpd.services`

用户添加[Ubuntu 16.04下vsftpd 安装配置实例（ftp服务器搭建）](https://blog.csdn.net/hanyuyang19940104/article/details/80421632?spm=1001.2101.3001.6650.1&utm_medium=distribute.pc_relevant.none-task-blog-2~default~CTRLIST~Rate-1-80421632-blog-79304076.pc_relevant_multi_platform_whitelistv3&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2~default~CTRLIST~Rate-1-80421632-blog-79304076.pc_relevant_multi_platform_whitelistv3&utm_relevant_index=2)

```sh
//添加用户
sudo useradd -d /home/admin/ftpdata -s /bin/bash wpy
//配置密码
sudo passwd wpy

//记得给ftpdata下的所有文件夹都符权限 不然会出现553 Could not create file.错误

```
ftp使用外网：[阿里云服务器使用FTP传输文件](https://blog.csdn.net/qq_38113006/article/details/105520125?spm=1001.2101.3001.6650.2&utm_medium=distribute.pc_relevant.none-task-blog-2~default~CTRLIST~Rate-2-105520125-blog-105364079.pc_relevant_3mothn_strategy_and_data_recovery&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2~default~CTRLIST~Rate-2-105520125-blog-105364079.pc_relevant_3mothn_strategy_and_data_recovery&utm_relevant_index=3)
> 搭建环境是否成功测试 : cd到ftp文件夹下 `go test`

#### Nginx - http服务

Nginx安装 : ` https://b23.tv/vbrNbjn`

```sh
//云服务器需要打开80端口
sudo ufw allow 80
```

修改配置

``` sh
sudo vim /usr/local/nginx/conf/nginx.conf

// 增加如下字段后保存
        // 表示访问 localhost:80/videos/ 则从 /home/hjg/ftpdata/videos/ 下面找
        location ^~ /videos/ {
            root /home/hjg/ftpdata;
        }

        location ^~ /images/ {
            root /home/hjg/ftpdata;
        }

// 切换到nginx可执行文件目录
/usr/local/nginx/sbin

// 重新加载配置
./nginx -s reload

```
权限问题：需要给根目录读写和可执行文件的命令
解法一：在nginx的nginx.conf 文件的顶部加上user root;指定操作的用户是root。
解法二： chmod -R 777 html/test给路径设置权限

./config 出错的时候：[Ubuntu16下PCRE库的安装与验证](https://blog.csdn.net/qq_40965507/article/details/117620466)

命令行安装查找nginx可执行文件看这个：[ubuntu安装nginx](https://blog.csdn.net/qq_41985134/article/details/117991218)

### APP连接云服务器

1. APP端口
2. 打开8080端口 sudo ufw allow 8080


---
## Version 2.0

项目结构和软件性能优化：
1. ...
2. redis缓存
3. ...

### 项目架构

...

### Redis
#### redis本地安装配置
1. 安装：\
https://blog.csdn.net/X_lsod/article/details/123263429
按照上述教程安装，并**启动**。
2. 配置：\
在redis.conf文件内:
搜索"bind",找到redis的ip地址，端口号默认为6379。将config.go内的"Redis_addr_port"替换。
搜索"requirepass"(未配置密码可忽略),找到redis设置的密码。将config.go内的"Redis_password"替换。
3. 测试：\
命令行可查看，也可使用redis可视化工具。\
redis可视化工具下载:https://github.com/lework/RedisDesktopManager-Windows/releases

4. 说明: \
目前只使用了一个redis库，其余设置均为默认设置。


#### redis开发规范
1. util/redisConstants.go \
该文件是对redis填入的key和TTL进行书写格式的规范\
以key为例:采用类似`const 模块_表名_key = "模块:表名:"`的分级结构，字符间的":"即为分级目录。\
例如：`const Relation_Follow_Key = "relation:follow:"`，在该文件内进行统一设置，使用时再加上对应的id或唯一性标识的字段。\
TTL类似。


    
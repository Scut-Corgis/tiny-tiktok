package config

// ssh
var Ssh_addr_port = "47.108.112.214:22"
var Ssh_username = "root"
var Ssh_password = "Wpy01104517"
var Ssh_max_taskCnt = 100

// ftp
var Ftp_addr_port = "172.19.32.144:21"
var Ftp_video_path = "/home/admin/ftpdata/videos/" //ftp服务器上的视频文件路径
var Ftp_image_path = "/home/admin/ftpdata/images/" //图片路径
var Ftp_username = "wpy"
var Ftp_password = "123456"

const Ftp_max_concurrent_cnt = 20 //Ftp并发处理的文件上限

// Url
var Url_addr_port = "127.0.0.1:80"
var Url_Play_prefix = "/videos/"
var Url_Image_prefix = "/images/"

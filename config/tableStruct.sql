SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`       bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '用户id，自增主键',
    `name`     varchar(255) NOT NULL COMMENT '用户名',
    `password` varchar(255) NOT NULL COMMENT '用户密码',
    PRIMARY KEY (`id`),
    UNIQUE INDEX (`name`, `password`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1000
  DEFAULT CHARSET = utf8 COMMENT ='用户表';

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`
(
    `id`           bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '评论id，自增主键',
    `user_id`      bigint(20)   NOT NULL COMMENT '评论发布用户id',
    `video_id`     bigint(20)   NOT NULL COMMENT '评论视频id',
    `comment_text` varchar(255) NOT NULL COMMENT '评论内容',
    `create_date`  datetime     NOT NULL COMMENT '评论发布时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1000
  DEFAULT CHARSET = utf8 COMMENT ='评论表';

-- ----------------------------
-- Table structure for follows
-- ----------------------------
DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`     bigint(20) NOT NULL COMMENT '用户id',
    `follower_id` bigint(20) NOT NULL COMMENT '关注的用户',
    PRIMARY KEY (`id`),
    INDEX (`user_id`),
    INDEX (`follower_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1000
  DEFAULT CHARSET = utf8 COMMENT ='关注表';

-- ----------------------------
-- Table structure for likes
-- ----------------------------
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes`
(
    `id`       bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `user_id`  bigint(20) NOT NULL COMMENT '点赞用户id',
    `video_id` bigint(20) NOT NULL COMMENT '被点赞的视频id',
    PRIMARY KEY (`id`),
    UNIQUE KEY `userIdtoVideoIdIdx` (`user_id`,`video_id`) USING BTREE,
    KEY `userIdIdx` (`user_id`) USING BTREE,
    KEY `videoIdx` (`video_id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 1000
  DEFAULT CHARSET = utf8 COMMENT ='点赞表';

-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`
(
    `id`           bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '自增主键，视频唯一id',
    `author_id`    bigint(20)   NOT NULL COMMENT '视频作者id',
    `play_url`     varchar(255) NOT NULL COMMENT '播放url',
    `cover_url`    varchar(255) NOT NULL COMMENT '封面url',
    `publish_time` datetime     NOT NULL COMMENT '发布时间戳',
    `title`        varchar(255) DEFAULT NULL COMMENT '视频名称',
    PRIMARY KEY (`id`),
    INDEX (`author_id`),
    INDEX (`publish_time`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1000
  DEFAULT CHARSET = utf8 COMMENT ='\r\n视频表';

-- ----------------------------
-- Table structure for messages
-- ----------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages`
(
    `id`           bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '消息id，自增主键',
    `to_user_id`   bigint(20)   NOT NULL COMMENT '接收用户id',
    `from_user_id` bigint(20)   NOT NULL COMMENT '发送用户id',
    `content`      varchar(255) NOT NULL COMMENT '消息内容',
    `create_time`  datetime NOT NULL COMMENT '消息发送时间',
    PRIMARY KEY (`id`),
    INDEX (`to_user_id`),
    INDEX (`from_user_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1000
  DEFAULT CHARSET = utf8 COMMENT ='消息表';
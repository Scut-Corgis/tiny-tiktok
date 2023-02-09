package util

import "time"

const Day = time.Hour * 24
const Month = Day * 30

// relation 模块
const Relation_Follow_Key = "relation:follow:"
const Relation_Follow_TTL = Day

const Relation_FollowerList_Key = "relation:followerList:"
const Relation_FollowerList_TTL = Day

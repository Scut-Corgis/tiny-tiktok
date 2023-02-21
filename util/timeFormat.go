package util

import "time"

/*
时间字符串转时间戳
*/
func TimeStrToUnix(timeStr string) int64 {
	loc, _ := time.LoadLocation("Local")
	timeDate, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return timeDate.Unix()
}

/*
时间字符串转时间戳
*/
func TimeStrToTime(timeStr string) time.Time {
	loc, _ := time.LoadLocation("Local")
	timeDate, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	return timeDate
}

/*
时间转时间戳
*/
func TimeToUnix(timeTime time.Time) int64 {
	loc, _ := time.LoadLocation("Local")
	timeDate, _ := time.ParseInLocation("2006-01-02 15:04:05", TimeToTimeStr(timeTime), loc)
	return timeDate.Unix()
}

/*
时间转时间字符串
*/
func TimeToTimeStr(timeTime time.Time) string {
	return timeTime.Format("2006-01-02 15:04:05")
}

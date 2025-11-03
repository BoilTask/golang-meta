package metatime

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "time/tzdata" // 必须引入
)

var beijingTimeLocation *time.Location

func Init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic("cannot load Asia/Shanghai location: " + err.Error())
	}
	beijingTimeLocation = loc
	// 设置当前时区为上海
	time.Local = beijingTimeLocation
}

// 获取常用时间文本，时区上海
func GetTimeStringByTime(t *time.Time) string {
	localTime := t.In(time.Local)
	return localTime.Format("2006-01-02 15:04:05")
}

// GetTimeNow 获取当前时间
func GetTimeNow() time.Time {
	return time.Now()
}

func GetTimeNowBeijing() time.Time {
	return time.Now().In(beijingTimeLocation)
}

func GetTimeNowString() string {
	currentTime := GetTimeNow()
	return GetTimeStringByTime(&currentTime)
}

// GetTimeDayStart 获取偏移时间后凌晨0点的时间
func GetTimeDayStart(offset time.Duration) time.Time {
	now := time.Now().Add(offset)
	year, month, day := now.Date()
	location := now.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, location)
}

// GetTimeDayEnd 获取偏移时间后23:59点的时间
func GetTimeDayEnd(offset time.Duration) time.Time {
	now := time.Now().Add(offset)
	year, month, day := now.Date()
	location := now.Location()
	return time.Date(year, month, day, 23, 59, 59, 0, location)
}

// GetTimeTodayStart 获取今日凌晨0点的时间
func GetTimeTodayStart() time.Time {
	return GetTimeDayStart(0)
}

// GetTimeTodayEnd 获取今日凌晨23:59点的时间
func GetTimeTodayEnd() time.Time {
	return GetTimeDayEnd(0)
}

// GetWeekDay 获取当前是星期几, 1 is Monday, 7 is Sunday
func GetWeekDay() int {
	currentTime := GetTimeNow()
	weekDay := int(currentTime.Weekday())
	if weekDay == 0 {
		weekDay = 7
	}
	return weekDay
}

func IsInWeekWorkDay() bool {
	weekDay := GetWeekDay()
	return weekDay >= 1 && weekDay <= 5
}

func GetBeijingTimeZone() *time.Location {
	return beijingTimeLocation
}

func GetTimeByUtcTimeString(utcTimeString string) (*time.Time, error) {
	utcTime, err := time.Parse(time.RFC3339, utcTimeString)
	if err != nil {
		return nil, err
	}
	utcTime = utcTime.In(GetBeijingTimeZone())
	return &utcTime, nil
}

func GetTimeByTimestampMs(timestampMs int64) time.Time {
	return time.Unix(0, timestampMs*int64(time.Millisecond))
}

func GetTimeByTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func GetTimeStampMsByDateString(dateString string) (int64, error) {
	// 定义时间格式
	layout := "2006-01-02 15:04:05"
	// 解析日期字符串
	localDateTime, err := time.ParseInLocation(layout, dateString, GetBeijingTimeZone())
	if err != nil {
		return 0, err
	}
	// 返回毫秒级时间戳
	return localDateTime.UnixMilli(), nil
}

func GetTimeStampByDateString(dateString string) (int64, error) {
	ms, err := GetTimeStampMsByDateString(dateString)
	if err != nil {
		return 0, err
	}
	// 返回秒级时间戳
	return ms / 1000, nil
}

func GetTimestampBySvnTime(timeString string) (int64, error) {
	re := regexp.MustCompile(`\s*\(.*?\)\s*`)
	timeString = re.ReplaceAllString(timeString, "")
	const layout = "2006-01-02 15:04:05 -0700"
	t, err := time.ParseInLocation(layout, strings.TrimSpace(timeString), GetBeijingTimeZone())
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

func GetTimestampByGitTime(timeString string) (int64, error) {
	re := regexp.MustCompile(`\s*\(.*?\)\s*`)
	timeString = re.ReplaceAllString(timeString, "")
	const layout = "2006-01-02T15:04:05Z07:00"
	t, err := time.ParseInLocation(layout, strings.TrimSpace(timeString), GetBeijingTimeZone())
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

func GetTimeByDateString(dateString string) (*time.Time, error) {
	layout := "2006-01-02 15:04:05"
	localDateTime, err := time.ParseInLocation(layout, dateString, GetBeijingTimeZone())
	if err != nil {
		return nil, err
	}
	return &localDateTime, nil
}

func IsCountdownOver(endTime time.Time) bool {
	currentTime := GetTimeNow()
	return currentTime.Compare(endTime) >= 0
}

func GetCountdownSeconds(endTime time.Time) float64 {
	currentTime := GetTimeNow()
	duration := endTime.Sub(currentTime)
	return duration.Seconds()
}

func GetPastSeconds(startTime time.Time) float64 {
	currentTime := GetTimeNow()
	duration := currentTime.Sub(startTime)
	return duration.Seconds()
}

func GetTimeStringBySeconds(seconds int) string {
	if seconds <= 0 {
		return "0秒"
	}
	timeString := ""
	hours := seconds / 3600
	if hours > 0 {
		timeString += strconv.Itoa(hours) + "小时"
	}
	minutes := (seconds % 3600) / 60
	if minutes > 0 || hours > 0 {
		if hours > 0 && minutes < 10 {
			timeString += "0"
		}
		timeString += strconv.Itoa(minutes) + "分钟"
	}
	seconds = seconds % 60
	if seconds > 0 || minutes > 0 || hours > 0 {
		if (hours > 0 || minutes > 0) && seconds < 10 {
			timeString += "0"
		}
		timeString += strconv.Itoa(seconds) + "秒"
	}
	return timeString
}

func GetCurrentYear() int {
	currentTime := GetTimeNow()
	return currentTime.Year()
}

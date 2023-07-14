package paramx

// 时间范围
type TimeRangeList []*TimeRange
type TimeRange struct {
	StartTime string `json:"startTime" comment:"开始时间"`
	EndTime   string `json:"endTime" comment:"结束时间"`
}

// 时间戳范围
type TimestampRangeList []*TimestampRange
type TimestampRange struct {
	StartTime int64 `json:"startTime" comment:"开始时间"`
	EndTime   int64 `json:"endTime" comment:"结束时间"`
}

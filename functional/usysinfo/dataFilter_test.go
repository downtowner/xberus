package usysinfo

import (
	"testing"
	"time"
)

func TestFilterByTimeInterval(t *testing.T) {
	nowTime := time.Now()
	collectData := []OneCollectData{
		OneCollectData{
			Metrics: 1,
			Time:    nowTime.Unix(),
		},
		OneCollectData{
			Metrics: 2,
			Time:    nowTime.Add(3 * time.Second).Unix(),
		},
		OneCollectData{
			Metrics: 3,
			Time:    nowTime.Add(time.Minute).Unix(),
		},
	}
	timeInterval := time.Minute
	chartPoints := FilterByTimeInterval("test", collectData, timeInterval)
	if len(chartPoints) != 2 {

		t.Errorf("FilterByTimeInterval error len(chartPoints) = %d", len(chartPoints))
	}

	if chartPoints[1].Value != 3 {
		t.Log(chartPoints)
		t.Errorf("FilterByTimeInterval error  chartPoints[1].Value  = %d", chartPoints[1].Value)
	}
}

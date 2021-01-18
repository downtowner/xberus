package usysinfo

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	testDate := time.Unix(1600000000, 0)
	formatdate := FormatTime(testDate)
	if formatdate != "2020/09/13 20:26:40" {
		t.Error("FormatTime error")
	}

}

func TestStrToTime(t *testing.T) {
	testTime, err := StrToTime("2020/09/13 20:26:40", "2006/01/02 15:04:05")
	if err != nil {
		t.Error(err)
	}
	if testTime.Unix() != 1600000000 {
		t.Error("StrToTime error")
	}
}

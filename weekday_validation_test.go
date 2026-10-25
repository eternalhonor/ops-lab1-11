package main

import "testing"

func TestValidAllWeekdays(t *testing.T) {
	for day := 1; day <= 7; day++ {
		s := Section{
			Weekday: day,
			Start:   1,
			End:     2,
			Weeks:   AllWeeks,
		}
		if !s.Valid() {
			t.Errorf("Weekday=%d 的合法上课时段应当通过校验", day)
		}
	}
}

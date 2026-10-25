package main

import "testing"

func TestValidSunday(t *testing.T) {
	s := Section{
		Weekday: 7,
		Start:   1,
		End:     2,
		Weeks:   AllWeeks,
	}
	if !s.Valid() {
		t.Fatal("星期日（Weekday=7）的合法上课时段应当通过校验")
	}
}

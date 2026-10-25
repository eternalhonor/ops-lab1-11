package main

import "testing"

func TestEvenWeeksSameSinglePeriod(t *testing.T) {
	a := Section{Weekday: 1, Start: 3, End: 3, Weeks: EvenWeeks}
	b := Section{Weekday: 1, Start: 3, End: 3, Weeks: EvenWeeks}

	if !Conflicts(a, b) {
		t.Fatal("两个双周班在同一天同一节上课，应当冲突")
	}
}

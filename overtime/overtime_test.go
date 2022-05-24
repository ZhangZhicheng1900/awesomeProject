package main

import (
	"testing"
)

func TestParseTime(t *testing.T) {
	tt := parseTime("20220524", "08:38:40")
	if tt.String() != "2022-05-24 08:38:40 +0800 CST" {
		t.Errorf("invalid time")
	}
}

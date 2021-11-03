package main

import (
	"testing"
)

func TestToCityCode(t *testing.T) {
	res := toCityCode("26102")
	if res != "26100" {
		t.Fatal("failed")
	}
}

func TestToClass15s(t *testing.T) {
	res := parentcode("2610000")
	if res != "260011" {
		t.Fatal("failed")
	}
}

func TestToClass10s(t *testing.T) {
	res := parentcode("260011")
	if res != "260010" {
		t.Fatal("failed")
	}
}

package main

import (
  "testing"
  evdev "github.com/gvalkov/golang-evdev"
)

func TestGetDigit(t *testing.T) {

  if getDigit(evdev.KEY_1) != 1 {
    t.Fail()
  }
}

func TestisKeyDownEvent(t *testing.T) {
}

func TestisKeyEventNumeric(t *testing.T) {
}

func TestAddDigit(t *testing.T) {

  var code *uint64

  if *addDigit(code, 1) != 1 {
    t.Fail()
  }

}


package main

import (
  "testing"
  evdev "github.com/gvalkov/golang-evdev"
)

func TestGetDigit(t *testing.T) {
  if getDigit(evdev.KEY_1) != 1 { t.Fail() }
}

func TestIsKeyDownEvent(t *testing.T) {
  if !isKeyDownEvent(evdev.EV_KEY, 0) { t.Fail() }
}

func TestIsKeyEventNumeric(t *testing.T) {
  if !isKeyEventNumeric(evdev.KEY_1) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_2) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_3) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_4) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_5) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_6) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_7) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_8) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_9) { t.Fail() }
  if !isKeyEventNumeric(evdev.KEY_0) { t.Fail() }
}

func TestAddDigit(t *testing.T) {
  var code *uint64

  code = addDigit(code, 1)
  if *code != 1 { t.Fail() }
  if *addDigit(code, 1) != 11 { t.Fail() }
}

func TestIsValidationNeeded(t *testing.T) {
  var validater Validater
  v := ValidatorData{validationUrl: "", code: 12345}
  validater = v

  if validater.isValidationNeeded() { t.Fail() }
}

func TestGetValidationLink(t *testing.T) {
  var code *uint64 = func(val uint64) *uint64 { return &val }(12345)

  var validater Validater
  v := ValidatorData{validationUrl: "test", code: *code}
  validater = v

  if validater.getValidationLink("url", code) != "url/12345" { t.Fail() }
}

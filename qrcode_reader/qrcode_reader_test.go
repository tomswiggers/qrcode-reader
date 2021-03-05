package main

import (
  "testing"
  evdev "github.com/gvalkov/golang-evdev"
)

func TestGetChar(t *testing.T) {
  if getChar(evdev.KEY_1) != "1" { t.Fail() }
  if getChar(evdev.KEY_A) != "a" { t.Fail() }
}

func TestAddChar(t *testing.T) {
  var code *string

  code = addChar(code, "a", false)
  if *code != "a" { t.Fail() }

  code = addChar(code, "a", true)
  if *code != "aA" { t.Fail() }

  code = addChar(code, "1", true)
  if *code != "aA1" { t.Fail() }

  code = addChar(code, "=", true)
  if *code != "aA1=" { t.Fail() }

  code = addChar(code, "/", true)
  if *code != "aA1=/" { t.Fail() }
}

func TestIsKeyUpperCase(t *testing.T) {
  if !isKeyUpperCase(evdev.KEY_LEFTSHIFT) { t.Fail() }
  if !isKeyUpperCase(evdev.KEY_RIGHTSHIFT) { t.Fail() }
}

func TestIsKeyDownEvent(t *testing.T) {
  if !isKeyDownEvent(evdev.EV_KEY, 0) { t.Fail() }
}

func TestIsTerminationKey(t *testing.T) {
  if !isTerminationKey(evdev.KEY_ENTER) { t.Fail() }
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

func TestIsValidationNeeded(t *testing.T) {
  var validater Validater
  v := ValidatorData{validationUrl: "", code: "12345"}
  validater = v

  if validater.isValidationNeeded() { t.Fail() }
}

func TestGetValidationLink(t *testing.T) {
  var code string
  code = "12345"

  var validater Validater
  v := ValidatorData{validationUrl: "url", code: code}
  validater = v

  if validater.getValidationLink() != "url/12345" { t.Fail() }
}

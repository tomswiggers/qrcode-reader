package main

import (
	"fmt"
  "flag"
  "net/http"
  "log"
	"os"
	evdev "github.com/gvalkov/golang-evdev"
)

type QrCodeReader interface {
  isReadingPossible() bool
  isKeyDownEvent(eventType uint16, eventValue int32) bool
  isKeyEventNumeric(code uint16) bool
  getDigit(code uint16) uint64
  addDigit(code *uint64, digit uint64) *uint64
}

type EventDevice struct {
  path string
}


func isKeyDownEvent(eventType uint16, eventValue int32) bool {
  if eventType == evdev.EV_KEY && eventValue == 0 {
    return true
  }

  return false
}

func isKeyEventNumeric(code uint16) bool {
  if code >= evdev.KEY_1 && code <= evdev.KEY_0 {
    return true
  }

  return false
}

func getDigit(code uint16) uint64 {
  var digits = map[uint16]uint64 {
    evdev.KEY_1 : 1,
    evdev.KEY_2 : 2,
    evdev.KEY_3 : 3,
    evdev.KEY_4 : 4,
    evdev.KEY_5 : 5,
    evdev.KEY_6 : 6,
    evdev.KEY_7 : 7,
    evdev.KEY_8 : 8,
    evdev.KEY_9 : 9,
    evdev.KEY_0 : 0,
  }

  return digits[code]
}

func getChar(code uint16) string {
  var chars = map[uint16]string {
    evdev.KEY_1 : "1",
    evdev.KEY_2 : "2",
    evdev.KEY_3 : "3",
    evdev.KEY_4 : "4",
    evdev.KEY_5 : "5",
    evdev.KEY_6 : "6",
    evdev.KEY_7 : "7",
    evdev.KEY_8 : "8",
    evdev.KEY_9 : "9",
    evdev.KEY_0 : "0",
  }

  return chars[code]
}

func addDigit(code *uint64, digit uint64) *uint64 {

  if code == nil {
    code = &digit
  } else {
    *code = *code * 10
    *code = *code + digit
  }

  return code
}

func addChar(code *string, digit string) *string {

  if code == nil {
    code = &digit
  } else {
    *code += digit
  }

  return code
}


type Validater interface {
  isValidationNeeded() bool
  getValidationLink() string
  validateQrCode() bool
}

type ValidatorData struct {
  validationUrl string
  code string
}

func (d ValidatorData) isValidationNeeded() bool {
  if d.validationUrl == "" {
    return false
  }

  return true
}

func (d ValidatorData) getValidationLink() string {
  return fmt.Sprintf("%s/%s", d.validationUrl, d.code)
}

func (d ValidatorData) validateQrCode() bool {
  fmt.Printf("validate qr code\n")
  validationUrl := d.getValidationLink()
  fmt.Printf("link: %s\n", validationUrl)
  resp, err := http.Get(validationUrl)

  if err != nil {
    log.Fatal(err)
  }

  defer resp.Body.Close()


  fmt.Print(resp.Status)

  return true
}

func main() {
  path := flag.String("inputDevice", "/dev/input/event14", "The input device")
  validationUrl := flag.String("validationUrl", "", "Validation url when QR code is scanned")
  flag.Parse()

  var code *string
  var key string
  var validator Validater

	if !evdev.IsInputDevice(*path) {
		os.Exit(1)
	}

	device, err := evdev.Open(*path)
	if err != nil {
		fmt.Printf("Unable to open input device: %s\nError: %v\n", *path, err)
		os.Exit(1)
	}

	fmt.Println(device)
  device.Grab()

	for {
		events, err := device.Read()
		if err != nil {
			fmt.Printf("device.Read() Error: %v\n", err)
			os.Exit(1)
		}
    for _, event := range events {
      if isKeyDownEvent(event.Type, event.Value) {

        if isKeyEventNumeric(event.Code) {
          key = getChar(event.Code)
          code = addChar(code, key)
        }

        if event.Code == evdev.KEY_ENTER {
          fmt.Printf("QR code is complete: %s\n", *code)
          v := ValidatorData{validationUrl: *validationUrl, code: *code}
          validator = v

          if validator.isValidationNeeded() {
            validator.validateQrCode()
          }
          code = nil
        }
      }
    }
	}
}

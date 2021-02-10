package main

import (
	"fmt"
  "flag"
  "net/http"
  "log"
	"os"
	evdev "github.com/gvalkov/golang-evdev"
)


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

func addDigit(code *uint64, digit uint64) *uint64 {

  if code == nil {
    code = &digit
  } else {
    *code = *code * 10
    *code = *code + digit
  }

  return code
}

type Validator interface {
  isValidationNeeded(validationUrl string) bool
  getValidationLink(validationUrl string, code uint64) string
  validateQrCode(validationUrl string, code uint64) bool
}

type ValidatorData struct {
  validationUrl string
  code uint64
}

func (d ValidatorData) isValidationNeeded() bool {
  if d.validationUrl == "" {
    return false
  }

  return true
}

func (d ValidatorData) getValidationLink(validationUrl string, code *uint64) string {
  return fmt.Sprintf("%s/%d", validationUrl, *code)
}

func (d ValidatorData) validateQrCode(validationUrl string, code *uint64) bool {
  if !d.isValidationNeeded() {
    fmt.Println("No validation needed")
    return false
  }

  validationUrl = d.getValidationLink(validationUrl, code)
  resp, err := http.Get(validationUrl)

  if err != nil {
    log.Fatal(err)
  }

  defer resp.Body.Close()

  return true
}

func main() {
  path := flag.String("inputDevice", "/dev/input/event14", "The input device") 
  validationUrl := flag.String("validationUrl", "", "Validation url when QR code is scanned") 
  flag.Parse()

  var code *uint64
  var key uint64

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
          key = getDigit(event.Code)
          code = addDigit(code, key)
        }

        if event.Code == evdev.KEY_ENTER {
          fmt.Printf("QR code is complete: %d\n", *code)
          validator := ValidatorData{validationUrl: *validationUrl, code: *code}
          validator.validateQrCode(*validationUrl, code)
          code = nil
        }
      }
    }
	}
}

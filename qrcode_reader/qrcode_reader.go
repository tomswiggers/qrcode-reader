package main

import (
	"fmt"
  "flag"
  "net/http"
  "log"
	"os"
  "strings"
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

func isKeyUpEvent(eventType uint16, eventValue int32) bool {
  if eventType == evdev.EV_KEY && eventValue == 1 {
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

func isKeyUpperCase(code uint16) bool {
  if code == evdev.KEY_LEFTSHIFT || code == evdev.KEY_RIGHTSHIFT {
    return true
  }

  return false
}

func isTerminationKey(code uint16) bool {
  if code == evdev.KEY_ENTER {
    return true
  }

  return false
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
    evdev.KEY_A : "a",
    evdev.KEY_B : "b",
    evdev.KEY_C : "c",
    evdev.KEY_D : "d",
    evdev.KEY_E : "e",
    evdev.KEY_F : "f",
    evdev.KEY_G : "g",
    evdev.KEY_H : "h",
    evdev.KEY_I : "i",
    evdev.KEY_J : "j",
    evdev.KEY_K : "k",
    evdev.KEY_L : "l",
    evdev.KEY_N : "n",
    evdev.KEY_M : "m",
    evdev.KEY_O : "o",
    evdev.KEY_P : "p",
    evdev.KEY_Q : "q",
    evdev.KEY_R : "r",
    evdev.KEY_S : "s",
    evdev.KEY_T : "t",
    evdev.KEY_U : "u",
    evdev.KEY_V : "v",
    evdev.KEY_W : "w",
    evdev.KEY_X : "x",
    evdev.KEY_Y : "y",
    evdev.KEY_Z : "z",
    evdev.KEY_SLASH : "/",
    evdev.KEY_EQUAL : "=",
    evdev.KEY_KPPLUS : "+",
  }

  return chars[code]
}

func addChar(code *string, char string, upper bool) *string {

  if upper {
    char = strings.ToUpper(char)
  }

  if code == nil {
    code = &char
  } else {
    *code += char
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
  var upper bool

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
      if isKeyUpEvent(event.Type, event.Value) {

        fmt.Printf("%d\n", event.Code)

        key = getChar(event.Code)
        code = addChar(code, key, upper)
        upper = isKeyUpperCase(event.Code)

        if isTerminationKey(event.Code) {
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

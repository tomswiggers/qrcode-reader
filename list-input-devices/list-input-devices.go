package main

import (
	"fmt"
	evdev "github.com/gvalkov/golang-evdev"
)

func main() {
  devices, _ := evdev.ListInputDevices()

  for _, dev := range devices {
    fmt.Printf(
      "%s %s %s %v %v\n",
      dev.Fn,
      dev.Name,
      dev.Phys,
      dev.Vendor,
      dev.Product)
  }
}

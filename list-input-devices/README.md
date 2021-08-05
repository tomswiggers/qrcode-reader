# Create module

```
$ go mod init inuits.eu/list-input-devices
go: creating new go.mod: module inuits.eu/list-input-devices
go: to add module requirements and sums:
	go mod tidy
```
  
```  
$ go mod tidy
go: finding module for package github.com/gvalkov/golang-evdev
go: found github.com/gvalkov/golang-evdev in github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
```

# Build

```
go build list-input-devices.go
```

package main

import (
	"fmt"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

type (
	HANDLE uintptr
	HWND   HANDLE
)

const (
	HWND_BROADCAST  = HWND(0xFFFF)
	WM_SYSCOMMAND   = 274
	SC_MONITORPOWER = 0xF170
	// DISPLAY_ON      = ^uint(0)
	// DISPLAY_ON      = -1 どっちもだめ。やり方が分からず。
	DISPLAY_OFF = 2
)

var (
	user32          = windows.NewLazySystemDLL("user32.dll")
	postMessageProc *windows.LazyProc
)

func init() {
	postMessageProc = user32.NewProc("PostMessageA")
}

func monitorSwitch(ms uint) {
	syscall.Syscall6(postMessageProc.Addr(), 4,
		uintptr(HWND_BROADCAST),
		uintptr(WM_SYSCOMMAND),
		uintptr(SC_MONITORPOWER),
		uintptr(ms),
		0, 0)
}

func main() {
	const start = 5
	const cycle = 60 * time.Second
	for i := 0; i < start; i++ {
		fmt.Printf("\r消灯まで%d秒前", start-i)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\r\n")
	var duration time.Duration
	for {
		monitorSwitch(DISPLAY_OFF)
		time.Sleep(cycle)
		duration += cycle
		fmt.Printf("\r%v経過", duration)
	}

	// これもだめ。
	// x := -1
	// fmt.Println("piyo")
	// monitorSwitch(uint(x))
}

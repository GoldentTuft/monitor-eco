package main

import (
	"fmt"
	"math"
	"time"
	"unsafe"

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
	user32           = windows.NewLazySystemDLL("user32.dll")
	postMessageProc  = user32.NewProc("PostMessageA")
	procGetCursorPos = user32.NewProc("GetCursorPos")
)

func monitorSwitch(ms uint) {
	postMessageProc.Call(
		uintptr(HWND_BROADCAST),
		uintptr(WM_SYSCOMMAND),
		uintptr(SC_MONITORPOWER),
		uintptr(ms),
		0, 0)
}

type Pos struct {
	X, Y int32
}

func getCursorPos() (Pos, error) {
	var pos Pos
	ret, _, err := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pos)))
	if ret == 0 {
		return pos, err
	}
	return pos, nil
}

func clearLine() {
	fmt.Print("\r                    \r")
}

func main() {
	// カウントダウン
	const countdownStart = 5
	for i := 0; i < countdownStart; i++ {
		clearLine()
		fmt.Printf("消灯まで%d秒前", countdownStart-i)
		time.Sleep(1 * time.Second)
	}

	// 指定間隔でモニターを消す
	go func() {
		const cycle = 60 * time.Second
		var duration time.Duration
		for {
			monitorSwitch(DISPLAY_OFF)
			var cycleCount time.Duration
			for cycleCount = 0; cycleCount < cycle; {
				clearLine()
				fmt.Printf("\r%v経過", duration)
				time.Sleep(1 * time.Second)
				cycleCount += 1 * time.Second
				duration += 1 * time.Second
			}
		}
	}()

	// カーソル指定距離以上移動したら終了
	exitDistance := 400.0
	startPos, err := getCursorPos()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for {
		pos, err := getCursorPos()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		xd := pos.X - startPos.X
		yd := pos.Y - startPos.Y
		d := math.Sqrt(math.Pow(float64(xd), 2) + math.Pow(float64(yd), 2))
		if d >= exitDistance {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

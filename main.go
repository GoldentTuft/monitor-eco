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

func (p Pos) distance(pos Pos) float64 {
	xd := p.X - pos.X
	yd := p.Y - pos.Y
	return math.Sqrt(math.Pow(float64(xd), 2) + math.Pow(float64(yd), 2))
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
	for i := range countdownStart {
		clearLine()
		fmt.Printf("消灯まで%d秒前", countdownStart-i)
		time.Sleep(1 * time.Second)
	}

	// 経過時間を表示
	go func() {
		const loopCycle = 1 * time.Second
		var elapsedTime time.Duration
		for {
			for {
				clearLine()
				fmt.Printf("\r%v経過", elapsedTime)
				time.Sleep(1 * time.Second)
				elapsedTime += 1 * time.Second
			}
		}
	}()

	const exitDistance = 400.0 // カーソルが移動したら終了する距離
	minDistance := 5.0         // ちょっとした移動は無視
	totalDistance := 0.0
	const loopCycle = 1 * time.Second
	var monitorOffCount time.Duration
	const monitorOffInterval = 60 * time.Second // n秒間カーソル移動がなければモニターを再度消す
	prevPos, err := getCursorPos()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	monitorSwitch(DISPLAY_OFF)
	for {
		time.Sleep(loopCycle)
		monitorOffCount += loopCycle
		if monitorOffCount >= monitorOffInterval {
			monitorSwitch(DISPLAY_OFF)
			monitorOffCount = 0
		}

		currentPos, err := getCursorPos()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		d := currentPos.distance(prevPos)
		// モニターをoffにするタイマーをリセットする
		if d >= minDistance {
			totalDistance += d
			monitorOffCount = 0
		}
		// 指定距離で終了
		if totalDistance >= exitDistance {
			return
		}
		prevPos = currentPos
	}
}

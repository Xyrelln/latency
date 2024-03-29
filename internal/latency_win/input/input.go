package input

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

// MouseClick ...
func MouseClick(x, y int) {
	// scale := robotgo.ScaleF()
	// fmt.Printf("scale: %f\n", scale)
	scale := 1.0

	robotgo.MouseSleep = 100
	robotgo.Move(x, y)

	robotgo.Toggle("left")
	defer robotgo.Toggle("left", "up")
	x1, y1 := robotgo.GetMousePos()
	fmt.Println("mouse pos: ", float64(x1)/scale, float64(y1)/scale)
	time.Sleep(time.Second)
}

// KeyboardPress ...
func KeyboardPress(key string) {
	robotgo.KeySleep = 100 // 100 millisecond
	robotgoKey, exist := keyCodeToRobotkey[key]
	if exist {
		log.Infof("KeyPress: %s", robotgoKey)
		robotgo.KeyTap(robotgoKey)
	}
}

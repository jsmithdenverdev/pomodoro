package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {

	b, err := getIcon("assets/tomato.ico")

	if err != nil {
		log.Fatal(err)
	}

	// set icon
	systray.SetIcon(b)

	// options
	pomo := systray.AddMenuItem("Pomodoro", "")
	short := systray.AddMenuItem("Short Break", "")
	long := systray.AddMenuItem("Long Break", "")
	reset := systray.AddMenuItem("Reset", "")

	systray.AddSeparator()

	quit := systray.AddMenuItem("Quit", "")

	resetChan := make(chan bool)
	go func() {
		for {
			select {
			case <-pomo.ClickedCh:
				go startCountdown(1500, resetChan)
			case <-short.ClickedCh:
				go startCountdown(300, resetChan)
			case <-long.ClickedCh:
				go startCountdown(600, resetChan)
			case <-reset.ClickedCh:
				// send a value to the reset channel
				resetChan <- true
			case <-quit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()

	systray.SetTooltip("Tomato Pomodoro Timer")
}

func onExit() {}

func getIcon(s string) ([]byte, error) {
	b, err := ioutil.ReadFile(s)

	if err != nil {
		return nil, fmt.Errorf("could not load icon: %w", err)
	}

	return b, nil
}

func startCountdown(seconds int, reset chan bool) {
	// stop the timer if we get a message in quit
	for {
		select {
		case <-reset:
			systray.SetTitle("")
			return
		default:
			systray.SetTitle(secondsToMinutes(seconds))
			seconds--
			if seconds == 0 {
				systray.SetTitle("")
				err := beeep.Notify("Buzzzzzz!", "Times up!", "")

				// log.Fatal?!
				if err != nil {
					log.Fatal(err)
				}

				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func secondsToMinutes(inSeconds int) string {
	var str string
	minutes := inSeconds / 60
	seconds := inSeconds % 60

	if seconds == 0 {
		str = fmt.Sprintf("%d:00", minutes)
	} else {
		str = fmt.Sprintf("%d:%d", minutes, seconds)
	}

	return str
}

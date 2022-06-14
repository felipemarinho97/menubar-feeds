package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/pkg/browser"
)

const shellToUse = "bash"

var displayI = 0
var displayCalls = 0
var currFeed = ""
var oldFeed = ""
var fr = NewReader()

func shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(shellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func startDbus() error {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		return err
	}
	defer conn.Close()

	if err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath("/org/darklyn/feeds"),
	); err != nil {
		return err
	}
	c := make(chan *dbus.Signal, 1)
	conn.Signal(c)

	for v := range c {
		switch v.Body[0] {
		case "OpenNews":
			browser.OpenURL(fr.GetURL())
			break
		case "NextItem":
			displayI = len([]rune(currFeed)) + 1
		case "PrevItem":
			fr.PrevItem()
			displayI = len([]rune(currFeed)) + 1
		}
	}

	return nil
}

func main() {
	args := os.Args[1:]
	go startDbus()

	panelWidth, _ := strconv.ParseInt(args[0], 10, 0)
	displayForever(int(panelWidth))
}

func getDisplayString() string {
	if displayI >= len([]rune(currFeed)) {
		displayI = 0
		currFeed = fr.GetFeed()
	}

	displayI++

	return currFeed
}

func displayForever(panelWidth int) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\033[?25h") // restore the cursor
		os.Exit(0)
	}()
	i := 0

	text := getDisplayString()
	runeText := []rune(text)

	go func() {
		for {
			newText := " " + getDisplayString() + " "
			if newText != text {
				text = newText
				i = 0
			} else {
				// text = newText
			}
			runeText = []rune(text)
			time.Sleep(450 * time.Millisecond)
		}

	}()

	fmt.Printf("\033[?25l") // hide the cursor
	for {
		if len(runeText) > panelWidth {
			for i = 0; i < len(runeText); i++ {
				// fmt.Print("\033[2J")              // clear terminal
				// fmt.Printf("\033[0;0H")           // place cursor at top left corner
				var toDisplay bytes.Buffer
				for j := 0; j < panelWidth; j++ { // character terminal width
					pos := (i + j) % len(runeText)

					toDisplay.WriteRune(runeText[pos])
				}
				fmt.Printf("%s", toDisplay.String())
				fmt.Printf("\n")

				time.Sleep(150 * time.Millisecond)
			}
		} else {
			// fmt.Print("\033[2J")    // clear terminal
			// fmt.Printf("\033[0;0H") // place cursor at top left corner
			// fmt.Printf("%s%s", getSpeaker(), text)
			fmt.Printf("\n")
			time.Sleep(250 * time.Millisecond)
		}

	}
}

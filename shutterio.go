package main

import (
	"encoding/json"
	"fmt"
	"github.com/oliverlj/go-rpio-sr595"
	"github.com/stianeikeland/go-rpio"
	"os"
	"strconv"
	"time"
)

type Configuration struct {
	Shutters []Shutter `json:"shutters"`
}

type Shutter struct {
	Up    int `json:"up,omitempty"`
	Stop  int `json:"stop,omitempty"`
	Down  int `json:"down,omitempty"`
	Delay int `json:"delay,omitempty"`
}

var (
	numberRegisterPins = 8
	dataPin            = rpio.Pin(22)
	clockPin           = rpio.Pin(17)
	latchPin           = rpio.Pin(27)
	vccPin             = rpio.Pin(23)
	configuration      Configuration
)

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parse configuration file
	file, _ := os.Open("shutterio.json")
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	shutterId, err := strconv.ParseInt(os.Args[1], 10, 8)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	shutterId = shutterId - 1

	direction := os.Args[2]
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	// Unmap gpio memory when done
	defer rpio.Close()

	// Power up the shift register and remotes
	vccPin.Output()
	vccPin.High()

	sr595.Setup(clockPin, dataPin, latchPin, numberRegisterPins)
	sr595.Reset()
	pin := getPin(shutterId, direction)
	sr595.SetRegisterPin(pin, rpio.High)
	sr595.WriteRegisters()
	time.Sleep(500 * time.Millisecond)
	sr595.SetRegisterPin(pin, rpio.Low)
	sr595.WriteRegisters()

	// Power down the shift register and remotes
	vccPin.Low()
}

func getPin(shutterId int64, direction string) int {
	shutter := configuration.Shutters[shutterId]
	if direction == "up" {
		return shutter.Up
	} else if direction == "down" {
		return shutter.Down
	} else if direction == "stop" {
		return shutter.Stop
	}
	os.Exit(2)
	return -1
}

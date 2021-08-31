package main

import (
	"bytes"
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	ledger "github.com/spacemeshos/go-ledger-sdk"
)

type speculosEvent struct {
	text string
	skip bool
	action func () error
}

type eventListener interface {
	OnEvent(event map[string]interface{}) bool
}

// Speculos Speculos
type Speculos struct {
	Info ledger.HidDeviceInfo
	step int
	events []speculosEvent
	done bool
	ready *sync.Cond
}

// NewSpeculos NewSpeculos
func NewSpeculos() *Speculos {
	return &Speculos {
		step: -1,
		done: false,
		ready: sync.NewCond(&sync.Mutex{}),
	}
}

// Open Open
func (device *Speculos) Open() error {
	return nil
}

// Close Close
func (device *Speculos) Close() {
}

// GetInfo GetInfo
func (device *Speculos) GetInfo() *ledger.HidDeviceInfo {
	return &device.Info
}

// OnEvent OnEvent
func (device *Speculos) OnEvent(data map[string]interface{}) bool {
	textField, ok := data["text"]
	if !ok {
		panic("No 'text' field")
	}
	text, ok := textField.(string)
	if !ok {
		panic("'text' field in not string")
	}
	if device.step == -1 {
		for i := 0; i < len(device.events); i++ {
			if device.events[i].text == text {
				if device.events[i].skip {
					return true
				}
				device.step = i
				break
			}
		}
		if device.step == -1 {
			panic("Unexpected event " + text)
		}
	}
	event := &device.events[device.step]
	if text != event.text {
		panic("Unexpected event " + text)
	}
	if event.action != nil {
		event.action()
	}
	device.step++;
	return device.step < len(device.events)
}

func post(url string, data string) (map[string]interface{}, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(data)))

	if err != nil {
		return nil, err
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	resp.Body.Close()

	return res, nil
}

func sendApdu(apdu string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:5000/apdu", bytes.NewBuffer([]byte("{\"data\": \"" + apdu + "\"}")))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	data, ok := res["data"].(string)
	if !ok {
		return "", fmt.Errorf("Wrong response")
	}
	return data, nil
}

func pressLeft() error {
	_, err := post("http://127.0.0.1:5000/button/left", "{\"action\":\"press-and-release\"}")
	return err
}

func pressBoth() error {
	_, err := post("http://127.0.0.1:5000/button/both", "{\"action\":\"press-and-release\"}")
	return err
}

func pressRight() error {
	_, err := post("http://127.0.0.1:5000/button/right", "{\"action\":\"press-and-release\"}")
	return err
}

// Exchange Exchange
func (device *Speculos) Exchange(apdu []byte) ([]byte, error) {
	hexData, err := sendApdu(hex.EncodeToString(apdu))
	if err != nil {
		return nil, err
	}
	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (device *Speculos) setupTest(ctx context.Context, events []speculosEvent) {
	device.step = -1
	device.events = events
	device.done = false

	go func() {
		// defer fmt.Printf("Speculos events pump done!\n")
		defer device.ready.Signal()

		// fmt.Printf("Speculos events pump start!\n")
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:5000/events?stream=true", nil)
		req = req.WithContext(ctx)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			device.done = false
			device.ready.Signal()
			fmt.Printf("Speculos events pump REQUEST: %v\n", err)
			return
		}
		defer resp.Body.Close()
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				device.done = false
				device.ready.Signal()
				fmt.Printf("Speculos events pump ERROR: %v\n", err)
				return
			}
			text := string(line)
			// fmt.Printf("Line: %v\n", text)
			if strings.HasPrefix(text, "data: ") {
				end := strings.LastIndexByte(text, '}')
				if end == -1 {
					fmt.Printf("Speculos events pump PANIC!\n")
					panic("Wrong event")
				}
				var event map[string]interface{}
				json.Unmarshal(line[6:end+1], &event)
				if (!device.OnEvent(event)) {
					// fmt.Printf("Speculos events pump DONE!\n")
					device.done = true
					device.ready.Signal()
					return
				}
			}
		}
	}()
}

func (device *Speculos) waitTestDone() bool {
	device.ready.L.Lock()
	device.ready.Wait()
	device.ready.L.Unlock()
	return device.done
}

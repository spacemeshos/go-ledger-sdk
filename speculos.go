package ledger

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Speculos event description struct
type speculosEvent struct {
	text   string
	skip   bool
	action func() error
}

// Speculos object struct
type speculos struct {
	Info   HidDeviceInfo
	step   int
	events []speculosEvent
	done   bool
	ready  *sync.Cond
}

// NewSpeculos creates a new speculos object
func NewSpeculos() *speculos {
	return &speculos{
		step:  -1,
		done:  false,
		ready: sync.NewCond(&sync.Mutex{}),
	}
}

// Open dummy method for Speculos
func (device *speculos) Open() error {
	return nil
}

// Close  dummy method for Speculos
func (device *speculos) Close() {
}

// GetInfo dummy method for Speculos
func (device *speculos) GetInfo() *HidDeviceInfo {
	return &device.Info
}

// Processing Speculos events
func (device *speculos) onEvent(data map[string]interface{}) bool {
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
	device.step++
	return device.step < len(device.events)
}

// HTTP Post method implementation for sending data to Speculos
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

// Send APDU packet to Speculos
func sendApdu(apdu string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:5001/apdu", bytes.NewBuffer([]byte("{\"data\": \""+apdu+"\"}")))
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

// Emulate of press left button on Ledger
func pressLeft() error {
	_, err := post("http://127.0.0.1:5001/button/left", "{\"action\":\"press-and-release\"}")
	return err
}

// Emulate of press both buttons on Ledger
func pressBoth() error {
	_, err := post("http://127.0.0.1:5001/button/both", "{\"action\":\"press-and-release\"}")
	return err
}

// Emulate of press right button on Ledger
func pressRight() error {
	_, err := post("http://127.0.0.1:5001/button/right", "{\"action\":\"press-and-release\"}")
	return err
}

// Exchange APDU packets with Speculos
func (device *speculos) Exchange(apdu []byte) ([]byte, error) {
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

// SetupTest prepares the device for testing and starts the Speculos event pump.
func (device *speculos) SetupTest(ctx context.Context, events []speculosEvent) {
	device.step = -1
	device.events = events
	device.done = false

	go func() {
		// defer t.Logf("Speculos events pump done!\n")
		defer device.ready.Signal()

		// t.Logf("Speculos events pump start!\n")
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:5001/events?stream=true", nil)
		req = req.WithContext(ctx)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			device.done = false
			device.ready.Signal()
			return
		}
		defer resp.Body.Close()
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				device.done = false
				device.ready.Signal()
				return
			}
			text := string(line)
			// t.Logf("Line: %v\n", text)
			if strings.HasPrefix(text, "data: ") {
				end := strings.LastIndexByte(text, '}')
				if end == -1 {
					panic("Wrong event")
				}
				var event map[string]interface{}
				json.Unmarshal(line[6:end+1], &event)
				if !device.onEvent(event) {
					// t.Logf("Speculos events pump DONE!\n")
					device.done = true
					device.ready.Signal()
					return
				}
			}
		}
	}()
}

// WaitTestDone waits for a test to complete
func (device *speculos) WaitTestDone() bool {
	device.ready.L.Lock()
	device.ready.Wait()
	device.ready.L.Unlock()
	return device.done
}

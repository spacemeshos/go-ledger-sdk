// +build speculos

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
	"testing"
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

// Create new Speculos object
func newSpeculos() *speculos {
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

	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:5000/apdu", bytes.NewBuffer([]byte("{\"data\": \""+apdu+"\"}")))
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
	_, err := post("http://127.0.0.1:5000/button/left", "{\"action\":\"press-and-release\"}")
	return err
}

// Emulate of press both buttons on Ledger
func pressBoth() error {
	_, err := post("http://127.0.0.1:5000/button/both", "{\"action\":\"press-and-release\"}")
	return err
}

// Emulate of press right button on Ledger
func pressRight() error {
	_, err := post("http://127.0.0.1:5000/button/right", "{\"action\":\"press-and-release\"}")
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

// Prepare the device for testing and start the Speculoos event pump.
func (device *speculos) setupTest(ctx context.Context, events []speculosEvent) {
	device.step = -1
	device.events = events
	device.done = false

	go func() {
		// defer t.Logf("Speculos events pump done!\n")
		defer device.ready.Signal()

		// t.Logf("Speculos events pump start!\n")
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:5000/events?stream=true", nil)
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

// Wait for a test to complete
func (device *speculos) waitTestDone() bool {
	device.ready.L.Lock()
	device.ready.Wait()
	device.ready.L.Unlock()
	return device.done
}

// Run tests on Speculos emilator
func doSpeculosTests(t *testing.T) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	ok := true
	speculos := newSpeculos()
	device := NewLedger(speculos)
	
	// run GetExtendedPublicKey test
	speculos.setupTest(ctx, []speculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Export public key"},
		{text: "m/44'/540'/0'/0/0", action: pressBoth},
		{text: "Confirm export"},
		{text: "public key?", action: pressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	publicKey, err := device.GetExtendedPublicKey(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		ok = false
		t.Logf("get public key ERROR: %v\n", err)
	} else {
		key := hex.EncodeToString(publicKey.PublicKey)
		t.Logf("public key: %v\n", key)
		if key != "a47a88814cecde42f2ad0d75123cf530fbe8e5940bbc44273014714df9a33e16" {
			ok = false
			t.Logf("WRONG public key\n")
		} else {
			t.Logf("Get public key: OK\n")
		}
	}

	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// run GetAddress test
	speculos.setupTest(ctx, []speculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Export address"},
		{text: "Path: m/44'/540'/", action: pressBoth},
		{text: "Confirm"},
		{text: "export address?", action: pressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	address, err := device.GetAddress(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		ok = false
		t.Logf("get address ERROR: %v\n", err)
	} else {
		addressStr := hex.EncodeToString(address)
		t.Logf("address: %v\n", addressStr)
		if addressStr != "a47a88814cecde42f2ad0d75123cf530fbe8e594" {
			ok = false
			t.Logf("WRONG address\n")
		} else {
			t.Logf("Get address: OK\n")
		}
	}

	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// run ShowAddress test
	speculos.setupTest(ctx, []speculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Verify address"},
		{text: "Make sure it agre", action: pressBoth},
		{text: "Address path"},
		{text: "m/44'/540'/0'/0/0", action: pressBoth},
		{text: "Address"},
		{text: "a47a88814cecde42f", action: pressBoth},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	err = device.ShowAddress(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		ok = false
		t.Logf("Show address ERROR: %v\n", err)
	} else {
		t.Logf("Show address: OK\n")
	}

	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// run Sign coin transaction test
	speculos.setupTest(ctx, []speculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Tx type:"},
		{text: "COIN ED", action: pressBoth},
		{text: "Send SMH"},
		{text: "1.0", action: pressBoth},
		{text: "To address"},
		{text: "a47a88814cecde42f", action: pressBoth},
		{text: "Max Tx Fee"},
		{text: "0.001", action: pressBoth},
		{text: "Confirm"},
		{text: "transaction?", action: pressRight},
		{text: "Signer"},
		{text: "a47a88814cecde42f", action: pressBoth},
		{text: "Sign using"},
		{text: "this signer?", action: pressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	ok = testTx(t, device, "coin.tx.json", "coin", publicKey.PublicKey, nil)
	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// run Sign app transaction test
	speculos.setupTest(ctx, []speculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Tx type:"},
		{text: "EXEC APP ED", action: pressBoth},
		{text: "Send SMH"},
		{text: "1.0", action: pressBoth},
		{text: "To address"},
		{text: "a47a88814cecde42f", action: pressBoth},
		{text: "Max Tx Fee"},
		{text: "0.001", action: pressBoth},
		{text: "Confirm"},
		{text: "transaction?", action: pressRight},
		{text: "Signer"},
		{text: "a47a88814cecde42f", action: pressBoth},
		{text: "Sign using"},
		{text: "this signer?", action: pressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	ok = testTx(t, device, "app.tx.json", "app", publicKey.PublicKey, nil)
	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// run Sign spawn transaction test
	speculos.setupTest(ctx, []speculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Tx type:"},
		{text: "SPAWN APP ED", action: pressBoth},
		{text: "Send SMH"},
		{text: "1.0", action: pressBoth},
		{text: "To address"},
		{text: "a47a88814cecde42f", action: pressBoth},
		{text: "Max Tx Fee"},
		{text: "0.001", action: pressBoth},
		{text: "Confirm"},
		{text: "transaction?", action: pressRight},
		{text: "Signer"},
		{text: "a47a88814cecde42f", action: pressBoth},
		{text: "Sign using"},
		{text: "this signer?", action: pressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	ok = testTx(t, device, "spawn.tx.json", "spawn", publicKey.PublicKey, nil)
	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	return ok
}

// Main Speculos test route
func TestSpeculos(t *testing.T) {
	if !doSpeculosTests(t) {
		t.FailNow()
	}
}

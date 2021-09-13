package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spacemeshos/ed25519"
	ledger "github.com/spacemeshos/go-ledger-sdk"
	"github.com/urfave/cli"
)

type speculosEvent struct {
	text   string
	skip   bool
	action func() error
}

// Speculos Speculos
type Speculos struct {
	Info   ledger.HidDeviceInfo
	step   int
	events []speculosEvent
	done   bool
	ready  *sync.Cond
}

// TxInfo TxInfo
type TxInfo struct {
	PublicKey []byte
	NetworkID []byte
	Type      byte
	Nonce     uint64
	To        []byte
	GasLimit  uint64
	GasPrice  uint64
	Amount    uint64
	Data      []byte
}

func uint64ToBuf(value uint64) []byte {
	data := make([]byte, 8)
	data[0] = byte((value >> 56) & 0xff)
	data[1] = byte((value >> 48) & 0xff)
	data[2] = byte((value >> 40) & 0xff)
	data[3] = byte((value >> 32) & 0xff)
	data[4] = byte((value >> 24) & 0xff)
	data[5] = byte((value >> 16) & 0xff)
	data[6] = byte((value >> 8) & 0xff)
	data[7] = byte((value) & 0xff)
	return data
}

func loadTxInfo(fileName string) (*TxInfo, error) {
	txInfo := &TxInfo{}
	jsonFile, err := os.Open(dataDirStringFlag + string(os.PathSeparator) + fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	json.Unmarshal([]byte(byteValue), &data)

	if field, ok := data["networkId"]; ok {
		txInfo.NetworkID, err = hex.DecodeString(field.(string))
		if err != nil {
			return nil, err
		}
	}

	if field, ok := data["type"]; ok {
		txInfo.Type = byte(field.(float64))
	}

	if field, ok := data["nonce"]; ok {
		txInfo.Nonce = uint64(field.(float64))
	}

	if field, ok := data["to"]; ok {
		txInfo.To, err = hex.DecodeString(field.(string))
		if err != nil {
			return nil, err
		}
	}

	if field, ok := data["gasLimit"]; ok {
		txInfo.GasLimit = uint64(field.(float64))
	}

	if field, ok := data["gasPrice"]; ok {
		txInfo.GasPrice = uint64(field.(float64))
	}

	if field, ok := data["amount"]; ok {
		txInfo.Amount = uint64(field.(float64))
	}

	if array, ok := data["data"]; ok {
		txInfo.Data = make([]byte, 0)
		items := array.([]interface{})
		for _, item := range items {
			bin, err := hex.DecodeString(item.(string))
			if err != nil {
				return nil, err
			}
			txInfo.Data = append(txInfo.Data, bin...)
		}
	}

	return txInfo, nil
}

func createTx(txInfo *TxInfo) []byte {
	tx := make([]byte, 0)
	tx = append(tx, txInfo.NetworkID...)
	tx = append(tx, txInfo.Type)
	tx = append(tx, uint64ToBuf(txInfo.Nonce)...)
	tx = append(tx, txInfo.To...)
	tx = append(tx, uint64ToBuf(txInfo.GasLimit)...)
	tx = append(tx, uint64ToBuf(txInfo.GasPrice)...)
	tx = append(tx, uint64ToBuf(txInfo.Amount)...)
	if txInfo.Data != nil {
		tx = append(tx, txInfo.Data...)
	}
	tx = append(tx, txInfo.PublicKey...)
	return tx
}

func testTx(device *ledger.Ledger, txInfoFileName string, txType string, publicKey []byte) bool {
	if txInfo, err := loadTxInfo(txInfoFileName); err == nil {
		txInfo.PublicKey = publicKey
		tx := createTx(txInfo)
		response, err := device.SignTx(ledger.StringToPath("44'/540'/0'/0/0'"), tx)
		if err == nil {
			hash := sha512.Sum512(tx)
			if ed25519.Verify(publicKey, hash[:], response[1:65]) {
				fmt.Printf("Verify %s tx: OK\n", txType)
			} else {
				fmt.Printf("Verify %s tx: FAILED\n", txType)
				return false
			}
		} else {
			fmt.Printf("Verify %s tx ERROR: %v\n", txType, err)
			return false
		}
	} else {
		fmt.Printf("Load %s tx info ERROR: %v\n", txType, err)
		return false
	}
	return true
}

func doLedgerTest(device *ledger.Ledger) bool {
	if err := device.Open(); err != nil {
		fmt.Printf("Open Ledger ERROR: %v\n", err)
		return false
	}
	defer device.Close()

	version, err := device.GetVersion()
	if err != nil {
		fmt.Printf("get version ERROR: %v\n", err)
		return false
	}
	fmt.Printf("version: %+v\n", version)

	publicKey, err := device.GetExtendedPublicKey(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("get public key ERROR: %v\n", err)
		return false
	}
	fmt.Printf("public key: %x\n", publicKey)

	address, err := device.GetAddress(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("get address ERROR: %v\n", err)
		return false
	}
	fmt.Printf("address: %x\n", address)

	err = device.ShowAddress(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("show address ERROR: %v\n", err)
		return false
	}
	fmt.Printf("show address: OK\n")

	if !testTx(device, "coin.tx.json", "coin", publicKey.PublicKey) {
		return false
	}

	if !testTx(device, "app.tx.json", "app", publicKey.PublicKey) {
		return false
	}

	if !testTx(device, "spawn.tx.json", "spawn", publicKey.PublicKey) {
		return false
	}

	return true
}

// NewSpeculos NewSpeculos
func NewSpeculos() *Speculos {
	return &Speculos{
		step:  -1,
		done:  false,
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
	device.step++
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
				if !device.OnEvent(event) {
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

func doSpeculosTests() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	ok := true
	speculos := NewSpeculos()
	device := ledger.NewLedger(speculos)
	// Test getExtendedPublicKey
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

	publicKey, err := device.GetExtendedPublicKey(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		ok = false
		fmt.Printf("get public key ERROR: %v\n", err)
	} else {
		key := hex.EncodeToString(publicKey.PublicKey)
		fmt.Printf("public key: %v\n", key)
		if key != "a47a88814cecde42f2ad0d75123cf530fbe8e5940bbc44273014714df9a33e16" {
			ok = false
			fmt.Printf("WRONG public key\n")
		} else {
			fmt.Printf("Get public key: OK\n")
		}
	}

	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// Test getAddress
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

	address, err := device.GetAddress(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		ok = false
		fmt.Printf("get address ERROR: %v\n", err)
	} else {
		addressStr := hex.EncodeToString(address)
		fmt.Printf("address: %v\n", addressStr)
		if addressStr != "a47a88814cecde42f2ad0d75123cf530fbe8e594" {
			ok = false
			fmt.Printf("WRONG address\n")
		} else {
			fmt.Printf("Get address: OK\n")
		}
	}

	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// Test showAddress
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

	err = device.ShowAddress(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		ok = false
		fmt.Printf("Show address ERROR: %v\n", err)
	} else {
		fmt.Printf("Show address: OK\n")
	}

	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// Test signCoinTx
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

	ok = testTx(device, "coin.tx.json", "coin", publicKey.PublicKey)
	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// Test signAppTx
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

	ok = testTx(device, "app.tx.json", "app", publicKey.PublicKey)
	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	// Test signSpawnTx
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

	ok = testTx(device, "spawn.tx.json", "spawn", publicKey.PublicKey)
	ok = ok && speculos.waitTestDone()
	if !ok {
		return false
	}

	return ok
}

var (
	targetStringFlag string
	dataDirStringFlag string
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "target",
		Usage:       "Run test on phisical device (ledger) or emulator (speculos)",
		Required:    false,
		Destination: &targetStringFlag,
		Value:       "ledger",
	},
	cli.StringFlag{
		Name:        "data",
		Usage:       "data directory for tx data",
		Required:    false,
		Destination: &dataDirStringFlag,
		Value:       ".",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "Spacemesh Ledger GO SDK test"
	app.Version = "0.1"
	app.Flags = flags
	app.Writer = os.Stderr

	app.Action = func(ctx *cli.Context) error {
		if targetStringFlag == "ledger" {
			devices := ledger.GetDevices(0)
			if devices == nil || len(devices) == 0 {
				fmt.Printf("No Ledger Devices Found\n")
				os.Exit(1)
			}
			for _, device := range devices {
				fmt.Printf("device: %+v\n", device.GetHidInfo())
				doLedgerTest(device)
			}
		} else if targetStringFlag == "speculos" {
			if !doSpeculosTests() {
				os.Exit(1)
			}
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

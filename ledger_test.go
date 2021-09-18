// +build !speculos

package ledger

import (
	"fmt"
	"testing"
)

func getTxTypeString(txType byte) string {
	switch txType {
	case 0:
		return "COIN ED"
	case 2:
		return "EXEC APP ED"
	case 4:
		return "SPAWN APP ED"
	default:
		return "UNWNOWN"
	}
}

func printTxInfo(txInfo *txInfo) {
	fmt.Printf("Check tx params on ledger:\n")
	fmt.Printf("\tTx type: %s\n", getTxTypeString(txInfo.Type))
	fmt.Printf("\tSend SMH: %v\n", float64(txInfo.Amount)/1000000000000.0)
	fmt.Printf("\tTo address: %x\n", txInfo.To)
	fmt.Printf("\tMax Tx Fee: %v\n", float64(txInfo.GasLimit*txInfo.GasPrice)/1000000000000.0)
	fmt.Printf("\tSigner: %x\n", txInfo.PublicKey[:20])
}

func doLedgerTest(t *testing.T, device *Ledger) bool {
	if err := device.Open(); err != nil {
		fmt.Printf("Open Ledger ERROR: %v\n", err)
		return false
	}
	defer device.Close()

	fmt.Printf("GetVersion test:\n")
	version, err := device.GetVersion()
	if err != nil {
		fmt.Printf("get version ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK, Version: %+v\n", version)

	fmt.Printf("GetExtendedPublicKey test: Follow Ledger display\n")
	publicKey, err := device.GetExtendedPublicKey(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("get public key ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK, public key: %x\n", publicKey)

	fmt.Printf("GetAddress test: Follow Ledger display\n")
	address, err := device.GetAddress(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("get address ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK, address: %x\n", address)

	fmt.Printf("ShowAddress test: Follow Ledger display\n")
	fmt.Printf("Expected address %x\n", address)
	err = device.ShowAddress(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("show address ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK\n")

	fmt.Printf("Sign coin transaction test: Follow Ledger display\n")
	if !testTx(t, device, "coin.tx.json", "coin", publicKey.PublicKey, printTxInfo) {
		return false
	}
	fmt.Printf("OK\n")

	fmt.Printf("Sign app transaction test: Follow Ledger display\n")
	if !testTx(t, device, "app.tx.json", "app", publicKey.PublicKey, printTxInfo) {
		return false
	}
	fmt.Printf("OK\n")

	fmt.Printf("Sign spawn transaction test: Follow Ledger display\n")
	if !testTx(t, device, "spawn.tx.json", "spawn", publicKey.PublicKey, printTxInfo) {
		return false
	}
	fmt.Printf("OK\n")

	return true
}

func TestLedger(t *testing.T) {
	devices := GetDevices(0)
	if devices == nil || len(devices) == 0 {
		t.Fatalf("No Ledger Devices Found\n")
	}
	for _, device := range devices {
		if !doLedgerTest(t, device) {
			t.FailNow()
		}
	}
}

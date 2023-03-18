//go:build !speculos
// +build !speculos

package ledger

import (
	"fmt"
	"testing"
)

// Return string representation of transaction type
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

// Display transaction info
func printTxInfo(txInfo *txInfo) {
	fmt.Printf("Check tx params on ledger:\n")
	fmt.Printf("\tTx type: %s\n", getTxTypeString(txInfo.Type))
	fmt.Printf("\tSend SMH: %v\n", float64(txInfo.Amount)/1000000000000.0)
	fmt.Printf("\tTo address: %x\n", txInfo.To)
	fmt.Printf("\tMax Tx Fee: %v\n", float64(txInfo.GasLimit*txInfo.GasPrice)/1000000000000.0)
	fmt.Printf("\tSigner: %x\n", txInfo.PublicKey[:20])
}

// Run tests on real Ledger device
func doLedgerTest(t *testing.T, device *Ledger) bool {
	// open Ledger device
	if err := device.Open(); err != nil {
		fmt.Printf("Open Ledger ERROR: %v\n", err)
		return false
	}
	defer device.Close()

	// run GetVersion test
	fmt.Printf("GetVersion test:\n")
	version, err := device.GetVersion()
	if err != nil {
		fmt.Printf("get version ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK, Version: %+v\n", version)

	bip44path := "44'/540'/0'/0/0'"
	path := StringToPath(bip44path)

	// run GetExtendedPublicKey test
	fmt.Printf("GetExtendedPublicKey test: export public key for account \"%v\"\n", bip44path)
	fmt.Printf("Please confirm exporting the public key on your Ledger.\n")
	publicKey, err := device.GetExtendedPublicKey(path)
	if err != nil {
		fmt.Printf("get public key ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK, public key: %x\n", publicKey)

	// run GetAddress test
	fmt.Printf("GetAddress test: get address for account \"%v\"\n", bip44path)
	fmt.Printf("Please confirm exporting the address on your Ledger.\n")
	address, err := device.GetAddress(path)
	if err != nil {
		fmt.Printf("get address ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK, address: %x\n", address)

	// run ShowAddress test
	fmt.Printf("ShowAddress test: show address for account \"%v\"\n", bip44path)
	fmt.Printf("Please make sure the following address agrees with your Ledger display.\n")
	fmt.Printf("Expected address is %x\n", address)
	err = device.ShowAddress(path)
	if err != nil {
		fmt.Printf("show address ERROR: %v\n", err)
		return false
	}
	fmt.Printf("OK\n")

	// run Sign coin transaction test
	fmt.Printf("Sign coin transaction test: Follow Ledger display\n")
	if !testTx(t, device, "coin.tx.json", "coin", publicKey.PublicKey, printTxInfo) {
		return false
	}
	fmt.Printf("OK\n")

	// run Sign app transaction test
	fmt.Printf("Sign app transaction test: Follow Ledger display\n")
	if !testTx(t, device, "app.tx.json", "app", publicKey.PublicKey, printTxInfo) {
		return false
	}
	fmt.Printf("OK\n")

	// run Sign spawn transaction test
	fmt.Printf("Sign spawn transaction test: Follow Ledger display\n")
	if !testTx(t, device, "spawn.tx.json", "spawn", publicKey.PublicKey, printTxInfo) {
		return false
	}
	fmt.Printf("OK\n")

	return true
}

// Main Ledger test route
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

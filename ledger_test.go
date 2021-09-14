// +build !speculos

package ledger

import (
	"testing"
)

func doLedgerTest(t *testing.T, device *Ledger) bool {
	if err := device.Open(); err != nil {
		t.Logf("Open Ledger ERROR: %v\n", err)
		return false
	}
	defer device.Close()

	version, err := device.GetVersion()
	if err != nil {
		t.Logf("get version ERROR: %v\n", err)
		return false
	}
	t.Logf("version: %+v\n", version)

	publicKey, err := device.GetExtendedPublicKey(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		t.Logf("get public key ERROR: %v\n", err)
		return false
	}
	t.Logf("public key: %x\n", publicKey)

	address, err := device.GetAddress(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		t.Logf("get address ERROR: %v\n", err)
		return false
	}
	t.Logf("address: %x\n", address)

	err = device.ShowAddress(StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		t.Logf("show address ERROR: %v\n", err)
		return false
	}
	t.Logf("show address: OK\n")

	if !testTx(t, device, "coin.tx.json", "coin", publicKey.PublicKey) {
		return false
	}

	if !testTx(t, device, "app.tx.json", "app", publicKey.PublicKey) {
		return false
	}

	if !testTx(t, device, "spawn.tx.json", "spawn", publicKey.PublicKey) {
		return false
	}

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

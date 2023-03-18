//go:build speculos
// +build speculos

package ledger

import (
	"context"
	"encoding/hex"
	"testing"
	"time"
)

// Run tests on Speculos emulator
func doSpeculosTests(t *testing.T) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	ok := true
	speculos := NewSpeculos()
	device := NewLedger(speculos)

	path := StringToPath("44'/540'/0'/0/0'")

	// run GetExtendedPublicKey test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{Text: "Spacemesh", Skip: true},
		{Text: "is ready", Skip: true},
		{Text: "Export public key"},
		{Text: "m/44'/540'/0'/0/0", Action: PressBoth},
		{Text: "Confirm export"},
		{Text: "public key?", Action: PressRight},
		{Text: "Spacemesh"},
		{Text: "is ready"},
	})

	publicKey, err := device.GetExtendedPublicKey(path)
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

	ok = ok && speculos.WaitTestDone()
	if !ok {
		return false
	}

	// run GetAddress test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{Text: "Spacemesh", Skip: true},
		{Text: "is ready", Skip: true},
		{Text: "Export address"},
		{Text: "Path: m/44'/540'/", Action: PressBoth},
		{Text: "Confirm"},
		{Text: "export address?", Action: PressRight},
		{Text: "Spacemesh"},
		{Text: "is ready"},
	})

	address, err := device.GetAddress(path)
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

	ok = ok && speculos.WaitTestDone()
	if !ok {
		return false
	}

	// run ShowAddress test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{Text: "Spacemesh", Skip: true},
		{Text: "is ready", Skip: true},
		{Text: "Verify address"},
		{Text: "Make sure it agre", Action: PressBoth},
		{Text: "Address path"},
		{Text: "m/44'/540'/0'/0/0", Action: PressBoth},
		{Text: "Address"},
		{Text: "a47a88814cecde42f", Action: PressBoth},
		{Text: "Spacemesh"},
		{Text: "is ready"},
	})

	err = device.ShowAddress(path)
	if err != nil {
		ok = false
		t.Logf("Show address ERROR: %v\n", err)
	} else {
		t.Logf("Show address: OK\n")
	}

	ok = ok && speculos.WaitTestDone()
	if !ok {
		return false
	}

	// run Sign coin transaction test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{Text: "Spacemesh", Skip: true},
		{Text: "is ready", Skip: true},
		{Text: "Tx type:"},
		{Text: "COIN ED", Action: PressBoth},
		{Text: "Send SMH"},
		{Text: "1.0", Action: PressBoth},
		{Text: "To address"},
		{Text: "a47a88814cecde42f", Action: PressBoth},
		{Text: "Max Tx Fee"},
		{Text: "0.001", Action: PressBoth},
		{Text: "Confirm"},
		{Text: "transaction?", Action: PressRight},
		{Text: "Signer"},
		{Text: "a47a88814cecde42f", Action: PressBoth},
		{Text: "Sign using"},
		{Text: "this signer?", Action: PressRight},
		{Text: "Spacemesh"},
		{Text: "is ready"},
	})

	ok = TestTx(t, device, "coin.tx.json", "coin", publicKey.PublicKey, nil)
	ok = ok && speculos.WaitTestDone()
	if !ok {
		return false
	}

	// run Sign app transaction test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{Text: "Spacemesh", Skip: true},
		{Text: "is ready", Skip: true},
		{Text: "Tx type:"},
		{Text: "EXEC APP ED", Action: PressBoth},
		{Text: "Send SMH"},
		{Text: "1.0", Action: PressBoth},
		{Text: "To address"},
		{Text: "a47a88814cecde42f", Action: PressBoth},
		{Text: "Max Tx Fee"},
		{Text: "0.001", Action: PressBoth},
		{Text: "Confirm"},
		{Text: "transaction?", Action: PressRight},
		{Text: "Signer"},
		{Text: "a47a88814cecde42f", Action: PressBoth},
		{Text: "Sign using"},
		{Text: "this signer?", Action: PressRight},
		{Text: "Spacemesh"},
		{Text: "is ready"},
	})

	ok = TestTx(t, device, "app.tx.json", "app", publicKey.PublicKey, nil)
	ok = ok && speculos.WaitTestDone()
	if !ok {
		return false
	}

	// run Sign spawn transaction test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{Text: "Spacemesh", Skip: true},
		{Text: "is ready", Skip: true},
		{Text: "Tx type:"},
		{Text: "SPAWN APP ED", Action: PressBoth},
		{Text: "Send SMH"},
		{Text: "1.0", Action: PressBoth},
		{Text: "To address"},
		{Text: "a47a88814cecde42f", Action: PressBoth},
		{Text: "Max Tx Fee"},
		{Text: "0.001", Action: PressBoth},
		{Text: "Confirm"},
		{Text: "transaction?", Action: PressRight},
		{Text: "Signer"},
		{Text: "a47a88814cecde42f", Action: PressBoth},
		{Text: "Sign using"},
		{Text: "this signer?", Action: PressRight},
		{Text: "Spacemesh"},
		{Text: "is ready"},
	})

	ok = TestTx(t, device, "spawn.tx.json", "spawn", publicKey.PublicKey, nil)
	ok = ok && speculos.WaitTestDone()
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

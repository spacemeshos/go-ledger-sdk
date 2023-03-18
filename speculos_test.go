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
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Export public key"},
		{text: "m/44'/540'/0'/0/0", action: PressBoth},
		{text: "Confirm export"},
		{text: "public key?", action: PressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
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
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Export address"},
		{text: "Path: m/44'/540'/", action: PressBoth},
		{text: "Confirm"},
		{text: "export address?", action: PressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
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
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Verify address"},
		{text: "Make sure it agre", action: PressBoth},
		{text: "Address path"},
		{text: "m/44'/540'/0'/0/0", action: PressBoth},
		{text: "Address"},
		{text: "a47a88814cecde42f", action: PressBoth},
		{text: "Spacemesh"},
		{text: "is ready"},
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
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Tx type:"},
		{text: "COIN ED", action: PressBoth},
		{text: "Send SMH"},
		{text: "1.0", action: PressBoth},
		{text: "To address"},
		{text: "a47a88814cecde42f", action: PressBoth},
		{text: "Max Tx Fee"},
		{text: "0.001", action: PressBoth},
		{text: "Confirm"},
		{text: "transaction?", action: PressRight},
		{text: "Signer"},
		{text: "a47a88814cecde42f", action: PressBoth},
		{text: "Sign using"},
		{text: "this signer?", action: PressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	ok = testTx(t, device, "coin.tx.json", "coin", publicKey.PublicKey, nil)
	ok = ok && speculos.WaitTestDone()
	if !ok {
		return false
	}

	// run Sign app transaction test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Tx type:"},
		{text: "EXEC APP ED", action: PressBoth},
		{text: "Send SMH"},
		{text: "1.0", action: PressBoth},
		{text: "To address"},
		{text: "a47a88814cecde42f", action: PressBoth},
		{text: "Max Tx Fee"},
		{text: "0.001", action: PressBoth},
		{text: "Confirm"},
		{text: "transaction?", action: PressRight},
		{text: "Signer"},
		{text: "a47a88814cecde42f", action: PressBoth},
		{text: "Sign using"},
		{text: "this signer?", action: PressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	ok = testTx(t, device, "app.tx.json", "app", publicKey.PublicKey, nil)
	ok = ok && speculos.WaitTestDone()
	if !ok {
		return false
	}

	// run Sign spawn transaction test
	speculos.SetupTest(ctx, []SpeculosEvent{
		{text: "Spacemesh", skip: true},
		{text: "is ready", skip: true},
		{text: "Tx type:"},
		{text: "SPAWN APP ED", action: PressBoth},
		{text: "Send SMH"},
		{text: "1.0", action: PressBoth},
		{text: "To address"},
		{text: "a47a88814cecde42f", action: PressBoth},
		{text: "Max Tx Fee"},
		{text: "0.001", action: PressBoth},
		{text: "Confirm"},
		{text: "transaction?", action: PressRight},
		{text: "Signer"},
		{text: "a47a88814cecde42f", action: PressBoth},
		{text: "Sign using"},
		{text: "this signer?", action: PressRight},
		{text: "Spacemesh"},
		{text: "is ready"},
	})

	ok = testTx(t, device, "spawn.tx.json", "spawn", publicKey.PublicKey, nil)
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

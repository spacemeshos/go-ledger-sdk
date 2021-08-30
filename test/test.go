package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/spacemeshos/ed25519"
	ledger "github.com/spacemeshos/go-ledger-sdk"
)

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

func doTest(device *ledger.HidDevice) {
	if err := device.Open(); err != nil {
		return
	}
	defer device.Close()

	version, err := device.GetVersion()
	if err != nil {
		fmt.Printf("get version ERROR: %v\n", err)
	} else {
		fmt.Printf("version: %+v\n", version)
	}
	publicKey, err := device.GetExtendedPublicKey(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("get public key ERROR: %v\n", err)
	} else {
		fmt.Printf("public key: %x\n", publicKey)
	}
	address, err := device.GetAddress(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("get address ERROR: %v\n", err)
	} else {
		fmt.Printf("address: %x\n", address)
	}
	err = device.ShowAddress(ledger.StringToPath("44'/540'/0'/0/0'"))
	if err != nil {
		fmt.Printf("show address ERROR: %v\n", err)
	} else {
		fmt.Printf("show address: OK\n")
	}

	tx := make([]byte, 0)
	var bin []byte
	bin, _ = hex.DecodeString("1835df3489b3a39e0f38a77d347f8327e8937c623543b84bd8734fc237ae3f33") // network id
	tx = append(tx, bin...)
	tx = append(tx, 0)                             // coin transaction with ed
	tx = append(tx, uint64ToBuf(1)...)             // nonce
	bin, _ = hex.DecodeString("a47a88814cecde42f2ad0d75123cf530fbe8e594") // recepient
	tx = append(tx, bin...)
	tx = append(tx, uint64ToBuf(1000000)...)       // gas limit
	tx = append(tx, uint64ToBuf(1000)...)          // gas price
	tx = append(tx, uint64ToBuf(1000000000000)...) // amount
	tx = append(tx, publicKey.PublicKey...)

	response, err := device.SignTx(ledger.StringToPath("44'/540'/0'/0/0'"), tx)
	if err != nil {
		fmt.Printf("Verify coin tx ERROR: %v\n", err)
	} else {
		hash := sha512.Sum512(tx)
		fmt.Printf("Verify coin tx: %v\n", ed25519.Verify(publicKey.PublicKey, hash[:], response[1:65]))
	}

	tx = make([]byte, 0)
	bin, _ = hex.DecodeString("1835df3489b3a39e0f38a77d347f8327e8937c623543b84bd8734fc237ae3f33") // network id
	tx = append(tx, bin...)
	tx = append(tx, 2)                             // exec app transaction with ed
	tx = append(tx, uint64ToBuf(1)...)             // nonce
	bin, _ = hex.DecodeString("a47a88814cecde42f2ad0d75123cf530fbe8e594") // app address
	tx = append(tx, bin...)
	tx = append(tx, uint64ToBuf(1000000)...)       // gas limit
	tx = append(tx, uint64ToBuf(1000)...)          // gas price
	tx = append(tx, uint64ToBuf(1000000000000)...) // amount
	// call data
	bin, _ = hex.DecodeString("f27596fadee9fd74a7745ab45d978e82e275e34d06e86f127bb1f18bfb8f1889") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("94b739420fe48a6fc6aa26e73e79ae743addb37c615c85bfd3688995be7ba7b1") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("ca459caa9b121df4a11125428e186ab633483f01ee14c7f70229153e39a873c6") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("c1a74d00ac388471cd99c6691b0168ee801028b393b66dc88b304e9c179617b2") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("13d834d17081701e7e47e25b8c7dbea6378603fb8d94e0912acb8d06692ffdc7") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("84a5219d52980f459b23b89ec975bd432acd041415edddd65f19ff5d7bfbb57a") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("72cac44b71f2a8728cb3f460c012592b5a2b7bb4530392498149e9fa70ada5fc") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("1881f6fee14a24eae3caab13c0dbafb9b544bf0f9949654445fa87158ff42642") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("4c6fcfc41e314c71bf7c372eaae41079e9f4b27d6e6bb61c5e1e22a7111eb54e") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("aa40133add17d2b9aef11781f65ab8fff5b1649cfd291f828384efa1c89584bd") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("0e71041e3750cf03daa482bda9106cb6f4112b2f4034279bd6abd6e54fedbbaf") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("d101ff66ef173d83769cf581e72853dd5cdb42f296f330f66df97100804386cd") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("6bf953ce0cd6e84a388f36be9ceee72d6bcc98fbfc9cfbc5b8f93ab7be5beed2") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("2ba72e494667f5548323398156377d930b1efd27971608dd29df8d1523b9ac41") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("a4e87c684ec7ee7790ca5899745b4be54026b3be233a39a454eff328c71995d4") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("20a460255a11e4508572915a8c394c764682392a1d1ecf43951119cf41ce5829") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("00918629e3b96b34b82396058659931d1904ef8e38cc84e717057d68e37bba33") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("e1f80e36a9031e39e84b6523f0470ffa41515f8847f5a0e7cd8cc31bf287639a") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("4d0249a50d866f6400c95f7f895f16ced4148f1b4c1ba2b63a7b4423adf96085") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("0432dea6b22bab7018bc20713348dbca2722a75d08f9ab60718355ba7063d337") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("3cd547a3a019b78ae83e6cf5fc9bcc3567ccff0d298a98cf2fa66770007428b6") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("4645f110521e83257bbbb686d23a89e36b92eeee61334d2ca644489478fb48b1") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("afd08737fcaacf637b9425ee514480aafe0af536fed8361a09503ccf512456c8") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("2cdf11e92f7d2e032a3c622cd113f1ea75f5ce1a6d93b67bfdb6354cf8496158") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("efbbe5b97f074e2c062a965bd44569040da5056a02ad03eeb9bd42c8fdae963a") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("ba569c7643fe2fbe106f5e3cce03e5bd773074e24dd26298f07f8a5a1b00b31d") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("cbaf0de81f0a2d7234d1ddd234d2093d66a515eb4c0e0c5e144e0cac300a151f") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("81a2f9a4c87dbd1b8a602ae80fba55913a0830c0d7f34a5fed47a5da7c776bb1") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("59ab8c2aa492f4c06c8519717ac6abc4413efaa8ff0e569768441a3780abb21c") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("71d81d9222c6d225d87744a60433c5741dd6fbf363fab3a93cabbf8b16c6ef58") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("0fc359fc5e7853bf24a7825f8aba7a1497a2d4b892ccd2d1ffa37ffcf4b4cc87") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("2b0170bde7f2bdefc226dd50d2aaf3f5a971593a0d52ee0c5979e233ab08dbf8") // bin data
	tx = append(tx, bin...)
	tx = append(tx, publicKey.PublicKey...)

	response, err = device.SignTx(ledger.StringToPath("44'/540'/0'/0/0'"), tx)
	if err != nil {
		fmt.Printf("Verify app tx ERROR: %v\n", err)
	} else {
		hash := sha512.Sum512(tx)
		fmt.Printf("Verify app tx: %v\n", ed25519.Verify(publicKey.PublicKey, hash[:], response[1:65]))
	}

	tx = make([]byte, 0)
	bin, _ = hex.DecodeString("1835df3489b3a39e0f38a77d347f8327e8937c623543b84bd8734fc237ae3f33") // network id
	tx = append(tx, bin...)
	tx = append(tx, 4)                             // spawn app + ed
	tx = append(tx, uint64ToBuf(1)...)             // nonce
	bin, _ = hex.DecodeString("a47a88814cecde42f2ad0d75123cf530fbe8e594") // template address
	tx = append(tx, bin...)
	tx = append(tx, uint64ToBuf(1000000)...)       // gas limit
	tx = append(tx, uint64ToBuf(1000)...)          // gas price
	tx = append(tx, uint64ToBuf(1000000000000)...) // amount
	// call data
	bin, _ = hex.DecodeString("f27596fadee9fd74a7745ab45d978e82e275e34d06e86f127bb1f18bfb8f1889") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("94b739420fe48a6fc6aa26e73e79ae743addb37c615c85bfd3688995be7ba7b1") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("ca459caa9b121df4a11125428e186ab633483f01ee14c7f70229153e39a873c6") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("c1a74d00ac388471cd99c6691b0168ee801028b393b66dc88b304e9c179617b2") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("13d834d17081701e7e47e25b8c7dbea6378603fb8d94e0912acb8d06692ffdc7") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("84a5219d52980f459b23b89ec975bd432acd041415edddd65f19ff5d7bfbb57a") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("72cac44b71f2a8728cb3f460c012592b5a2b7bb4530392498149e9fa70ada5fc") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("1881f6fee14a24eae3caab13c0dbafb9b544bf0f9949654445fa87158ff42642") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("4c6fcfc41e314c71bf7c372eaae41079e9f4b27d6e6bb61c5e1e22a7111eb54e") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("aa40133add17d2b9aef11781f65ab8fff5b1649cfd291f828384efa1c89584bd") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("0e71041e3750cf03daa482bda9106cb6f4112b2f4034279bd6abd6e54fedbbaf") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("d101ff66ef173d83769cf581e72853dd5cdb42f296f330f66df97100804386cd") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("6bf953ce0cd6e84a388f36be9ceee72d6bcc98fbfc9cfbc5b8f93ab7be5beed2") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("2ba72e494667f5548323398156377d930b1efd27971608dd29df8d1523b9ac41") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("a4e87c684ec7ee7790ca5899745b4be54026b3be233a39a454eff328c71995d4") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("20a460255a11e4508572915a8c394c764682392a1d1ecf43951119cf41ce5829") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("00918629e3b96b34b82396058659931d1904ef8e38cc84e717057d68e37bba33") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("e1f80e36a9031e39e84b6523f0470ffa41515f8847f5a0e7cd8cc31bf287639a") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("4d0249a50d866f6400c95f7f895f16ced4148f1b4c1ba2b63a7b4423adf96085") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("0432dea6b22bab7018bc20713348dbca2722a75d08f9ab60718355ba7063d337") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("3cd547a3a019b78ae83e6cf5fc9bcc3567ccff0d298a98cf2fa66770007428b6") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("4645f110521e83257bbbb686d23a89e36b92eeee61334d2ca644489478fb48b1") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("afd08737fcaacf637b9425ee514480aafe0af536fed8361a09503ccf512456c8") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("2cdf11e92f7d2e032a3c622cd113f1ea75f5ce1a6d93b67bfdb6354cf8496158") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("efbbe5b97f074e2c062a965bd44569040da5056a02ad03eeb9bd42c8fdae963a") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("ba569c7643fe2fbe106f5e3cce03e5bd773074e24dd26298f07f8a5a1b00b31d") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("cbaf0de81f0a2d7234d1ddd234d2093d66a515eb4c0e0c5e144e0cac300a151f") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("81a2f9a4c87dbd1b8a602ae80fba55913a0830c0d7f34a5fed47a5da7c776bb1") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("59ab8c2aa492f4c06c8519717ac6abc4413efaa8ff0e569768441a3780abb21c") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("71d81d9222c6d225d87744a60433c5741dd6fbf363fab3a93cabbf8b16c6ef58") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("0fc359fc5e7853bf24a7825f8aba7a1497a2d4b892ccd2d1ffa37ffcf4b4cc87") // bin data
	tx = append(tx, bin...)
	bin, _ = hex.DecodeString("2b0170bde7f2bdefc226dd50d2aaf3f5a971593a0d52ee0c5979e233ab08dbf8") // bin data
	tx = append(tx, bin...)
	tx = append(tx, publicKey.PublicKey...)

	response, err = device.SignTx(ledger.StringToPath("44'/540'/0'/0/0'"), tx)
	if err != nil {
		fmt.Printf("Verify spawn tx ERROR: %v\n", err)
	} else {
		hash := sha512.Sum512(tx)
		fmt.Printf("Verify spawn tx: %v\n", ed25519.Verify(publicKey.PublicKey, hash[:], response[1:65]))
	}
}

func main() {
	devices := ledger.GetDevices(0)
	for _, device := range devices {
		fmt.Printf("device: %+v\n", device.Info)
		doTest(device)
	}
}

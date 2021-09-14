package ledger

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spacemeshos/ed25519"
)

type txInfo struct {
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

func loadTxInfo(fileName string) (*txInfo, error) {
	txInfo := &txInfo{}
	jsonFile, err := os.Open("./test" + string(os.PathSeparator) + fileName)
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

func createTx(txInfo *txInfo) []byte {
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

func testTx(t *testing.T, device *Ledger, txInfoFileName string, txType string, publicKey []byte) bool {
	if txInfo, err := loadTxInfo(txInfoFileName); err == nil {
		txInfo.PublicKey = publicKey
		tx := createTx(txInfo)
		response, err := device.SignTx(StringToPath("44'/540'/0'/0/0'"), tx)
		if err == nil {
			hash := sha512.Sum512(tx)
			if ed25519.Verify(publicKey, hash[:], response[1:65]) {
				t.Logf("Verify %s tx: OK\n", txType)
			} else {
				t.Logf("Verify %s tx: FAILED\n", txType)
				return false
			}
		} else {
			t.Logf("Verify %s tx ERROR: %v\n", txType, err)
			return false
		}
	} else {
		t.Logf("Load %s tx info ERROR: %v\n", txType, err)
		return false
	}
	return true
}

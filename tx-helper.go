package ledger

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spacemeshos/ed25519"
)

// Transaction info struct
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

// Load transactoin info from JSON file
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
	json.Unmarshal(byteValue, &data)

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

// Convert transaction info to byte array
func createTx(txInfo *txInfo) []byte {
	tx := make([]byte, 32+1+8+20+8+8+8)
	copy(tx, txInfo.NetworkID)
	tx[32] = txInfo.Type
	binary.BigEndian.PutUint64(tx[33:], txInfo.Nonce)
	copy(tx[41:], txInfo.To)
	binary.BigEndian.PutUint64(tx[61:], txInfo.GasLimit)
	binary.BigEndian.PutUint64(tx[69:], txInfo.GasPrice)
	binary.BigEndian.PutUint64(tx[77:], txInfo.Amount)
	if txInfo.Data != nil {
		tx = append(tx, txInfo.Data...)
	}
	tx = append(tx, txInfo.PublicKey...)
	return tx
}

// Do transaction test
func testTx(t *testing.T, device *Ledger, txInfoFileName string, txType string, publicKey []byte, callback func(txInfo *txInfo)) bool {
	if txInfo, err := loadTxInfo(txInfoFileName); err == nil {
		txInfo.PublicKey = publicKey
		if callback != nil {
			callback(txInfo)
		}
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

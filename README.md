# go-ledger-sdk
A golang SDK for using Ledger wallets from golang Spacemesh wallets such as smrepl

# Note: this repository is no longer maintained. Please see [spacemesh-sdk](https://github.com/spacemeshos/spacemesh-sdk) instead.

## API

Enumerate Ledger devices.
```
/**
 * @param {int} productID USB Product ID filter, 0 - all.
 * @return {[]*HidDevice} Discovered Ledger devices.
 *
 * @example
 * devices := ledger.GetDevices(0)
 * if devices != nil && len(devices) > 0 {
 * 	device := devices[0]
 * 	if err := device.Open(); err == nil {
 * 		...
 * 		device.Close()
 * 	} else {
 * 		fmt.Printf("Open device ERROR: %v\n", err)
 * 	}
 * }
 */
func GetDevices(productID int) []*HidDevice
```

Open Ledger device for communication.
```
/**
 * @return {error} Error value.
 *
 * @example
 * devices := ledger.GetDevices(0)
 * if devices != nil && len(devices) > 0 {
 * 	device := devices[0]
 * 	if err := device.Open(); err == nil {
 * 		...
 * 		device.Close()
 * 	} else {
 * 		fmt.Printf("Open device ERROR: %v\n", err)
 * 	}
 * }
 */
func (device *HidDevice) Open() error
```

Close communication with Ledger device.
```
/**
 * @example
 * devices := ledger.GetDevices(0)
 * if devices != nil && len(devices) > 0 {
 * 	device := devices[0]
 * 	if err := device.Open(); err == nil {
 * 		...
 * 		device.Close()
 * 	} else {
 * 		fmt.Printf("Open device ERROR: %v\n", err)
 * 	}
 * }
 */
func (device *HidDevice) Close()
```

Get the ledger app version.
```
/**
 * @returns {Version} Result object containing the application version number.
 * @return {error} Error value.
 *
 * @example
 * version, err := device.GetVersion()
 * if err != nil {
 * 	fmt.Printf("get version ERROR: %v\n", err)
 * } else {
 * 	fmt.Printf("version: %+v\n", version)
 * }
 */
func (device *HidDevice) GetVersion() (*Version, error)
```

Get a public key from the specified BIP 32 path.
```
/**
 * @param {BipPath} path The BIP 32 path indexes. Path must begin with `44'/540'/n'`, and shuld be 5 indexes long.
 * @return {ExtendedPublicKey} The public key with chaincode for the given path.
 * @return {error} Error value.
 *
 * @example
 * publicKey, err := device.GetExtendedPublicKey(ledger.StringToPath("44'/540'/0'/0/0'"))
 * if err != nil {
 * 	fmt.Printf("get public key ERROR: %v\n", err)
 * } else {
 * 	fmt.Printf("public key: %+v\n", publicKey)
 * }
 */
func (device *HidDevice) GetExtendedPublicKey(path BipPath) (*ExtendedPublicKey, error)
```

Gets an address from the specified BIP 32 path.
```
/**
 * @param {BipPath} path The BIP 32 path indexes. Path must begin with `44'/540'/0'/0/i`
 * @return {[]byte} The address for the given path.
 * @return {error} Error value.
 *
 * @example
 * address, err := device.GetAddress(ledger.StringToPath("44'/540'/0'/0/0'"))
 * if err != nil {
 * 	fmt.Printf("get address ERROR: %v\n", err)
 * } else {
 * 	fmt.Printf("address: %+v\n", address)
 * }
 */
func (device *HidDevice) GetAddress(path BipPath) ([]byte, error)
```

Show an address from the specified BIP 32 path for verify.
```
/**
 * @param {BipPath} indexes The path indexes. Path must begin with `44'/540'/0'/0/i`
 * @return {error} Error value.
 *
 * @example
 * err := device.ShowAddress(ledger.StringToPath("44'/540'/0'/0/1'"))
 * if err != nil {
 * 	fmt.Printf("show address ERROR: %v\n", err)
 * } else {
 * 	fmt.Printf("show address: OK\n")
 * }
 */
func (device *HidDevice) ShowAddress(path BipPath) error
```

Sign a transaction by the specified BIP 32 path account address.
```
/**
 * @param {BipPath} path The BIP 32 path indexes. Path must begin with `44'/540'/0'/0/i`
 * @param {[]byte} tx The XDR encoded transaction data, include transaction type
 * @return {[]byte} Signed transaction.
 * @return {error} Error value
 *
 * @example
 * tx := make([]byte, 0)
 * var bin []byte
 * bin, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000") // network id
 * tx = append(tx, bin...)
 * tx = append(tx, 0) // coin transaction with ed
 * tx = append(tx, uint64_to_buf(1)...) // nonce
 * bin, _ = hex.DecodeString("0000000000000000000000000000000000000000") // recepient
 * tx = append(tx, bin...)
 * tx = append(tx, uint64_to_buf(1000000)...) // gas limit
 * tx = append(tx, uint64_to_buf(1000)...) // gas price
 * tx = append(tx, uint64_to_buf(1000000000000)...) // amount
 * tx = append(tx, publicKey.PublicKey...)
 * 
 * response, err := device.SignTx(ledger.StringToPath("44'/540'/0'/0/0'"), tx)
 * if err != nil {
 * 	fmt.Printf("Verify coin tx ERROR: %v\n", err)
 * } else {
 * 	hash := sha512.Sum512(tx)
 * 	fmt.Printf("Verify coin tx: %v\n", ed25519.Verify(publicKey.PublicKey, hash[:], response[1:65]))
 * }
 */
func (device *HidDevice) SignTx(path BipPath, tx []byte) ([]byte, error)
```

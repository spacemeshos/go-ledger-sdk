package ledger

import (
	"fmt"
//	"strconv"
//	"strings"
)

const CLA = 0x30
const MAX_PACKET_LENGTH = 240

const INS_GET_VERSION		= 0x00
const INS_GET_EXT_PUBLIC_KEY	= 0x10
const INS_GET_ADDRESS		= 0x11
const INS_SIGN_TX		= 0x20

const P1_UNUSED 		= 0x00
const P1_RETURN			= 0x01
const P1_DISPLAY		= 0x02
const P1_HAS_HEADER		= 0x01
const P1_HAS_DATA		= 0x02
const P1_IS_LAST		= 0x04

const P2_UNUSED 		= 0x00

type Version struct {
  Major	byte
  Minor	byte
  Patch	byte
  Flags	byte
}

type ExtendedPublicKey struct {
  PublicKey	[]byte
  ChainCode	[]byte
}

type BipPath []uint32

func stripRetcodeFromResponse(response []byte) []byte {
	L := len(response)
	if L < 2 {
		return nil
	}
	if response[L-2] != 0x90 || response[L-1] != 0x00 { // OK code 0x9000
		fmt.Printf("error code %x\n", (uint32(response[L-2]) << 8) + uint32(response[L-1]))
		return nil
	}
	return response[0 : L-2]
}

/**
 * wrapper on top of exchange to simplify work of the implementation.
 * @param cla
 * @param ins
 * @param p1
 * @param p2
 * @param data
 * @param statusList is a list of accepted status code (shorts). [0x9000] by default
 * @return a Promise of response buffer
 */
func (device *HidDevice) send(cla byte, ins byte, p1 byte, p2 byte, data []byte) []byte {
	if len(data) >= 256 {
//		throw new TransportError("data.length exceed 256 bytes limit. Got: " + data.length, "DataLengthTooBig")
		return nil
	}
	buffer := make([]byte, 5 + len(data))
	buffer[0] = cla
	buffer[1] = ins
	buffer[2] = p1
	buffer[3] = p2
	buffer[4] = byte(len(data))
	copy(buffer[5:], data)
	response := device.exchange(buffer)
	if response != nil {
		response = stripRetcodeFromResponse(response)
	}

	return response;
}

/**
 * Returns an object containing the app version.
 *
 * @returns {Promise<GetVersionResponse>} Result object containing the application version number.
 *
 * @example
 * const { major, minor, patch, flags } = await smesh.getVersion();
 * console.log(`App version ${major}.${minor}.${patch}`);
 *
 */
func (device *HidDevice) GetVersion() *Version {
	response := device.send(CLA, INS_GET_VERSION, P1_UNUSED, P2_UNUSED, []byte{})
	if response == nil || len(response) != 4 {
		return nil
	}
	return &Version {
		Major: response[0],
		Minor: response[1], 
		Patch: response[2],
		Flags: response[3],
	}
}

/**
 * @description Get a public key from the specified BIP 32 path.
 *
 * @param {BIP32Path} indexes The path indexes. Path must begin with `44'/540'/n'`, and shuld be 5 indexes long.
 * @return {Promise<GetExtendedPublicKeyResponse>} The public key with chaincode for the given path.
 *
 * @throws 0x6E07 - Some part of request data is invalid
 * @throws 0x6E09 - User rejected the action
 * @throws 0x6E11 - Pin screen
 *
 * @example
 * const { publicKey, chainCode } = await smesh.getExtendedPublicKey([ HARDENED + 44, HARDENED + 540, HARDENED + 1 ]);
 * console.log(publicKey);
 *
 */
func (device *HidDevice) GetExtendedPublicKey(path BipPath) *ExtendedPublicKey {

	data := pathToBytes(path)
	response := device.send(CLA, INS_GET_EXT_PUBLIC_KEY, P1_UNUSED, P2_UNUSED, data)
	fmt.Printf("result len %v\n", len(response))
	if response == nil || len(response) != (32 + 32) {
		return nil
	}

	return &ExtendedPublicKey {
		PublicKey: response[:32],
		ChainCode: response[32:],
	}
}

/**
 * @description Gets an address from the specified BIP 32 path.
 *
 * @param {BIP32Path} indexes The path indexes. Path must begin with `44'/540'/0'/0/i`
 * @return {Promise<GetAddressResponse>} The address for the given path.
 *
 * @throws 0x6E07 - Some part of request data is invalid
 * @throws 0x6E09 - User rejected the action
 * @throws 0x6E11 - Pin screen
 *
 * @example
 * const { address } = await smesh.getAddress([ HARDENED + 44, HARDENED + 540, HARDENED + 0, 0, 2 ]);
 *
 */
func (device *HidDevice) GetAddress(path BipPath) []byte {

	data := pathToBytes(path)
	response := device.send(CLA, INS_GET_ADDRESS, P1_RETURN, P2_UNUSED, data)
	fmt.Printf("result len %v\n", len(response))
	if response == nil || len(response) != 32 {
		return nil
	}

	return response
}

/**
 * @description Show an address from the specified BIP 32 path for verify.
 *
 * @param {BIP32Path} indexes The path indexes. Path must begin with `44'/540'/0'/0/i`
 * @return {Promise<void>} No return.
 *
 * @throws 0x6E07 - Some part of request data is invalid
 * @throws 0x6E11 - Pin screen
 *
 * @example
 * await smesh.showAddress([ HARDENED + 44, HARDENED + 540, HARDENED + 0, 0, HARDENED + 2 ]);
 *
 */
func (device *HidDevice) ShowAddress(path BipPath) bool {

	data := pathToBytes(path)
	response := device.send(CLA, INS_GET_ADDRESS, P1_DISPLAY, P2_UNUSED, data)
	fmt.Printf("result len %v\n", len(response))
	if response == nil || len(response) != 0 {
		return false
	}

	return true
}

/**
 * @description Sign an transaction by the specified BIP 32 path account address.
 *
 * @param {BIP32Path} indexes The path indexes. Path must begin with `44'/540'/0'/0/i`
 * @param {Buffer} tx The XDR encoded transaction data, include transaction type
 * @return {Promise<Buffer>} Signed transaction.
 *
 * @throws 0x6E05 - P1, P2 or payload is invalid
 * @throws 0x6E06 - Request is not valid in the context of previous calls
 * @throws 0x6E07 - Some part of request data is invalid
 * @throws 0x6E09 - User rejected the action
 * @throws 0x6E11 - Pin screen
 *
 * @example
 * const { signature } = await smesh.signTx([ HARDENED + 44, HARDENED + 540, HARDENED + 0, 0, 2 ], txData);
 *
 */
func (device *HidDevice) SignTx(path BipPath, tx []byte) []byte {

	data := pathToBytes(path)
	data = append(data, tx...)
	var response []byte

	fmt.Printf("data length %v\n", len(data))
	fmt.Printf("data %v\n", data)

	if len(data) <= MAX_PACKET_LENGTH {
		response = device.send(CLA, INS_SIGN_TX, P1_HAS_HEADER | P1_IS_LAST, P2_UNUSED, data)
	} else {
		dataSize := len(data)
		chunkSize := MAX_PACKET_LENGTH
		offset := 0
		// Send tx header + tx data
		response = device.send(CLA, INS_SIGN_TX, P1_HAS_HEADER | P1_HAS_DATA, P2_UNUSED, data[offset : offset + chunkSize])
	fmt.Printf("response 2 len %v\n", len(response))
		if response == nil || len(response) != 0 {
			return nil
		}
		dataSize -= chunkSize
		offset += chunkSize
		// Send tx data
		for ; dataSize > MAX_PACKET_LENGTH;  {
			response = device.send(CLA, INS_SIGN_TX, P1_HAS_DATA, P2_UNUSED, data[offset : offset + chunkSize])
	fmt.Printf("response 3 len %v\n", len(response))
			if response == nil || len(response) != 0 {
				return nil
			}
			dataSize -= chunkSize
			offset += chunkSize
		}
		response = device.send(CLA, INS_SIGN_TX, P1_IS_LAST, P2_UNUSED, data[offset:])
	}

	fmt.Printf("response len %v\n", len(response))
	if response == nil || len(response) != (64 + 32) {
		return nil
	}

	result := make([]byte, 64 + len(tx))
	result[0] = tx[0]
	copy(result[1:], response[:64])
	copy(result[65:], tx[1:])

	return result
}

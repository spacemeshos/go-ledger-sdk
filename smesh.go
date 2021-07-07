package ledger

import (
//	"fmt"
)

const CLA = 0x30
const MAX_PACKET_LENGTH = 240

const INS_GET_VERSION		= 0x00
const INS_GET_EXT_PUBLIC_KEY	= 0x10
const INS_GET_ADDRESS		= 0x11

const INS_SIGN_TX		= 0x20

const P1_UNUSED 		= 0x00
const P2_UNUSED 		= 0x00

type Version struct {
  Major	byte
  Minor	byte
  Patch	byte
  Flags	byte
}

func stripRetcodeFromResponse(response []byte) []byte {
	L := len(response)
	if L < 2 {
		return nil
	}
	if response[L-2] != 0x90 || response[L-1] != 0x00 { // OK code 0x9000
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

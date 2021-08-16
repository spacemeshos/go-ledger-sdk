package ledger

import (
	"fmt"
)

const cTag = 0x05
const cPacketSize = 64

type Frame struct {
	data       []byte
	dataLength int
	sequence   int
}

// add chink to APDU response
func (frame *Frame) add(channel int, chunk []byte) (*Frame, error) {

	if chunk[0] != byte((channel>>8)&0xff) || chunk[1] != byte(channel&0xff) {
		return nil, fmt.Errorf("Invalid channel")
	}
	if chunk[2] != cTag {
		return nil, fmt.Errorf("Invalid tag")
	}
	if chunk[3] != byte((frame.sequence>>8)&0xff) || chunk[4] != byte(frame.sequence&0xff) {
		return nil, fmt.Errorf("Invalid sequence")
	}

	if frame.data == nil {
		frame.dataLength = (int(chunk[5]) << 8) + int(chunk[6])
		frame.data = make([]byte, 0)
	}
	if frame.sequence == 0 {
		frame.data = append(frame.data, chunk[7:]...)
	} else {
		frame.data = append(frame.data, chunk[5:]...)
	}
	frame.sequence++
	if len(frame.data) > frame.dataLength {
		frame.data = frame.data[:frame.dataLength]
	}

	return frame, nil
}

// returns result if ready
func (frame *Frame) getResult() []byte {
	if frame != nil && frame.data != nil && frame.dataLength == len(frame.data) {
		return frame.data
	}
	return nil
}

/**
 * Exchange with the device using APDU protocol.
 * @param apdu
 * @return {[]byte} apdu response
 * @return {error} Error value.
 */
func (device *HidDevice) exchange(apdu []byte) ([]byte, error) {
	message := make([]byte, cPacketSize+1)
	dataLength := len(apdu)
	chunkLength := dataLength
	offset := 0

	message[0] = 0
	// Channel
	message[1] = byte((device.channel >> 8) & 0xff)
	message[2] = byte(device.channel & 0xff)
	// Tag
	message[3] = cTag
	// Sequence index for first APDU packet
	message[4] = 0
	message[5] = 0
	// Data length
	message[6] = byte((dataLength >> 8) & 0xff)
	message[7] = byte(dataLength & 0xff)

	if chunkLength > cPacketSize-7 {
		chunkLength = cPacketSize - 7
	}
	dataLength -= chunkLength

	copy(message[8:], apdu[offset:chunkLength])
	// Send first APDU packet
	if writeLength := device.write(message, chunkLength+8); writeLength != (chunkLength+8) && writeLength != (cPacketSize+1) {
		return nil, fmt.Errorf("writeHID error %v", writeLength)
	}
	offset += chunkLength

	for i := 1; dataLength > 0; i++ {
		// Sequence index for this APDU packet
		message[4] = byte((i >> 8) & 0xff)
		message[5] = byte(i & 0xff)

		chunkLength = dataLength
		if chunkLength > cPacketSize-5 {
			chunkLength = cPacketSize - 5
		}
		dataLength -= chunkLength

		copy(message[6:], apdu[offset:offset+chunkLength])
		// Send this APDU packet
		if writeLength := device.write(message, chunkLength+6); writeLength != (chunkLength+6) && writeLength != (cPacketSize+1) {
			return nil, fmt.Errorf("writeHID error %v", writeLength)
		}
		offset += chunkLength
	}

	// Read response
	var result []byte
	var err error
	frame := &Frame{}
	for result = frame.getResult(); result == nil; result = frame.getResult() {
		buffer := device.read()
		if buffer == nil {
			return nil, fmt.Errorf("Buffer is nil")
		}
		frame, err = frame.add(device.channel, buffer)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

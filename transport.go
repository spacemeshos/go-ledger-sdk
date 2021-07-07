package ledger

const LedgerUSBVendorId = 0x2c97

const cTag = 0x05
const cPacketSize = 64

type Frame struct {
	data		[]byte
	dataLength	int
	sequence	int
}

func (frame *Frame) reduceResponse(channel int, chunk []byte) *Frame {

	if chunk[0] != byte((channel >> 8) & 0xff) || chunk[1] != byte(channel & 0xff) {
//		throw new TransportError("Invalid channel", "InvalidChannel");
		return nil
	}
	if chunk[2] != cTag {
//		throw new TransportError("Invalid tag", "InvalidTag");
		return nil
	}
	if chunk[3] != byte((frame.sequence >> 8) & 0xff) || chunk[4] != byte(frame.sequence & 0xff) {
//		throw new TransportError("Invalid sequence", "InvalidSequence");
		return nil
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

	return frame
}

func (frame *Frame) getReducedResult() []byte {
	if frame != nil && frame.data != nil && frame.dataLength == len(frame.data) {
		return frame.data
	}
	return nil
}


func (device *HidDevice) writeHID(message []byte) int {
	return device.write(message)
}

func (device *HidDevice) readHID() []byte {
	data := device.read()
        if data == nil {
          return nil
        }
	return data
}

/**
 * Exchange with the device using APDU protocol.
 * @param apdu
 * @returns a promise of apdu response
 */
func (device *HidDevice) exchange(apdu []byte) []byte {
	message := make([]byte, cPacketSize + 1)
	dataLength := len(apdu)
	chunkLength := dataLength
	offset := 0

	// Write...

	message[0] = 0

        message[1] = byte((device.channel >> 8) & 0xff)
	message[2] = byte(device.channel & 0xff)
	message[3] = cTag

	message[4] = 0
	message[5] = 0

	message[6] = byte((dataLength >> 8) & 0xff)
	message[7] = byte(dataLength & 0xff)
	if chunkLength > cPacketSize - 7 {
		chunkLength = cPacketSize - 7
	}
	dataLength -= chunkLength

	copy(message[8:], apdu[offset:chunkLength])
        device.writeHID(message);
	offset += chunkLength

	for i := 1 ; dataLength > 0; i++ {
		message[4] = byte((i >> 8) & 0xff)
		message[5] = byte(i & 0xff)

		chunkLength = dataLength
		if chunkLength > cPacketSize - 7 {
			chunkLength = cPacketSize - 7
		}
		dataLength -= chunkLength

		copy(message[8:], apdu[offset:offset + chunkLength])
	        device.writeHID(message);
		offset += chunkLength
	}

	// Read...
	var result []byte
	frame := &Frame{}
	for result = frame.getReducedResult(); result == nil; result = frame.getReducedResult() {
		buffer := device.readHID();
		if buffer == nil {
			return nil
		}
		frame = frame.reduceResponse(device.channel, buffer);
		if frame == nil {
			return nil
		}
	}

	return result;
}

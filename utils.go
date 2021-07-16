package ledger

import (
	"strconv"
	"strings"
)

// Parse string to BIP32 path
func StringToPath(pathStr string) BipPath {
	if len(pathStr) == 0 {
		return nil
	}

	items := strings.Split(pathStr, "/")
	path := make(BipPath, len(items))
	for i := 0; i < len(items); i++ {
		var p uint64
		var base uint32
		var err error
		if strings.HasSuffix(items[i], "'") {
			p, err = strconv.ParseUint(items[i][:len(items[i])-1], 10, 32)
			base = 0x80000000
		} else {
			p, err = strconv.ParseUint(items[i], 10, 32)
		}
		if err != nil {
			return nil
		}
		path[i] = base + uint32(p);
	}

	return path
}

// Convert PIB32 path to BE bytes array
func pathToBytes(path BipPath) []byte {
	data := make([]byte, 1 + 4 * len(path))
	data[0] = byte(len(path))

	for i, p := range(path) {
		data[1 + i * 4] = byte((p >> 24) & 0xff)
		data[2 + i * 4] = byte((p >> 16) & 0xff)
		data[3 + i * 4] = byte((p >>  8) & 0xff)
		data[4 + i * 4] = byte( p        & 0xff)
	}
	return data;
}

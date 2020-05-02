// Generate Version 1 or Version 4 UUID's encoded as a 22 alphanumeric
// string. Useful for situations where shorter strings are more desirable.
// 22 character strings are only marginally longer than the 16 byte
//underlying representation, as opposed to the 36 character length UUID.
// Base62 strings are more easily copied and less prone to being split by
// email clients, and when being viewed in text editors.
//
// This library only encodes and decodes 16 byte UUID's. It is not designed
// for general base62 encoding and decoding, and it will not decode standard
// base62 encoded strings.
package base62

import (
	"github.com/google/uuid"
)

const BASE62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var BASE62_REVERSE = []uint64{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0, 0,
	0, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
	51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 0, 0, 0, 0, 0,
	0, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

/* // Make a reverse map of the BASE62 array for decoding
BASE62_REVERSE := make([]uint64, 256)
for i, v := range BASE62 {
	BASE62_REVERSE[v] = uint64(i)
}
	for i, v := range BASE62_REVERSE {
		if i%16 == 0 {
			fmt.Printf("\n    ")
		}
		fmt.Printf("%d, ", v)
	}
	fmt.Println()
*/

// NewUuid generates a 22 character string representation of
// a time based uuid (Version 1).
func NewUuid() string {
	data, err := uuid.NewUUID()
	if err != nil {
		data, err = uuid.NewRandom()
		if err != nil {
			return ""
		}
	}
	return UuidToBase62([16]byte(data))
}

// NewUuid generates a 22 character string representation of
// a random uuid (Version 4)
func NewRandomUuid() string {
	data, err := uuid.NewRandom()
	if err != nil {
		data, err = uuid.NewUUID()
		if err != nil {
			return ""
		}
	}
	return UuidToBase62([16]byte(data))
}

// UuidToBase62 converts a 16 character byte array into a
// base62 encoded string. The first bytes in a time based UUID
// are usually the most likely to vary/change, therefore to ensure
// database sharding operates efficiently the most frequently changing
// time based bytes are interspersed throught the start of the string.
func UuidToBase62(data [16]byte) string {
	if len(data) != 16 {
		return ""
	}

	out := make([]byte, 22)

	p := uint64(data[0]) | uint64(data[1])<<8 |
		uint64(data[2])<<16 | uint64(data[3])<<24 |
		uint64(data[4])<<32 | uint64(data[5])<<40 |
		uint64(data[6])<<48 | uint64(data[7])<<56

	p2 := uint64(data[8]) | uint64(data[9])<<8 |
		uint64(data[10])<<16 | uint64(data[11])<<24 |
		uint64(data[12])<<32 | uint64(data[13])<<40 |
		uint64(data[14])<<48 | uint64(data[15])<<56

	i := 0
	j := 21
	for i < 22 {
		r := p % 62
		out[i] = BASE62[r]
		p = p / 62

		r = p2 % 62
		out[j] = BASE62[r]
		p2 = p2 / 62

		i = i + 2
		j = j - 2
	}

	return string(out)
}

// Base62ToUuid converts a base62 encoded string into a
// 16 byte character array.
func Base62ToUuid(data string) *[16]byte {
	if len(data) != 22 {
		return nil
	}

	var p uint64
	var p2 uint64

	p = BASE62_REVERSE[data[0]] +
		BASE62_REVERSE[data[2]]*62 +
		BASE62_REVERSE[data[4]]*62*62 +
		BASE62_REVERSE[data[6]]*62*62*62 +
		BASE62_REVERSE[data[8]]*62*62*62*62 +
		BASE62_REVERSE[data[10]]*62*62*62*62*62 +
		BASE62_REVERSE[data[12]]*62*62*62*62*62*62 +
		BASE62_REVERSE[data[14]]*62*62*62*62*62*62*62 +
		BASE62_REVERSE[data[16]]*62*62*62*62*62*62*62*62 +
		BASE62_REVERSE[data[18]]*62*62*62*62*62*62*62*62*62 +
		BASE62_REVERSE[data[20]]*62*62*62*62*62*62*62*62*62*62

	p2 = BASE62_REVERSE[data[21]] +
		BASE62_REVERSE[data[19]]*62 +
		BASE62_REVERSE[data[17]]*62*62 +
		BASE62_REVERSE[data[15]]*62*62*62 +
		BASE62_REVERSE[data[13]]*62*62*62*62 +
		BASE62_REVERSE[data[11]]*62*62*62*62*62 +
		BASE62_REVERSE[data[9]]*62*62*62*62*62*62 +
		BASE62_REVERSE[data[7]]*62*62*62*62*62*62*62 +
		BASE62_REVERSE[data[5]]*62*62*62*62*62*62*62*62 +
		BASE62_REVERSE[data[3]]*62*62*62*62*62*62*62*62*62 +
		BASE62_REVERSE[data[1]]*62*62*62*62*62*62*62*62*62*62

	var out [16]byte
	out[0] = byte(p)
	out[1] = byte(p >> 8)
	out[2] = byte(p >> 16)
	out[3] = byte(p >> 24)
	out[4] = byte(p >> 32)
	out[5] = byte(p >> 40)
	out[6] = byte(p >> 48)
	out[7] = byte(p >> 56)
	out[8] = byte(p2)
	out[9] = byte(p2 >> 8)
	out[10] = byte(p2 >> 16)
	out[11] = byte(p2 >> 24)
	out[12] = byte(p2 >> 32)
	out[13] = byte(p2 >> 40)
	out[14] = byte(p2 >> 48)
	out[15] = byte(p2 >> 56)

	return &out
}

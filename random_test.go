package base62

import (
	crand "crypto/rand"
	"encoding/binary"

	"math/rand"
)

var seeded = false

func RandomString(size int) string {
	random := rand.New(NewCryptoSeededSource())

	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		b := uint8(random.Int31n(62))
		if b < 26 {
			bytes[i] = 'A' + b
		} else if b < 52 {
			bytes[i] = 'a' + (b - 26)
		} else {
			bytes[i] = '0' + (b - 52)
		}
	}
	return string(bytes)
}

const rst = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"

func NewCryptoSeededSource() rand.Source {
	var seed int64
	binary.Read(crand.Reader, binary.BigEndian, &seed)
	return rand.NewSource(seed)
}

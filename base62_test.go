package base62

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkBase62Encode(b *testing.B) {
	data := [16]byte{230, 2, 3, 4, 5, 6, 7, 8, 101, 102, 103, 104, 105, 106, 107, 108}
	for n := 0; n < b.N; n++ {
		UuidToBase62(data)
	}
}

func BenchmarkBase62Decode(b *testing.B) {
	data := "1234qwer1234zxcvitow41"
	for n := 0; n < b.N; n++ {
		Base62ToUuid(data)
	}
}

func BenchmarkBase64Decode(b *testing.B) {
	data := "091375D5-8A48-4F80-B5AB-EBB5AD2E36E4"
	for n := 0; n < b.N; n++ {
		uuid.Parse(data)
	}
}

func BenchmarkBase62EncodeDecode(b *testing.B) {
	data := [16]byte{230, 2, 3, 4, 5, 6, 7, 8, 101, 102, 103, 104, 105, 106, 107, 108}
	for n := 0; n < b.N; n++ {
		out := UuidToBase62(data)
		Base62ToUuid(out)
	}
}

func BenchmarkBase62Generate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewUuid()
	}
}

func BenchmarkBase64Generate(b *testing.B) {
	for n := 0; n < b.N; n++ {
		data, err := uuid.NewUUID()
		if err != nil {
			continue
		}
		data.String()
	}
}

func TestEncodeDecode(t *testing.T) {

	{
		data := [16]byte{230, 2, 3, 4, 5, 6, 7, 8, 101, 102, 103, 104, 105, 106, 107, 108}
		out := UuidToBase62(data)
		data2 := Base62ToUuid(out)
		if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", *data2) {
			t.Fatalf("encode/decode failed")

		}
	}

	{
		data := [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
		out := UuidToBase62(data)
		data2 := Base62ToUuid(out)
		if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", *data2) {
			t.Fatalf("encode/decode failed")
		}

	}

	{
		data := [16]byte{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160}
		out := UuidToBase62(data)
		data2 := Base62ToUuid(out)
		if fmt.Sprintf("%v", data) != fmt.Sprintf("%v", *data2) {
			t.Fatalf("encode/decode failed")
		}
	}

	for i := 0; i < 10; i++ {
		v := NewUuid()
		data := Base62ToUuid(v)
		v2 := UuidToBase62(*data)
		if v != v2 {
			t.Fatalf("encode/decode failed")
		}
	}

	for i := 0; i < 10; i++ {
		// 22 characters of base 62 can hold a little over what can be
		// held in a 128 bit number, thus our random strings must
		// be zeroed out at points to prevent that overflow
		v := RandomString(1) + "0" +
			RandomString(18) + "0" +
			RandomString(1)
		data := Base62ToUuid(v)
		v2 := UuidToBase62(*data)
		//fmt.Printf("Test base62 uuid using random string: %s\n", v)
		//fmt.Printf("    %v\n", *data)
		//fmt.Printf("    %v\n", v2)
		//data2 := base62_to_uuid(v2)
		//fmt.Printf("    %v\n", *data2)
		if v != v2 {
			t.Fatalf("encode/decode failed %s != %s", v, v2)
		}
	}
}

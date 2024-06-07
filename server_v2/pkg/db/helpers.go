package db

import (
	"encoding/binary"
	"math"
	"time"

	"github.com/google/uuid"
)

func FLOAT_BYTES(f float64) []byte {
	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], math.Float64bits(f))
	return buffer[:]
}

func INT_BYTES(i int) []byte {
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], uint32(i))
	return buffer[:]
}

func STRING_BYTES(s string) []byte {
	return []byte(s)
}

func BOOL_BYTES(bo bool) []byte {
	var b byte
	if bo {
		b = 1
	}
	return []byte{b}
}

func UUID_BYTES(id uuid.UUID) []byte {
	return id[:]
}

func TIME_BYTES(t time.Time) []byte {
	var buffer [8]byte
	binary.BigEndian.PutUint64(buffer[:], uint64(t.UnixNano()))
	return buffer[:]
}

func BYTES_FLOAT(b []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b))
}

func BYTES_INT(b []byte) int {
	return int(binary.BigEndian.Uint32(b))
}

func BYTES_STRING(b []byte) string {
	return string(b)
}

func BYTES_BOOL(b []byte) bool {
	if b[0] == 0x01 {
		return true
	} else {
		return false
	}
}

func BYTES_UUID(b []byte) uuid.UUID {
	return uuid.UUID([16]byte(b))
}

func BYTES_TIME(b []byte) time.Time {
	return time.Unix(0, int64(binary.BigEndian.Uint64(b)))
}

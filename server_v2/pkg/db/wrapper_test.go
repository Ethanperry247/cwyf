package db

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"
)

type mockStruct struct {
	Int   int
	Str   string
	Bool  bool
	Float float64
}

const (
	mockId    string = "my_id"
	mockKind  Kind   = "kind"
	mockTag   string = "my_tag"
	mockIdTwo string = "my_id_two"
)

var (
	mockMappings = Mapping[mockStruct]{
		"Int": AssembleDisassemble[mockStruct]{
			A: func(ms *mockStruct, b []byte) {
				ms.Int = BYTES_INT(b)
			}, D: func(ms *mockStruct) []byte {
				return INT_BYTES(ms.Int)
			},
		},
		"Str": AssembleDisassemble[mockStruct]{
			A: func(ms *mockStruct, b []byte) {
				ms.Str = BYTES_STRING(b)
			}, D: func(ms *mockStruct) []byte {
				return STRING_BYTES(ms.Str)
			},
		},
		"Bool": AssembleDisassemble[mockStruct]{
			A: func(ms *mockStruct, b []byte) {
				ms.Bool = BYTES_BOOL(b)
			}, D: func(ms *mockStruct) []byte {
				return BOOL_BYTES(ms.Bool)
			},
		},
		"Float": AssembleDisassemble[mockStruct]{
			A: func(ms *mockStruct, b []byte) {
				ms.Float = BYTES_FLOAT(b)
			}, D: func(ms *mockStruct) []byte {
				return FLOAT_BYTES(ms.Float)
			},
		},
	}
	mockObject = &mockStruct{
		Int:   1,
		Str:   "1",
		Bool:  true,
		Float: 1.1,
	}
	mockKvs = []KV{{K: []uint8{0x6b, 0x69, 0x6e, 0x64, 0x2e, 0x6d, 0x79, 0x5f, 0x69, 0x64, 0x2e, 0x49, 0x6e, 0x74}, V: []uint8{0x0, 0x0, 0x0, 0x1}}, {K: []uint8{0x6b, 0x69, 0x6e, 0x64, 0x2e, 0x6d, 0x79, 0x5f, 0x69, 0x64, 0x2e, 0x53, 0x74, 0x72}, V: []uint8{0x31}}, {K: []uint8{0x6b, 0x69, 0x6e, 0x64, 0x2e, 0x6d, 0x79, 0x5f, 0x69, 0x64, 0x2e, 0x42, 0x6f, 0x6f, 0x6c}, V: []uint8{0x1}}, {K: []uint8{0x6b, 0x69, 0x6e, 0x64, 0x2e, 0x6d, 0x79, 0x5f, 0x69, 0x64, 0x2e, 0x46, 0x6c, 0x6f, 0x61, 0x74}, V: []uint8{0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}}}
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCreateGetAndDelete(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(t, err)
	defer db.Close()

	wrapper := New[mockStruct](db, mockKind, mockMappings)
	require.NoError(t, err)

	err = wrapper.Create(mockId, mockObject)
	require.NoError(t, err)

	res, err := wrapper.Get(mockId)
	require.NoError(t, err)
	require.Equal(t, mockObject, res)

	err = wrapper.Delete(mockId)
	require.NoError(t, err)

	_, err = wrapper.Get(mockId)
	require.Equal(t, &NotFoundError{id: mockId}, err)
}

func TestTagCreateAndListByTag(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(t, err)
	defer db.Close()

	wrapper := New[mockStruct](db, mockKind, mockMappings)
	require.NoError(t, err)

	err = wrapper.Create(mockId, mockObject, mockTag)
	require.NoError(t, err)

	res, err := wrapper.ListByTag(mockTag)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{mockId}, res)

	err = wrapper.Create(mockIdTwo, mockObject, mockTag)
	require.NoError(t, err)

	res, err = wrapper.ListByTag(mockTag)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{mockId, mockIdTwo}, res)

	err = wrapper.Delete(mockId)
	require.NoError(t, err)

	res, err = wrapper.ListByTag(mockTag)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{mockIdTwo}, res)

	err = wrapper.Delete(mockIdTwo)
	require.NoError(t, err)

	res, err = wrapper.ListByTag(mockTag)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{}, res)
}

func BenchmarkDBCreate(b *testing.B) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(b, err)
	defer db.Close()

	wrapper := New[mockStruct](db, mockKind, mockMappings)

	err = wrapper.Create(mockId, mockObject)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.Create(mockId, mockObject)
	}
}

func BenchmarkDBCreateAndDelete(b *testing.B) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(b, err)
	defer db.Close()

	wrapper := New[mockStruct](db, mockKind, mockMappings)

	err = wrapper.Create(mockId, mockObject)
	require.NoError(b, err)

	err = wrapper.Delete(mockId)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.Create(mockId, mockObject)
		wrapper.Delete(mockId)
	}
}

func BenchmarkDBGet(b *testing.B) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(b, err)
	defer db.Close()

	wrapper := New[mockStruct](db, mockKind, mockMappings)

	err = wrapper.Create(mockId, mockObject)
	require.NoError(b, err)

	_, err = wrapper.Get(mockId)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.Get(mockId)
	}
}

func BenchmarkDBTagAndList(b *testing.B) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	require.NoError(b, err)
	defer db.Close()

	wrapper := New[mockStruct](db, mockKind, mockMappings)

	err = wrapper.Create(mockId, mockObject, mockTag)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.ListByTag(mockTag)
	}
}

func BenchmarkNoopDBCreate(b *testing.B) {

	wrapper := New[mockStruct](&NoopBadgerDatabase{}, mockKind, mockMappings)

	err := wrapper.Create(mockId, mockObject)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.Create(mockId, mockObject)
	}
}

func BenchmarkNoopDBCreateAndDelete(b *testing.B) {

	wrapper := New[mockStruct](&NoopBadgerDatabase{}, mockKind, mockMappings)

	err := wrapper.Create(mockId, mockObject)
	require.NoError(b, err)

	err = wrapper.Delete(mockId)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.Create(mockId, mockObject)
		wrapper.Delete(mockId)
	}
}

func BenchmarkNoopDBGet(b *testing.B) {
	wrapper := New[mockStruct](&NoopBadgerDatabase{}, mockKind, mockMappings)

	err := wrapper.Create(mockId, mockObject)
	require.NoError(b, err)

	_, err = wrapper.Get(mockId)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.Get(mockId)
	}
}

func BenchmarkNoopDBTagAndList(b *testing.B) {
	wrapper := New[mockStruct](&NoopBadgerDatabase{}, mockKind, mockMappings)

	err := wrapper.Create(mockId, mockObject, mockTag)
	require.NoError(b, err)

	for index := 0; index < b.N; index++ {
		wrapper.ListByTag(mockTag)
	}
}

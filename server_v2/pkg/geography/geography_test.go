package geography

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	point = &Point{
		Latitude: 47.597400889267455, Longitude: -122.12061416585233,
	}
	hashed = &Point{
		Latitude: 47.5974, Longitude: -122.12061,
	}
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestPointHash(t *testing.T) {
	geo := New()

	res, err := geo.Hash(point)
	require.NoError(t, err)
	require.Equal(t, hashed, res)
}

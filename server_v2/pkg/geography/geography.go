package geography

import "fmt"

const (
	LAT_DEGREES  = 180
	LONG_DEGREES = 360
)

type Point struct {
	Latitude  float64
	Longitude float64
}

type Geography struct {
	Divisions int
}

func New(options ...func(*Geography)) *Geography {
	geo := &Geography{
		Divisions: 100000,
	}

	for _, option := range options {
		option(geo)
	}

	return geo
}

func (geo *Geography) validate(point *Point) error {
	if point.Latitude > 90 || point.Latitude < -90 {
		return &InvalidLatitudeError{
			value: point.Latitude,
		}
	}

	if point.Longitude > 180 || point.Longitude < -180 {
		return &InvalidLongitudeError{
			value: point.Longitude,
		}
	}

	return nil
}

func (geo *Geography) Subdivisions(point *Point, len int) ([]*Point, error) {
	err := geo.validate(point)
	if err != nil {
		return nil, err
	}

	points := make([]*Point, len*len)

	hash := geo.hash(point)
	for latIndex := 0; latIndex < len; latIndex++ {
		for longIndex := 0; longIndex < len; longIndex++ {
			points[latIndex*len+longIndex] = &Point{
				Latitude:  hash.Latitude + float64(latIndex)*(float64(geo.Divisions)/LAT_DEGREES),
				Longitude: hash.Longitude + float64(longIndex)*(float64(geo.Divisions)/LONG_DEGREES),
			}
		}
	}

	return points, nil
}

func (geo *Geography) Hash(point *Point) (*Point, error) {
	err := geo.validate(point)
	if err != nil {
		return nil, err
	}

	return geo.hash(point), nil
}

func (geo *Geography) hash(point *Point) *Point {
	return &Point{
		Latitude:  float64(int(point.Latitude*float64(geo.Divisions))) / float64(geo.Divisions),
		Longitude: float64(int(point.Longitude*float64(geo.Divisions*2))) / float64(geo.Divisions*2),
	}
}

func (geo *Geography) AsString(point *Point) string {
	return fmt.Sprintf("%f_%f", point.Latitude, point.Longitude)
}

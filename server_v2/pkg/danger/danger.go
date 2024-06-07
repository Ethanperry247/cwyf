package danger

type (
	Point struct {
		Latitude  float64
		Longitude float64
	}
	AltitudePoint struct {
		Point
		Altitude float64
	}
	HazardPoint struct {
		AltitudePoint
		Hazard int
	}
)

const (
	DEFAULT_HAZARD_MIN = 0
	DEFAULT_HAZARD_MAX = 100
)
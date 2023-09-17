package internal

import (
	"danger-dodgers/pkg/db"
	"time"
)

type (
	ID struct {
		UserID string
	}
	Token struct {
		Token string
	}
	User struct {
		Name     string
		Email    string
		Username string
		Password string
	}
	Report struct {
		ID          string
		Description string
		Tag         string
		Title       string
		Timestamp   time.Time
		UserID      string
		Latitude  float64
		Longitude float64
	}
	Point struct {
		Latitude  float64
		Longitude float64
	}
	Position struct {
		ID string
		Point
	}
	AltitudePosition struct {
		Position
		Altitude float64
	}
	Area struct {
		Point
		Radius float64
	}
	Route struct {
		ID string
	}
)

const (
	USER              db.Kind = "user"
	REPORT            db.Kind = "report"
	POINT             db.Kind = "point"
	POSITION          db.Kind = "position"
	ALTITUDE_POSITION db.Kind = "altitudeposition"
	AREA              db.Kind = "area"
	ROUTE             db.Kind = "route"
)

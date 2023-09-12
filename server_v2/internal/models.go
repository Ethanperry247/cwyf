package internal

import (
	"danger-dodgers/pkg/db"
	"time"

	"github.com/google/uuid"
)

type (
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
		ID          uuid.UUID
		Tag         string
		Description string
		Title       string
		Timestamp   time.Time
		UserID      string
		PositionID  string
	}
	Point struct {
		Latitude  float64
		Longitude float64
	}
	Position struct {
		ID uuid.UUID
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
		ID uuid.UUID
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

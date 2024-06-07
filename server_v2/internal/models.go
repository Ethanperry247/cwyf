package internal

import (
	"time"
)

type (
	ID struct {
		UserID string
	}
	Token struct {
		Token string `json:"token"`
	}
	User struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	ModifyUser struct {
		CurrentUser *User `json:"currentUser"`
		NewUser     *User `json:"newUser"`
	}
	Report struct {
		ID          string    `json:"id"`
		Description string    `json:"description"`
		Tag         string    `json:"tag"`
		Title       string    `json:"title"`
		Timestamp   time.Time `json:"timestamp"`
		UserID      string    `json:"userID"`
		ActivityID  string    `json:"activityID"`
		Latitude    float64   `json:"latitude"`
		Longitude   float64   `json:"longitude"`
	}
	Point struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	Position struct {
		ID string `json:"id"`
		Point
	}
	AltitudePosition struct {
		Position
		Altitude float64 `json:"altitude"`
	}
	IndexedPosition struct {
		AltitudePosition
		Index int `json:"index"`
	}
	ActivityPosition struct {
		ActivityID string `json:"activityID"`
		UserID     string `json:"userID"`
		IndexedPosition
		Timestamp time.Time `json:"timestamp"`
	}
	Activity struct {
		ID     string `json:"id"`
		UserID string `json:"userID"`
	}
	Directions struct {
		Positions []IndexedPosition `json:"positions"`
		Reports   []Report          `json:"reports"`
	}
	DirectionsRequest struct {
	}
	Feedback struct {
		Message   string    `json:"positions"`
		Timestamp time.Time `json:"timestamp"`
		UserID    string    `json:"userID"`
	}
)

const (
	USER              string = "user"
	REPORT            string = "report"
	POINT             string = "point"
	POSITION          string = "position"
	ALTITUDE_POSITION string = "altitudeposition"
	ACTIVITY_POSITION string = "activityposition"
	AREA              string = "area"
	ROUTE             string = "route"
	ACTIVITY          string = "activity"
)

const (
	MIN_OPTIONAL_FIELD_LENGTH = 0
	MIN_REQUIRED_FIELD_LENGTH = 1
	MAX_STANDARD_FIELD_LENGTH = 100
	MAX_EXTENDED_FIELD_LENGTH = 1000
	MIN_PASSWORD_FIELD_LENGTH = 8
)

var (
	REQUIRED_STANDARD_VERIFIER = func(v FieldVerifier) FieldVerifier {
		return v.WithMin(MIN_REQUIRED_FIELD_LENGTH).WithMax(MAX_STANDARD_FIELD_LENGTH)
	}
	OPTIONAL_STANDARD_VERIFIER = func(v FieldVerifier) FieldVerifier {
		return v.WithMin(MIN_OPTIONAL_FIELD_LENGTH).WithMax(MAX_STANDARD_FIELD_LENGTH)
	}
	REQUIRED_EXTENDED_VERIFIER = func(v FieldVerifier) FieldVerifier {
		return v.WithMin(MIN_REQUIRED_FIELD_LENGTH).WithMax(MAX_EXTENDED_FIELD_LENGTH)
	}
	OPTIONAL_EXTENDED_VERIFIER = func(v FieldVerifier) FieldVerifier {
		return v.WithMin(MIN_OPTIONAL_FIELD_LENGTH).WithMax(MAX_EXTENDED_FIELD_LENGTH)
	}
	PASSWORD_VERIFIER = func(v FieldVerifier) FieldVerifier {
		return v.WithMin(MIN_PASSWORD_FIELD_LENGTH).WithMax(MAX_STANDARD_FIELD_LENGTH).WithField("password")
	}
)

package internal

import (
	"danger-dodgers/pkg/db"
)

var (
	UserMapping = db.Mapping[User]{
		"Name": db.AssembleDisassemble[User]{
			A: func(t *User, b []byte) {
				t.Name = db.BYTES_STRING(b)
			},
			D: func(t *User) []byte {
				return db.STRING_BYTES(t.Name)
			},
		}, "Email": db.AssembleDisassemble[User]{
			A: func(t *User, b []byte) {
				t.Email = db.BYTES_STRING(b)
			},
			D: func(t *User) []byte {
				return db.STRING_BYTES(t.Email)
			},
		}, "Username": db.AssembleDisassemble[User]{
			A: func(t *User, b []byte) {
				t.Username = db.BYTES_STRING(b)
			},
			D: func(t *User) []byte {
				return db.STRING_BYTES(t.Username)
			},
		}, "Password": db.AssembleDisassemble[User]{
			A: func(t *User, b []byte) {
				t.Password = db.BYTES_STRING(b)
			},
			D: func(t *User) []byte {
				return db.STRING_BYTES(t.Password)
			},
		},
	}

	ReportMapping = db.Mapping[Report]{
		"ID": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.ID = db.BYTES_STRING(b)
			}, D: func(t *Report) []byte {
				return db.STRING_BYTES(t.ID)
			},
		}, "Tag": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.Tag = db.BYTES_STRING(b)
			},
			D: func(t *Report) []byte {
				return db.STRING_BYTES(t.Tag)
			},
		}, "Description": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.Description = db.BYTES_STRING(b)
			},
			D: func(t *Report) []byte {
				return db.STRING_BYTES(t.Description)
			},
		}, "Title": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.Title = db.BYTES_STRING(b)
			},
			D: func(t *Report) []byte {
				return db.STRING_BYTES(t.Title)
			},
		}, "Timestamp": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.Timestamp = db.BYTES_TIME(b)
			}, D: func(t *Report) []byte {
				return db.TIME_BYTES(t.Timestamp)
			},
		}, "UserID": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.UserID = db.BYTES_STRING(b)
			},
			D: func(t *Report) []byte {
				return db.STRING_BYTES(t.UserID)
			},
		}, "Latitude": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.Latitude = db.BYTES_FLOAT(b)
			}, D: func(t *Report) []byte {
				return db.FLOAT_BYTES(t.Latitude)
			},
		}, "Longitude": db.AssembleDisassemble[Report]{
			A: func(t *Report, b []byte) {
				t.Longitude = db.BYTES_FLOAT(b)
			}, D: func(t *Report) []byte {
				return db.FLOAT_BYTES(t.Longitude)
			},
		},
	}

	PositionMapping = db.Mapping[Position]{
		"ID": db.AssembleDisassemble[Position]{
			A: func(t *Position, b []byte) {
				t.ID = db.BYTES_STRING(b)
			}, D: func(t *Position) []byte {
				return db.STRING_BYTES(t.ID)
			},
		}, "Latitude": db.AssembleDisassemble[Position]{
			A: func(t *Position, b []byte) {
				t.Latitude = db.BYTES_FLOAT(b)
			}, D: func(t *Position) []byte {
				return db.FLOAT_BYTES(t.Latitude)
			},
		}, "Longitude": db.AssembleDisassemble[Position]{
			A: func(t *Position, b []byte) {
				t.Longitude = db.BYTES_FLOAT(b)
			}, D: func(t *Position) []byte {
				return db.FLOAT_BYTES(t.Longitude)
			},
		},
	}

	AltitudePositionMapping = db.Mapping[AltitudePosition]{
		"ID": db.AssembleDisassemble[AltitudePosition]{
			A: func(t *AltitudePosition, b []byte) {
				t.ID = db.BYTES_STRING(b)
			}, D: func(t *AltitudePosition) []byte {
				return db.STRING_BYTES(t.ID)
			},
		}, "Latitude": db.AssembleDisassemble[AltitudePosition]{
			A: func(t *AltitudePosition, b []byte) {
				t.Latitude = db.BYTES_FLOAT(b)
			}, D: func(t *AltitudePosition) []byte {
				return db.FLOAT_BYTES(t.Latitude)
			},
		}, "Longitude": db.AssembleDisassemble[AltitudePosition]{
			A: func(t *AltitudePosition, b []byte) {
				t.Longitude = db.BYTES_FLOAT(b)
			}, D: func(t *AltitudePosition) []byte {
				return db.FLOAT_BYTES(t.Longitude)
			},
		}, "Altitude": db.AssembleDisassemble[AltitudePosition]{
			A: func(t *AltitudePosition, b []byte) {
				t.Altitude = db.BYTES_FLOAT(b)
			}, D: func(t *AltitudePosition) []byte {
				return db.FLOAT_BYTES(t.Altitude)
			},
		},
	}

	AreaMapping = db.Mapping[Area]{
		"Latitude": db.AssembleDisassemble[Area]{
			A: func(t *Area, b []byte) {
				t.Latitude = db.BYTES_FLOAT(b)
			}, D: func(t *Area) []byte {
				return db.FLOAT_BYTES(t.Latitude)
			},
		}, "Longitude": db.AssembleDisassemble[Area]{
			A: func(t *Area, b []byte) {
				t.Longitude = db.BYTES_FLOAT(b)
			}, D: func(t *Area) []byte {
				return db.FLOAT_BYTES(t.Longitude)
			},
		}, "Radius": db.AssembleDisassemble[Area]{
			A: func(t *Area, b []byte) {
				t.Radius = db.BYTES_FLOAT(b)
			}, D: func(t *Area) []byte {
				return db.FLOAT_BYTES(t.Radius)
			},
		},
	}

	RouteMapping = db.Mapping[Route]{
		"ID": db.AssembleDisassemble[Route]{
			A: func(t *Route, b []byte) {
				t.ID = db.BYTES_STRING(b)
			}, D: func(t *Route) []byte {
				return db.STRING_BYTES(t.ID)
			},
		},
	}
)

package geography

import "fmt"

type InvalidLongitudeError struct {
	value float64
}

func (err *InvalidLongitudeError) Error() string {
	return fmt.Sprintf("longitude is invalid -- provided value is %0.2f but value must be between -180 and 180", err.value)
}

type InvalidLatitudeError struct {
	value float64
}

func (err *InvalidLatitudeError) Error() string {
	return fmt.Sprintf("latitude is invalid -- provided value is %0.2f but value must be between -90 and 90", err.value)
}

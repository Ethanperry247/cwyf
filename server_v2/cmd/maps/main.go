package main

import (
	"context"
	"fmt"
	"os"

	"googlemaps.github.io/maps"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func GetMapsAPIKey() string {
	return os.Getenv("MAPS_API_KEY")
}

func run() error {
	gMaps, err := maps.NewClient(maps.WithAPIKey(GetMapsAPIKey()))
	if err != nil {
		return err
	}

	r := &maps.DirectionsRequest{
		Origin:      "Sydney",
		Destination: "Perth",
		Mode:        maps.TravelModeBicycling,
	}

	route, _, err := gMaps.Directions(context.Background(), r)
	if err != nil {
		return err
	}

	fmt.Printf("%v", route)

	return nil
}

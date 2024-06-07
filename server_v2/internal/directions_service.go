package internal

import "googlemaps.github.io/maps"

type DirectionsService struct {
	client *maps.Client
}

func NewDirectionsService(client *maps.Client) *DirectionsService {
	return &DirectionsService{
		client: client,
	}
}

func (service *DirectionsService) Route(request *DirectionsRequest) (*Directions, error) {

	return nil, nil
}
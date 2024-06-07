package internal

import "danger-dodgers/pkg/db"

var (
	PositionVerifier         = FieldVerifier{}.WithModel(ACTIVITY_POSITION)
	VerifyPositionID         = REQUIRED_STANDARD_VERIFIER(PositionVerifier).WithField("ID").Build()
	VerifyPositionUserID     = REQUIRED_STANDARD_VERIFIER(PositionVerifier).WithField("userID").Build()
	VerifyPositionActivityID = REQUIRED_STANDARD_VERIFIER(PositionVerifier).WithField("activityID").Build()
)

type PositionService struct {
	positionDB db.Database[ActivityPosition]
	activityDB db.Database[Activity]
}

func NewPositionService(positionDB db.Database[ActivityPosition], activityDB db.Database[Activity]) *PositionService {
	return &PositionService{
		positionDB: positionDB,
		activityDB: activityDB,
	}
}

func (service *PositionService) Create(position *ActivityPosition) error {
	err := VerifyPositionID(position.ID)
	if err != nil {
		return err
	}

	err = Peek(position.ID, ACTIVITY_POSITION, service.positionDB)
	if err != nil {
		return err
	}

	err = VerifyPositionUserID(position.UserID)
	if err != nil {
		return err
	}

	err = VerifyPositionActivityID(position.ActivityID)
	if err != nil {
		return err
	}

	return service.positionDB.Create(position.ID, position, position.ActivityID)
}

func (service *PositionService) ListByActivity(activity *Activity) ([]*ActivityPosition, error) {
	err := VerifyActivityID(activity.ID)
	if err != nil {
		return nil, err
	}

	res, err := service.activityDB.Get(activity.ID)
	if err != nil {
		return nil, err
	}

	err = CompareUserIDs(res.UserID, activity.UserID)
	if err != nil {
		return nil, err
	}

	positionIds, err := service.positionDB.ListByTag(activity.ID)
	if err != nil {
		return nil, err
	}

	positions := make([]*ActivityPosition, len(positionIds))
	for idx, id := range positionIds {
		position, err := service.positionDB.Get(id)
		if err != nil {
			return nil, err
		}
		positions[idx] = position
	}

	return positions, nil
}

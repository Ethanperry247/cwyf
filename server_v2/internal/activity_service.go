package internal

import "danger-dodgers/pkg/db"

var (
	ActivityVerifier     = FieldVerifier{}.WithModel(ACTIVITY)
	VerifyActivityID     = REQUIRED_STANDARD_VERIFIER(ActivityVerifier).WithField("ID").Build()
	VerifyActivityUserID = REQUIRED_STANDARD_VERIFIER(ActivityVerifier).WithField("userID").Build()
)

type ActivityService struct {
	db db.Database[Activity]
}

func NewActivityService(activityDB db.Database[Activity]) *ActivityService {
	return &ActivityService{
		db: activityDB,
	}
}

func (service *ActivityService) Create(activity *Activity) error {
	
	err := VerifyActivityID(activity.ID)
	if err != nil {
		return err
	}

	err = Peek(activity.ID, ACTIVITY, service.db)
	if err != nil {
		return err
	}

	err = VerifyActivityUserID(activity.UserID)
	if err != nil {
		return err
	}

	return service.db.Create(activity.ID, activity)
}

func (service *ActivityService) Delete(activity *Activity) error {
	err := VerifyActivityID(activity.ID)
	if err != nil {
		return err
	}

	res, err := service.db.Get(activity.ID)
	if err != nil {
		return err
	}

	err = CompareUserIDs(activity.UserID, res.UserID)
	if err != nil {
		return err
	}
	
	return service.db.Delete(activity.ID)
}

func (service *ActivityService) ListByUser(id string) ([]string, error) {
	return service.db.ListByTag(id)
}

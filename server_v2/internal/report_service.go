package internal

import (
	"danger-dodgers/pkg/db"
	"danger-dodgers/pkg/geography"
	"time"
)

var (
	ReportVerifier    = FieldVerifier{}.WithModel(REPORT)
	VerifyReportID          = REQUIRED_STANDARD_VERIFIER(ReportVerifier).WithField("ID").Build()
	VerifyReportDescription = OPTIONAL_EXTENDED_VERIFIER(ReportVerifier).WithField("description").Build()
	VerifyReportTitle       = REQUIRED_STANDARD_VERIFIER(ReportVerifier).WithField("title").Build()
	VerifyReportUserID      = REQUIRED_STANDARD_VERIFIER(ReportVerifier).WithField("userID").Build()
	VerifyReportActivityID  = OPTIONAL_STANDARD_VERIFIER(ReportVerifier).WithField("activityID").Build()
)

type ReportService struct {
	db  db.Database[Report]
	geo *geography.Geography
}

func NewReportService(db db.Database[Report]) *ReportService {
	return &ReportService{
		db: db,
	}
}

func (service *ReportService) Create(report *Report) error {

	err := VerifyReportID(report.ID)
	if err != nil {
		return err
	}

	err = Peek(report.ID, REPORT, service.db)
	if err != nil {
		return err
	}

	err = VerifyReportTitle(report.Title)
	if err != nil {
		return err
	}

	err = VerifyReportUserID(report.UserID)
	if err != nil {
		return err
	}

	err = VerifyReportDescription(report.Description)
	if err != nil {
		return err
	}

	err = VerifyReportActivityID(report.ActivityID)
	if err != nil {
		return err
	}

	point, err := service.geo.Hash(&geography.Point{
		Latitude:  report.Latitude,
		Longitude: report.Longitude,
	})
	if err != nil {
		return err
	}

	report.Timestamp = time.Now()

	tags := []string{service.geo.AsString(point), report.UserID}

	if report.Tag == "" {
		tags = append(tags, report.Tag)
	}

	if report.ActivityID == "" {
		tags = append(tags, report.ActivityID)
	}

	return service.db.Create(report.ID, report, tags...)
}

func (service *ReportService) Get(report *Report) (*Report, error) {
	err := VerifyReportID(report.ID)
	if err != nil {
		return nil, err
	}

	res, err := service.db.Get(report.ID)
	if err != nil {
		return nil, err
	}

	err = CompareUserIDs(report.UserID, res.UserID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (service *ReportService) Delete(report *Report) error {
	err := VerifyReportID(report.ID)
	if err != nil {
		return err
	}

	res, err := service.db.Get(report.ID)
	if err != nil {
		return err
	}

	err = CompareUserIDs(report.UserID, res.UserID)
	if err != nil {
		return err
	}

	return service.db.Delete(report.ID)
}

func (service *ReportService) ListByUser(id string) ([]string, error) {
	err := VerifyReportID(id)
	if err != nil {
		return nil, err
	}

	return service.db.ListByTag(id)
}

func (service *ReportService) ListByGeography(latitude float64, longitude float64, len int) ([]string, error) {
	var ids []string
	points, err := service.geo.Subdivisions(&geography.Point{
		Latitude:  latitude,
		Longitude: longitude,
	}, len)
	if err != nil {
		return nil, err
	}
	for _, point := range points {
		res, err := service.db.ListByTag(service.geo.AsString(point))
		if err != nil {
			return nil, err
		}
		ids = append(ids, res...)
	}

	return ids, nil
}

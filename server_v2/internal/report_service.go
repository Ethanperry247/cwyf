package internal

import (
	"danger-dodgers/pkg/db"
	"danger-dodgers/pkg/geography"
	"time"
)

const (
	MAX_REPORT_ID_LENGTH          = 100
	MAX_REPORT_TITLE_LENGTH       = 100
	MAX_REPORT_TAG_LENGTH         = 100
	MAX_REPORT_DESCRIPTION_LENGTH = 1000
	MAX_REPORT_USER_ID_LENGTH     = 100
	MAX_REPORT_POSITION_ID_LENGTH = 100
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

func (service *ReportService) check(report *Report) error {
	if report.ID == "" {
		return &ReportFieldBlankError{
			field: "id",
		}
	}

	if len(report.ID) > MAX_REPORT_ID_LENGTH {
		return &ReportFieldTooLargeError{
			max:   MAX_REPORT_ID_LENGTH,
			field: "id",
		}
	}

	_, err := service.db.Get(report.ID)
	if err == nil {
		return &ReportAlreadyExistsError{}
	}

	_, ok := err.(*db.NotFoundError)
	if !ok {
		return err
	}

	return nil
}

func (service *ReportService) checkStringField(fieldName string, field string, max int) error {
	if field == "" {
		return &ReportFieldBlankError{
			field: fieldName,
		}
	}

	if len(field) > max {
		return &ReportFieldTooLargeError{
			field: fieldName,
			max:   max,
		}
	}

	return nil
}

func (service *ReportService) Create(report *Report) error {
	err := service.check(report)
	if err != nil {
		return err
	}

	err = service.checkStringField("id", report.ID, MAX_REPORT_ID_LENGTH)
	if err != nil {
		return err
	}

	err = service.checkStringField("title", report.Title, MAX_REPORT_TITLE_LENGTH)
	if err != nil {
		return err
	}

	err = service.checkStringField("userID", report.UserID, MAX_REPORT_USER_ID_LENGTH)
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

	if report.Tag == "" {
		return service.db.Create(report.ID, report, service.geo.AsString(point), report.UserID)
	}

	return service.db.Create(report.ID, report, report.Tag, service.geo.AsString(point), report.UserID)
}

func (service *ReportService) Get(report *Report) (*Report, error) {
	err := service.checkStringField("id", report.ID, MAX_REPORT_ID_LENGTH)
	if err != nil {
		return nil, err
	}

	res, err := service.db.Get(report.ID)
	if err != nil {
		return nil, err
	}

	if res.UserID != report.UserID {
		return nil, &InvalidUserIDError{}
	}

	return res, nil
}

func (service *ReportService) Delete(report *Report) error {
	err := service.checkStringField("id", report.ID, MAX_REPORT_ID_LENGTH)
	if err != nil {
		return err
	}

	res, err := service.db.Get(report.ID)
	if err != nil {
		return err
	}

	if res.UserID != report.UserID {
		return &InvalidUserIDError{}
	}

	return service.db.Delete(report.ID)
}

func (service *ReportService) ListByUser(id string) ([]string, error) {
	err := service.checkStringField("id", id, MAX_REPORT_ID_LENGTH)
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
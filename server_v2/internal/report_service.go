package internal

import "danger-dodgers/pkg/db"

type ReportService struct {
	db db.Database[Report]
}

func NewReportService(db db.Database[Report]) *ReportService {
	return &ReportService{
		db: db,
	}
}

func (service *ReportService) Create(report *Report) error {
	
}
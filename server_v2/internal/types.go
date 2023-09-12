package internal

type Service interface {
	CreateUser(user User, id string) error
	UpdateUser(user User, id string) error
	GetUser(id string) (User, error)
	DeleteUser(id string) error

	CreateReport(report Report, id string) error
	UpdateReport(report Report, id string) error
	GetReport(id string) (Report, error)
	DeleteReport(id string) error
	ListReport() ([]Report, error)
	ListReportByGeography(Area) ([]Report, error)

	CreateRoute(route Route, id string) error
	UpdateRoute(route Route, id string) error
	GetRoute(id string) (Route, error)
	DeleteRoute(id string) error
}

package models

type Doc struct {
	ID                    string
	DateReported          string
	ProjectName           string
	Description           string
	Category              string
	Borough               string
	ManagingAgency        string
	ClientAgency          string
	CurrentPhase          string
	DesignStart           string
	BudgetForecast        string
	LatestBudgetChanges   string
	TotalBudgetChanges    string
	ForecastCompletion    string
	LatestScheduleChanges string
	TotalScheduleChanges  string
	Text                  string
}

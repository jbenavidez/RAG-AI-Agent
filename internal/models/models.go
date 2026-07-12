package models

import "time"

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]any
}

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type UploadedFile struct {
	ID               string
	OriginalFileName string
	StoredFileName   string
	FilePath         string
	Description      string
	ContentType      string
	Status           string
	Size             int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ErrorMessage     string
}

type CapitalProject struct {
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

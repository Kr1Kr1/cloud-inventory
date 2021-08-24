package ovh

type Project struct {
	ID          string `json:"project_id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

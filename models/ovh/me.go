package ovh

type Me struct {
	ID           string `json:"id"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Country      string `json:"country"`
	CustomerCode string `json:"customerCode"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	Name         string `json:"name"`
	Nichandle    string `json:"nichandle"`
	Zip          string `json:"zip"`
	Currency     struct {
		Code string `json:"code"`
	} `json:"currency"`
}

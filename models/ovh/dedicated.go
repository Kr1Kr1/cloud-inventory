package ovh

type DedicatedServer struct {
	ServerID         int         `json:"serverId"`
	Reverse          string      `json:"reverse"`
	NewUpgradeSystem bool        `json:"newUpgradeSystem"`
	Os               string      `json:"os"`
	IP               string      `json:"ip"`
	SupportLevel     string      `json:"supportLevel"`
	RescueMail       string      `json:"rescueMail"`
	Rack             string      `json:"rack"`
	CommercialRange  string      `json:"commercialRange"`
	Datacenter       string      `json:"datacenter"`
	LinkSpeed        int         `json:"linkSpeed"`
	Monitoring       bool        `json:"monitoring"`
	BootID           int         `json:"bootId"`
	ProfessionalUse  bool        `json:"professionalUse"`
	RootDevice       interface{} `json:"rootDevice"`
	NoIntervention   bool        `json:"noIntervention"`
	Name             string      `json:"name"`
	State            string      `json:"state"`
}

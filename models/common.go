package models

type Server struct {
	CloudName  string `json:"cloudName"`
	CloudType  string `json:"cloudType"`
	Datacenter string `json:"datacenter"`
	Name       string `json:"name"`
	State      string `json:"state"`
}

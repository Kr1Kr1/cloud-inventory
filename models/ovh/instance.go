package ovh

import "time"

type Instance struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IPAddresses []struct {
		IP        string `json:"ip"`
		Type      string `json:"type"`
		Version   int    `json:"version"`
		NetworkID string `json:"networkId"`
		GatewayIP string `json:"gatewayIp"`
	} `json:"ipAddresses"`
	FlavorID       string    `json:"flavorId"`
	ImageID        string    `json:"imageId"`
	SSHKeyID       string    `json:"sshKeyId"`
	Created        time.Time `json:"created"`
	Region         string    `json:"region"`
	MonthlyBilling struct {
		Since  time.Time `json:"since"`
		Status string    `json:"status"`
	} `json:"monthlyBilling"`
	Status       string        `json:"status"`
	PlanCode     string        `json:"planCode"`
	OperationIds []interface{} `json:"operationIds"`
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ovh/go-ovh/ovh"
)

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

type Project struct {
	ID          string `json:"project_id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type vRack struct {
	ID string `json:"id"`
}

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

// PrettyPrint Prettify the display
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}

func main() {
	client, _ := ovh.NewDefaultClient()

	// Me
	me := Me{}
	client.Get("/me", &me)
	fmt.Printf("Welcome %s!\n", me.FirstName)

	// Projects
	projects := make([]Project, 0)
	ids := []string{}
	err := client.Get("/cloud/project", &ids)
	if err != nil {
		fmt.Printf("[Error] %s!\n", err)
	}

	for _, id := range ids {
		project := Project{}
		client.Get(fmt.Sprintf("/cloud/project/%s", id), &project)
		projects = append(projects, project)
	}

	// vRacks
	for _, project := range projects {
		vrack := vRack{}
		client.Get(fmt.Sprintf("/cloud/project/%s/vrack", project.ID), &vrack)
		fmt.Printf("Project %s has vRack %s\n", project.Description, vrack.ID)
	}

	// Instances
	listOfInstances := make([]Instance, 0)
	for _, project := range projects {
		instances := make([]Instance, 0)
		// test := make(json.RawMessage, 0)
		client.Get(fmt.Sprintf("/cloud/project/%s/instance", project.ID), &instances)
		// client.Get(fmt.Sprintf("/cloud/project/%s/instance", project.ID), &test)
		// PrettyPrint(test)

		for _, instance := range instances {
			instance.PlanCode = strings.Replace(instance.PlanCode, ".consumption", "", -1)
			listOfInstances = append(listOfInstances, instance)
		}
	}

	// Pretty Print
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	fmt.Fprintln(w, "REGION\tNAME")
	for _, instance := range listOfInstances {
		w.Write([]byte(fmt.Sprintf("%s\t%s\n", instance.Region, instance.Name)))
	}
	w.Flush()
}

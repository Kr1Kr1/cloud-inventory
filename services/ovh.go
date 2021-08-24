package services

import (
	"fmt"
	"strings"

	"cloud-inventory/models"
	"cloud-inventory/models/ovh"

	govh "github.com/ovh/go-ovh/ovh"
)

func OvhListAllServers() ([]models.Server, error) {
	client, _ := govh.NewDefaultClient()

	// Projects
	projects := make([]ovh.Project, 0)
	ids := []string{}
	err := client.Get("/cloud/project", &ids)
	if err != nil {
		fmt.Printf("[Error] %s!\n", err)
	}

	for _, id := range ids {
		project := ovh.Project{}
		client.Get(fmt.Sprintf("/cloud/project/%s", id), &project)
		projects = append(projects, project)
	}

	// Instances
	listOfInstances := make([]ovh.Instance, 0)
	for _, project := range projects {
		instances := make([]ovh.Instance, 0)
		client.Get(fmt.Sprintf("/cloud/project/%s/instance", project.ID), &instances)

		for _, instance := range instances {
			instance.PlanCode = strings.Replace(instance.PlanCode, ".consumption", "", -1)
			listOfInstances = append(listOfInstances, instance)
		}
	}

	// Bare Metal Cloud > Dedicated servers
	dedicatedServices := []string{}
	listOfDedicatedServers := make([]ovh.DedicatedServer, 0)
	client.Get("/dedicated/server", &dedicatedServices)
	for _, dedicatedService := range dedicatedServices {
		dedicatedServer := ovh.DedicatedServer{}
		client.Get(fmt.Sprintf("/dedicated/server/%s", dedicatedService), &dedicatedServer)
		listOfDedicatedServers = append(listOfDedicatedServers, dedicatedServer)
	}

	// Servers
	servers := make([]models.Server, 0)
	for _, instance := range listOfInstances {
		servers = append(servers, models.Server{CloudName: "OVH", CloudType: "Public Cloud", Datacenter: instance.Region, Name: instance.Name, State: instance.Status})
	}
	for _, dedicatedServer := range listOfDedicatedServers {
		servers = append(servers, models.Server{CloudName: "OVH", CloudType: "Dedicated Cloud", Datacenter: dedicatedServer.Datacenter, Name: dedicatedServer.Name, State: dedicatedServer.State})
	}

	return servers, nil
}

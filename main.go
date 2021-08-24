package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"cloud-inventory/api"
	"cloud-inventory/models/ovh"

	govh "github.com/ovh/go-ovh/ovh"
	"github.com/spf13/cobra"
)

// PrettyPrint Prettify the display
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return err
}

func main() {

	var cmdOVH = &cobra.Command{
		Use:   "ovh",
		Short: "OVH inventory",
		Long:  `Listing OVH servers (cloud & bare metal)`,
		// Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// sandbox.Run()
			ovhListing()
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "cloud-inventory",
		Short: "cloud-inventory go api",
		Long:  "cloud-inventory go api",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting api...")
			api.LaunchServer()
		},
	}

	rootCmd.AddCommand(cmdOVH)
	rootCmd.Execute()
}

func ovhListing() {
	client, _ := govh.NewDefaultClient()

	// Me
	me := ovh.Me{}
	client.Get("/me", &me)
	fmt.Printf("Welcome %s!\n", me.FirstName)

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

	// vRacks
	for _, project := range projects {
		vrack := ovh.VRack{}
		client.Get(fmt.Sprintf("/cloud/project/%s/vrack", project.ID), &vrack)
		fmt.Printf("Project %s has vRack %s\n", project.Description, vrack.ID)
	}

	// Instances
	listOfInstances := make([]ovh.Instance, 0)
	for _, project := range projects {
		instances := make([]ovh.Instance, 0)
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
	fmt.Fprintln(w, "CLOUD TYPE\tREGION\tNAME\tSTATE")
	for _, instance := range listOfInstances {
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t%s\n", "Public Cloud", instance.Region, instance.Name, instance.Status)))
	}
	// w.Flush()

	// Bare Metal Cloud > Dedicated servers
	dedicatedServices := []string{}
	listOfDedicatedServers := make([]ovh.DedicatedServer, 0)
	client.Get("/dedicated/server", &dedicatedServices)
	for _, dedicatedService := range dedicatedServices {
		dedicatedServer := ovh.DedicatedServer{}
		client.Get(fmt.Sprintf("/dedicated/server/%s", dedicatedService), &dedicatedServer)
		listOfDedicatedServers = append(listOfDedicatedServers, dedicatedServer)
	}

	// w = tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)
	// fmt.Fprintln(w, "CLOUD TYPE\tREGION\tNAME")
	for _, server := range listOfDedicatedServers {
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t%s\n", "Bare Metal Cloud", strings.ToUpper(server.Datacenter), server.Reverse, strings.ToUpper(server.State))))
	}
	w.Flush()
}

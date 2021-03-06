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

// All prices are per month and HT
var prices = map[string]string{
	"b2-7":      "22",
	"b2-15":     "42",
	"b2-30":     "85",
	"b2-60":     "165",
	"b2-120":    "325",
	"c2-7":      "32",
	"c2-15":     "62",
	"c2-30":     "125",
	"c2-60":     "245",
	"c2-120":    "485",
	"d2-2":      "5",
	"d2-4":      "10",
	"d2-8":      "18",
	"i1-45":     "200",
	"i1-90":     "400",
	"i1-180":    "800",
	"r2-15":     "32",
	"r2-30":     "37",
	"r2-60":     "72",
	"r2-120":    "145",
	"r2-240":    "285",
	"t1-45":     "799",
	"t1-90":     "1599",
	"t1-180":    "3199",
	"t2-45":     "825",
	"t2-90":     "1649",
	"t2-180":    "3299",
	"advance-5": "190",
	"fs-72t":    "390",
	"infra-1":   "95",
	"infra-4":   "219",
	"infra-3":   "209",
	"hgr-sds-2": "576",
	"fs-48t":    "314",
	"advance-3": "81",
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
	fmt.Fprintln(w, "CLOUD TYPE\tREGION\tNAME\tIP\tSTATE\tBILLING (Month, HT)")
	for _, instance := range listOfInstances {
		// Get monthly price
		var price string
		billingName := strings.Split(instance.PlanCode, ".")
		if val, ok := prices[strings.ToLower(billingName[0])]; ok {
			price = val
		} else {
			price = instance.PlanCode
		}
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n", "Public Cloud", instance.Region, instance.Name, instance.IPAddresses[0].IP, instance.Status, price)))
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
		// Get monthly price
		var price string
		if val, ok := prices[strings.ToLower(server.CommercialRange)]; ok {
			price = val
		} else {
			price = server.CommercialRange
		}
		w.Write([]byte(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n", "Bare Metal Cloud", strings.ToUpper(server.Datacenter), server.Reverse, server.IP, strings.ToUpper(server.State), price)))
	}
	w.Flush()
}

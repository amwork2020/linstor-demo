package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/LINBIT/golinstor/client"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.TODO()

	u, err := url.Parse("http://10.1.5.11:3370")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u)

	c, err := client.NewClient(client.BaseURL(u), client.Log(log.StandardLogger()))
	if err != nil {
		log.Fatal(err)
	}
	version, err := c.Controller.GetVersion(ctx)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(version)
	
	nodes, err := c.Nodes.GetAll(ctx)
	if err != nil {
		panic(err.Error())
	}

	// println(nodes)
	for i, node := range nodes {
		fmt.Printf("%v\t %v\t %v\t %v\t %v\n", i+1, node.Name, node.Type, node.ConnectionStatus, node.NetInterfaces)
	}
	fmt.Println("------------------------------------------------------")

	rsall, err := c.Resources.GetResourceView(ctx)
	if err != nil {
		panic(err.Error())
	}

	for _, r := range rsall {
		kind := ""
		for _, volume := range r.Volumes {
			// kind = string(volume.ProviderKind)
			kind = volume.State.DiskState
		}
		fmt.Printf("%v\t %v\t %v\t %v\n", r.Name, r.NodeName, *r.State.InUse, kind)
	}
	fmt.Println("------------------------------------------------------")

	rs, err := c.Resources.GetAll(ctx, "linstor_db")
	if errs, ok := err.(client.ApiCallError); ok {
		log.Error("A LINSTOR API error occurred:")
		for i, e := range errs {
			log.Errorf("  Message #%d:", i)
			log.Errorf("    Code: %d", e.RetCode)
			log.Errorf("    Message: %s", e.Message)
			log.Errorf("    Cause: %s", e.Cause)
			log.Errorf("    Details: %s", e.Details)
			log.Errorf("    Correction: %s", e.Correction)
			log.Errorf("    Error Reports: %v", e.ErrorReportIds)
		}
		return
	}
	if err != nil {
		log.Fatalf("Some other error occurred: %s", err.Error())
	}

	for _, r := range rs {
		fmt.Printf("Resource with name '%s' on node with name '%s'\n", r.Name, r.NodeName)
	}
}

package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/casbin/casbin"
	etcdadapter "github.com/ntk148v/etcd-adapter"
)

func main() {
	// Initialize a casbin etcd adapter and use it in a Casbin enforcer:
	// The adapter will use the ETCD and a named path with the key you give.
	// If not provided, the adapter will use the default value casbin_policy.
	a := etcdadapter.NewAdapter([]string{"http://127.0.0.1:2379"}, "casbin_policy_test") // Your etcd endpoints and the path key.

	e := casbin.NewEnforcer("rbac_model.conf", a)
	log.Println("connected to Etcd and init an Enforcer")

	// Load the policy from ETCD.
	if err := e.LoadPolicy(); err != nil {
		log.Println(err)
	}

	// Check the permission.
	if ok := e.Enforce("alice", "data1", "read"); !ok {
		log.Println("alice doesn't permission to read data1")
	} else {
		log.Println("alice is allowed to read data1")
	}

	// Modify the policy.
	csvFile, err := os.Open("rbac_policy.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Println(err)
	}
	for _, line := range csvLines {
		if len(line) != 4 {
			continue
		}
		if ok := e.AddPolicy(line[1], line[2], line[3]); ok {
			log.Println("added new rule to Etcd")
		}
	}

	// e.RemovePolicy(...)

	// Save the policy back to DB.
	if err := e.SavePolicy(); err != nil {
		log.Fatal(err)
	}
}

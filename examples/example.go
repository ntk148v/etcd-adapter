package main

import (
	"log"
	
	"github.com/casbin/casbin/v2"
	etcdadapter "github.com/ntk148v/etcd-adapter"
)

func main() {
	// Initialize a casbin etcd adapter and use it in a Casbin enforcer:
	// The adapter will use the ETCD and a named path with the key you give.
	// If not provided, the adapter will use the default value casbin_policy.
	a := etcdadapter.NewAdapter([]string{"http://127.0.0.1:6379"}, "casbin_policy_test") // Your etcd endpoints and the path key.

	e, err := casbin.NewEnforcer("rbac_model.conf", a)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to Etcd and init an Enforcer")

	// Load the policy from ETCD.
	if err := e.LoadPolicy(); err != nil {
		log.Println(err)
	}

	// Modify the policy.
	rules := [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"data2_admin", "data2", "read"},
		{"data2_admin", "data2", "write"},
	}
	if ok, _ := e.AddPolicies(rules); ok {
		log.Println("added new policies")
	}

	// Check the permission.
	if ok, err := e.Enforce("alice", "data1", "read"); !ok || err != nil {
		log.Println("alice doesn't permission to read data1")
	} else {
		log.Println("alice is allowed to read data1")
	}

	// e.RemovePolicy(...)

	// Save the policy back to DB.
	if err := e.SavePolicy(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"

	xctl ".."
)

func main() {

	c := xctl.NewServiceClient("127.0.0.1", 10085)
	fmt.Println("--------------------------------------------")
	ret := c.QueryStats("rand", false)
	for key, val := range ret {
		fmt.Printf("%s -> %d\n", key, val)
	}
	// fmt.Println("--------------------------------------------")
	// c.GetStats("user>>>test@test.com>>>traffic>>>downlink", false)
	// c.GetStats("user>>>test@test.com>>>traffic>>>uplink", false)
	// fmt.Println("--------------------------------------------")

	// email := fmt.Sprintf("rand%d@test.com", rand.Int31())
	// uuid := v2ctl.GenerateUUID()

	// fmt.Printf("generate user: %s ... %s\n", uuid, email)
	// c.AddUser("vmess", email, 0, uuid, 64)

	// c.RemoveUser("vmess", "test@test.com")
	fmt.Println("--------------------------------------------")
}

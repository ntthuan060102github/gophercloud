package routers

import "github.com/vnpaycloud-console/gophercloud/v2"

const resourcePath = "routers"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func addInterfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "add_router_interface")
}

func removeInterfaceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "remove_router_interface")
}

func listl3AgentsURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "l3-agents")
}

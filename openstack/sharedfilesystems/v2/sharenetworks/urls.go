package sharenetworks

import "github.com/vnpaycloud-console/gophercloud/v2"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-networks")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-networks", id)
}

func listDetailURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("share-networks", "detail")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return deleteURL(c, id)
}

func updateURL(c *gophercloud.ServiceClient, id string) string {
	return deleteURL(c, id)
}

func addSecurityServiceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-networks", id, "action")
}

func removeSecurityServiceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("share-networks", id, "action")
}

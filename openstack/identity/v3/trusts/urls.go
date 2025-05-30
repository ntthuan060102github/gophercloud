package trusts

import "github.com/vnpaycloud-console/gophercloud/v2"

const resourcePath = "OS-TRUST/trusts"

func rootURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func resourceURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id)
}

func createURL(c *gophercloud.ServiceClient) string {
	return rootURL(c)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return resourceURL(c, id)
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func listRolesURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL(resourcePath, id, "roles")
}

func getRoleURL(c *gophercloud.ServiceClient, id, roleID string) string {
	return c.ServiceURL(resourcePath, id, "roles", roleID)
}

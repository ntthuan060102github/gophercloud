package availabilityzones

import "github.com/vnpaycloud-console/gophercloud/v2"

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("os-availability-zone")
}

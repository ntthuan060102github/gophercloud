package RESOURCE

import "github.com/vnpaycloud-console/gophercloud/v2"

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("resource")
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("resource", id)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("resource")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("resource", id)
}

func updateURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("resource", id)
}

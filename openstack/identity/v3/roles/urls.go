package roles

import "github.com/vnpaycloud-console/gophercloud/v2"

const (
	rolePath = "roles"
)

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func getURL(client *gophercloud.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL(rolePath)
}

func updateURL(client *gophercloud.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func deleteURL(client *gophercloud.ServiceClient, roleID string) string {
	return client.ServiceURL(rolePath, roleID)
}

func listAssignmentsURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("role_assignments")
}

func listAssignmentsOnResourceURL(client *gophercloud.ServiceClient, targetType, targetID, actorType, actorID string) string {
	return client.ServiceURL(targetType, targetID, actorType, actorID, rolePath)
}

func assignURL(client *gophercloud.ServiceClient, targetType, targetID, actorType, actorID, roleID string) string {
	return client.ServiceURL(targetType, targetID, actorType, actorID, rolePath, roleID)
}

func createRoleInferenceRuleURL(client *gophercloud.ServiceClient, priorRoleID, impliedRoleID string) string {
	return client.ServiceURL(rolePath, priorRoleID, "implies", impliedRoleID)
}

func getRoleInferenceRuleURL(client *gophercloud.ServiceClient, priorRoleID, impliedRoleID string) string {
	return client.ServiceURL(rolePath, priorRoleID, "implies", impliedRoleID)
}

func listRoleInferenceRulesURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("role_inferences")
}

func deleteRoleInferenceRuleURL(client *gophercloud.ServiceClient, priorRoleID, impliedRoleID string) string {
	return client.ServiceURL(rolePath, priorRoleID, "implies", impliedRoleID)
}

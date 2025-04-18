package v2

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/sharedfilesystems/v2/securityservices"
)

// CreateSecurityService will create a security service with a random name. An
// error will be returned if the security service was unable to be created.
func CreateSecurityService(t *testing.T, client *gophercloud.ServiceClient) (*securityservices.SecurityService, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires share network creation in short mode.")
	}

	securityServiceName := tools.RandomString("ACPTTEST", 16)
	securityServiceDescription := tools.RandomString("ACPTTEST-DESC", 16)
	t.Logf("Attempting to create security service: %s", securityServiceName)

	createOpts := securityservices.CreateOpts{
		Name:        securityServiceName,
		Description: securityServiceDescription,
		Type:        "kerberos",
	}

	securityService, err := securityservices.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return securityService, err
	}

	return securityService, nil
}

// DeleteSecurityService will delete a security service. An error will occur if
// the security service was unable to be deleted.
func DeleteSecurityService(t *testing.T, client *gophercloud.ServiceClient, securityService *securityservices.SecurityService) {
	err := securityservices.Delete(context.TODO(), client, securityService.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Failed to delete security service %s: %v", securityService.ID, err)
	}

	t.Logf("Deleted security service: %s", securityService.ID)
}

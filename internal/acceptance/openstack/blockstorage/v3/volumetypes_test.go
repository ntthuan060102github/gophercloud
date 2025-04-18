//go:build acceptance || blockstorage || volumetypes

package v3

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	identity "github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/openstack/identity/v3"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestVolumeTypes(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	vt, err := CreateVolumeType(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, vt)

	allPages, err := volumetypes.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allVolumeTypes, err := volumetypes.ExtractVolumeTypes(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allVolumeTypes {
		tools.PrintResource(t, v)
		if v.ID == vt.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	isPublic := false
	name := vt.Name + "-UPDATED"
	description := ""
	updateOpts := volumetypes.UpdateOpts{
		Name:        &name,
		Description: &description,
		IsPublic:    &isPublic,
	}

	newVT, err := volumetypes.Update(context.TODO(), client, vt.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newVT)
	th.AssertEquals(t, name, newVT.Name)
	th.AssertEquals(t, description, newVT.Description)
	th.AssertEquals(t, isPublic, newVT.IsPublic)
}

func TestVolumeTypesExtraSpecs(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	vt, err := CreateVolumeTypeNoExtraSpecs(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, vt)

	createOpts := volumetypes.ExtraSpecsOpts{
		"capabilities":        "gpu",
		"volume_backend_name": "ssd",
	}

	createdExtraSpecs, err := volumetypes.CreateExtraSpecs(context.TODO(), client, vt.ID, createOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, createdExtraSpecs)

	th.AssertEquals(t, len(createdExtraSpecs), 2)
	th.AssertEquals(t, createdExtraSpecs["capabilities"], "gpu")
	th.AssertEquals(t, createdExtraSpecs["volume_backend_name"], "ssd")

	err = volumetypes.DeleteExtraSpec(context.TODO(), client, vt.ID, "volume_backend_name").ExtractErr()
	th.AssertNoErr(t, err)

	updateOpts := volumetypes.ExtraSpecsOpts{
		"capabilities": "gpu-2",
	}
	updatedExtraSpec, err := volumetypes.UpdateExtraSpec(context.TODO(), client, vt.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedExtraSpec)

	th.AssertEquals(t, updatedExtraSpec["capabilities"], "gpu-2")

	allExtraSpecs, err := volumetypes.ListExtraSpecs(context.TODO(), client, vt.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, allExtraSpecs)

	th.AssertEquals(t, len(allExtraSpecs), 1)
	th.AssertEquals(t, allExtraSpecs["capabilities"], "gpu-2")

	singleSpec, err := volumetypes.GetExtraSpec(context.TODO(), client, vt.ID, "capabilities").Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, singleSpec)

	th.AssertEquals(t, singleSpec["capabilities"], "gpu-2")
}

func TestVolumeTypesAccess(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	vt, err := CreatePrivateVolumeType(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, vt)

	project, err := identity.CreateProject(t, identityClient, nil)
	th.AssertNoErr(t, err)
	defer identity.DeleteProject(t, identityClient, project.ID)

	addAccessOpts := volumetypes.AddAccessOpts{
		Project: project.ID,
	}

	err = volumetypes.AddAccess(context.TODO(), client, vt.ID, addAccessOpts).ExtractErr()
	th.AssertNoErr(t, err)

	allPages, err := volumetypes.ListAccesses(client, vt.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	accessList, err := volumetypes.ExtractAccesses(allPages)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, accessList)

	th.AssertEquals(t, len(accessList), 1)
	th.AssertEquals(t, accessList[0].ProjectID, project.ID)
	th.AssertEquals(t, accessList[0].VolumeTypeID, vt.ID)

	removeAccessOpts := volumetypes.RemoveAccessOpts{
		Project: project.ID,
	}

	err = volumetypes.RemoveAccess(context.TODO(), client, vt.ID, removeAccessOpts).ExtractErr()
	th.AssertNoErr(t, err)

	allPages, err = volumetypes.ListAccesses(client, vt.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	accessList, err = volumetypes.ExtractAccesses(allPages)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, accessList)

	th.AssertEquals(t, len(accessList), 0)
}

func TestEncryptionVolumeTypes(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	vt, err := CreateVolumeType(t, client)
	th.AssertNoErr(t, err)
	defer DeleteVolumeType(t, client, vt)

	createEncryptionOpts := volumetypes.CreateEncryptionOpts{
		KeySize:         256,
		Provider:        "luks",
		ControlLocation: "front-end",
		Cipher:          "aes-xts-plain64",
	}

	eVT, err := volumetypes.CreateEncryption(context.TODO(), client, vt.ID, createEncryptionOpts).Extract()
	th.AssertNoErr(t, err)
	defer volumetypes.DeleteEncryption(context.TODO(), client, eVT.VolumeTypeID, eVT.EncryptionID)

	geVT, err := volumetypes.GetEncryption(context.TODO(), client, vt.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, geVT)

	key := "cipher"
	gesVT, err := volumetypes.GetEncryptionSpec(context.TODO(), client, vt.ID, key).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, gesVT)

	updateEncryptionOpts := volumetypes.UpdateEncryptionOpts{
		ControlLocation: "back-end",
	}

	newEVT, err := volumetypes.UpdateEncryption(context.TODO(), client, vt.ID, eVT.EncryptionID, updateEncryptionOpts).Extract()
	tools.PrintResource(t, newEVT)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "back-end", newEVT.ControlLocation)
}

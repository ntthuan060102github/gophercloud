//go:build acceptance || compute || usage

package v2

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/compute/v2/usage"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestUsageSingleTenant(t *testing.T) {
	// TODO(emilien): This test is failing for now
	t.Skip("This is not passing now, will fix later")

	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	DeleteServer(t, client, server)

	endpointParts := strings.Split(client.Endpoint, "/")
	tenantID := endpointParts[4]

	end := time.Now()
	start := end.AddDate(0, -1, 0)
	opts := usage.SingleTenantOpts{
		Start: &start,
		End:   &end,
	}

	err = usage.SingleTenant(client, tenantID, opts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		tenantUsage, err := usage.ExtractSingleTenant(page)
		th.AssertNoErr(t, err)

		tools.PrintResource(t, tenantUsage)
		if tenantUsage.TotalHours == 0 {
			t.Fatalf("TotalHours should not be 0")
		}

		return true, nil
	})

	th.AssertNoErr(t, err)
}

func TestUsageAllTenants(t *testing.T) {
	t.Skip("This is not passing in OpenLab. Works locally")

	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	DeleteServer(t, client, server)

	end := time.Now()
	start := end.AddDate(0, -1, 0)
	opts := usage.AllTenantsOpts{
		Detailed: true,
		Start:    &start,
		End:      &end,
	}

	err = usage.AllTenants(client, opts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		allUsage, err := usage.ExtractAllTenants(page)
		th.AssertNoErr(t, err)

		tools.PrintResource(t, allUsage)

		if len(allUsage) == 0 {
			t.Fatalf("No usage returned")
		}

		if allUsage[0].TotalHours == 0 {
			t.Fatalf("TotalHours should not be 0")
		}
		return true, nil
	})

	th.AssertNoErr(t, err)
}

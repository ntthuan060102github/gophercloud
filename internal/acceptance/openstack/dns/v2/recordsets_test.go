//go:build acceptance || dns || recordsets

package v2

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/dns/v2/recordsets"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestRecordSetsListByZone(t *testing.T) {
	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	zone, err := CreateZone(t, client)
	th.AssertNoErr(t, err)
	defer DeleteZone(t, client, zone)

	allPages, err := recordsets.ListByZone(client, zone.ID, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allRecordSets, err := recordsets.ExtractRecordSets(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, recordset := range allRecordSets {
		tools.PrintResource(t, &recordset)

		if recordset.ZoneID == zone.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	listOpts := recordsets.ListOpts{
		Limit: 1,
	}

	pager := recordsets.ListByZone(client, zone.ID, listOpts)
	err = pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		rr, err := recordsets.ExtractRecordSets(page)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, len(rr), 1)
		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestRecordSetsCRUD(t *testing.T) {
	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	zone, err := CreateZone(t, client)
	th.AssertNoErr(t, err)
	defer DeleteZone(t, client, zone)

	tools.PrintResource(t, &zone)

	rs, err := CreateRecordSet(t, client, zone)
	th.AssertNoErr(t, err)
	defer DeleteRecordSet(t, client, rs)

	tools.PrintResource(t, &rs)

	description := ""
	updateOpts := recordsets.UpdateOpts{
		Description: &description,
	}

	newRS, err := recordsets.Update(context.TODO(), client, rs.ZoneID, rs.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, &newRS)

	th.AssertEquals(t, newRS.Description, description)

	records := []string{"10.1.0.3"}
	updateOpts = recordsets.UpdateOpts{
		Records: records,
	}

	newRS, err = recordsets.Update(context.TODO(), client, rs.ZoneID, rs.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, &newRS)

	th.AssertDeepEquals(t, newRS.Records, records)
	th.AssertEquals(t, newRS.TTL, 3600)

	ttl := 0
	updateOpts = recordsets.UpdateOpts{
		TTL: &ttl,
	}

	newRS, err = recordsets.Update(context.TODO(), client, rs.ZoneID, rs.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, &newRS)

	th.AssertDeepEquals(t, newRS.Records, records)
	th.AssertEquals(t, newRS.TTL, ttl)
}

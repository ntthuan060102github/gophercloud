package testing

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/loadbalancer/v2/amphorae"
	fake "github.com/vnpaycloud-console/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestListAmphorae(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAmphoraListSuccessfully(t)

	pages := 0
	err := amphorae.List(fake.ServiceClient(), amphorae.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := amphorae.ExtractAmphorae(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 amphorae, got %d", len(actual))
		}

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllAmphorae(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAmphoraListSuccessfully(t)

	allPages, err := amphorae.List(fake.ServiceClient(), amphorae.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := amphorae.ExtractAmphorae(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(actual))
	th.AssertDeepEquals(t, ExpectedAmphoraeSlice, actual)
}

func TestGetAmphora(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAmphoraGetSuccessfully(t)

	client := fake.ServiceClient()
	actual, err := amphorae.Get(context.TODO(), client, "45f40289-0551-483a-b089-47214bc2a8a4").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, FirstAmphora, *actual)
}

func TestFailoverAmphora(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleAmphoraFailoverSuccessfully(t)

	res := amphorae.Failover(context.TODO(), fake.ServiceClient(), "36e08a3e-a78f-4b40-a229-1e7e23eee1ab")
	th.AssertNoErr(t, res.Err)
}

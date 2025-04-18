//go:build acceptance || workflow || crontriggers

package v2

import (
	"testing"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/workflow/v2/crontriggers"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestCronTriggersCreateGetDelete(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)

	trigger, err := CreateCronTrigger(t, client, workflow)
	th.AssertNoErr(t, err)
	defer DeleteCronTrigger(t, client, trigger)

	gettrigger, err := GetCronTrigger(t, client, trigger.ID)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, trigger.ID, gettrigger.ID)

	tools.PrintResource(t, trigger)
}

func TestCronTriggersList(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)
	workflow, err := CreateWorkflow(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkflow(t, client, workflow)
	trigger, err := CreateCronTrigger(t, client, workflow)
	th.AssertNoErr(t, err)
	defer DeleteCronTrigger(t, client, trigger)
	list, err := ListCronTriggers(t, client, &crontriggers.ListOpts{
		Name: &crontriggers.ListFilter{
			Filter: crontriggers.FilterEQ,
			Value:  trigger.Name,
		},
		Pattern: &crontriggers.ListFilter{
			Value: "0 0 1 1 *",
		},
		CreatedAt: &crontriggers.ListDateFilter{
			Filter: crontriggers.FilterGT,
			Value:  time.Now().AddDate(-1, 0, 0),
		},
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(list))
	tools.PrintResource(t, list)
}

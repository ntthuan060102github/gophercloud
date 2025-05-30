package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2"
	fake "github.com/vnpaycloud-console/gophercloud/v2/openstack/networking/v2/common"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/rules"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "firewall_rules": [
        {
            "protocol": "tcp",
            "description": "ssh rule",
            "source_port": null,
            "source_ip_address": null,
            "destination_ip_address": "192.168.1.0/24",
            "firewall_policy_id": ["e2a5fb51-698c-4898-87e8-f1eee6b50919"],
            "destination_port": "22",
            "id": "f03bd950-6c56-4f5e-a307-45967078f507",
            "name": "ssh_form_any",
            "tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
			"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
            "enabled": true,
            "action": "allow",
            "ip_version": 4,
            "shared": false
        },
        {
            "protocol": "udp",
            "description": "udp rule",
            "source_port": null,
            "source_ip_address": null,
            "destination_ip_address": null,
            "firewall_policy_id": ["98d7fb51-698c-4123-87e8-f1eee6b5ab7e"],
            "destination_port": null,
            "id": "ab7bd950-6c56-4f5e-a307-45967078f890",
            "name": "deny_all_udp",
            "tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
			"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
            "enabled": true,
            "action": "deny",
            "ip_version": 4,
            "shared": false
        }
    ]
}
        `)
	})

	count := 0

	err := rules.List(fake.ServiceClient(), rules.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := rules.ExtractRules(page)
		if err != nil {
			t.Errorf("Failed to extract members: %v", err)
			return false, err
		}

		expected := []rules.Rule{
			{
				Protocol:             string(rules.ProtocolTCP),
				Description:          "ssh rule",
				SourcePort:           "",
				SourceIPAddress:      "",
				DestinationIPAddress: "192.168.1.0/24",
				FirewallPolicyID:     []string{"e2a5fb51-698c-4898-87e8-f1eee6b50919"},
				DestinationPort:      "22",
				ID:                   "f03bd950-6c56-4f5e-a307-45967078f507",
				Name:                 "ssh_form_any",
				TenantID:             "80cf934d6ffb4ef5b244f1c512ad1e61",
				ProjectID:            "80cf934d6ffb4ef5b244f1c512ad1e61",
				Enabled:              true,
				Action:               string(rules.ActionAllow),
				IPVersion:            4,
				Shared:               false,
			},
			{
				Protocol:             "udp",
				Description:          "udp rule",
				SourcePort:           "",
				SourceIPAddress:      "",
				DestinationIPAddress: "",
				FirewallPolicyID:     []string{"98d7fb51-698c-4123-87e8-f1eee6b5ab7e"},
				DestinationPort:      "",
				ID:                   "ab7bd950-6c56-4f5e-a307-45967078f890",
				Name:                 "deny_all_udp",
				TenantID:             "80cf934d6ffb4ef5b244f1c512ad1e61",
				ProjectID:            "80cf934d6ffb4ef5b244f1c512ad1e61",
				Enabled:              true,
				Action:               "deny",
				IPVersion:            4,
				Shared:               false,
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"firewall_rule": {
		"protocol": "tcp",
		"description": "ssh rule",
		"destination_ip_address": "192.168.1.0/24",
		"destination_port": "22",
		"name": "ssh_form_any",
		"action": "allow",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61"
	}
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
	"firewall_rule":{
		"protocol": "tcp",
		"description": "ssh rule",
		"source_port": null,
		"source_ip_address": null,
		"destination_ip_address": "192.168.1.0/24",
		"firewall_policy_id": ["e2a5fb51-698c-4898-87e8-f1eee6b50919"],
		"position": 2,
		"destination_port": "22",
		"id": "f03bd950-6c56-4f5e-a307-45967078f507",
		"name": "ssh_form_any",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"enabled": true,
		"action": "allow",
		"ip_version": 4,
		"shared": false
	}
}
        `)
	})

	options := rules.CreateOpts{
		TenantID:             "80cf934d6ffb4ef5b244f1c512ad1e61",
		ProjectID:            "80cf934d6ffb4ef5b244f1c512ad1e61",
		Protocol:             rules.ProtocolTCP,
		Description:          "ssh rule",
		DestinationIPAddress: "192.168.1.0/24",
		DestinationPort:      "22",
		Name:                 "ssh_form_any",
		Action:               "allow",
	}

	_, err := rules.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestCreateAnyProtocol(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"firewall_rule": {
		"protocol": null,
		"description": "any to 192.168.1.0/24",
		"destination_ip_address": "192.168.1.0/24",
		"name": "any_to_192.168.1.0/24",
		"action": "allow",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61"
	}
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `
{
	"firewall_rule":{
		"protocol": null,
		"description": "any to 192.168.1.0/24",
		"source_port": null,
		"source_ip_address": null,
		"destination_ip_address": "192.168.1.0/24",
		"firewall_policy_id": ["e2a5fb51-698c-4898-87e8-f1eee6b50919"],
		"position": 2,
		"destination_port": null,
		"id": "f03bd950-6c56-4f5e-a307-45967078f507",
		"name": "any_to_192.168.1.0/24",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"enabled": true,
		"action": "allow",
		"ip_version": 4,
		"shared": false
	}
}
        `)
	})

	options := rules.CreateOpts{
		TenantID:             "80cf934d6ffb4ef5b244f1c512ad1e61",
		ProjectID:            "80cf934d6ffb4ef5b244f1c512ad1e61",
		Protocol:             rules.ProtocolAny,
		Description:          "any to 192.168.1.0/24",
		DestinationIPAddress: "192.168.1.0/24",
		Name:                 "any_to_192.168.1.0/24",
		Action:               "allow",
	}

	_, err := rules.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules/f03bd950-6c56-4f5e-a307-45967078f507", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
	"firewall_rule":{
		"protocol": "tcp",
		"description": "ssh rule",
		"source_port": null,
		"source_ip_address": null,
		"destination_ip_address": "192.168.1.0/24",
		"firewall_policy_id": ["e2a5fb51-698c-4898-87e8-f1eee6b50919"],
		"position": 2,
		"destination_port": "22",
		"id": "f03bd950-6c56-4f5e-a307-45967078f507",
		"name": "ssh_form_any",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"enabled": true,
		"action": "allow",
		"ip_version": 4,
		"shared": false
	}
}
        `)
	})

	rule, err := rules.Get(context.TODO(), fake.ServiceClient(), "f03bd950-6c56-4f5e-a307-45967078f507").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "tcp", rule.Protocol)
	th.AssertEquals(t, "ssh rule", rule.Description)
	th.AssertEquals(t, "192.168.1.0/24", rule.DestinationIPAddress)
	th.AssertEquals(t, 1, len(rule.FirewallPolicyID))
	th.AssertEquals(t, "e2a5fb51-698c-4898-87e8-f1eee6b50919", rule.FirewallPolicyID[0])
	th.AssertEquals(t, "22", rule.DestinationPort)
	th.AssertEquals(t, "f03bd950-6c56-4f5e-a307-45967078f507", rule.ID)
	th.AssertEquals(t, "ssh_form_any", rule.Name)
	th.AssertEquals(t, "80cf934d6ffb4ef5b244f1c512ad1e61", rule.TenantID)
	th.AssertEquals(t, "80cf934d6ffb4ef5b244f1c512ad1e61", rule.ProjectID)
	th.AssertEquals(t, true, rule.Enabled)
	th.AssertEquals(t, "allow", rule.Action)
	th.AssertEquals(t, 4, rule.IPVersion)
	th.AssertEquals(t, false, rule.Shared)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules/f03bd950-6c56-4f5e-a307-45967078f507", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
	"firewall_rule":{
		"protocol": "tcp",
		"description": "ssh rule",
		"destination_ip_address": "192.168.1.0/24",
		"destination_port": "22",
		"name": "ssh_form_any",
		"action": "allow",
		"enabled": false
	}
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
	"firewall_rule":{
		"protocol": "tcp",
		"description": "ssh rule",
		"destination_ip_address": "192.168.1.0/24",
		"firewall_policy_id": ["e2a5fb51-698c-4898-87e8-f1eee6b50919"],
		"position": 2,
		"destination_port": "22",
		"id": "f03bd950-6c56-4f5e-a307-45967078f507",
		"name": "ssh_form_any",
		"tenant_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"project_id": "80cf934d6ffb4ef5b244f1c512ad1e61",
		"enabled": false,
		"action": "allow",
		"ip_version": 4,
		"shared": false
	}
}
		`)
	})

	newProtocol := rules.ProtocolTCP
	newDescription := "ssh rule"
	newDestinationIP := "192.168.1.0/24"
	newDestintionPort := "22"
	newName := "ssh_form_any"
	newAction := rules.ActionAllow

	options := rules.UpdateOpts{
		Protocol:             &newProtocol,
		Description:          &newDescription,
		DestinationIPAddress: &newDestinationIP,
		DestinationPort:      &newDestintionPort,
		Name:                 &newName,
		Action:               &newAction,
		Enabled:              gophercloud.Disabled,
	}

	_, err := rules.Update(context.TODO(), fake.ServiceClient(), "f03bd950-6c56-4f5e-a307-45967078f507", options).Extract()
	th.AssertNoErr(t, err)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_rules/4ec89077-d057-4a2b-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := rules.Delete(context.TODO(), fake.ServiceClient(), "4ec89077-d057-4a2b-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}

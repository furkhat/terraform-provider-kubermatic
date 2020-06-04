package kubermatic

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kubermatic/go-kubermatic/models"
)

func TestFlattenClusterSpec(t *testing.T) {
	cases := []struct {
		Input          *models.ClusterSpec
		ExpectedOutput []interface{}
	}{
		{
			&models.ClusterSpec{
				Version:         "1.15.6",
				MachineNetworks: nil,
				AuditLogging:    &models.AuditLoggingSettings{},
				Cloud: &models.CloudSpec{
					DatacenterName: "eu-west-1",
					Bringyourown:   map[string]interface{}{},
				},
			},
			[]interface{}{
				map[string]interface{}{
					"version":             "1.15.6",
					"audit_logging":       false,
					"pod_security_policy": false,
					"cloud": []interface{}{
						map[string]interface{}{
							"bringyourown": []interface{}{map[string]interface{}{}},
						},
					},
				},
			},
		},
		{
			&models.ClusterSpec{},
			[]interface{}{
				map[string]interface{}{
					"audit_logging":       false,
					"pod_security_policy": false,
				},
			},
		},
		{
			nil,
			[]interface{}{},
		},
	}

	for _, tc := range cases {
		output := flattenClusterSpec(clusterPreserveValues{}, tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestFlattenClusterCloudSpec(t *testing.T) {
	cases := []struct {
		Input          *models.CloudSpec
		ExpectedOutput []interface{}
	}{
		{
			&models.CloudSpec{
				Aws: &models.AWSCloudSpec{},
			},
			[]interface{}{
				map[string]interface{}{
					"aws": []interface{}{
						map[string]interface{}{},
					},
				},
			},
		},
		{
			&models.CloudSpec{
				Bringyourown: map[string]interface{}{},
			},
			[]interface{}{
				map[string]interface{}{
					"bringyourown": []interface{}{
						map[string]interface{}{},
					},
				},
			},
		},
		{
			&models.CloudSpec{},
			[]interface{}{
				map[string]interface{}{},
			},
		},
		{
			nil,
			[]interface{}{},
		},
	}

	for _, tc := range cases {
		output := flattenClusterCloudSpec(clusterPreserveValues{}, tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestFlattenAWSCloudSpec(t *testing.T) {
	cases := []struct {
		Input          *models.AWSCloudSpec
		ExpectedOutput []interface{}
	}{
		{
			&models.AWSCloudSpec{
				AccessKeyID:         "AKIAIOSFODNN7EXAMPLE",
				ControlPlaneRoleARN: "default",
				InstanceProfileName: "default",
				RouteTableID:        "rtb-09ba434c1bEXAMPLE",
				SecretAccessKey:     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				SecurityGroupID:     "sg-51530134",
				VPCID:               "e5e4b2ef2fe",
			},
			[]interface{}{
				map[string]interface{}{
					"access_key_id":         "AKIAIOSFODNN7EXAMPLE",
					"role_arn":              "default",
					"instance_profile_name": "default",
					"route_table_id":        "rtb-09ba434c1bEXAMPLE",
					"secret_access_key":     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"security_group_id":     "sg-51530134",
					"vpc_id":                "e5e4b2ef2fe",
				},
			},
		},
		{
			&models.AWSCloudSpec{},
			[]interface{}{
				map[string]interface{}{},
			},
		},
		{
			nil,
			[]interface{}{},
		},
	}

	for _, tc := range cases {
		output := flattenAWSCloudSpec(tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestFlattenOpenstackCloudSpec(t *testing.T) {
	cases := []struct {
		Input          *models.OpenstackCloudSpec
		PreserveValues clusterPreserveValues
		ExpectedOutput []interface{}
	}{
		{
			&models.OpenstackCloudSpec{
				FloatingIPPool: "FloatingIPPool",
				Network:        "Network",
				Password:       "",
				RouterID:       "RouterID",
				SecurityGroups: "SecurityGroups",
				SubnetID:       "SubnetID",
				Tenant:         "",
				TenantID:       "TenantID",
				Username:       "",
			},
			clusterPreserveValues{
				openstackUsername: "Username",
				openstackPassword: "Password",
				openstackTenant:   "Tenant",
			},
			[]interface{}{
				map[string]interface{}{
					"username":         "Username",
					"password":         "Password",
					"tenant":           "Tenant",
					"floating_ip_pool": "FloatingIPPool",
				},
			},
		},
		{
			&models.OpenstackCloudSpec{},
			clusterPreserveValues{},
			[]interface{}{
				map[string]interface{}{},
			},
		},
		{
			nil,
			clusterPreserveValues{},
			[]interface{}{},
		},
	}

	for _, tc := range cases {
		output := flattenOpenstackSpec(tc.PreserveValues, tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestFlattenMachineNetwork(t *testing.T) {
	cases := []struct {
		Input          []*models.MachineNetworkingConfig
		ExpectedOutput []interface{}
	}{
		{
			[]*models.MachineNetworkingConfig{
				{
					CIDR:    "192.168.0.0/24",
					Gateway: "192.168.1.1",
					DNSServers: []string{
						"192.200.200.1",
						"192.200.200.201",
					},
				},
			},
			[]interface{}{
				map[string]interface{}{
					"cidr":    "192.168.0.0/24",
					"gateway": "192.168.1.1",
					"dns_servers": []interface{}{
						"192.200.200.1",
						"192.200.200.201",
					},
				},
			},
		},
		{
			[]*models.MachineNetworkingConfig{},
			[]interface{}{},
		},
		{
			nil,
			[]interface{}{},
		},
	}

	for _, tc := range cases {
		output := flattenMachineNetworks(tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestExpandClusterSpec(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *models.ClusterSpec
		DCName         string
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"version":             "1.15.6",
					"machine_networks":    []interface{}{},
					"audit_logging":       false,
					"pod_security_policy": true,
					"cloud": []interface{}{
						map[string]interface{}{
							"bringyourown": []interface{}{
								map[string]interface{}{},
							},
						},
					},
				},
			},
			&models.ClusterSpec{
				Version:                             "1.15.6",
				MachineNetworks:                     nil,
				AuditLogging:                        &models.AuditLoggingSettings{},
				UsePodSecurityPolicyAdmissionPlugin: true,
				Cloud: &models.CloudSpec{
					DatacenterName: "eu-west-1",
					Bringyourown:   map[string]interface{}{},
				},
			},
			"eu-west-1",
		},
		{
			[]interface{}{
				map[string]interface{}{},
			},
			&models.ClusterSpec{},
			"",
		},
		{
			[]interface{}{},
			nil,
			"",
		},
	}

	for _, tc := range cases {
		output := expandClusterSpec(tc.Input, tc.DCName)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestExpandClusterCloudSpec(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *models.CloudSpec
		DCName         string
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"bringyourown": []interface{}{
						map[string]interface{}{},
					},
				},
			},
			&models.CloudSpec{
				DatacenterName: "eu-west-1",
				Bringyourown:   map[string]interface{}{},
			},
			"eu-west-1",
		},
		{
			[]interface{}{
				map[string]interface{}{
					"aws": []interface{}{
						map[string]interface{}{},
					},
				},
			},
			&models.CloudSpec{
				DatacenterName: "eu-west-1",
				Aws:            &models.AWSCloudSpec{},
			},
			"eu-west-1",
		},
		{
			[]interface{}{
				map[string]interface{}{},
			},
			&models.CloudSpec{
				DatacenterName: "eu-west-1",
			},
			"eu-west-1",
		},
		{
			[]interface{}{},
			nil,
			"eu-west-1",
		},
	}

	for _, tc := range cases {
		output := expandClusterCloudSpec(tc.Input, tc.DCName)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestExpandBringYourOwnCloud(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput models.BringYourOwnCloudSpec
	}{
		{
			[]interface{}{
				map[string]interface{}{},
			},
			map[string]interface{}{},
		},
		{
			[]interface{}{},
			nil,
		},
	}

	for _, tc := range cases {
		output := expandBringYourOwnCloudSpec(tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestExpandAWSCloudSpec(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *models.AWSCloudSpec
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"access_key_id":         "AKIAIOSFODNN7EXAMPLE",
					"role_arn":              "default",
					"instance_profile_name": "default",
					"route_table_id":        "rtb-09ba434c1bEXAMPLE",
					"secret_access_key":     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
					"security_group_id":     "sg-51530134",
					"vpc_id":                "e5e4b2ef2fe",
				},
			},
			&models.AWSCloudSpec{
				AccessKeyID:         "AKIAIOSFODNN7EXAMPLE",
				ControlPlaneRoleARN: "default",
				InstanceProfileName: "default",
				RouteTableID:        "rtb-09ba434c1bEXAMPLE",
				SecretAccessKey:     "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
				SecurityGroupID:     "sg-51530134",
				VPCID:               "e5e4b2ef2fe",
			},
		},
		{
			[]interface{}{
				map[string]interface{}{},
			},
			&models.AWSCloudSpec{},
		},
		{
			[]interface{}{},
			nil,
		},
	}

	for _, tc := range cases {
		output := expandAWSCloudSpec(tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestExpandOpenstackCloudSpec(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput *models.OpenstackCloudSpec
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"tenant":           "Tenant",
					"floating_ip_pool": "FloatingIPPool",
					"username":         "Username",
					"password":         "Password",
				},
			},
			&models.OpenstackCloudSpec{
				Domain:         "Default",
				FloatingIPPool: "FloatingIPPool",
				Password:       "Password",
				Tenant:         "Tenant",
				Username:       "Username",
			},
		},
		{
			[]interface{}{
				map[string]interface{}{},
			},
			&models.OpenstackCloudSpec{
				Domain: "Default",
			},
		},
		{
			[]interface{}{},
			nil,
		},
	}

	for _, tc := range cases {
		output := expandOpenstackCloudSpec(tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestExpandMachineNetwork(t *testing.T) {
	cases := []struct {
		Input          []interface{}
		ExpectedOutput []*models.MachineNetworkingConfig
	}{
		{
			[]interface{}{
				map[string]interface{}{
					"cidr":    "192.168.0.0/24",
					"gateway": "192.168.1.1",
					"dns_servers": []interface{}{
						"192.200.200.1",
						"192.200.200.201",
					},
				},
			},
			[]*models.MachineNetworkingConfig{
				{
					CIDR:    "192.168.0.0/24",
					Gateway: "192.168.1.1",
					DNSServers: []string{
						"192.200.200.1",
						"192.200.200.201",
					},
				},
			},
		},
		{
			[]interface{}{
				map[string]interface{}{},
			},
			[]*models.MachineNetworkingConfig{{}},
		},
		{
			[]interface{}{},
			nil,
		},
	}

	for _, tc := range cases {
		output := expandMachineNetworks(tc.Input)
		if diff := cmp.Diff(tc.ExpectedOutput, output); diff != "" {
			t.Fatalf("Unexpected output from expander: mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestExpandAuditLogging(t *testing.T) {
	want := &models.AuditLoggingSettings{
		Enabled: true,
	}
	got := expandAuditLogging(true)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %+v, got %+v", want, got)
	}
}

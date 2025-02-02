package unittest

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

const bitbucketModule = "products/bitbucket"

func TestBitbucketVariablesPopulatedWithValidValues(t *testing.T) {
	t.Parallel()

	tfOptions := GenerateTFOptions(BitbucketCorrectVariables, t, bitbucketModule)
	plan := terraform.InitAndPlanAndShowWithStruct(t, tfOptions)

	// verify Bitbucket
	bitbucketKey := "helm_release.bitbucket"
	terraform.RequirePlannedValuesMapKeyExists(t, plan, bitbucketKey)
	bitbucket := plan.ResourcePlannedValuesMap[bitbucketKey]
	assert.Equal(t, "deployed", bitbucket.AttributeValues["status"])
	assert.Equal(t, "bitbucket", bitbucket.AttributeValues["chart"])
	assert.Equal(t, "https://atlassian.github.io/data-center-helm-charts", bitbucket.AttributeValues["repository"])
}

func TestBitbucketVariablesPopulatedWithInvalidValues(t *testing.T) {
	t.Parallel()

	tfOptions := GenerateTFOptions(BitbucketInvalidVariables, t, bitbucketModule)
	_, err := terraform.InitAndPlanAndShowWithStructE(t, tfOptions)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid value for variable")
	assert.Contains(t, err.Error(), "Invalid environment name. Valid name is up to 25 characters starting with")
	assert.Contains(t, err.Error(), "Bitbucket configuration is not valid.")
	assert.Contains(t, err.Error(), "Bitbucket administrator configuration is not valid.")
	assert.Contains(t, err.Error(), "Invalid elasticsearch replicas. Valid replicas is a positive integer in")
	assert.Contains(t, err.Error(), "Bitbucket display name must be a non-empty value less than 255 characters.")
}

func TestBitbucketVariablesNotProvided(t *testing.T) {
	t.Parallel()

	tfOptions := GenerateTFOptions(nil, t, bitbucketModule)

	_, err := terraform.InitAndPlanAndShowWithStructE(t, tfOptions)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "No value for required variable")
	assert.Contains(t, err.Error(), "\"environment_name\" is not set")
	assert.Contains(t, err.Error(), "\"namespace\" is not set")
	assert.Contains(t, err.Error(), "\"vpc\" is not set")
	assert.Contains(t, err.Error(), "\"eks\" is not set")
	assert.Contains(t, err.Error(), "\"db_major_engine_version\" is not set")
	assert.Contains(t, err.Error(), "\"db_allocated_storage\" is not set")
	assert.Contains(t, err.Error(), "\"db_instance_class\" is not set")
	assert.Contains(t, err.Error(), "\"db_iops\" is not set")
	assert.Contains(t, err.Error(), "\"bitbucket_configuration\" is not set")
	assert.Contains(t, err.Error(), "\"admin_configuration\" is not set")
	assert.Contains(t, err.Error(), "\"elasticsearch_cpu\" is not set")
	assert.Contains(t, err.Error(), "\"elasticsearch_mem\" is not set")
	assert.Contains(t, err.Error(), "\"elasticsearch_storage\" is not set")
	assert.Contains(t, err.Error(), "\"elasticsearch_replicas\" is not set")
	assert.NotContains(t, err.Error(), "display_name")
}

const nfsModule = "products/bitbucket/nfs"

func TestNfsVariablesNotProvided(t *testing.T) {
	t.Parallel()

	tfOptions := GenerateTFOptions(nil, t, nfsModule)

	_, err := terraform.InitAndPlanAndShowWithStructE(t, tfOptions)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "No value for required variable")
	assert.Contains(t, err.Error(), "\"namespace\" is not set")
}

func TestNfsVariablesPopulatedWithValidValues(t *testing.T) {
	t.Parallel()

	tfOptions := GenerateTFOptions(NfsValidVariable, t, nfsModule)
	plan := terraform.InitAndPlanAndShowWithStruct(t, tfOptions)

	terraform.RequirePlannedValuesMapKeyExists(t, plan, "helm_release.nfs")
	helmRelease := plan.ResourcePlannedValuesMap["helm_release.nfs"]
	values := helmRelease.AttributeValues["values"].([]interface{})[0].(string)

	expectedHelmValues := fmt.Sprintf("\"nameOverride\": \"%s\"\n\"persistence\":\n  \"size\": \"%s\"\n\"resources\":\n  \"limits\":\n    \"cpu\": \"%s\"\n    \"memory\": \"%s\"\n  \"requests\":\n    \"cpu\": \"%s\"\n    \"memory\": \"%s\"\n",
		nfsVarChartNameOverride, nfsVarCapacity, nfsLimitsCpu, nfsLimitsMemory, nfsRequestsCpu, nfsRequestsMemory)

	expectedNamespace := nfsVarNamespace

	assert.Equal(t, "bitbucket-nfs", helmRelease.AttributeValues["name"])
	assert.Equal(t, expectedNamespace, helmRelease.AttributeValues["namespace"])
	assert.Equal(t, expectedHelmValues, values)
}

// Variables
var BitbucketCorrectVariables = map[string]interface{}{
	"environment_name": "dummy-environment",
	"namespace":        "dummy-namespace",
	"eks": map[string]interface{}{
		"kubernetes_provider_config": map[string]interface{}{
			"host":                   "dummy-host",
			"token":                  "dummy-token",
			"cluster_ca_certificate": "dummy-certificate",
		},
		"cluster_security_group": "dummy-sg",
		"cluster_size":           2,
	},
	"vpc":                     VpcDefaultModuleVariable,
	"db_major_engine_version": "13",
	"db_allocated_storage":    5,
	"db_instance_class":       "dummy_db_instance_class",
	"db_iops":                 1000,
	"admin_configuration": map[string]interface{}{
		"admin_username":      "dummy_admin_username",
		"admin_password":      "dummy_admin_password",
		"admin_display_name":  "dummy_admin_display_name",
		"admin_email_address": "dummy_admin_email_address",
	},
	"display_name":  "dummy_display_name",
	"ingress":       map[string]interface{}{},
	"replica_count": 1,
	"bitbucket_configuration": map[string]interface{}{
		"helm_version": "1.2.0",
		"cpu":          "1",
		"mem":          "1Gi",
		"min_heap":     "256m",
		"max_heap":     "512m",
		"license":      "dummy_license",
	},
	"nfs_requests_cpu":       "0.25",
	"nfs_requests_memory":    "256Mi",
	"nfs_limits_cpu":         "0.25",
	"nfs_limits_memory":      "256Mi",
	"elasticsearch_cpu":      "1",
	"elasticsearch_mem":      "1Gi",
	"elasticsearch_storage":  10,
	"elasticsearch_replicas": 2,
}

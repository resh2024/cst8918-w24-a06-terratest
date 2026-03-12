package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "d4f264fa-08cf-4cd1-bfcb-e8ea4367743b"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "resh0004",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))
	// --- NIC connected to VM ---
	t.Run("NIC exists", func(t *testing.T) {
		nicName := terraform.Output(t, terraformOptions, "nic_name")

		exists := azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID)
		assert.True(t, exists, "NIC should exist")
	})
	t.Run("VM has correct Ubuntu version", func(t *testing.T) {
		vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)

		imagePublisher := *vm.StorageProfile.ImageReference.Publisher
		imageOffer := *vm.StorageProfile.ImageReference.Offer
		imageSKU := *vm.StorageProfile.ImageReference.Sku

		assert.Equal(t, "Canonical", imagePublisher)
		assert.Contains(t, imageOffer, "ubuntu")
		assert.Contains(t, imageSKU, "22")
	})
}

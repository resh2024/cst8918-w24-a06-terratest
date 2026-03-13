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
	terraform.InitAndApply(t, terraformOptions)

	// Collect Terraform outputs
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	rgName := terraform.Output(t, terraformOptions, "resource_group_name")

	// ----------------------
	// Check 1: VM existence
	// ----------------------
	exists := azure.VirtualMachineExists(t, vmName, rgName, subscriptionID)
	assert.True(t, exists, "Expected VM to exist in the resource group")

	// Retrieve VM details for further validation
	vmDetails, err := azure.GetVirtualMachineE(vmName, rgName, subscriptionID)
	assert.NoError(t, err)

	// ----------------------
	// Check 2: Ubuntu version
	// ----------------------
	skuValue := *vmDetails.StorageProfile.ImageReference.Sku
	expectedUbuntuSku := "22_04-lts-gen2"

	assert.Equal(
		t,
		expectedUbuntuSku,
		skuValue,
		"VM image SKU does not match the expected Ubuntu version",
	)

	// ----------------------
	// Check 3: NIC attachment
	// ----------------------
	nicList, err := FetchVMNicNames(vmName, rgName, subscriptionID)
	assert.NoError(t, err)

	assert.Greater(
		t,
		len(nicList),
		0,
		"VM should have at least one network interface attached",
	)

	// Verify each NIC exists in Azure
	for _, nic := range nicList {
		assert.True(
			t,
			azure.NetworkInterfaceExists(t, nic, rgName, subscriptionID),
			"Network interface "+nic+" should exist",
		)
	}
}

func FetchVMNicNames(vmName string, resourceGroup string, subscriptionID string) ([]string, error) {

	vm, err := azure.GetVirtualMachineE(vmName, resourceGroup, subscriptionID)
	if err != nil {
		return nil, err
	}

	if vm.NetworkProfile == nil || vm.NetworkProfile.NetworkInterfaces == nil {
		return []string{}, nil
	}

	vmInterfaces := *vm.NetworkProfile.NetworkInterfaces
	nicNames := make([]string, 0)

	for _, nic := range vmInterfaces {
		name, err := azure.GetNameFromResourceIDE(*nic.ID)
		if err == nil {
			nicNames = append(nicNames, name)
		}
	}

	return nicNames, nil
}

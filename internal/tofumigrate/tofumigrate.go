package tofumigrate

import (
	tfaddr "github.com/opentofu/registry-address"

	"github.com/opentofu/opentofu/internal/configs"
	"github.com/opentofu/opentofu/internal/states"
	"github.com/opentofu/opentofu/internal/tfdiags"
)

// MigrateStateProviderAddresses can be used to update the in-memory view of the state to use registry.opentofu.org
// provider addresses. This only applies for providers which are *not* explicitly referenced in the configuration in full form.
// For example, if the configuration contains a provider block like this:
//
//	terraform {
//	 required_providers {
//	   random = {}
//	 }
//	}
//
// we will migrate the in-memory view of the statefile to use registry.opentofu.org/hashicorp/random.
// However, if the configuration contains a provider block like this:
//
//	terraform {
//	 required_providers {
//	   random = {
//	     source = "registry.terraform.io/hashicorp/random"
//	   }
//	 }
//	}
//
// then we keep the old address.
func MigrateStateProviderAddresses(config *configs.Config, state *states.State) (*states.State, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics

	stateCopy := state.DeepCopy()

	providers, hclDiags := config.ProviderRequirements()
	diags = diags.Append(hclDiags)
	if hclDiags.HasErrors() {
		return nil, diags
	}

	for _, module := range stateCopy.Modules {
		for _, resource := range module.Resources {
			_, referencedInConfig := providers[resource.ProviderConfig.Provider]
			if resource.ProviderConfig.Provider.Hostname == "registry.terraform.io" && !referencedInConfig {
				resource.ProviderConfig.Provider.Hostname = tfaddr.DefaultProviderRegistryHost
			}
		}
	}

	return stateCopy, diags
}

// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tofu

import (
	"context"
	"fmt"
	"log"

	"github.com/opentofu/opentofu/internal/instances"
	"github.com/opentofu/opentofu/internal/plans"
	"github.com/opentofu/opentofu/internal/tfdiags"

	"github.com/opentofu/opentofu/internal/addrs"
	"github.com/opentofu/opentofu/internal/configs"
	"github.com/opentofu/opentofu/internal/states"
)

// NodeDestroyResourceInstance represents a resource instance that is to be
// destroyed.
type NodeDestroyResourceInstance struct {
	*NodeAbstractResourceInstance

	// If DeposedKey is set to anything other than states.NotDeposed then
	// this node destroys a deposed object of the associated instance
	// rather than its current object.
	DeposedKey states.DeposedKey
}

var (
	_ GraphNodeModuleInstance      = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeConfigResource      = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeResourceInstance    = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeDestroyer           = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeDestroyerCBD        = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeReferenceable       = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeReferencer          = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeExecutable          = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeProviderConsumer    = (*NodeDestroyResourceInstance)(nil)
	_ GraphNodeProvisionerConsumer = (*NodeDestroyResourceInstance)(nil)
)

func (n *NodeDestroyResourceInstance) Name() string {
	if n.DeposedKey != states.NotDeposed {
		return fmt.Sprintf("%s (destroy deposed %s)", n.ResourceInstanceAddr(), n.DeposedKey)
	}
	return n.ResourceInstanceAddr().String() + " (destroy)"
}

func (n *NodeDestroyResourceInstance) ProvidedBy() RequestedProvider {
	if n.Addr.Resource.Resource.Mode == addrs.DataResourceMode {
		// indicate that this node does not require a configured provider
		return RequestedProvider{}
	}
	return n.NodeAbstractResourceInstance.ProvidedBy()
}

// GraphNodeDestroyer
func (n *NodeDestroyResourceInstance) DestroyAddr() *addrs.AbsResourceInstance {
	addr := n.ResourceInstanceAddr()
	return &addr
}

// GraphNodeDestroyerCBD
func (n *NodeDestroyResourceInstance) CreateBeforeDestroy() bool {
	// State takes precedence during destroy.
	// If the resource was removed, there is no config to check.
	// If CBD was forced from descendent, it should be saved in the state
	// already.
	if s := n.instanceState; s != nil {
		if s.Current != nil {
			return s.Current.CreateBeforeDestroy
		}
	}

	if n.Config != nil && n.Config.Managed != nil {
		return n.Config.Managed.CreateBeforeDestroy
	}

	return false
}

// GraphNodeDestroyerCBD
func (n *NodeDestroyResourceInstance) ModifyCreateBeforeDestroy(v bool) error {
	return nil
}

// GraphNodeReferenceable, overriding NodeAbstractResource
func (n *NodeDestroyResourceInstance) ReferenceableAddrs() []addrs.Referenceable {
	normalAddrs := n.NodeAbstractResourceInstance.ReferenceableAddrs()
	destroyAddrs := make([]addrs.Referenceable, len(normalAddrs))

	phaseType := addrs.ResourceInstancePhaseDestroy
	if n.CreateBeforeDestroy() {
		phaseType = addrs.ResourceInstancePhaseDestroyCBD
	}

	for i, normalAddr := range normalAddrs {
		switch ta := normalAddr.(type) {
		case addrs.Resource:
			destroyAddrs[i] = ta.Phase(phaseType)
		case addrs.ResourceInstance:
			destroyAddrs[i] = ta.Phase(phaseType)
		default:
			destroyAddrs[i] = normalAddr
		}
	}

	return destroyAddrs
}

// GraphNodeReferencer, overriding NodeAbstractResource
func (n *NodeDestroyResourceInstance) References() []*addrs.Reference {
	// If we have a config, then we need to include destroy-time dependencies
	if c := n.Config; c != nil && c.Managed != nil {
		var result []*addrs.Reference

		// We include conn info and config for destroy time provisioners
		// as dependencies that we have.
		for _, p := range c.Managed.Provisioners {
			schema := n.ProvisionerSchemas[p.Type]

			if p.When == configs.ProvisionerWhenDestroy {
				if p.Connection != nil {
					result = append(result, ReferencesFromConfig(p.Connection.Config, connectionBlockSupersetSchema)...)
				}
				result = append(result, ReferencesFromConfig(p.Config, schema)...)
			}
		}

		return result
	}

	return nil
}

// GraphNodeExecutable
func (n *NodeDestroyResourceInstance) Execute(evalCtx EvalContext, _ walkOperation) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	// TODO: In future we'll change GraphNodeExecutable to include context.Context as the first
	// argument to this method, but for now we'll just stub it.
	ctx := context.TODO()

	addr := n.ResourceInstanceAddr()

	// Eval info is different depending on what kind of resource this is
	var execute func(ctx context.Context, evalCtx EvalContext) tfdiags.Diagnostics
	switch addr.Resource.Resource.Mode {
	case addrs.ManagedResourceMode:
		diags = n.resolveProvider(evalCtx, false)
		if diags.HasErrors() {
			return diags
		}
		execute = n.managedResourceExecute
	case addrs.DataResourceMode:
		execute = n.dataResourceExecute
	default:
		panic(fmt.Errorf("unsupported resource mode %s", n.Config.Mode))
	}

	// TODO: The entire "execute" function is a very coarse way to define "I/O";
	// it would be better to limit this only to the individual provider calls
	// and the associated log/hook events about them, but we've implemented it
	// this way at first just to minimize the size of the change compared to
	// the old design where the entire call to Execute was the unit used for
	// concurrency limit, previously handled automatically by the graph walker.
	return evalCtx.PerformIO(ctx, func(ctx context.Context) tfdiags.Diagnostics {
		return execute(ctx, evalCtx)
	})
}

func (n *NodeDestroyResourceInstance) managedResourceExecute(_ context.Context, evalCtx EvalContext) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	addr := n.ResourceInstanceAddr()

	// Get our state
	is := n.instanceState
	if is == nil {
		log.Printf("[WARN] NodeDestroyResourceInstance for %s with no state", addr)
	}

	// These vars are updated through pointers at various stages below.
	var changeApply *plans.ResourceInstanceChange
	var state *states.ResourceInstanceObject

	_, providerSchema, err := getProvider(evalCtx, n.ResolvedProvider.ProviderConfig, n.ResolvedProviderKey)
	diags = diags.Append(err)
	if diags.HasErrors() {
		return diags
	}

	changeApply, err = n.readDiff(evalCtx, providerSchema)
	diags = diags.Append(err)
	if changeApply == nil || diags.HasErrors() {
		return diags
	}

	changeApply = reducePlan(addr.Resource, changeApply, true)
	// reducePlan may have simplified our planned change
	// into a NoOp if it does not require destroying.
	if changeApply == nil || changeApply.Action == plans.NoOp {
		return diags
	}

	state, readDiags := n.readResourceInstanceState(evalCtx, addr)
	diags = diags.Append(readDiags)
	if diags.HasErrors() {
		return diags
	}

	// Exit early if the state object is null after reading the state
	if state == nil || state.Value.IsNull() {
		return diags
	}

	diags = diags.Append(n.preApplyHook(evalCtx, changeApply))
	if diags.HasErrors() {
		return diags
	}

	// Run destroy provisioners if not tainted
	if state.Status != states.ObjectTainted {
		applyProvisionersDiags := n.evalApplyProvisioners(evalCtx, state, false, configs.ProvisionerWhenDestroy)
		diags = diags.Append(applyProvisionersDiags)
		// keep the diags separate from the main set until we handle the cleanup

		if diags.HasErrors() {
			// If we have a provisioning error, then we just call
			// the post-apply hook now.
			diags = diags.Append(n.postApplyHook(evalCtx, state, diags.Err()))
			return diags
		}
	}

	// Managed resources need to be destroyed, while data sources
	// are only removed from state.
	// we pass a nil configuration to apply because we are destroying
	s, d := n.apply(evalCtx, state, changeApply, nil, instances.RepetitionData{}, false)
	state, diags = s, diags.Append(d)
	// we don't return immediately here on error, so that the state can be
	// finalized

	err = n.writeResourceInstanceState(evalCtx, state, workingState)
	if err != nil {
		return diags.Append(err)
	}

	// create the err value for postApplyHook
	diags = diags.Append(n.postApplyHook(evalCtx, state, diags.Err()))
	diags = diags.Append(updateStateHook(evalCtx))
	return diags
}

func (n *NodeDestroyResourceInstance) dataResourceExecute(_ context.Context, evalCtx EvalContext) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics

	log.Printf("[TRACE] NodeDestroyResourceInstance: removing state object for %s", n.Addr)
	evalCtx.State().SetResourceInstanceCurrent(n.Addr, nil, n.ResolvedProvider.ProviderConfig, n.ResolvedProviderKey)
	return diags.Append(updateStateHook(evalCtx))
}

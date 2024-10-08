// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package states

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/opentofu/opentofu/internal/addrs"
)

// Resource represents the state of a resource.
type Resource struct {
	// Addr is the absolute address for the resource this state object
	// belongs to.
	Addr addrs.AbsResource

	// Instances contains the potentially-multiple instances associated with
	// this resource. This map can contain a mixture of different key types,
	// but only the ones of InstanceKeyType are considered current.
	Instances map[addrs.InstanceKey]*ResourceInstance

	// ProviderConfig is the absolute address for the provider configuration that
	// most recently managed this resource and all its instances. This is used to
	// connect a resource with a provider configuration when the resource
	// configuration block is not available, such as if it has been removed from
	// configuration altogether.
	// ProviderConfig can be not set if the provider is set in the instance level,
	// in cases we have count / for_each in provider reference of the resource /
	// module. Then, instead, ProviderInstanceConfig will be set on each instance.
	ProviderConfig addrs.AbsProviderConfig
}

// Instance returns the state for the instance with the given key, or nil
// if no such instance is tracked within the state.
func (rs *Resource) Instance(key addrs.InstanceKey) *ResourceInstance {
	return rs.Instances[key]
}

// InstanceProvider returns the calculated provider for the instance with the given key. It also returns an indication
// whether the provider was set on the resource level or the instance level
func (rs *Resource) InstanceProvider(key addrs.InstanceKey) (provider addrs.AbsProviderConfig, isSetOnInstanceLevel bool) {
	var resourceProvider addrs.AbsProviderConfig
	var instanceProvider addrs.AbsProviderConfig

	resourceProvider = rs.ProviderConfig

	// If the provider is set on the instance level, we can't get it from the resource
	instance := rs.Instances[key]

	if instance.Current != nil && instance.Current.InstanceProvider.IsSet() {
		instanceProvider = instance.Current.InstanceProvider
	}

	if !resourceProvider.IsSet() && !instanceProvider.IsSet() {
		// At this point we are trying to find any provider
		// If instance.Current is not set, then maybe the resource has deposed instances instead
		for _, deposedInstance := range instance.Deposed {
			// We are assuming that all the deposed instances should have the same instance provider, so we can get it
			// from the first deposed instance we stumble upon.
			// This assumption might cause a bug, in the scenario where the provider of the resource got changed between
			// runs, and the deposed object still requires the old provider configuration. This bug exists not only for
			// the InstanceProvider, but also always existed for the resourceProvider, and we should solve it in a
			// holistic approach.
			if deposedInstance.InstanceProvider.IsSet() {
				// Found one, let's assume it's good enough for now
				instanceProvider = deposedInstance.InstanceProvider
				break // Exit after the first iteration
			}
		}
	}

	if !resourceProvider.IsSet() && !instanceProvider.IsSet() {
		panic(fmt.Sprintf("InstanceProvider for %s (instance key %s) failed to read provider from the state", rs.Addr, key.String()))
	}

	if resourceProvider.IsSet() && instanceProvider.IsSet() {
		panic(fmt.Sprintf("InstanceProvider for %s (instance key %s) found two providers in state for the instance", rs.Addr, key.String()))
	}

	if resourceProvider.IsSet() {
		return resourceProvider, true
	} else {
		return instanceProvider, false
	}
}

// CreateInstance creates an instance and adds it to the resource
func (rs *Resource) CreateInstance(key addrs.InstanceKey) *ResourceInstance {
	is := NewResourceInstance()
	rs.Instances[key] = is
	return is
}

// EnsureInstance returns the state for the instance with the given key,
// creating a new empty state for it if one doesn't already exist.
//
// Because this may create and save a new state, it is considered to be
// a write operation.
func (rs *Resource) EnsureInstance(key addrs.InstanceKey) *ResourceInstance {
	ret := rs.Instance(key)
	if ret == nil {
		ret = NewResourceInstance()
		rs.Instances[key] = ret
	}
	return ret
}

// ResourceInstance represents the state of a particular instance of a resource.
type ResourceInstance struct {
	// Current, if non-nil, is the remote object that is currently represented
	// by the corresponding resource instance.
	Current *ResourceInstanceObjectSrc

	// Deposed, if len > 0, contains any remote objects that were previously
	// represented by the corresponding resource instance but have been
	// replaced and are pending destruction due to the create_before_destroy
	// lifecycle mode.
	Deposed map[DeposedKey]*ResourceInstanceObjectSrc
}

// NewResourceInstance constructs and returns a new ResourceInstance, ready to
// use.
func NewResourceInstance() *ResourceInstance {
	return &ResourceInstance{
		Deposed: map[DeposedKey]*ResourceInstanceObjectSrc{},
	}
}

// HasCurrent returns true if this resource instance has a "current"-generation
// object. Most instances do, but this can briefly be false during a
// create-before-destroy replace operation when the current has been deposed
// but its replacement has not yet been created.
func (i *ResourceInstance) HasCurrent() bool {
	return i != nil && i.Current != nil
}

// HasDeposed returns true if this resource instance has a deposed object
// with the given key.
func (i *ResourceInstance) HasDeposed(key DeposedKey) bool {
	return i != nil && i.Deposed[key] != nil
}

// HasAnyDeposed returns true if this resource instance has one or more
// deposed objects.
func (i *ResourceInstance) HasAnyDeposed() bool {
	return i != nil && len(i.Deposed) > 0
}

// HasObjects returns true if this resource has any objects at all, whether
// current or deposed.
func (i *ResourceInstance) HasObjects() bool {
	return i.Current != nil || len(i.Deposed) != 0
}

// deposeCurrentObject is part of the real implementation of
// SyncState.DeposeResourceInstanceObject. The exported method uses a lock
// to ensure that we can safely allocate an unused deposed key without
// collision.
func (i *ResourceInstance) deposeCurrentObject(forceKey DeposedKey) DeposedKey {
	if !i.HasCurrent() {
		return NotDeposed
	}

	key := forceKey
	if key == NotDeposed {
		key = i.findUnusedDeposedKey()
	} else {
		if _, exists := i.Deposed[key]; exists {
			panic(fmt.Sprintf("forced key %s is already in use", forceKey))
		}
	}
	i.Deposed[key] = i.Current
	i.Current = nil
	return key
}

// GetGeneration retrieves the object of the given generation from the
// ResourceInstance, or returns nil if there is no such object.
//
// If the given generation is nil or invalid, this method will panic.
func (i *ResourceInstance) GetGeneration(gen Generation) *ResourceInstanceObjectSrc {
	if gen == CurrentGen {
		return i.Current
	}
	if dk, ok := gen.(DeposedKey); ok {
		return i.Deposed[dk]
	}
	if gen == nil {
		panic("get with nil Generation")
	}
	// Should never fall out here, since the above covers all possible
	// Generation values.
	panic(fmt.Sprintf("get invalid Generation %#v", gen))
}

// FindUnusedDeposedKey generates a unique DeposedKey that is guaranteed not to
// already be in use for this instance at the time of the call.
//
// Note that the validity of this result may change if new deposed keys are
// allocated before it is used. To avoid this risk, instead use the
// DeposeResourceInstanceObject method on the SyncState wrapper type, which
// allocates a key and uses it atomically.
func (i *ResourceInstance) FindUnusedDeposedKey() DeposedKey {
	return i.findUnusedDeposedKey()
}

// findUnusedDeposedKey generates a unique DeposedKey that is guaranteed not to
// already be in use for this instance.
func (i *ResourceInstance) findUnusedDeposedKey() DeposedKey {
	for {
		key := NewDeposedKey()
		if _, exists := i.Deposed[key]; !exists {
			return key
		}
		// Spin until we find a unique one. This shouldn't take long, because
		// we have a 32-bit keyspace and there's rarely more than one deposed
		// instance.
	}
}

// DeposedKey is a 8-character hex string used to uniquely identify deposed
// instance objects in the state.
type DeposedKey string

// NotDeposed is a special invalid value of DeposedKey that is used to represent
// the absence of a deposed key. It must not be used as an actual deposed key.
const NotDeposed = DeposedKey("")

var deposedKeyRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// NewDeposedKey generates a pseudo-random deposed key. Because of the short
// length of these keys, uniqueness is not a natural consequence and so the
// caller should test to see if the generated key is already in use and generate
// another if so, until a unique key is found.
func NewDeposedKey() DeposedKey {
	v := deposedKeyRand.Uint32()
	return DeposedKey(fmt.Sprintf("%08x", v))
}

func (k DeposedKey) String() string {
	return string(k)
}

func (k DeposedKey) GoString() string {
	ks := string(k)
	switch {
	case ks == "":
		return "states.NotDeposed"
	default:
		return fmt.Sprintf("states.DeposedKey(%s)", ks)
	}
}

// Generation is a helper method to convert a DeposedKey into a Generation.
// If the receiver is anything other than NotDeposed then the result is
// just the same value as a Generation. If the receiver is NotDeposed then
// the result is CurrentGen.
func (k DeposedKey) Generation() Generation {
	if k == NotDeposed {
		return CurrentGen
	}
	return k
}

// generation is an implementation of Generation.
func (k DeposedKey) generation() {}

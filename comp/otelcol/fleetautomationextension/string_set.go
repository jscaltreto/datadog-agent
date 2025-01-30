// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package fleetautomationextension

// StringSet is a small dumb map-backed implementation of a set. Not
// thread-safe.
type StringSet map[string]struct{}

// NewStringSet returns a new empty string set.
func NewStringSet() StringSet {
	return make(StringSet)
}

// NewStringSetFromSlice returns a new string set.
func NewStringSetFromSlice(sl []string) StringSet {
	ss := make(StringSet)
	for _, s := range sl {
		ss.Add(s)
	}
	return ss
}

// NewStringSetWithCapacity returns a new empty string set with the given expected capacity.
func NewStringSetWithCapacity(c int) StringSet {
	return make(StringSet, c)
}

// Add adds a key ot the set.
func (ss StringSet) Add(key string) {
	ss[key] = struct{}{}
}

// Remove removes a key from the set.
func (ss StringSet) Remove(key string) {
	delete(ss, key)
}

// Contains returns true if the key is in the set, false otherwise
func (ss StringSet) Contains(key string) bool {
	_, ok := ss[key]
	return ok
}

// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package fleetautomationextension

import (
	"time"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"
)

var file_domains_redapl_shared_resourcespb_proto_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 20)
var file_domains_redapl_shared_resourcespb_proto_async_intake_proto_msgTypes = make([]protoimpl.MessageInfo, 1)

type RedaplEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the source is required and cannot be empty
	// if the source isn't recognized then the payload is assumed to be a RawResourceV3
	Source string `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	// this can be either a RawResourceV3 or a client-specific payload
	// supported by redapl async intake
	Message []byte `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *RedaplEvent) Reset() {
	*x = RedaplEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domains_redapl_shared_resourcespb_proto_async_intake_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedaplEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (x *RedaplEvent) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (*RedaplEvent) ProtoMessage() {}

func (x *RedaplEvent) ProtoReflect() protoreflect.Message {
	mi := &file_domains_redapl_shared_resourcespb_proto_async_intake_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*RawResourceV3) ProtoMessage() {}

func (x *RawResourceV3) Reset() {
	*x = RawResourceV3{}
	if protoimpl.UnsafeEnabled {
		mi := &file_domains_redapl_shared_resourcespb_proto_resource_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawResourceV3) ProtoReflect() protoreflect.Message {
	mi := &file_domains_redapl_shared_resourcespb_proto_resource_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *RawResourceV3) String() string {
	return protoimpl.X.MessageStringOf(x)
}

// BuildMetadata holds metadata for the migrator, when available, it provides additional information of the migration
// that was ran over the resource
type MigrationBuildMetadata struct {
	BuildCommit string
	BuildDate   time.Time
}

// // MigrationRecord is an internal type within the metadata to reflect a migration applied to a resource
type MigrationRecordMetadata struct {
	RevisionID        string
	RevisionCreatedAt time.Time
	Timestamp         time.Time
	BuildMetadata     *MigrationBuildMetadata
}

// MigrationDerived describes that the current resource was derived due to a migration from  a separate resource
type MigrationDerived struct {
	FromResourceType string
	FromResourceName string
	FromRevisionId   string
}

// MigrationTestMetadata describes a shadow resource
type MigrationTestMetadata struct {
	OriginalResourceType string
}

// MigrationsMetadata is an internal type within the metadata to reflect whether
// it is a derived resource, and all applied migrations
type MigrationsMetadata struct {
	MigrationRecords []MigrationRecordMetadata
	Derived          *MigrationDerived
	TestMetadata     *MigrationTestMetadata
}

// ResourceStoreProxyMetadata stores metadata useful for resource store proxy to
// store internally and to expose to callers.
type ResourceStoreProxyMetadata struct {
	// Populated on retrieval - indicates the backend that returned the
	// information.
	Backend string

	// Populated on retrieval if false / nil; otherwise respect the previous
	// value from the backend.
	//
	// When a backend retrieves a resources with non-zero Missing/Stale/Err, it
	// respects the existing value. This indicates that a resource with one of
	// these characteristics was written back to the current backend, so it
	// doesn't make sense to overwrite the value. If the values
	// are zero, the backend sets the values based on its own semantics.
	Missing bool
	Stale   bool
	Err     error

	// Never populated on retrieval; stored before write, either in the normal
	// write path or as part of a write-back to a datastore on retrieval.
	Provenance string // most recent backend this value was retrieved from (i.e. not original provenance) or "client"
	WrittenAt  time.Time
}

// InternalMetadata contains data for the resource that isn't relevant
// semantically to the resource, but is needed for the resource pipeline
type InternalMetadata struct {
	EnableDebugTracking bool
	ResourceStoreProxy  ResourceStoreProxyMetadata
	Migrations          MigrationsMetadata
	IsCanary            bool
}

// SchemaVersion is a version type for schemas
// Deprecated : Schemas semantics do not hold versions
type SchemaVersion int32

// RawResourceV3 is a raw resourceV3
type RawResourceV3 struct {
	OrgID        int64
	Type         string
	Name         string
	FieldsByName map[string]*structpb.Value
	SeenAt       time.Time
	ExpireAt     time.Time
	Meta         InternalMetadata
	Version      SchemaVersion
	Scope        string
	FieldsJson   []byte
}

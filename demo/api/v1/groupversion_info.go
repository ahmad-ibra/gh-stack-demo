// Package v1 contains the API schema definitions for the backup v1 API group.
// +kubebuilder:object:generate=true
// +groupName=backup.example.com
package v1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is the group/version used to register these objects.
	GroupVersion = schema.GroupVersion{Group: "backup.example.com", Version: "v1"}

	// SchemeBuilder registers the Go types with the GroupVersionKind scheme.
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group/version to a runtime scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

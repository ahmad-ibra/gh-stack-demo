// Package features tracks which optional operator capabilities are compiled
// into this build. Each stack layer that adds a controller or webhook appends
// its feature here, so the manager can wire up exactly what is enabled.
package features

// Feature is the name of an optional operator capability.
type Feature string

const (
	// FeatureClusterBackup enables the scheduled cluster-backup controller.
	FeatureClusterBackup Feature = "cluster-backup"
)

// registered is the set of features compiled into this operator, in
// registration order.
var registered = []Feature{
	FeatureClusterBackup,
}

// Enabled reports whether the given feature is registered in this build.
func Enabled(f Feature) bool {
	for _, r := range registered {
		if r == f {
			return true
		}
	}
	return false
}

// All returns the registered features in registration order.
func All() []Feature {
	out := make([]Feature, len(registered))
	copy(out, registered)
	return out
}

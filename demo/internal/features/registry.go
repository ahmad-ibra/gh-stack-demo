package features

// registeredFeatures lists feature keys enabled in this build.
//
// NOTE: the demo intentionally has multiple stack layers append to this slice
// at the same spot, to produce a controlled rebase conflict on `gs sync`.
var registeredFeatures = []string{
	"clusters",
}

package mixin

import (
	// Import the original pkg path, generator will replace this string

	"get.porter.sh/porter/pkg/mixin"
	"get.porter.sh/porter/pkg/pkgmgmt"
	"get.porter.sh/porter/pkg/porter/version"

	"github.com/getporter/skeletor/pkg"
)

func (m *Mixin) PrintVersion(opts version.Options) error {
	metadata := mixin.Metadata{
		Name: "skeletor", // Use original name, generator will replace this string
		VersionInfo: pkgmgmt.VersionInfo{
			Version: pkg.Version, // Use original pkg alias
			Commit:  pkg.Commit,
			Author:  "YOURNAME", // Use placeholder, generator will replace this string
		},
	}
	return version.PrintVersion(m.Context, opts, metadata)
}

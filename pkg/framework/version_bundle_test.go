package framework

import (
	"sort"
	"testing"

	"github.com/giantswarm/versionbundle"
)

func Test_SortBySemver(t *testing.T) {
	tcs := []struct {
		inputIndexReleases  []versionbundle.IndexRelease
		outputIndexReleases []versionbundle.IndexRelease
		description         string
	}{
		{
			description: "right order, no need to change",
			inputIndexReleases: []versionbundle.IndexRelease{
				versionbundle.IndexRelease{Version: "0.2.0"},
				versionbundle.IndexRelease{Version: "0.1.0"},
			},
			outputIndexReleases: []versionbundle.IndexRelease{
				versionbundle.IndexRelease{Version: "0.2.0"},
				versionbundle.IndexRelease{Version: "0.1.0"},
			},
		},
		{
			description: "reversed order",
			inputIndexReleases: []versionbundle.IndexRelease{
				versionbundle.IndexRelease{Version: "0.1.0"},
				versionbundle.IndexRelease{Version: "0.2.0"},
			},
			outputIndexReleases: []versionbundle.IndexRelease{
				versionbundle.IndexRelease{Version: "0.2.0"},
				versionbundle.IndexRelease{Version: "0.1.0"},
			},
		},
		{
			description: "multiple items",
			inputIndexReleases: []versionbundle.IndexRelease{
				versionbundle.IndexRelease{Version: "0.2.6"},
				versionbundle.IndexRelease{Version: "1.1.3"},
				versionbundle.IndexRelease{Version: "0.2.5"},
			},
			outputIndexReleases: []versionbundle.IndexRelease{
				versionbundle.IndexRelease{Version: "1.1.3"},
				versionbundle.IndexRelease{Version: "0.2.6"},
				versionbundle.IndexRelease{Version: "0.2.5"},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.description, func(t *testing.T) {
			actual := SortReleasesBySemver(tc.inputIndexReleases)
			sort.Sort(sort.Reverse(actual))

			for i := 0; i < len(tc.outputIndexReleases); i++ {
				if tc.outputIndexReleases[i].Version != actual[i].Version {
					t.Errorf("at index %d want %q, got %q", i, tc.outputIndexReleases[i].Version, actual[i].Version)
				}
			}
		})
	}
}

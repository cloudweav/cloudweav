package version

import (
	"fmt"
	"strings"

	semverv3 "github.com/Masterminds/semver/v3"
)

type CloudweavVersion struct {
	version      *semverv3.Version
	rawVersion   string
	isDev        bool
	isPrerelease bool
}

func NewCloudweavVersion(versionStr string) (*CloudweavVersion, error) {
	var (
		v                         *semverv3.Version
		isDev, isReleaseCandidate bool
	)

	trimmedVersionStr := strings.TrimPrefix(versionStr, "v")

	// CloudweavVersion accepts any form of version strings except the empty string.
	if trimmedVersionStr == "" {
		return nil, ErrInvalidVersion
	}

	// Version strings other than semantic versions will be treated as dev versions.
	// For example, v1.2-head, f024f49a, etc. are dev versions.
	v, err := semverv3.StrictNewVersion(trimmedVersionStr)
	if err != nil {
		isDev = true
	}

	if !isDev && v.Prerelease() != "" {
		isReleaseCandidate = true
	}

	return &CloudweavVersion{
		version:      v,
		rawVersion:   versionStr,
		isDev:        isDev,
		isPrerelease: isReleaseCandidate,
	}, nil
}

// GetStableVersion returns the Version object without the suffix. For example, given "v1.2.2-rc1", it returns "v1.2.2"
func (v *CloudweavVersion) GetStableVersion() *CloudweavVersion {
	sv := semverv3.New(v.version.Major(), v.version.Minor(), v.version.Patch(), "", "")
	return &CloudweavVersion{
		version:      sv,
		rawVersion:   sv.Original(),
		isDev:        false,
		isPrerelease: false,
	}
}

func (v *CloudweavVersion) IsNewer(version *CloudweavVersion) (bool, error) {
	if version == nil {
		return false, ErrInvalidCloudweavVersion
	}

	if v.isDev || version.isDev {
		return false, ErrIncomparableVersion
	}

	var constraint string

	if v.isPrerelease && !version.isPrerelease {
		constraint = fmt.Sprintf("> %s-z", version.rawVersion)
	} else {
		constraint = fmt.Sprintf("> %s", version.rawVersion)
	}

	c, err := semverv3.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	if v.version == nil {
		return false, ErrInvalidVersion
	}

	return c.Check(v.version), nil
}

func (v *CloudweavVersion) IsEqual(version *CloudweavVersion) (bool, error) {
	if version == nil {
		return false, ErrInvalidCloudweavVersion
	}

	if v.isDev || version.isDev {
		if v.rawVersion == version.rawVersion {
			return true, nil
		}
		return false, ErrIncomparableVersion
	}

	constraint := version.rawVersion

	c, err := semverv3.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	if v.version == nil {
		return false, ErrInvalidVersion
	}

	return c.Check(v.version), nil
}

func (v *CloudweavVersion) IsOlder(version *CloudweavVersion) (bool, error) {
	if version == nil {
		return false, ErrInvalidCloudweavVersion
	}

	if v.isDev || version.isDev {
		return false, ErrIncomparableVersion
	}

	isNewer, err := v.IsNewer(version)
	if err != nil {
		return false, err
	}
	isEqual, err := v.IsEqual(version)
	if err != nil {
		return false, err
	}

	return !(isNewer || isEqual), nil
}

func (v *CloudweavVersion) String() string {
	if v.isDev {
		return v.rawVersion
	}

	return fmt.Sprintf("v%s", v.version.String())
}

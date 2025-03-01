package version

type CloudweavUpgradeVersion struct {
	currentVersion       *CloudweavVersion
	upgradeVersion       *CloudweavVersion
	minUpgradableVersion *CloudweavVersion
}

func NewCloudweavUpgradeVersion(cv, uv, mv *CloudweavVersion) *CloudweavUpgradeVersion {
	return &CloudweavUpgradeVersion{
		currentVersion:       cv,
		upgradeVersion:       uv,
		minUpgradableVersion: mv,
	}
}

func (u *CloudweavUpgradeVersion) IsUpgrade() error {
	if u.currentVersion == nil || u.upgradeVersion == nil {
		return ErrInvalidCloudweavVersion
	}

	isDowngrade, err := u.currentVersion.IsNewer(u.upgradeVersion)
	if err != nil {
		return err
	}

	if isDowngrade {
		return ErrDowngrade
	}

	return nil
}

// IsUpgradable checks that whether the current version satisfies the minimum upgrade requirement.
func (u *CloudweavUpgradeVersion) IsUpgradable() error {
	// If minUpgradableVersion is nil means there's no restriction.
	if u.minUpgradableVersion == nil {
		return nil
	}

	isOlderThanMinUpgradableVersion, err := u.currentVersion.IsOlder(u.minUpgradableVersion)
	if err != nil {
		return err
	}

	if isOlderThanMinUpgradableVersion {
		return ErrMinUpgradeRequirement
	}

	return nil
}

func (u *CloudweavUpgradeVersion) CheckUpgradeEligibility(strictMode bool) error {
	// Upgrading to dev versions is always allowed
	if u.upgradeVersion.isDev {
		return nil
	}

	// Upgrading from dev versions is restricted if strict mode is enabled
	if u.currentVersion.isDev {
		if strictMode {
			return ErrDevUpgrade
		}
		return nil
	}

	// Same-version upgrade is always allowed
	isSameVersion, err := u.currentVersion.IsEqual(u.upgradeVersion)
	if err != nil {
		return err
	}
	if isSameVersion {
		return nil
	}

	// Check if it's effectively a downgrade
	if err := u.IsUpgrade(); err != nil {
		return err
	}

	// Check the minimum upgradable version
	if err := u.IsUpgradable(); err != nil {
		return err
	}

	// Check if it's a prerelease cross-version upgrade
	if u.currentVersion.isPrerelease {
		currentStableVersion := u.currentVersion.GetStableVersion()
		upgradeStableVersion := u.upgradeVersion.GetStableVersion()
		if isSameStableVersion, _ := currentStableVersion.IsEqual(upgradeStableVersion); !isSameStableVersion {
			return ErrPrereleaseCrossVersionUpgrade
		}
	}

	return nil
}

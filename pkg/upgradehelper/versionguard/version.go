package versionguard

import (
	"errors"

	"github.com/cloudweav/go-common/version"
	"github.com/sirupsen/logrus"

	"github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
)

func Check(upgrade *v1beta1.Upgrade, strictMode bool, minUpgradableVersionStr string) error {

	repoInfo, err := getRepoInfo(upgrade)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"namespace": upgrade.Namespace,
			"name":      upgrade.Name,
		}).Error("failed to retrieve repo info")
		return err
	}

	upgradeVersion, err := version.NewCloudweavVersion(repoInfo.Release.Cloudweav)
	if err != nil {
		return err
	}

	currentVersion, err := version.NewCloudweavVersion(upgrade.Status.PreviousVersion)
	if err != nil {
		return err
	}

	var minUpgradableVersion *version.CloudweavVersion
	if minUpgradableVersionStr != "" {
		minUpgradableVersion, err = version.NewCloudweavVersion(minUpgradableVersionStr)
		if err != nil {
			return err
		}
	} else {
		minUpgradableVersion, err = version.NewCloudweavVersion(repoInfo.Release.MinUpgradableVersion)
		// When the error is ErrInvalidVersion, let the nil minUpgradableVersion slip through the check since it's a
		// valid scenario. It implies "upgrade with no restrictions."
		if err != nil && !errors.Is(err, version.ErrInvalidVersion) {
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"namespace":            upgrade.Namespace,
		"name":                 upgrade.Name,
		"currentVersion":       currentVersion,
		"upgradeVersion":       upgradeVersion,
		"minUpgradableVersion": minUpgradableVersion,
	}).Info("upgrade eligibility check")

	cloudweavUpgradeVersion := version.NewCloudweavUpgradeVersion(currentVersion, upgradeVersion, minUpgradableVersion)

	return cloudweavUpgradeVersion.CheckUpgradeEligibility(strictMode)
}

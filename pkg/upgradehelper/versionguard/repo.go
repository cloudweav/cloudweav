package versionguard

import (
	"github.com/cloudweav/cloudweav/pkg/apis/cloudweavhci.io/v1beta1"
	"github.com/cloudweav/cloudweav/pkg/controller/master/upgrade/repoinfo"
)

func getRepoInfo(upgrade *v1beta1.Upgrade) (*repoinfo.RepoInfo, error) {
	repoInfo := &repoinfo.RepoInfo{}
	if err := repoInfo.Load(upgrade.Status.RepoInfo); err != nil {
		return nil, err
	}
	return repoInfo, nil
}

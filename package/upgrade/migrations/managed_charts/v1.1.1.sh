#!/bin/bash -ex

CHART_NAME=$1
CHART_MANIFEST=$2

patch_cloudweav_ignore_default_sc()
{
	# add ignoring resources when upgrading to match the pr (https://github.com/cloudweav/cloudweav-installer/pull/481)
	yq e '.spec.diff.comparePatches += [{"apiVersion": "storage.k8s.io/v1", "kind": "StorageClass", "name": "cloudweav-longhorn", "jsonPointers":["/metadata/annotations"]}]' $CHART_MANIFEST -i
}

case $CHART_NAME in
  cloudweav)
    patch_cloudweav_ignore_default_sc
    ;;
esac

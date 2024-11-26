# Cloudweav-Network-Controller Helm Chart

[Cloudweav Network Contrller](https://github.com/cloudweav/network-controller-cloudweav) is a network controller that helps to manage the host network configuration of the Cloudweav cluster.

Introduction
------------

This chart installs the network-controller daemonset on the [cloudweav](https://github.com/cloudweav/cloudweav) cluster using the [Helm](https://helm.sh) package manager.

Prerequisites
-------------
- [multus-cni](https://github.com/intel/multus-cni) v3.6+
- Vlan filtering support on bridge
- Switch to support `trunk` mode

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm repo add cloudweav https://charts.cloudweavhci.io
helm repo update
helm install -n $namespace my-release cloudweav/cloudweav-network-controller
```

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
helm uninstall -n $namespace my-release
```

package utils

import "github.com/cloudweav/cloudweav-network-controller/pkg/apis/network.cloudweavhci.io"

const (
	KeyNetworkConf         = network.GroupName + "/route"
	KeyVlanLabel           = network.GroupName + "/vlan-id"
	KeyVlanConfigLabel     = network.GroupName + "/vlanconfig"
	KeyClusterNetworkLabel = network.GroupName + "/clusternetwork"
	KeyNodeLabel           = network.GroupName + "/node"
	KeyNetworkType         = network.GroupName + "/type"

	KeyMatchedNodes = network.GroupName + "/matched-nodes"

	ValueTrue = "true"
)

# Cloudweav Seeder Embedded Mode Integration

## Summary

We have been running and developing [Seeder](https://github.com/cloudweav/seeder) for provisioning and managing the lifecycle of bare-metal clusters in our internal labs.

With version 1.2.0, we plan to introduce an embedded mode of Seeder, which will allow deployment of Seeder to an existing cluster via **Cloudweav Addons**.

In embedded mode, the addon will enable end-users to define inventory objects that map the Cloudweav nodes to the corresponding bare-metal nodes. The bare-metal interaction will be through IPMI and Redfish protocols.

Once defined, Seeder in embedded mode leverages the cluster event controllers to query underlying hardware and propagate that information to the corresponding node objects.

In addition, users will be able to perform underlying hardware power cycle operations via the Cloudweav UI.


### Related Issues

https://github.com/cloudweav/cloudweav/issues/2318

## Motivation

### Goals

- Allow propagation of hardware information into Cloudweav.

- Allow propagation of hardware events into Cloudweav events.

- Allow users to power cycle hardware via Cloudweav.

### Non-goals [optional]

Provisioning of additional hardware.

## Proposal

### User Stories

#### Discover hardware details via Cloudweav

A Cloudweav user wants to discover underlying hardware information from Cloudweav.

Discovering underlying hardware information from Cloudweav is not currently possible. However, once the seeder addon is enabled, the user can define inventory objects for their Cloudweav nodes to map individual Cloudweav nodes to the underlying hardware.

Once the addon is enabled and the inventory object is defined, the end user can see underlying hardware information in the Cloudweav UI.

This will include information such as (but not limited to):
* Hardware manufacture
* Hardware serial number
* Hardware model
* Hardware events

![](./20230327-cloudweav-seeder-integration/hardware-discovery.png)

![](./20230327-cloudweav-seeder-integration/events.png)


#### Power cycle nodes via Cloudweav

A Cloudweav user wishes to reboot, shutdown, or power-on underlying hardware from Cloudweav.

Once the Seeder addon is enabled, users can power cycle the underlying nodes by using the additional power action options available in the Cloudweav UI.

### API changes
Like PCIDeviceClaims, the Cloudweav UI will allow end users to define a secret for BMC credentials and associated inventory objects.

The two changes explained above can be handled directly by the Cloudweav UI.

Once an inventory is associated with a node, the associated nodeObject needs to be annotated with additional information:

`cloudweavhci.io/inventoryName: inventoryName`

There will be changes to the node API: https://github.com/cloudweav/cloudweav/blob/master/pkg/api/node/formatter.go#L40

Changes will involve additional resource actions:

* If the node is in `maintenanceMode`, users can invoke `powerActionPossible` to check if power actions are possible on this node.
* If an API call `powerActionPossible` returns `HTTP status 200`, the user can invoke `powerAction` to trigger the correct `powerAction` on the node.


## Design
The Seeder addon will introduce a few new CRDs. However, only the following are used in embedded mode:

* inventories.metal.cloudweavhci.io
* clusters.metal.cloudweavhci.io
* jobs.bmc.tinkerbell.org
* machines.bmc.tinkerbell.org
* tasks.bmc.tinkerbell.org

### Implementation Overview

Seeder will run as a deployment in the `cloudweav-system` namespace in the cluster.

The addon will deploy Seeder with `EMBEDDED_MODE` set to `true`.

In this mode, the provisioning controllers are not bootstrapped. However, the following controllers are bootstrapped:

* local cluster controller
* local node contorller
* rufio machine controller
* rufio job controller
* rufio task controller
* cluster event controller

When booting in embedded mode, Seeder will create a `local` cluster for `cluster.metal` objects pointing to the K8s default service as the cluster endpoint address.

This is a placeholder for `inventory` objects as they are added to the cluster.

```yaml
apiVersion: metal.cloudweavhci.io/v1alpha1
kind: Cluster
metadata:
  creationTimestamp: "2023-03-26T22:50:15Z"
  generation: 4
  name: local
  namespace: cloudweav-system
  resourceVersion: "31519751"
  uid: 44ad4855-65c7-46ba-a321-04c54ab69337
spec:
  clusterConfig: {}
  nodes:
  - addressPoolReference:
      name: ""
      namespace: ""
    inventoryReference:
      name: cloudweav-659jw
      namespace: cloudweav-system
  version: local
  vipConfig:
    addressPoolReference:
      name: ""
      namespace: ""
status:
  clusterAddress: 10.53.0.1
  status: clusterRunning
```

The Cloudweav UI will allow users to create secrets and inventory definition corresponding to existing nodes.

sample secret creation call to endpoint: `${ENDPOINT}/v1/cloudweav/secrets`
```json
{
  "type": "Opaque",
  "metadata": {
    "namespace": "cloudweav-system",
    "name": "cloudweav-659jw"
  },
  "_type": "Opaque",
  "data": {
    "password": "TjBJTE80dSEh",
    "username": "RXRob3M="
  }
}
```

sample inventory creation call to endpoint: `${ENDPOINT}/v1/cloudweav/metal.cloudweavhci.io.inventories`

```json
{
  "kind": "Inventory",
  "metadata": {
    "name": "cloudweav-659jw",
    "namespace": "cloudweav-system",
    "annotations": {
      "metal.cloudweavhci.io/localInventory": "true",
      "metal.cloudweavhci.io/localNodeName": "cloudweav-659jw"
    }
  },
  "spec": {
    "baseboardSpec": {
      "connection": {
        "authSecretRef": {
          "name": "cloudweav-659jw",
          "namespace": "cloudweav-system"
        },
        "host": "172.19.1.172",
        "insecureTLS": true,
        "port": 623
      }
    },
    "events": {
      "enabled": true,
      "pollingInterval": "1h"
    }
  }
}
```

Since no actual provisioning is performed by Seeder, the node name the inventory is associated with is passed via the annotation `metal.cloudweavhci.io/localNodeName`. The local cluster controller in Seeder will use this annotation to query the k8s node object and generate the correct inventory status.

```json
{
  "ownerCluster": {
    "name": "local",
    "namespace": "cloudweav-system"
  },
  "pxeBootConfig": {
    "address": "172.19.109.10",
    "gateway": "",
    "netmask": ""
  },
  "status": "inventoryNodeReady",
  "generatedPassword": "",
  "hardwareID": ""
}

```

Once the inventory object is created, the `local cluster controller`, will add the inventory to the `local` cluster object in the `cloudweav-system` namespace.

The `cluster event controller` will now regularly reconcile inventory objects, query underlying hardware for hardware details and events, and trigger updates to the Cloudweav nodes.

The `inventory controller` will also watch nodes for power action requests via updates to the `powerActionStatus` in the status resource for the inventory object.

Once a power action is complete, the associated `powerStatus` field, `LastPowerAction`, will be updated with the associated status.

### Test plan

* Enable Seeder addon integration.
* Define valid inventory and secrets for node.
* Wait for additional node labels to be propagated on the node.
* Wait for additional hardware events to be generated for the node.
* Perform node power actions via Cloudweav UI.


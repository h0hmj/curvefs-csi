# Curvefs CSI Driver

## Introduction

The Curvefs CSI Driver implements the CSI specification for container orchestrators to manager Curvefs File Systems.

## Prerequisites

- Kubernetes 1.14+

## CSI Interface Implemented

- ControllerServer: CreateVolume, DeleteVolume, ValidateVolumeCapabilities
- NodeServer: NodePublishVolume, NodeUnpublishVolume, NodeGetInfo, NodeGetCapabilities
- IdentityServer: GetPluginInfo, Probe, GetPluginCapabilities

## How to use

### via kubectl

1. add label to node
    ```bash
    kubectl label node <nodename> curvefs-csi-controller=enabled
    kubectl label node <nodename> curvefs-csi-node=enabled
    ```
2. deploy csi driver
    ```bash
    kubectl apply -f deploy/csi-driver.yaml
    kubectl apply -f deploy/csi-rbac.yaml
    kubectl apply -f deploy/csi-controller-deployment.yaml
    kubectl apply -f deploy/csi-node-daemonset.yaml
    ```
3. create storage class
   ```bash
   # copy and fill in the blanks in storageclass-default.yaml
   kubectl apply -f /path/to/sc.yaml
   ```
4. now you can use this storageclass to create pvc and bind to a pod

## Build Status

| Curvefs CSI Driver Version | Curvefs Version Commit                    | Curvefs CSI Driver Image                          |
|----------------------------|-------------------------------------------|---------------------------------------------------|
| v1.0.0                     | c93730d4d0ba0e4b7843e30e384dd44640e4d547  | harbor.cloud.netease.com/curve/curvefs:csi-v1.0.0 |

## Follow-up Work

- more create/mount options support (require future curvefs support)
- move sensitive info like s3 ak/sk to secret
- subpath mount support
- resource limitation of single mount point
- quota(bytes) support
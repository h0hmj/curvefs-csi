# Curvefs CSI Driver

## Introduction

The Curvefs CSI Driver implements the CSI specification for container orchestrators to manager Curvefs File Systems.

## Prerequisites

- Kubernetes 1.18+

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
   attention: if you want to enable DiskCache, read the related section below
3. create storage class and pvc
   ```bash
   # copy and fill in the blanks in storageclass-default.yaml
   kubectl apply -f storageclass.yaml
   # copy and modify the pvc-default.yaml
   kubectl apply -f pvc.yaml
   ```
4. now you can bind this pvc to a pod

#### DiskCache related

what is DiskCache? A disk based cache used by client to increase the io performance
of client.

If you want to enable it:
1. check out content in csi-node-daemonset-enable-cache.yaml to bind the cache dir on curvefs-csi-node to pod's /curvefs/client/data/cache
2. add "diskCache.diskCacheType=2" or "diskCache.diskCacheType=1" to your mountOptions section of storageclass.yaml, 2 for read and write, 1 for read

Know Issue:

With discache enabled (type=2, write), metadata in metadatasever will be newer than data in s3 storage,
if the csi node pod crash but write cache is not fully uploaded to s3 storage,
you may lose this part of data. Remount will crash, because you only have meta but without data (haven't been flushed to s3).


## Build Status

| Curvefs CSI Driver Version | Curvefs Version | Curvefs CSI Driver Image                          |
|----------------------------|-----------------|---------------------------------------------------|
| v1.0.0                     | v2.3.0-rc0      | harbor.cloud.netease.com/curve/curvefs:csi-v1.0.0 |

## Follow-up Work

- more create/mount options support (require future curvefs support)
- move sensitive info like s3 ak/sk to secret
- subpath mount support (require future curvefs support)
- move every mount into a seperate pod, inspired by juicefs-csi
- quota(bytes) support (require future curvefs support)

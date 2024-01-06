# node-shell

A `kubectl` plugin to access k8s node. It starts a pod with mounted host filesystem and escalated privileges

## Intro

This plugin is an alternative to ```kubectl debug``` to run a pre-configured pod. Currently ```kubectl debug``` creates a debug pod with a container mounted on the root of the host file system. The following options are also set for this container:

```
hostIPC: true
hostNetwork: true
hostPID: true
securityContext: {}
```

However, ```kubectl debug``` currently lacks the flexibility to customize the pod/container that is launched. For example, we cannot specify the capabilities we are interested in. ```Node-shell``` plugin allows you to configure some options for launching a debug pod

## Installation

Manual installation via [krew](https://github.com/kubernetes-sigs/krew):

1) Download the [latest release archive](https://github.com/Tinkoff/node-shell/releases) for your os/arch via browser
2) Go to repo root directory
3) Install ```kubectl krew install --manifest=deploy/krew/plugin.yaml --archive=<ABSOLUTE-path-to-downloaded-archive>```


## Quick Start

```
kubectl node-shell --debug-image=<debug-image> --nodename=<node-to-run-on>
```

## Usage

```
kubectl node-shell [--debug-image <DEBUG_IMAGE>] [-n <NAMESPACE_NAME>] [--nodename <NODE_NAME>]  [--cpu <CPU>] [--mem <MEM>] [--ips <IPS>] [--caps <CAPS>]

DEBUG_IMAGE: Required. Image to run as debug container. No default value
NAMESPACE_NAME: Optional. Namespace to run debug pod in. Default: default
NODE_NAME: Optional. Name of a cluster node to run debug pod on. If ommited, pod will run on node determined by scheduler
CPU: Optional. Set cpu request and limit for debug container equal to its value. Default: 500m
MEM: Optional. Set memory request and limit for debug container equal to its value. Default: 256Mi
IPS: Optional. Set ImagePullSecrets for pod. Required for image pull authorization. No Default value
CAPS: Optional. List of POSIX capabilities for container separated by comma. Default: "NET_ADMIN", "SYS_ADMIN", "SYS_PTRACE"
```

## Examples

### Run on specific node 

```kubectl node-shell --debug-image="nikolaka/netshoot" --nodename="k8s-worker-1"```

### Run on specific node and specific namespace

```kubectl node-shell --debug-image="nikolaka/netshoot" --nodename="k8s-worker-1" -n="kube-system```

### Run with configured limits/requests

```kubectl node-shell --debug-image="nikolaka/netshoot" --nodename="k8s-worker-1" --cpu="50m" -n="kube-system ``` 


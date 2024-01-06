package plugin

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	debuggerName = "node-shell"
)

var t = true

var userToRun = int64(0)

var hostDirType = v1.HostPathDirectory
var hostSocketType = v1.HostPathSocket

var podMeta = metav1.ObjectMeta{
	Name:   debuggerName,
	Labels: map[string]string{"app": debuggerName},
}

var podSpec = v1.PodSpec{
	AutomountServiceAccountToken: &t,
	HostNetwork:                  true,
	HostPID:                      true,
	HostIPC:                      true,
	Containers:                   []v1.Container{},
	Tolerations:                  []v1.Toleration{},
	Volumes:                      []v1.Volume{},
}

var debugContainer = v1.Container{
	Name:            debuggerName,
	Stdin:           true,
	TTY:             true,
	ImagePullPolicy: "IfNotPresent",
	SecurityContext: &v1.SecurityContext{},
}

var volumeMounts = []v1.VolumeMount{
	{
		Name:      "host-root",
		MountPath: "/host",
	},
	{
		Name:      "cri-socket",
		MountPath: "/var/run/containerd/containerd.sock",
	},
	{
		Name:      "cgroup",
		MountPath: "/sys/fs/cgroup",
	},
}

var debuggerProbe = v1.Probe{
	ProbeHandler: v1.ProbeHandler{
		Exec: &v1.ExecAction{
			Command: []string{"/bin/true"},
		},
	},
}

var podTolerations = []v1.Toleration{
	{
		Key:      "CriticalAddonsOnly",
		Operator: "Exists",
	},
	{
		Key:      "NoExecute",
		Operator: "Exists",
	},
}

var podVolumes = []v1.Volume{
	{
		Name: "host-root",
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: "/",
				Type: &hostDirType,
			},
		},
	},
	{
		Name: "cri-socket",
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: "/var/run/containerd/containerd.sock",
				Type: &hostSocketType,
			},
		},
	},
	{
		Name: "cgroup",
		VolumeSource: v1.VolumeSource{
			HostPath: &v1.HostPathVolumeSource{
				Path: "/sys/fs/cgroup",
				Type: &hostDirType,
			},
		},
	},
}

func (p *Plugin) getDebugPod(n string) *v1.Pod {
	dp := v1.Pod{}
	podMeta.Namespace = n
	dp.ObjectMeta = podMeta
	dp.Spec = podSpec
	dp.Spec.ImagePullSecrets = append(dp.Spec.ImagePullSecrets, v1.LocalObjectReference{
		Name: p.Ips,
	})
	dp.Spec.NodeName = p.NodeName
	dp.Spec.Tolerations = podTolerations
	dp.Spec.Volumes = podVolumes
	debugContainer.Image = p.DebugImage
	resources := map[v1.ResourceName]resource.Quantity{v1.ResourceCPU: resource.MustParse(p.Cpu), v1.ResourceMemory: resource.MustParse(p.Mem)}
	debugContainer.Resources = v1.ResourceRequirements{
		Limits:   resources,
		Requests: resources,
	}
	debugContainer.VolumeMounts = volumeMounts
	debugContainer.ReadinessProbe = &debuggerProbe
	debugContainer.LivenessProbe = &debuggerProbe
	debugContainer.SecurityContext.Capabilities = p.getContainerCaps()
	debugContainer.SecurityContext.RunAsUser = &userToRun
	dp.Spec.Containers = append(dp.Spec.Containers, debugContainer)

	return &dp
}

func (p *Plugin) getContainerCaps() *v1.Capabilities {
	cs := v1.Capabilities{}
	for _, c := range p.Caps {
		cs.Add = append(cs.Add, v1.Capability(c))
	}
	return &cs
}

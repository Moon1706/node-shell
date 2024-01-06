package plugin

import (
	"context"
	"fmt"

	runshell "github.com/Moon1706/node-shell/pkg/utils"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

type Plugin struct {
	DebugImage string
	NodeName   string
	Ips        string
	Cpu        string
	Mem        string
	Caps       []string
}

func NewPlugin(pluginFlags *pflag.FlagSet) (*Plugin, error) {
	p := Plugin{}
	var err error
	p.DebugImage, err = pluginFlags.GetString("debug-image")
	if err != nil {
		return &p, err
	}
	p.Ips, err = pluginFlags.GetString("ips")
	if err != nil {
		return &p, err
	}
	p.NodeName, err = pluginFlags.GetString("nodename")
	if err != nil {
		return &p, err
	}
	p.Caps, err = pluginFlags.GetStringSlice("caps")
	if err != nil {
		return &p, err
	}
	p.Cpu, err = pluginFlags.GetString("cpu")
	if err != nil {
		return &p, err
	}
	p.Mem, err = pluginFlags.GetString("mem")
	if err != nil {
		return &p, err
	}

	return &p, nil
}

func (p *Plugin) RunPlugin(configFlags *genericclioptions.ConfigFlags) error {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("failed to read kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	if *configFlags.Namespace == "" {
		*configFlags.Namespace = "default"
	}

	namespace, err := clientset.CoreV1().Namespaces().Get(context.TODO(), *configFlags.Namespace, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}

	podObj := p.getDebugPod(namespace.Name)

	debugPod, err := clientset.CoreV1().Pods(namespace.Name).Create(context.TODO(), podObj, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}

	podLabelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": debuggerName}}

	podReady, err := runshell.WaitForPodReadiness(clientset, namespace.Name, podLabelSelector)
	if err != nil {
		return err
	}
	if podReady {
		err = runshell.ExecIntoPod(clientset, config, namespace.Name, debugPod)
		if err != nil {
			return err
		}
	}

	return nil
}

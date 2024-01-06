package runshell

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/cli/cli/streams"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

func WaitForPodReadiness(c *kubernetes.Clientset, namespace string, selector metav1.LabelSelector) (bool, error) {
	watchTimeout := int64(60)
	watcher, err := c.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{
		LabelSelector:  labels.Set(selector.MatchLabels).String(),
		Watch:          true,
		TimeoutSeconds: &watchTimeout,
	})
	if err != nil {
		return false, err
	}
	for event := range watcher.ResultChan() {
		pod := event.Object.(*v1.Pod)
		for _, s := range pod.Status.Conditions {
			if s.Type == v1.PodReady && s.Status == v1.ConditionTrue {
				fmt.Printf("debug pod %v is ready\n", pod.ObjectMeta.Name)
				watcher.Stop()
				return true, nil
			}
		}
		fmt.Printf("waiting for pod %v to be ready\n", pod.ObjectMeta.Name)
	}
	return false, fmt.Errorf("debug pod is not ready after %d seconds", watchTimeout)
}

func ExecIntoPod(c *kubernetes.Clientset, conf *rest.Config, ns string, p *v1.Pod) error {
	request := c.CoreV1().RESTClient().
		Post().
		Namespace(ns).
		Resource("pods").
		Name(p.Name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: []string{"chroot", "/host"},
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(conf, "POST", request.URL())
	if err != nil {
		return err
	}
	in := streams.NewIn(os.Stdin)
	if err := in.SetRawTerminal(); err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  in,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err == os.ErrProcessDone {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%w failed getting shell on %v/%v", err, p.Namespace, p.Name)
	}
	return err

}

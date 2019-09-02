package podinfo

import (
	"errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

// http://localhost:8001/api/v1/pods

func GetPodInfo(podname string) (*corev1.Pod, error) {
	// creates the in-cluster config
	var pod1 *corev1.Pod
	var err error
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// find namespace for pod
	PL, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range PL.Items {
		if pod.Name == podname {
			pod1 = &pod
			break
		}
	}

	// wait pod to be ready
	if pod1 != nil {
		for i := 0; i < 3; i++ {
			for _, cond := range pod1.Status.Conditions {
				//	log.Printf("Pod status: %v = %v\n", string(cond.Type), string(cond.Status))
				if cond.Type == corev1.PodReady &&
					cond.Status == corev1.ConditionTrue {
					return pod1, nil
				}
			}

			time.Sleep(5 * time.Second)

			pod1, err = clientset.CoreV1().Pods(pod1.Namespace).Get(podname, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}
		}
		return pod1, nil
	}
	return nil, errors.New(podname + " not found.")
}

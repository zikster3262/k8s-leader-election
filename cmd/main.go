package main

import (
	"context"
	"flag"
	"os"

	leader "github.com/zikster3262/k8s-leader-election/pkg/leader"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

var (
	client *clientset.Clientset
)

func main() {
	var (
		leaseLockName      string
		leaseLockNamespace string
		podName            = os.Getenv("POD_NAME")
	)
	flag.StringVar(&leaseLockName, "lease-name", "", "Name of lease lock")
	flag.StringVar(&leaseLockNamespace, "lease-namespace", "default", "Name of lease lock namespace")
	flag.Parse()

	if leaseLockName == "" {
		klog.Fatal("missing lease-name flag")
	}
	if leaseLockNamespace == "" {
		klog.Fatal("missing lease-namespace flag")
	}

	config, err := rest.InClusterConfig()
	client = clientset.NewForConfigOrDie(config)

	if err != nil {
		klog.Fatalf("failed to get kubeconfig")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lock := leader.GetNewLock(leaseLockName, podName, leaseLockNamespace)
	leader.RunLeaderElection(lock, ctx, podName)
}

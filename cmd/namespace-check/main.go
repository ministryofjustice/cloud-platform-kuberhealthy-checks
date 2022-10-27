package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/kuberhealthy/kuberhealthy/v2/pkg/checks/external/checkclient"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	checkclient.Debug = true
}

func main() {
	// K8s config file for the client.
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	// We have to explicitly list of namespaces that we want to look for
	namespaces := []string{
		"cert-manager",
		"default",
		"ingress-controllers",
		"kube-system",
		"logging",
		"monitoring",
		"opa",
		"velero",
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := doExpectedNamespacesExist(context.Background(), clientset, namespaces); err != nil {
		reportErr := checkclient.ReportFailure([]string{"Namespace check failed:" + err.Error()})
		if reportErr != nil {
			log.Fatalln("Unable to communicate with kuberhealthy", reportErr.Error())
		}
		log.Fatalln("Error checking for namespaces", err)
	}

	// report success to Kuberhealthy. If it fails, fail the check.
	reportErr := checkclient.ReportSuccess()
	if reportErr != nil {
		log.Fatalln("error reporting to kuberhealthy:", err.Error())
	}
}

// doExpectedNamespacesExist checks if the expected namespaces exist in the cluster.
func doExpectedNamespacesExist(ctx context.Context, client kubernetes.Interface, expectedNamespaces []string) error {
	var missing []string
	for _, ns := range expectedNamespaces {
		if checkclient.Debug {
			log.Println("Checking for namespace", ns)
		}
		_, err := client.CoreV1().Namespaces().Get(ctx, ns, metav1.GetOptions{})
		if apierrors.IsNotFound(err) {
			missing = append(missing, ns)
		} else if err != nil {
			log.Println("Getting namespace from cluster failed:", err)
			return fmt.Errorf("failed getting namespace %s from cluster: %w", ns, err)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing namespaces: %s", strings.Join(missing, ", "))
	}
	return nil
}

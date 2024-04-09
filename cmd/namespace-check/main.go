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
	"k8s.io/utils/strings/slices"

	"github.com/kuberhealthy/kuberhealthy/v2/pkg/checks/external/checkclient"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	checkclient.Debug = true
}

func main() {
	// K8s config file for the client.
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Unable to get user home directory", err.Error())
	}

	kubeconfig := filepath.Join(userHomeDir, ".kube", "config")
	if kubeconfig == "" {
		log.Fatalln("kubeconfig: No kubeconfig file found in $HOME/.kube/config")
	}

	// We have to explicitly list of namespaces that we want to look for
	namespaces := []string{
		"calico-apiserver",
		"calico-system",
		"cert-manager",
		"default",
		"external-secrets-operator",
		"gatekeeper-system",
		"ingress-controllers",
		"kube-system",
		"kuberos",
		"logging",
		"monitoring",
		"overprovision",
		"tigera-operator",
		"trivy-system",
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

	currentEnv := os.Getenv("CLUSTER_ENV")

	if err := doExpectedNamespacesExist(context.Background(), clientset, namespaces, currentEnv); err != nil {
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
func doExpectedNamespacesExist(ctx context.Context, client kubernetes.Interface, expectedNamespaces []string, currentEnv string) error {
	var missing []string
	prodEnvs := []string{"manager", "live-2", "live"}
	liveEnvs := []string{"live"}

	isProd := slices.Contains(prodEnvs, currentEnv)

	isLive := slices.Contains(liveEnvs, currentEnv)

	prodOnlyNamespaces := []string{
		"velero",
	}

	liveOnlyNamespaces := []string{
		"overprovision",
	}

	for _, ns := range expectedNamespaces {
		if !isProd && slices.Contains(prodOnlyNamespaces, ns) {
			fmt.Printf("skipping namespace %s because we are running in a non-prod cluster\n", ns)
			continue
		}

		if !isLive && slices.Contains(liveOnlyNamespaces, ns) {
			fmt.Printf("skipping namespace %s because we are running in a non-live cluster\n", ns)
			continue
		}

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

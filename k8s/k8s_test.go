package k8s

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testNamespace = "sched-ops-test"
)

var (
	commonLabels = map[string]string{
		"owner": "sched-ops-test",
	}
)

func TestK8s(t *testing.T) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		t.Skipf("KUBECONFIG not defined. Skipping tests")
	}

	testInstance := Instance()
	err := testInstance.initK8sClient()
	require.NoError(t, err, "initK8sClient returned error")

	ns, err := testInstance.CreateNamespace(testNamespace, commonLabels)
	require.NotNil(t, ns, "CreateNamespace returned nil namespace")
	require.NoError(t, err, "CreateNamespace returned error")

	t.Run("pvcTests", pvcTests)

	err = testInstance.DeleteNamespace(testNamespace)
	require.NoError(t, err, "DeleteNamespace returned error")
}

func pvcTests(t *testing.T) {

}

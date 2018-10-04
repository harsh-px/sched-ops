package k8s

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/api/core/v1"
	storage_api "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testNamespace = "sched-ops-test"
)

var (
	testInstance Ops
	trueVar      = true
	commonLabels = map[string]string{
		"owner": "sched-ops-test",
		"foo":   "bar",
	}
)

func TestK8s(t *testing.T) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		t.Skipf("KUBECONFIG not defined. Skipping tests")
	}

	testInstance = Instance()
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
	scName := "sched-ops-sc"
	scParams := map[string]string{
		"repl": "1",
	}

	scObj := &storage_api.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:   scName,
			Labels: commonLabels,
		},
		Provisioner:          "kubernetes.io/portworx-volume",
		AllowVolumeExpansion: &trueVar,
		Parameters:           scParams,
	}

	sc, err := testInstance.CreateStorageClass(scObj)
	require.NoError(t, err, "create sc returned err")
	require.NotNil(t, sc, "got empty sc")

	pvc1Name := "sched-ops-pvc1"
	createPVC(t, pvc1Name, scName)

	pvc2Name := "sched-ops-pvc2"
	createPVC(t, pvc2Name, scName)

	pvcs, err := testInstance.GetPersistentVolumeClaims(testNamespace, map[string]string{
		"name": pvc1Name,
		"foo":  "bar",
	})
	require.NoError(t, err, "get pvcs returned err")
	require.NotNil(t, pvcs, "got nil pvcs")
	require.Len(t, pvcs.Items, 1, "expected only one pvc")

	fmt.Printf("List of pvc1 returned: %v\n", pvcs.Items)

	pvcs, err = testInstance.GetPersistentVolumeClaims(testNamespace, map[string]string{
		"name": pvc2Name,
		"foo":  "bar",
	})
	require.NoError(t, err, "get pvcs returned err")
	require.NotNil(t, pvcs, "got nil pvcs")
	require.Len(t, pvcs.Items, 1, "expected only one pvc")
	fmt.Printf("List of pvcc returned: %v\n", pvcs.Items)

	// List all
	pvcs, err = testInstance.GetPersistentVolumeClaims(testNamespace, commonLabels)
	require.NoError(t, err, "get pvcs returned err")
	require.NotNil(t, pvcs, "got nil pvcs")
	require.Len(t, pvcs.Items, 1, "expected only one pvc")
	fmt.Printf("List of all pvcs returned: %v\n", pvcs.Items)

	err = testInstance.DeletePersistentVolumeClaim(pvc1Name, testNamespace)
	require.NoError(t, err, "delete pvc returned err")

	err = testInstance.DeletePersistentVolumeClaim(pvc2Name, testNamespace)
	require.NoError(t, err, "delete pvc returned err")

	err = testInstance.DeleteStorageClass(scName)
	require.NoError(t, err, "create sc returned err")
}

func createPVC(t *testing.T, pvcName, scName string) {
	pvcLabels := commonLabels
	pvcLabels["name"] = pvcName

	pvc1Obj := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: testNamespace,
			Labels:    pvcLabels,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse("1Gi"),
				},
			},
			StorageClassName: &scName,
		},
	}

	pvc, err := testInstance.CreatePersistentVolumeClaim(pvc1Obj)
	require.NoError(t, err, "create pvc returned err")
	require.NotNil(t, pvc, "got empty pvc")
}

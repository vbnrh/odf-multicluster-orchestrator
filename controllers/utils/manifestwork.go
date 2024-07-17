package utils

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	workv1 "open-cluster-management.io/api/work/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func CreateOrUpdateManifestWork(ctx context.Context, c client.Client, name string, namespace string, objJson []byte, ownerRef metav1.OwnerReference) (controllerutil.OperationResult, error) {
	mw := workv1.ManifestWork{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				ownerRef,
			},
		},
	}

	operationResult, err := controllerutil.CreateOrUpdate(ctx, c, &mw, func() error {
		mw.Spec = workv1.ManifestWorkSpec{
			Workload: workv1.ManifestsTemplate{
				Manifests: []workv1.Manifest{
					{
						RawExtension: runtime.RawExtension{
							Raw: objJson,
						},
					},
				},
			},
		}
		return nil
	})

	if err != nil {
		return operationResult, fmt.Errorf("failed to create and update ManifestWork %s for namespace %s. error %w", name, namespace, err)
	}

	return operationResult, nil
}

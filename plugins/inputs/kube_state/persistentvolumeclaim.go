package kube_state

import (
	"context"
	"strings"

	"github.com/ericchiang/k8s/apis/core/v1"

	"github.com/influxdata/telegraf"
)

func collectPersistentVolumeClaims(ctx context.Context, acc telegraf.Accumulator, ks *KubernetesState) {
	list, err := ks.client.getPersistentVolumeClaims(ctx)
	if err != nil {
		acc.AddError(err)
		return
	}
	for _, pvc := range list.Items {
		if err = ks.gatherPersistentVolumeClaim(*pvc, acc); err != nil {
			acc.AddError(err)
			return
		}
	}
}

func (ks *KubernetesState) gatherPersistentVolumeClaim(pvc v1.PersistentVolumeClaim, acc telegraf.Accumulator) error {
	phaseType := 3
	switch strings.ToLower(pvc.Status.GetPhase()) {
	case "bound":
		phaseType = 0
	case "lost":
		phaseType = 1
	case "pending":
		phaseType = 2
	}
	fields := map[string]interface{}{
		"phase_type": phaseType,
	}
	tags := map[string]string{
		"pvc_name":     pvc.Metadata.GetName(),
		"namespace":    pvc.Metadata.GetNamespace(),
		"phase":        pvc.Status.GetPhase(),
		"storageclass": pvc.Spec.GetStorageClassName(),
	}

	acc.AddFields(persistentVolumeClaimMeasurement, fields, tags)

	return nil
}

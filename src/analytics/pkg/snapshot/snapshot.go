package snapshot

import (
	"encoding/json"
	"io"

	corev1 "k8s.io/api/core/v1"
	schedulingv1 "k8s.io/api/scheduling/v1"
	storagev1 "k8s.io/api/storage/v1"
)

// ResourcesForSnap indicates all resources and scheduler configuration to be snapped.
type SimulatorSnapshot struct {
	Pods            []corev1.Pod                   `json:"pods"`
	Nodes           []corev1.Node                  `json:"nodes"`
	Pvs             []corev1.PersistentVolume      `json:"pvs"`
	Pvcs            []corev1.PersistentVolumeClaim `json:"pvcs"`
	StorageClasses  []storagev1.StorageClass       `json:"storageClasses"`
	PriorityClasses []schedulingv1.PriorityClass   `json:"priorityClasses"`
	Namespaces      []corev1.Namespace             `json:"namespaces"`
}

func FromReader(src io.Reader) (snapshot SimulatorSnapshot, err error) {
	decoder := json.NewDecoder(src)
	return snapshot, decoder.Decode(&snapshot)
}

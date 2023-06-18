package insights

import (
	"github.com/defaulterrr/scheduler-playground/pkg/snapshot"
	v1 "k8s.io/api/core/v1"
)

type NodeResourceRequested map[string]float64

func CalculateRequestsCPU(snap snapshot.SimulatorSnapshot, filter ...func(node *v1.Node) float64) NodeResourceRequested {
	nodeCPUUsage := make(NodeResourceRequested, len(snap.Nodes))
	nodeCPURequestedmCPUSTotal := make(map[string]int64, len(snap.Nodes))

	for _, pod := range snap.Pods {
		if pod.Spec.NodeName == "" {
			continue
		}

		nodeCPURequestedmCPUSTotal[pod.Spec.NodeName] = nodeCPURequestedmCPUSTotal[pod.Spec.NodeName] + totalPodCPURequestsmCPU(pod)
	}

	for _, node := range snap.Nodes {
		usage := float64(nodeCPURequestedmCPUSTotal[node.Name]) / float64(node.Status.Allocatable.Cpu().MilliValue())
		nodeCPUUsage[node.Name] = usage
	}

	return nodeCPUUsage
}

func CalculateRequestsCPUFromStaticSource(
	snap snapshot.SimulatorSnapshot,
	CPUsource AppPodCPURequests,
) NodeResourceRequested {
	nodeCPUUsage := make(NodeResourceRequested, len(snap.Nodes))
	nodeCPURequestedmCPUSTotal := make(map[string]int64, len(snap.Nodes))

	for _, pod := range snap.Pods {
		if pod.Spec.NodeName == "" {
			continue
		}

		nodeCPURequestedmCPUSTotal[pod.Spec.NodeName] = nodeCPURequestedmCPUSTotal[pod.Spec.NodeName] + totalPodCPURequestmCPUFromSource(pod, CPUsource)
	}

	for _, node := range snap.Nodes {
		usage := float64(nodeCPURequestedmCPUSTotal[node.Name]) / float64(node.Status.Allocatable.Cpu().MilliValue())
		nodeCPUUsage[node.Name] = usage
	}

	return nodeCPUUsage
}

func totalPodCPURequestsmCPU(pod v1.Pod) (total int64) {
	for _, container := range pod.Spec.Containers {
		total += container.Resources.Requests.Cpu().MilliValue()
	}
	return total
}

func totalPodCPURequestmCPUFromSource(pod v1.Pod, source AppPodCPURequests) int64 {
	app := ""

	if len(pod.Labels) != 0 {
		app = pod.Labels["app"]
	}

	return source.GetCPUForApp(app)
}

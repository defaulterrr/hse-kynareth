package insights

import (
	"github.com/defaulterrr/scheduler-playground/pkg/snapshot"
	v1 "k8s.io/api/core/v1"
)

type AppPodCPURequests map[string]int64

func (r AppPodCPURequests) GetCPUForApp(app string) int64 {
	v, ok := r[app]
	if !ok {
		v = 0
	}

	return v
}

func GetAppPodCPURequests(snap snapshot.SimulatorSnapshot) AppPodCPURequests {
	totalRequests := AppPodCPURequests{}

	for _, pod := range snap.Pods {
		totalRequests[AppFromPod(pod)] = totalPodCPURequestsmCPU(pod)
	}

	return totalRequests
}

func AppFromPod(pod v1.Pod) string {
	app := ""

	if len(pod.Labels) != 0 {
		app = pod.Labels["app"]
	}

	return app
}

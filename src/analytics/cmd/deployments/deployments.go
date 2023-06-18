package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"gitlab.ozon.ru/platform/errgroup/v2"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var kubeconfig *string
	var deploymentsDir *string

	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	deploymentsDir = flag.String("deploydir", "", "absolute path to the nodes manifests")
	flag.Parse()

	if *kubeconfig == "" {
		panic("-kubeconfig flag must be explicitly set")
	}

	if *deploymentsDir == "" {
		panic("-deploydir must be explicitly set")
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	config.QPS = 200
	config.Burst = 400

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	go func() {

		group, ctx := errgroup.WithContext(ctx)
		group.SetLimit(1)

		deploymentFiles := make([]string, 0, 8000)

		if err := filepath.Walk(*deploymentsDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(info.Name(), ".json") {
				return nil
			}

			deploymentFiles = append(deploymentFiles, path)

			return nil
		}); err != nil {
			panic(err)
		}

		rand.Seed(time.Now().Unix())
		rand.Shuffle(len(deploymentFiles), func(i, j int) { deploymentFiles[i], deploymentFiles[j] = deploymentFiles[j], deploymentFiles[i] })

		// twoCPUS, _ := resource.ParseQuantity("2")
		// fourCPUS, _ := resource.ParseQuantity("4")
		// fourRAMs, _ := resource.ParseQuantity("4G")
		// eightRAMs, _ := resource.ParseQuantity("8G")

		// resourceMap := map[corev1.ResourceName]resource.Quantity{
		// 	corev1.ResourceCPU:    twoCPUS,
		// 	corev1.ResourceMemory: fourRAMs,
		// }

		// limitsMap := map[corev1.ResourceName]resource.Quantity{
		// 	corev1.ResourceCPU:    fourCPUS,
		// 	corev1.ResourceMemory: eightRAMs,
		// }

		for i, path := range deploymentFiles {
			i := i

			group.Go(func() error {
				var newDeployment appsv1.Deployment
				manifest, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if err := json.Unmarshal(manifest, &newDeployment); err != nil {
					return err
				}

				// reset resource version
				newDeployment.ResourceVersion = ""
				newDeployment.Spec.Template.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = nil
				// newDeployment.Spec.Template.Spec.Containers = []corev1.Container{newDeployment.Spec.Template.Spec.Containers[0]}
				// newDeployment.Spec.Template.Spec.Containers[0].Resources.Requests = resourceMap
				// newDeployment.Spec.Template.Spec.Containers[0].Resources.Limits = limitsMap

				deployment := clientset.AppsV1().Deployments(newDeployment.Namespace)
				dplmnt, err := deployment.Create(ctx, &newDeployment, v1.CreateOptions{})
				if err != nil {
					// return err
					fmt.Println(err)
					// fmt.Println(err)
					// if errors.Is(context.Canceled, err) {
					// 	return err
					// }
				}
				fmt.Println(i, "out of \t", len(deploymentFiles), dplmnt.Name, "was created")
				time.Sleep(time.Millisecond * 10)

				return nil
			})
			// fmt.Scanln()
		}

		if err := group.Wait(); err != nil {
			panic(err)
		}

		cancel()
	}()

	<-ctx.Done()

	fmt.Println("Done!")
}

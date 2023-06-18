package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"o3.ru/svpetrov/scheduler-playground/pkg/insights"
	"o3.ru/svpetrov/scheduler-playground/pkg/snapshot"
)

var snapshotFile = pflag.String("snapshot", "snapshot.json", "filepath to the snapshot file in JSON format")
var metadataFile = pflag.String("out", "out.json", "filepath to the intermediate file to be consumed by python script")

func main() {
	pflag.Parse()

	pathToSnapshot := *snapshotFile
	if pathToSnapshot == "" {
		panic("--snapshot cannot be empty")
	}

	fmt.Println("reading and unmarshaling snapshot from ", pathToSnapshot)
	sourceFile, err := os.Open(pathToSnapshot)
	if err != nil {
		panic(err)
	}
	defer sourceFile.Close()
	snapshot, err := snapshot.FromReader(sourceFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("calculating usage")
	nodeCPUUsage := insights.CalculateRequestsCPU(snapshot)

	out, err := json.Marshal(nodeCPUUsage)
	if err != nil {
		panic(err)
	}

	filename := *metadataFile
	if err := os.WriteFile(filename, out, 0755); err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

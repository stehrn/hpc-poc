package main

import (
	"log"
	"math"

	"github.com/stehrn/hpc-poc/internal/utils"
)

const defaultTaskLoadFactor = 0.2
const defaultMaxPodsPerJob = 100

// init task load factor
func init() {
	taskLoadFactor = utils.EnvAsFloat("TASK_LOAD_FACTOR", defaultTaskLoadFactor)
	maxPodsPerJob = utils.EnvAsInt("MAX_PODS_PER_JOB", defaultMaxPodsPerJob)
}

func parallelism(numTasks int) int32 {
	if numTasks == 1 {
		return 1
	}
	parallelism := int32(math.Max(float64(numTasks)*float64(taskLoadFactor), 1.0))
	log.Printf("Parallelism set to %d, (numtasks * taskLoadFactor) = (%d * %f)", parallelism, numTasks, taskLoadFactor)
	return parallelism
}

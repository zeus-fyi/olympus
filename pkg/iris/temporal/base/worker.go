package temporal_base

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

/*
The Worker is where Workflow Functions and Activity Functions are executed.

Each Worker in the Worker Process must register the exact Workflow Types and Activity Types it may execute.
Each Worker Entity must also associate itself with exactly one Task QueueLink preview icon.
Each Worker Entity polling the same Task Queue must be registered with the same Workflow Types and Activity Types.
A Worker is the component within a Worker Process that listens to a specific Task Queue.

Although multiple Worker Entities can be in a single Worker Process, a single Worker Entity Worker Process
may be perfectly sufficient. For more information, see the Worker tuning guide.
*/
type Worker struct {
	worker.Worker
	TemporalClient
	WorkerOpts    worker.Options
	TaskQueueName string
	Workflows     []interface{}
	Activities    []interface{}
}

func NewWorker(taskQueueName string) Worker {
	return Worker{
		TaskQueueName: taskQueueName,
		WorkerOpts:    worker.Options{},
		Workflows:     []interface{}{},
		Activities:    []interface{}{},
	}
}

func (w *Worker) RegisterWorker(c client.Client) {
	w.Worker = worker.New(c, w.TaskQueueName, w.WorkerOpts)
	for _, workflow := range w.Workflows {
		w.Worker.RegisterWorkflow(workflow)
	}
	for _, activity := range w.Activities {
		w.Worker.RegisterActivity(activity)
	}
}

func (w *Worker) AddWorkflow(wf interface{}) {
	w.Workflows = append(w.Workflows, wf)
}

func (w *Worker) AddWorkflows(wfs []interface{}) {
	w.Workflows = append(w.Workflows, wfs...)
}

func (w *Worker) AddActivities(a []interface{}) {
	w.Activities = append(w.Activities, a...)
}

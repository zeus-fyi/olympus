package temporal_base

import "go.temporal.io/sdk/worker"

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
	TemporalClient

	worker.Worker
	TaskQueueName string
	Workflows     []Workflow
	Activities    []interface{}
}

func NewWorker(tc TemporalClient, taskQueueName string) Worker {
	return Worker{
		TemporalClient: tc,
		TaskQueueName:  taskQueueName,
	}
}

func (w *Worker) RegisterWorker(wrk Worker) error {
	err := w.Connect()
	if err != nil {
		return err
	}
	defer w.Close()

	wrk.Worker = worker.New(w.Client, wrk.TaskQueueName, worker.Options{})
	for _, workflow := range wrk.Workflows {
		wrk.Worker.RegisterWorkflow(workflow)
	}
	for _, activity := range wrk.Activities {
		wrk.Worker.RegisterActivity(activity)
	}
	return nil
}

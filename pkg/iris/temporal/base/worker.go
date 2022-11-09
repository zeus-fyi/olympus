package temporal_base

import (
	"github.com/rs/zerolog/log"
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
	TemporalClient

	worker.Worker
	WorkerOpts    worker.Options
	TaskQueueName string
	Workflows     []interface{}
	Activities    []interface{}
}

func NewWorker(tc TemporalClient, taskQueueName string) Worker {
	return Worker{
		TemporalClient: tc,
		TaskQueueName:  taskQueueName,
		WorkerOpts:     worker.Options{},
		Workflows:      []interface{}{},
		Activities:     []interface{}{},
	}
}

func (w *Worker) RegisterWorker() error {
	err := w.ConnectTemporalClient()
	if err != nil {
		log.Err(err).Msg("RegisterWorker: ConnectTemporalClient failed")
		return err
	}
	defer w.Close()
	w.Worker = worker.New(w.Client, w.TaskQueueName, w.WorkerOpts)
	for _, workflow := range w.Workflows {
		w.Worker.RegisterWorkflow(workflow)
	}
	for _, activity := range w.Activities {
		w.Worker.RegisterActivity(activity)
	}
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Err(err).Msg("NewTopologyWorker: Run worker failed")
	}
	return nil
}

func (w *Worker) AddWorkflow(wf interface{}) {
	w.Workflows = append(w.Workflows, wf)
}

func (w *Worker) AddActivities(a []interface{}) {
	w.Activities = append(w.Activities, a...)
}

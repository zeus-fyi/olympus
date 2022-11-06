package temporal_base

/*
The Worker is where Workflow Functions and Activity Functions are executed.

Each Worker in the Worker Process must register the exact Workflow Types and Activity Types
it may execute.

Each Worker Entity must also associate itself with exactly one Task Queue.

Each Worker Entity polling the same Task Queue must be registered with the same Workflow Types
and Activity Types.
*/

type WorkersGroup struct {
	WorkerSlice []Worker
}

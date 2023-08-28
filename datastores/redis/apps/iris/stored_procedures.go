package iris_redis

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog/log"
	iris_programmable_proxy_v1_beta "github.com/zeus-fyi/zeus/zeus/iris_programmable_proxy/v1beta"
)

func (m *IrisCache) SetStoredProcedure(ctx context.Context, orgID int, procedure iris_programmable_proxy_v1_beta.IrisRoutingProcedure) error {
	procedureKey := getProcedureKey(orgID, procedure.Name)
	procedureStepsKey := getProcedureStepsKey(orgID, procedure.Name)

	var steps []iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep
	for procedure.OrderedSteps.Len() > 0 {
		step := procedure.OrderedSteps.PopFront()
		steps = append(steps, step.(iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep))
	}

	// Serialize the procedure
	data, err := json.Marshal(procedure)
	if err != nil {
		log.Err(err).Msg("Failed to serialize the procedure")
		return err
	}

	// Serialize the procedure
	stepsData, err := json.Marshal(steps)
	if err != nil {
		log.Err(err).Msg("Failed to serialize the procedure")
		return err
	}

	// Use Redis transaction (pipeline) to perform all operations atomically
	pipe := m.Writer.TxPipeline()

	// Add serialized data to the pipeline
	pipe.Set(ctx, procedureKey, data, 0)
	pipe.Set(ctx, procedureStepsKey, stepsData, 0)

	// Execute the transaction
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("Failed to set stored procedure in Redis")
		return err
	}
	return nil
}

func (m *IrisCache) GetStoredProcedure(ctx context.Context, orgID int, procedureName string) (iris_programmable_proxy_v1_beta.IrisRoutingProcedure, error) {
	procedureKey := getProcedureKey(orgID, procedureName)
	procedureStepsKey := getProcedureStepsKey(orgID, procedureName)

	pipe := m.Reader.TxPipeline()

	// Get the values from Redis
	procedureCmd := pipe.Get(ctx, procedureKey)
	procedureStepsKeyCmd := pipe.Get(ctx, procedureStepsKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Err(err).Msgf("Failed to set stored procedure in Redis")
		return iris_programmable_proxy_v1_beta.IrisRoutingProcedure{}, err
	}
	data, err := procedureCmd.Bytes()
	if err != nil {
		log.Err(err).Msg("Failed to get procedure from Redis")
		return iris_programmable_proxy_v1_beta.IrisRoutingProcedure{}, err
	}
	// Deserialize the procedure
	var procedure iris_programmable_proxy_v1_beta.IrisRoutingProcedure
	err = json.Unmarshal(data, &procedure)
	if err != nil {
		log.Err(err).Msg("failed to deserialize the procedure")
		return iris_programmable_proxy_v1_beta.IrisRoutingProcedure{}, err
	}
	stepsBytes, err := procedureStepsKeyCmd.Bytes()
	if err != nil {
		log.Err(err).Msg("Failed to get procedure steps from Redis")
		return iris_programmable_proxy_v1_beta.IrisRoutingProcedure{}, err
	}
	var steps []iris_programmable_proxy_v1_beta.IrisRoutingProcedureStep
	err = json.Unmarshal(stepsBytes, &steps)
	for _, step := range steps {
		procedure.OrderedSteps.PushBack(step)
	}
	return procedure, nil
}

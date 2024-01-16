import * as React from 'react';
import {useEffect, useState} from 'react';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Container from '@mui/material/Container';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import {
    Card,
    CardActions,
    CardContent,
    Collapse,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack,
    Tab,
    Tabs,
    TextareaAutosize
} from "@mui/material";
import authProvider from "../../redux/auth/auth.actions";
import MainListItems from "../dashboard/listItems";
import {WorkflowTable} from "./WorkflowTable";
import {ZeusCopyright} from "../copyright/ZeusCopyright";
import TextField from "@mui/material/TextField";
import {AppBar, Drawer} from "../dashboard/Dashboard";
import {RootState} from "../../redux/store";
import {
    removeAggregationFromWorkflowBuilderTaskMap,
    removeEvalFnFromWorkflowBuilderEvalMap,
    setAddAggregateTasks,
    setAddAggregationView,
    setAddAnalysisTasks,
    setAddAnalysisView,
    setAddAssistantsView,
    setAddEvalFns,
    setAddEvalFnsView,
    setAddRetrievalTasks,
    setAddRetrievalView,
    setAddSchemasView,
    setAddTriggerActionsView,
    setAddTriggersToEvalFnView,
    setAnalysisRetrievalsMap,
    setEditAggregateTask,
    setEditAnalysisTask,
    setEval,
    setEvalFns,
    setEvalMap,
    setEvalMetric,
    setEvalsTaskMap,
    setRetrieval,
    setSchema,
    setSchemas,
    setSelectedMainTabBuilder,
    setSelectedWorkflows,
    setTaskMap,
    setTriggerAction,
    setTriggerActions,
    setWorkflowBuilderTaskMap,
    setWorkflowGroupName,
    setWorkflowName,
    updateEvalMetrics,
} from "../../redux/ai/ai.reducer";
import {aiApiGateway} from "../../gateway/ai";
import {
    DeleteWorkflowsActionRequest,
    EvalFn,
    EvalMetric,
    PostWorkflowsRequest,
    TaskModelInstructions
} from "../../redux/ai/ai.types";
import {TasksTable} from "./TasksTable";
import {isValidLabel} from "../clusters/wizard/builder/AddComponentBases";
import {RetrievalsTable} from "./RetrievalsTable";
import {loadBalancingApiGateway} from "../../gateway/loadbalancing";
import {setEndpoints, setGroupEndpoints} from "../../redux/loadbalancing/loadbalancing.reducer";
import {ExpandLess, ExpandMore} from "@mui/icons-material";
import {EvalsTable} from "./EvalsTable";
import {Assistants} from "./Assistants";
import {AssistantsTable} from "./AssistantsTable";
import {ActionsTable} from "./ActionsTable";
import {Retrieval, TriggerAction} from "../../redux/ai/ai.types2";
import {SchemasTable} from "./SchemasTable";
import {Schemas} from "./Schemas";
import {JsonSchemaDefinition, JsonSchemaField} from "../../redux/ai/ai.types.schemas";

const mdTheme = createTheme();

function WorkflowEngineBuilder(props: any) {
    const addSchemasView = useSelector((state: RootState) => state.ai.addSchemasView);
    const schema = useSelector((state: RootState) => state.ai.schema);
    const schemas = useSelector((state: RootState) => state.ai.schemas);
    const schemaField = useSelector((state: RootState) => state.ai.schemaField);
    const [open, setOpen] = useState(true);
    const assistants = useSelector((state: RootState) => state.ai.assistants);
    const evalMap = useSelector((state: RootState) => state.ai.evalMap);
    const evalFns = useSelector((state: RootState) => state.ai.evalFns);
    const actions = useSelector((state: RootState) => state.ai.triggerActions);
    const addedEvalFns = useSelector((state: RootState) => state.ai.addedEvalFns);
    const groups = useSelector((state: RootState) => state.loadBalancing.groups);
    const [loading, setIsLoading] = useState(false);
    const selectedWorkflows = useSelector((state: any) => state.ai.selectedWorkflows);
    const selectedMainTabBuilder = useSelector((state: any) => state.ai.selectedMainTabBuilder);
    const [selected, setSelected] = useState<{ [key: number]: boolean }>({});
    const addAnalysisView = useSelector((state: RootState) => state.ai.addAnalysisView);
    const addAggregateView = useSelector((state: RootState) => state.ai.addAggregationView);
    const addRetrievalView = useSelector((state: RootState) => state.ai.addRetrievalView);
    const allTasks = useSelector((state: any) => state.ai.tasks);
    const retrievalStages = useSelector((state: RootState) => state.ai.addedRetrievals);
    const [taskType, setTaskType] = useState('analysis');
    const analysisStages = useSelector((state: RootState) => state.ai.addedAnalysisTasks);
    const [selectedRetrievalForAnalysis, setSelectedRetrievalForAnalysis] = useState('');
    const [selectedAnalysisStageForRetrieval, setSelectedAnalysisStageForRetrieval] = useState('');
    const [selectedAnalysisStageForAggregation, setSelectedAnalysisStageForAggregation] = useState('');
    const [selectedAggregationStageForAnalysis, setSelectedAggregationStageForAnalysis] = useState('');
    const [selectedAnalysisStageForEval, setSelectedAnalysisStageForEval] = useState('');
    const [selectedAggregationStageForEval, setSelectedAggregationStageForEval] = useState('');
    const [selectedEvalStage, setSelectedEvalStage] = useState('');
    const addTriggersToEvalFnView = useSelector((state: RootState) => state.ai.addTriggersToEvalFnView);
    const evalFnStages = useSelector((state: RootState) => state.ai.addedEvalFns);
    const addEvalsView = useSelector((state: RootState) => state.ai.addEvalFnsView);
    const addAssistantsView = useSelector((state: RootState) => state.ai.addAssistantsView);
    const addTriggerActionsView = useSelector((state: RootState) => state.ai.addTriggerActionsView);
    const aggregationStages = useSelector((state: RootState) => state.ai.addedAggregateTasks);
    const [tasks, setTasks] = useState(allTasks && allTasks.filter((task: TaskModelInstructions) => task.taskType === taskType));
    const retrievals = useSelector((state: RootState) => state.ai.retrievals);
    const workflowBuilderTaskMap = useSelector((state: RootState) => state.ai.workflowBuilderTaskMap);
    const workflowAnalysisRetrievalsMap = useSelector((state: RootState) => state.ai.workflowAnalysisRetrievalsMap);
    const workflowBuilderEvalsTaskMap = useSelector((state: RootState) => state.ai.workflowBuilderEvalsTaskMap);
    const taskMap = useSelector((state: RootState) => state.ai.taskMap);
    const retrievalsMap = useSelector((state: RootState) => state.ai.retrievalsMap);
    const retrieval = useSelector((state: RootState) => state.ai.retrieval);
    const workflowName = useSelector((state: RootState) => state.ai.workflowName);
    const workflowGroupName = useSelector((state: RootState) => state.ai.workflowGroupName);
    const workflows = useSelector((state: any) => state.ai.workflows);
    const action = useSelector((state: any) => state.ai.triggerAction);
    const evalMetric = useSelector((state: any) => state.ai.evalMetric);
    const evalFn = useSelector((state: any) => state.ai.evalFn);
    //const actionMetric = useSelector((state: any) => state.ai.actionMetric);
    //const actionsEvalTrigger = useSelector((state: any) => state.ai.actionsEvalTrigger);
    // const actionPlatformAccount = useSelector((state: any) => state.ai.actionPlatformAccount);
    const assistant = useSelector((state: RootState) => state.ai.assistant);
    const [openRetrievals, setOpenRetrievals] = useState<boolean>(true); // Or use an object/array for multiple sections
    const [openAnalysis, setOpenAnalysis] = useState<boolean>(true); // Or use an object/array for multiple sections
    const [openAggregation, setOpenAggregation] = useState<boolean>(true); // Or use an object/array for multiple sections
    const [openActions, setActions] = useState<boolean>(true); // Or use an object/array for multiple sections
    const [openEvals, setOpenEvals] = useState<boolean>(true); // Or use an object/array for multiple sections
    const [toggleEvalToTaskType, setToggleEvalToTaskType] = useState<boolean>(false); // Or use an object/array for multiple sections
    const editAnalysisTask = useSelector((state: any) => state.ai.editAnalysisTask);
    const editAggregateTask = useSelector((state: any) => state.ai.editAggregateTask);
    const setToggleEvalTaskType = () => {
        setToggleEvalToTaskType(!toggleEvalToTaskType);
    };
    const handleEditAnalysisTaskResponseFormat = (event: any) => {
        const responseFormat = event.target.value;
        dispatch(setEditAnalysisTask({ ...editAnalysisTask, responseFormat: responseFormat }))
    }
    const handleEditAggTaskResponseFormat = (event: any) => {
        const responseFormat = event.target.value;
        dispatch(setEditAggregateTask({ ...editAggregateTask, responseFormat: responseFormat }))
    }

    const handleEditAggregateTaskModel = (event: any) => {
        dispatch(setEditAggregateTask({ ...editAggregateTask, model: event.target.value }))
    }
    const editEvalMetricRow = (index: number) => {
        const updatedMetrics = evalFn.evalMetrics.filter((_: EvalMetric, i: number) => i === index);
        if (updatedMetrics.length > 0) {
            dispatch(setEvalMetric(updatedMetrics[0]));
        }
    };

    const removeEvalMetricRow = (index: number) => {
        const updatedMetrics = evalFn.evalMetrics.filter((_: EvalMetric, i: number) => i !== index);
        dispatch(updateEvalMetrics(updatedMetrics));
        setRequestEvalCreateOrUpdateStatus('')
        setRequestEvalCreateOrUpdateStatusError('')
    };

    const removeSchemaField = (index: number) => {
        const updatedFields = schema.fields.filter((_: JsonSchemaField, i: number) => i !== index);
        dispatch(setSchema({
            ...schema,
            fields: updatedFields
        }))
        setRequestStatusSchema('')
        setRequestStatusSchemaError('')
    };

    const clearEvalMetricRow = () => {
        dispatch(setEvalMetric({
            evalMetricName: '',
            evalModelPrompt: '',
            evalComparisonNumber: 1,
            evalComparisonString: '',
            evalComparisonBoolean: false,
            evalMetricDataType: '',
            evalOperator: '',
            evalState: 'info',
            evalMetricResult: '',
        }));
    }
    const addJsonSchemaFieldRow = () => {
        if (schemaField.fieldName.length <= 0) {
            setRequestStatusSchema('Field name is empty')
            setRequestStatusSchemaError('error')
            return;
        }
        if (schemaField.dataType.length <= 0) {
            setRequestStatusSchema('Data type must be set')
            setRequestStatusSchemaError('error')
            return;
        }
        if (schemaField.fieldDescription.length <= 0) {
            setRequestStatusSchema('Field description must be set')
            setRequestStatusSchemaError('error')
            return;
        }

        setRequestStatusSchema('')
        setRequestStatusSchemaError('')

        if (schema.fields && schema.fields.length > 0) {
            const updatedFields = updateFieldByName(schema.fields, schemaField);
            dispatch(setSchema(
                {...schema, fields: updatedFields}))
        } else {
            dispatch(setSchema(
                {...schema, fields: [schemaField]}))
        }
    };

    const updateFieldByName = (fields: JsonSchemaField[], newField: JsonSchemaField) => {
        // Filter out the old field if it exists
        const fieldsWithoutOld = fields.filter(field => field.fieldName !== newField.fieldName);
        // Add the new field to the array
        return [...fieldsWithoutOld, newField];
    };

    const updateSchemaByName = (schemas: JsonSchemaDefinition[], newSchema: JsonSchemaDefinition) => {
        // Filter out the old field if it exists
        const schemasWithoutOld = schemas.filter(schema => schema.schemaName !== newSchema.schemaName);
        // Add the new field to the array
        return [...schemasWithoutOld, schema];
    };


    const addEvalMetricRow = () => {
        if (!isValidLabel(evalMetric.evalMetricName)){
            setRequestEvalCreateOrUpdateStatus('Metric name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
            setRequestEvalCreateOrUpdateStatusError('error')
            return;
        }
        if (evalFn.evalType === 'model' && evalMetric.evalModelPrompt.length <= 0){
            setRequestEvalCreateOrUpdateStatus('Prompt is empty')
            setRequestEvalCreateOrUpdateStatusError('error')
            return;
        }
        setRequestEvalCreateOrUpdateStatus('')
        setRequestEvalCreateOrUpdateStatusError('')

        if (evalFn.evalMetrics && evalFn.evalMetrics.length > 0){
            const updatedMetrics = updateMetricByName(evalFn.evalMetrics, evalMetric);
            dispatch(updateEvalMetrics(updatedMetrics));
        } else {
            dispatch(updateEvalMetrics([evalMetric]));
        }
    };
    const updateMetricByName = (metrics: EvalMetric[], newMetric: EvalMetric) => {
        const metricsWithoutOld = metrics.filter(metric => metric.evalMetricName !== newMetric.evalMetricName);
        return [...metricsWithoutOld, newMetric];
    };
    // const addActionMetricRow = () => {
    //     if (!isValidLabel(actionMetric.metricName)){
    //         setRequestMetricActionCreateOrUpdateStatus('Metric name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
    //         setRequestMetricActionCreateOrUpdateStatusError('error')
    //         return;
    //     }
    //     // Check if the metric name already exists in actionMetrics
    //     const existingMetric = action.actionMetrics.some((metric: { metricName: string; }) => metric.metricName === actionMetric.metricName);
    //     if (existingMetric) {
    //         setRequestMetricActionCreateOrUpdateStatus('Metric name already exists.');
    //         setRequestMetricActionCreateOrUpdateStatusError('error');
    //         return;
    //     }
    //     setRequestMetricActionCreateOrUpdateStatus('')
    //     setRequestMetricActionCreateOrUpdateStatusError('')
    //     const updatedMetrics = [...action.actionMetrics,actionMetric];
    //     // dispatch(updateActionMetrics(updatedMetrics));
    // };
    const handleAddRetrievalToAnalysis = () => {
        if (selectedRetrievalForAnalysis.length <= 0 || selectedRetrievalForAnalysis.length <= 0) {
            return;
        }
        const retKey = Number(selectedRetrievalForAnalysis);
        const analysisKey = Number(selectedAnalysisStageForRetrieval);
        const payload = {
            key: retKey,
            subKey: analysisKey,
            value: true
        };
        dispatch(setAnalysisRetrievalsMap(payload));
    };
    const handleRemoveRetrievalRelationshipFromWorkflow = async (event: any, keystr: string, value: number) => {
        const key = Number(keystr);
        const payload = {
            key: key,
            subKey: value,
            value: false
        };
        dispatch(setAnalysisRetrievalsMap(payload));
    }

    const handleRemoveEvalFnRelationshipFromWorkflow = async (event: any, evalID: number, value: number) => {
        const payload = {
            evalID: evalID,
            evalTaskID: value,
            value: false
        };
        dispatch(setEvalsTaskMap(payload));
    }

    const handleAddEvalToSubTask = async (e: any) => {
        if (!toggleEvalToTaskType){
            if (selectedAnalysisStageForEval.length <= 0 && selectedEvalStage.length <= 0){
                return;
            }
            const evalID = Number(selectedEvalStage);
            const analysisKey = Number(selectedAnalysisStageForEval);
            const payload = {
                evalID: evalID,
                evalTaskID: analysisKey,
                value: true
            }
            dispatch(setEvalsTaskMap(payload));
        } else {
            if (selectedAggregationStageForEval.length <= 0 && selectedEvalStage.length <= 0){
                return;
            }
            const evalID = Number(selectedEvalStage);
            const aggKey = Number(selectedAggregationStageForEval);
            const payload = {
                evalID: evalID,
                evalTaskID: aggKey,
                value: true
            }
            dispatch(setEvalsTaskMap(payload));
        }
    };
    const handleAddSubTaskToAggregate = () => {
        if (selectedAggregationStageForAnalysis.length <= 0 || selectedAnalysisStageForAggregation.length <= 0) {
            return;
        }
        const aggKey = Number(selectedAggregationStageForAnalysis);
        const analysisKey = Number(selectedAnalysisStageForAggregation);
        const payload = {
            key: aggKey,
            subKey: analysisKey,
            value: true
        };
        dispatch(setWorkflowBuilderTaskMap(payload));
    };
    const handleRemoveRetrievalFromWorkflow = async (event: any, retrievalRemove: Retrieval) => {
        dispatch(setAddRetrievalTasks(retrievalStages.filter((ret: Retrieval) => ret.retrievalID !== retrievalRemove.retrievalID)));
    }
    const handleRemoveTriggerFromEvalFn = async (event: any, triggerRemove: TriggerAction) => {
        const newTriggerFunctions = evalFn.triggerFunctions.filter((trigger: TriggerAction) => trigger.triggerID !== triggerRemove.triggerID);
        dispatch(setEval({
            ...evalFn, // Spread the existing action properties
            triggerFunctions: newTriggerFunctions// Update the actionName
        }))
    }
    const handleAddTriggersToEvalFn = async (event: any) => {
        setIsLoading(true);
        if (addTriggersToEvalFnView) {
            // Get the currently selected triggers
            const selectedTriggers: TriggerAction[] = Object.keys(selected)
                .filter(key => selected[Number(key)])
                .map(key => actions[Number(key)]);


            if (evalFn && evalFn.triggerFunctions && selectedTriggers.length > 0) {
                // Filter out triggers that already exist in evalFn.triggerFunctions
                const newTriggersToAdd: TriggerAction[] = selectedTriggers.filter((st: TriggerAction) =>
                    !evalFn.triggerFunctions.some((tf: TriggerAction) => tf.triggerID === st.triggerID));

                // Combine the existing triggers with the new, non-duplicate triggers
                const updatedTriggerFunctions: TriggerAction[] = [...evalFn.triggerFunctions, ...newTriggersToAdd];

                // Dispatch the updated evalFn
                dispatch(setEval({
                    ...evalFn, // Spread the existing action properties
                    triggerFunctions: updatedTriggerFunctions // Update with combined triggers
                }));
            } else {
                // Dispatch the updated evalFn
                dispatch(setEval({
                    ...evalFn, // Spread the existing action properties
                    triggerFunctions: selectedTriggers // Update with combined triggers
                }));
            }
        }
        setIsLoading(false);
    };

    const handleRemoveAnalysisFromWorkflow = async (event: any, taskRemove: TaskModelInstructions) => {
        Object.entries(workflowBuilderTaskMap).map(([key, value], index) => {
            const payloadRemove = {
                key: Number(key),
                subKey: taskRemove.taskID? Number(taskRemove.taskID) : 0,
                value: false
            };
            dispatch(setWorkflowBuilderTaskMap(payloadRemove));
        })
        dispatch(setAddAnalysisTasks(analysisStages.filter((task: TaskModelInstructions) => task.taskID !== taskRemove.taskID)));
    }
    const handleRemoveAggregationFromWorkflow = async (event: any, taskRemove: TaskModelInstructions) => {
        const payload = {
            key: Number(taskRemove.taskID),
            subKey: 0,
            value: true
        };
        dispatch(removeAggregationFromWorkflowBuilderTaskMap(payload));
        dispatch(setAddAggregateTasks(aggregationStages.filter((task: TaskModelInstructions) => task.taskID !== taskRemove.taskID)));
    }
    const handleRemoveEvalFnFromWorkflow = async (event: any, evalFn: EvalFn) => {
        const payload = {
            key: evalFn.evalID? evalFn.evalID : 0,
            subKey: 0,
            value: false
        }
        dispatch(removeEvalFnFromWorkflowBuilderEvalMap(payload));
        dispatch(setAddEvalFns(addedEvalFns.filter(efn => efn.evalID !== evalFn.evalID)));
    }
    const handleRemoveTaskRelationshipFromWorkflow = async (event: any, keystr: string, value: number) => {
        const key = Number(keystr);
        const payload = {
            key: key,
            subKey: value,
            value: false
        };
        dispatch(setWorkflowBuilderTaskMap(payload));
    }
    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                setIsLoading(true); // Set loading to true
                const response = await loadBalancingApiGateway.getEndpoints();
                dispatch(setEndpoints(response.data.routes));
                dispatch(setGroupEndpoints(response.data.orgGroupsRoutes));
            } catch (error) {
                console.log("error", error);
            } finally {
                setIsLoading(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData({});
    }, []);


    useEffect(() => {
    }, [addEvalsView, addAggregateView, addAnalysisView,addRetrievalView, selectedMainTabBuilder, addSchemasView,
        analysisStages, aggregationStages, evalFnStages, retrievals, retrievalStages,workflowBuilderEvalsTaskMap, workflowBuilderTaskMap,
        workflowAnalysisRetrievalsMap, evalMap, taskMap]);
    const dispatch = useDispatch();
    const handleTaskCycleCountChange = (val: number, task: TaskModelInstructions) => {
        if (val <= 0) {
            val = 1
        }
        if (task && task.taskID) {
            const payload = {
                key: task.taskID,
                count: val
            };
            dispatch(setTaskMap(payload));
        }
    };

    const handleEvalCycleCountChange = (val: number, ef: EvalFn) => {
        if (val <= 0) {
            val = 1
        }
        if (ef && ef.evalID) {
            const payload = {
                key: ef.evalID,
                count: val
            };
            dispatch(setEvalMap(payload));
        }
    };

    const handleDeleteWorkflows = async (event: any) => {
        const params: DeleteWorkflowsActionRequest = {
            workflows: selectedWorkflows.map((ind: string) => {
                return workflows[Number(ind)]
            })
        }
        if (params.workflows.length === 0) {
            return
        }
        try {
            setIsLoading(true)
            const response = await aiApiGateway.deleteWorkflowsActionRequest(params);
            const statusCode = response.status;
            if (statusCode < 400) {
                dispatch(setSelectedWorkflows([]));
                const data = response.data;
            } else {
                console.log('Failed to delete', response);
            }
        } catch (e) {
        } finally {
            setIsLoading(false);
        }
    };

    const handleAddTasksToWorkflow = async (event: any) => {
        setIsLoading(true)
        if (addAnalysisView){
            const selectedTasks: TaskModelInstructions[] = Object.keys(selected)
                .filter(key => selected[Number(key)])
                .map(key => tasks[Number(key)]);
            dispatch(setAddAnalysisTasks(selectedTasks.filter((task: TaskModelInstructions) => task.taskType === 'analysis')));
        } else if (addAggregateView){
            const selectedTasks: TaskModelInstructions[] = Object.keys(selected)
                .filter(key => selected[Number(key)])
                .map(key => tasks[Number(key)]);
            dispatch(setAddAggregateTasks(selectedTasks.filter((task: TaskModelInstructions) => task.taskType === 'aggregation')));
        } else if (addRetrievalView) {
            const selectedTasks: Retrieval[] = Object.keys(selected)
                .filter(key => selected[Number(key)])
                .map(key => retrievals[Number(key)]);
            dispatch(setAddRetrievalTasks(selectedTasks));
        } else if (addEvalsView) {
            const selectedEvals: EvalFn[] = Object.keys(selected)
                .filter(key => selected[Number(key)])
                .map(key => evalFns[Number(key)]);
            dispatch(setAddEvalFns(selectedEvals));
        }
        setIsLoading(false)
    }

    const addSchemaToTask = async (event: any) => {
        const selectedSchemas: JsonSchemaDefinition[] = Object.keys(selected)
            .filter(key => selected[Number(key)])
            .map(key => schemas[Number(key)]);

        if (taskType === 'analysis') {
            dispatch(setEditAnalysisTask({ ...editAnalysisTask, schemas: selectedSchemas }))
        } else if (taskType === 'aggregation') {
            dispatch(setEditAggregateTask({ ...editAggregateTask, schemas: selectedSchemas }))
        }
    }

    const addSchemasViewToggle = async (event: any) => {
        const toggle = !addSchemasView;
        dispatch(setAddSchemasView(toggle));
        if (toggle) {
            if (taskType === 'analysis') {
                dispatch(setSelectedMainTabBuilder(1))
            } else if (taskType === 'aggregation') {
                dispatch(setSelectedMainTabBuilder(2))
            }
            setSelected({});
        } else {
            setSelected({});
        }

    }

    const removeSchemasViewToggle = async (event: any, index: number) => {
        if (taskType === 'analysis') {
            const updatedSchemas = editAnalysisTask.schemas.filter((_: JsonSchemaDefinition, i: number) => i !== index);
            dispatch(setEditAnalysisTask({ ...editAnalysisTask, schemas: updatedSchemas }))
        } else if (taskType === 'aggregation') {
            const updatedSchemas = editAggregateTask.schemas.filter((_: JsonSchemaDefinition, i: number) => i !== index);
            dispatch(setEditAggregateTask({ ...editAggregateTask, schemas: updatedSchemas }))
        }
    }

    const addEvalsStageView = async () => {
        const toggle = !addEvalsView;
        dispatch(setAddAnalysisView(false));
        dispatch(setAddAggregationView(false));
        dispatch(setAddRetrievalView(false));
        dispatch(setAddEvalFnsView(toggle));

        if (toggle) {
            dispatch(setSelectedMainTabBuilder(4))
            setSelected({});
        } else {
            dispatch(setSelectedMainTabBuilder(0))
        }
    }

    const addRetrievalStageView = async () => {
        const toggle = !addRetrievalView;
        dispatch(setAddAnalysisView(false));
        dispatch(setAddAggregationView(false));
        dispatch(setAddRetrievalView(toggle));
        dispatch(setAddEvalFnsView(false));
        if (toggle) {
            dispatch(setSelectedMainTabBuilder(3))
            setSelected({});
        } else {
            dispatch(setSelectedMainTabBuilder(0))
        }
    }

    const addAnalysisStageView = async () => {
        const toggle = !addAnalysisView;
        dispatch(setAddAnalysisView(toggle));
        dispatch(setAddAggregationView(false));
        dispatch(setAddRetrievalView(false));
        dispatch(setAddEvalFnsView(false));
        if (toggle) {
            dispatch(setSelectedMainTabBuilder(1))
            setSelected({});
            setTaskType('analysis');
            setTasks(allTasks && allTasks.filter((task: any) => task.taskType === 'analysis'));
        } else {
            dispatch(setSelectedMainTabBuilder(0))
        }
    }
    const addAggregationStageView = async () => {
        const toggle = !addAggregateView;
        dispatch(setAddAnalysisView(false));
        dispatch(setAddRetrievalView(false));
        dispatch(setAddAggregationView(toggle));
        dispatch(setAddEvalFnsView(false));
        if (toggle) {
            setSelected({});
            setTaskType('aggregation');
            setTasks(allTasks && allTasks.filter((task: any) => task.taskType === 'aggregation'));
            dispatch(setSelectedMainTabBuilder(2))
        } else {
            dispatch(setSelectedMainTabBuilder(0))
        }
    }
    const [stepSize, setStepSize] = useState(5);
    const handleTimeStepChange = (event: any) => {
        setStepSize(event.target.value);
    };
    const [stepSizeUnit, setStepSizeUnit] = React.useState('minutes');
    const handleUpdateStepSizeUnit = (event: any) => {
        setStepSizeUnit(event.target.value);
    };
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();

    const handleLogout = async (event: any) => {
        event.preventDefault();
        await authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }
    const [requestAggStatus, setRequestAggStatus] = useState('');
    const [requestAggStatusError, setRequestAggStatusError] = useState('');
    const [requestAnalysisStatus, setRequestAnalysisStatus] = useState('');
    const [requestAnalysisStatusError, setRequestAnalysisStatusError] = useState('');
    const [requestRetrievalStatus, setRequestRetrievalStatus] = useState('');
    const [requestRetrievalStatusError, setRequestRetrievalStatusError] = useState('');
    const [requestStatus, setRequestStatus] = useState('');
    const [requestStatusError, setRequestStatusError] = useState('');
    const [requestActionStatus, setRequestActionStatus] = useState('');
    const [requestActionStatusError, setRequestActionStatusError] = useState('');
    const [requestEvalCreateOrUpdateStatus, setRequestEvalCreateOrUpdateStatus] = useState('');
    const [requestEvalCreateOrUpdateStatusError, setRequestEvalCreateOrUpdateStatusError] = useState('');
    const [requestMetricActionCreateOrUpdateStatus, setRequestMetricActionCreateOrUpdateStatus] = useState('');
    const [requestMetricActionCreateOrUpdateStatusError, setRequestMetricActionCreateOrUpdateStatusError] = useState('');
    const [requestStatusAssistant, setRequestStatusAssistant] = useState('');
    const [requestStatusAssistantError, setRequestStatusAssistantError] = useState('');
    const [requestStatusSchema, setRequestStatusSchema] = useState('');
    const [requestStatusSchemaError, setRequestStatusSchemaError] = useState('');

    const createOrUpdateWorkflow = async () => {
        try {
            const allMappedRetrievalIDs: Set<number> = new Set();
            Object.entries(workflowAnalysisRetrievalsMap).forEach(([retrievalID, innerMap])=> {
                Object.entries(innerMap).forEach(([analysisID, isAdded], subInd) => {
                    if (isAdded) {
                        //allMappedAnalysisIDs.add(parseInt(analysisID, 10));
                        allMappedRetrievalIDs.add(parseInt(retrievalID, 10));
                    }
                });
            });
            if (allMappedRetrievalIDs.size < retrievalStages.length) {
                setRequestStatus('All retrieval stages must be connected to at least one analysis stage')
                setRequestStatusError('error')
                return;
            }
            if (analysisStages.length <= 0) {
                setRequestStatus('Workflow must have at least one analysis stage')
                setRequestStatusError('error')
                return;
            }
            if (aggregationStages.length > 0 && analysisStages.length <= 0) {
                setRequestStatus('Workflows with aggregation stages must have at least one connected analysis stage')
                setRequestStatusError('error')
                return;
            }
            if (Object.keys(workflowBuilderTaskMap).length <= 0 && analysisStages.length <= 0) {
                setRequestStatus('Workflows with aggregation stages must have at least one connected analysis stage')
                setRequestStatusError('error')
                return;
            }
            if (Object.keys(workflowBuilderTaskMap).length < aggregationStages.length) {
                setRequestStatus('Workflows with aggregation stages must have at least one connected analysis stage')
                setRequestStatusError('error')
                return;
            }
            if (Object.keys(workflowAnalysisRetrievalsMap).length > 0 && retrievalStages.length <= 0) {
                setRequestStatus('Workflows with retrieval stages must have at least one connected analysis stage')
                setRequestStatusError('error')
                return;
            }
            if (stepSize <= 0) {
                setRequestStatus('Step size must be greater than 0')
                setRequestStatusError('error')
                return;
            }
            if (!isValidLabel(workflowName)) {
                setRequestStatus('Workflow name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestStatusError('error')
                return;
            }
            if (!isValidLabel(workflowGroupName)) {
                setRequestStatus('Workflow group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestStatusError('error')
                return;
            }
            const payload: PostWorkflowsRequest = {
                workflowName: workflowName,
                workflowGroupName: workflowGroupName,
                stepSize: Number(stepSize),
                stepSizeUnit: stepSizeUnit,
                models: taskMap,
                aggregateSubTasksMap: workflowBuilderTaskMap,
                analysisRetrievalsMap: workflowAnalysisRetrievalsMap,
                evalsMap: evalMap,
                evalTasksMap: workflowBuilderEvalsTaskMap
            }
            setIsLoading(true)
            const response = await aiApiGateway.createAiWorkflowRequest(payload);
            const statusCode = response.status
            if (statusCode < 400) {
                const data = response.data;
                setSelected({});
                dispatch(setAddAnalysisView(false));
                dispatch(setAddAggregationView(false));
                dispatch(setAddEvalFnsView(false));
                dispatch(setAddRetrievalView(false));
                dispatch(setSelectedMainTabBuilder(0));
                setRequestStatus('Workflow created successfully')
                setRequestStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }
    const createOrUpdateSchema = async () => {
        try {
            setIsLoading(true)
            if (!isValidLabel(schema.schemaName)) {
                setRequestStatusSchema('Schema name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestStatusSchemaError('error')
                return;
            }
            if (!isValidLabel(schema.schemaGroup)) {
                setRequestStatusSchema('Schema group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestStatusSchemaError('error')
                return;
            }
            if (schema.fields.length <= 0) {
                setRequestStatusSchema('Schema must have at least one field')
                setRequestStatusSchemaError('error')
                return;
            }
            const response = await aiApiGateway.createOrUpdateJsonSchema(schema);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data as JsonSchemaDefinition;
                const updatedSchemas = updateSchemaByName(schemas, data);
                dispatch(setSchemas(updatedSchemas))
                setRequestStatusSchema('Schema created or updated successfully')
                setRequestStatusSchemaError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestStatusSchema('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestStatusSchemaError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }
    const createOrUpdateAssistant = async () => {
        try {
            setIsLoading(true)
            if (assistant.model.length <= 0) {
                setRequestStatusAssistant('Assistant model must be set')
                setRequestStatusAssistantError('error')
                return;
            }
            const response = await aiApiGateway.createOrUpdateAssistant(assistant);
            const statusCode = response.status;
            if (statusCode < 400) {
                // const data = response.data as Retrieval;
                // dispatch(setRetrievals([...retrievals, data]))
                setRequestStatusAssistant('Assistant created or updated successfully')
                setRequestStatusAssistantError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestStatusAssistant('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestStatusAssistantError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }

    const createOrUpdateRetrieval = async () => {
        try {
            setIsLoading(true)
            if (!isValidLabel(retrieval.retrievalName)) {
                setRequestRetrievalStatus('Retrieval name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestRetrievalStatusError('error')
                return;
            }
            if (!isValidLabel(retrieval.retrievalGroup)) {
                setRequestRetrievalStatus('Retrieval group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestRetrievalStatusError('error')
                return;
            }
            if ((retrieval.retrievalItemInstruction.retrievalKeywords && retrieval.retrievalItemInstruction.retrievalKeywords.length <= 0 && retrieval.retrievalItemInstruction.retrievalPrompt && retrieval.retrievalItemInstruction.retrievalPrompt.length <= 0 && retrieval.retrievalGroup.length <= 0 ))  {
                setRequestRetrievalStatus('At least one of retrieval keywords or prompt or group must be set')
                setRequestRetrievalStatusError('error')
                return;
            }
            if (retrieval.retrievalItemInstruction.retrievalPlatform.length <= 0) {
                setRequestRetrievalStatus('Retrieval platform must be set')
                setRequestRetrievalStatusError('error')
                return;
            }
            const response = await aiApiGateway.createOrUpdateRetrieval(retrieval);
            const statusCode = response.status;
            if (statusCode < 400) {
                // const data = response.data as Retrieval;
                // dispatch(setRetrievals([...retrievals, data]))
                setRequestRetrievalStatus('Retrieval created successfully')
                setRequestRetrievalStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestRetrievalStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestRetrievalStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }

    const createOrUpdateTask = async (taskType: string) => {
        try {
            setIsLoading(true)
            const tn = (taskType === 'analysis' ? editAnalysisTask.taskName : editAggregateTask.taskName);

            if (!isValidLabel(tn)) {
                if (taskType === 'analysis') {
                    setRequestAnalysisStatus('Analysis task name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                    setRequestAnalysisStatusError('error')
                } else if (taskType === 'aggregation') {
                    setRequestAggStatus('Aggregation task name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                    setRequestAggStatusError('error')
                }
                return;
            }
            const taskGn = (taskType === 'analysis' ? editAnalysisTask.taskGroup : editAggregateTask.taskGroup);
            if (!isValidLabel(taskGn)) {
                if (taskType === 'analysis') {
                    setRequestAnalysisStatus('Analysis task group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                    setRequestAnalysisStatusError('error')
                } else if (taskType === 'aggregation') {
                    setRequestAggStatus('Aggregation task group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                    setRequestAggStatusError('error')
                }
                return;
            }
            const prompt = (taskType === 'analysis' ? editAnalysisTask.prompt : editAggregateTask.prompt);
            if (prompt.length <= 0) {
                if (taskType === 'analysis') {
                    setRequestAnalysisStatus('Analysis task prompt is empty')
                    setRequestAnalysisStatusError('error')
                } else if (taskType === 'aggregation') {
                    setRequestAggStatus('Aggregation task prompt is empty')
                    setRequestAggStatusError('error')
                }
                return;
            }
            const task: TaskModelInstructions = {
                taskID: (taskType === 'analysis' ? editAnalysisTask.taskID : editAggregateTask.taskID),
                taskType: taskType,
                taskGroup:taskGn,
                taskName: tn,
                schemas: (taskType === 'analysis' ? editAnalysisTask.schemas : editAggregateTask.schemas),
                responseFormat: (taskType === 'analysis' ? editAnalysisTask.responseFormat : editAggregateTask.responseFormat),
                model: (taskType === 'analysis' ? editAnalysisTask.model : editAggregateTask.model),
                prompt: (taskType === 'analysis' ? editAnalysisTask.prompt : editAggregateTask.prompt),
                maxTokens:  (taskType === 'analysis' ? editAnalysisTask.maxTokens : editAggregateTask.maxTokens),
                cycleCount: (taskType === 'analysis' ? 1 : 1),
                tokenOverflowStrategy: (taskType === 'analysis' ? editAnalysisTask.tokenOverflowStrategy : editAggregateTask.tokenOverflowStrategy),
            };
            const response = await aiApiGateway.createOrUpdateTaskRequest(task);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data as TaskModelInstructions;
                if (taskType === 'analysis') {
                    const at = allTasks.filter((task: any) => task.taskType === 'analysis' && task.taskID !== data.taskID)
                    setTasks([data, ...at]);
                    setRequestAnalysisStatus('Task created successfully')
                    setRequestAnalysisStatusError('success')
                } else if (taskType === 'aggregation') {
                    const at = allTasks.filter((task: any) => task.taskType === 'aggregation' && task.taskID !== data.taskID)
                    setTasks([data, ...at]);
                    setRequestAggStatus('Task created successfully')
                    setRequestAggStatusError('success')
                }
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                if (taskType === 'analysis') {
                    setRequestAnalysisStatus('Billing setup required. Please configure your billing information to continue using this service.');
                    setRequestAnalysisStatusError('error')
                } else if (taskType === 'aggregation') {
                    setRequestAggStatus('Billing setup required. Please configure your billing information to continue using this service.');
                    setRequestAggStatusError('error')
                }
            }
        } finally {
            setIsLoading(false);
        }
    }
    if (loading) {
        return <div>Loading...</div>;
    }

    const createOrUpdateAction = async () => {
        try {
            setIsLoading(true)
            if (!isValidLabel(action.triggerName)) {
                setRequestActionStatus('Action name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestActionStatusError('error')
                return;
            }
            if (!isValidLabel(action.triggerGroup)) {
                setRequestActionStatus('Action group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestActionStatusError('error')
                return;
            }

            if (action.triggerAction.length <= 0) {
                setRequestActionStatus('Trigger action environment must be set')
                setRequestActionStatusError('error')
                return;
            }
            const response = await aiApiGateway.createOrUpdateAction(action);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data as TriggerAction;
                const at = actions.filter((act: TriggerAction) => act.triggerName !== data.triggerName)
                dispatch(setTriggerActions([data, ...at]));
                setRequestActionStatus('Action created or updated successfully')
                setRequestActionStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestActionStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestActionStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }
    // const addRetrievalStageView = async () => {
    //     const toggle = !addRetrievalView;
    //     dispatch(setAddAnalysisView(false));
    //     dispatch(setAddAggregationView(false));
    //     dispatch(setAddRetrievalView(toggle));
    //     dispatch(setAddEvalFnsView(false));
    //     if (toggle) {
    //         dispatch(setSelectedMainTabBuilder(3))
    //         setSelected({});
    //     } else {
    //         dispatch(setSelectedMainTabBuilder(0))
    //     }
    // }
    const addTriggersToEvalFn = async () => {
        const toggle = !addTriggersToEvalFnView;
        dispatch(setAddTriggersToEvalFnView(toggle));
        if (toggle) {
            dispatch(setSelectedMainTabBuilder(5))
            setSelected({});
        } else {
            dispatch(setSelectedMainTabBuilder(4))
            setSelected({});
        }
        return
    }

    const createOrUpdateEval= async () => {
        try {
            setIsLoading(true)
            if (!isValidLabel(evalFn.evalName)) {
                setRequestEvalCreateOrUpdateStatus('Eval name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestEvalCreateOrUpdateStatusError('error')
                return;
            }
            if (!isValidLabel(evalFn.evalGroupName)) {
                setRequestEvalCreateOrUpdateStatus('Eval group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestEvalCreateOrUpdateStatusError('error')
                return;
            }
            if (evalFn.evalMetrics.length <= 0){
                setRequestEvalCreateOrUpdateStatus('You must add at least one metric create an eval')
                setRequestEvalCreateOrUpdateStatusError('error')
                return;
            }
            const response = await aiApiGateway.createOrUpdateEval(evalFn);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data as EvalFn;
                const ae = evalFns.filter((ef: EvalFn) =>  ef.evalID !== data.evalID)
                dispatch(setEvalFns([data, ...ae]));
                setRequestEvalCreateOrUpdateStatus('Eval created or updated successfully')
                setRequestEvalCreateOrUpdateStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestEvalCreateOrUpdateStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestEvalCreateOrUpdateStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }

    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        if (newValue === 1) {
            setSelected({});
            setTaskType('analysis');
            dispatch(setSelectedWorkflows([]));
            setTasks(allTasks && allTasks.filter((task: any) => task.taskType === 'analysis'));
        } else if (newValue === 2) {
            dispatch(setSelectedWorkflows([]));
            setSelected({});
            setTaskType('aggregation');
            setTasks(allTasks && allTasks.filter((task: any) => task.taskType === 'aggregation'));
        } else if (newValue === 3) {
            dispatch(setSelectedWorkflows([]));
            setSelected({});
        } else if (newValue === 4) {
            dispatch(setSelectedWorkflows([]));
            setSelected({});
        } else if (newValue === 5) {
            dispatch(setSelectedWorkflows([]));
            setSelected({});
        } else if (newValue === 6) {
            dispatch(setSelectedWorkflows([]));
            setSelected({});
        } else if (newValue === 7 && addSchemasView) {
            dispatch(setAddSchemasView(!addSchemasView));
            dispatch(setSelectedWorkflows([]));
            setSelected({});
        }
        setSelected({});

        if (addAssistantsView && newValue !== 6) {
            dispatch(setAddAssistantsView(false));
        }
        if (addTriggersToEvalFnView && newValue !== 4) {
            dispatch(setAddTriggersToEvalFnView(false));
        }
        if (addEvalsView && newValue !== 4) {
            dispatch(setAddEvalFnsView(false));
        }
        if (addTriggerActionsView && newValue !== 5) {
            dispatch(setAddTriggerActionsView(false));
        }
        if (addAggregateView && newValue !== 2) {
            dispatch(setAddAggregationView(false));
        }
        if (addAnalysisView && newValue !== 1) {
            dispatch(setAddAnalysisView(false));
        }
        if (addRetrievalView && newValue !== 3) {
            dispatch(setAddRetrievalView(false));
        }
        setRequestStatus('');
        setRequestStatusError('');
        setRequestAnalysisStatus('');
        setRequestAnalysisStatusError('');
        setRequestAggStatus('');
        setRequestAggStatusError('');
        setRequestRetrievalStatus('');
        setRequestRetrievalStatusError('');
        setRequestActionStatus('');
        setRequestActionStatusError('');
        setRequestEvalCreateOrUpdateStatus('');
        setRequestEvalCreateOrUpdateStatusError('');
        setRequestMetricActionCreateOrUpdateStatus('');
        setRequestMetricActionCreateOrUpdateStatusError('');
        setRequestStatusAssistant('');
        setRequestStatusAssistantError('');
        setRequestStatusSchemaError('');
        setRequestStatusSchema('');
        dispatch(setAddAnalysisView(false));
        dispatch(setSelectedMainTabBuilder(newValue));
    };
    const handleClick = (index: number) => {
        setSelected((prevSelected) => ({
            ...prevSelected,
            [index]: !prevSelected[index]
        }));
    }

    const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
        const isChecked = event.target.checked;
        if (addRetrievalView || selectedMainTabBuilder === 3 && !addSchemasView) {
            const newSelection = retrievals.reduce((acc: { [key: number]: boolean }, task: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        } else if ( selectedMainTabBuilder === 5)  {
            const newSelection = actions.reduce((acc: { [key: number]: boolean }, task: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        } else if (selectedMainTabBuilder === 4)  {
            const newSelection = evalFns.reduce((acc: { [key: number]: boolean }, task: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        } else if (selectedMainTabBuilder === 5)  {
            const newSelection = actions.reduce((acc: { [key: number]: boolean }, action: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        } else if (selectedMainTabBuilder === 6) {
            const newSelection = assistants.reduce((acc: { [key: number]: boolean }, assistant: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        } else if (selectedMainTabBuilder === 7) {
            const newSelection = schemas.reduce((acc: { [key: number]: boolean }, schema: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        } else if (addSchemasView) {
            const newSelection = schemas.reduce((acc: { [key: number]: boolean }, schema: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        } else {
            const newSelection = tasks.reduce((acc: { [key: number]: boolean }, task: any, index: number) => {
                acc[index] = isChecked;
                return acc;
            }, {});
            setSelected(newSelection);
        }
    };

    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBar position="absolute" open={open} style={{ backgroundColor: '#333'}}>
                    <Toolbar
                        sx={{
                            pr: '24px', // keep right padding when drawer closed
                        }}
                    >
                        <IconButton
                            edge="start"
                            color="inherit"
                            aria-label="open drawer"
                            onClick={toggleDrawer}
                            sx={{
                                marginRight: '36px',
                                ...(open && { display: 'none' }),
                            }}
                        >
                            <MenuIcon />
                        </IconButton>
                        <Typography
                            component="h1"
                            variant="h6"
                            color="inherit"
                            noWrap
                            sx={{ flexGrow: 1 }}
                        >
                            Mockingbird — An Intelligently Designed AI Systems Coordinator
                        </Typography>
                        <Button
                            color="inherit"
                            onClick={handleLogout}
                        >Logout
                        </Button>
                    </Toolbar>
                </AppBar>
                <Drawer variant="permanent" open={open}>
                    <Toolbar
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'flex-end',
                            px: [1],
                        }}
                    >
                        <IconButton onClick={toggleDrawer}>
                            <ChevronLeftIcon />
                        </IconButton>
                    </Toolbar>
                    <Divider />
                    <List component="nav">
                        <MainListItems />
                        <Divider sx={{ my: 1 }} />
                    </List>
                </Drawer>
                <Box
                    component="main"
                    sx={{
                        backgroundColor: (theme) =>
                            theme.palette.mode === 'light'
                                ? theme.palette.grey[100]
                                : theme.palette.grey[900],
                        flexGrow: 1,
                        height: '100vh',
                        overflow: 'auto',
                    }}
                >
                    <Toolbar />
                    <Container maxWidth="xl" sx={{ mt: 2, mb: 4 }}>
                        <Stack direction="row" spacing={2}>
                            <Card sx={{ minWidth: 700, maxWidth: 1200 }}>
                                {( selectedMainTabBuilder === 0 || addAnalysisView || addAggregateView || addRetrievalView || addEvalsView) &&
                                    <div>
                                        <CardContent>
                                            <Box sx={{ml: 2, mr: 2, mb: 2}}>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Workflow Generation
                                                </Typography>
                                            </Box>
                                            <Stack direction={"row"} >
                                                <Box sx={{ width: '50%',ml:2, mb: 0, mt: 2 }}>
                                                    <TextField
                                                        label={`Workflow Name`}
                                                        variant="outlined"
                                                        value={workflowName}
                                                        onChange={(event) => dispatch(setWorkflowName(event.target.value))}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box sx={{ width: '50%', mb: 0, mt: 2, ml: 2, mr:2 }}>
                                                    <TextField
                                                        label={`Workflow Group Name`}
                                                        variant="outlined"
                                                        value={workflowGroupName}
                                                        onChange={(event) => dispatch(setWorkflowGroupName(event.target.value))}
                                                        fullWidth
                                                    />
                                                </Box>
                                            </Stack>
                                            <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                <Divider/>
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 2}}>
                                                <Stack direction={"row"}>
                                                    <IconButton
                                                        onClick={() => setOpenRetrievals(!openRetrievals)}
                                                        aria-expanded={openRetrievals}
                                                        aria-label="show more"
                                                    >
                                                        {openRetrievals ? <ExpandLess /> : <ExpandMore />}
                                                    </IconButton>
                                                    <Typography gutterBottom variant="h5" component="div">
                                                        Retrievals
                                                    </Typography>
                                                </Stack>
                                            </Box>
                                            <Collapse in={openRetrievals} timeout="auto" unmountOnExit>
                                            <Box flexGrow={2} sx={{mt: 2}}>
                                                {retrievalStages && retrievalStages.map((ret, subIndex) => (
                                                    <Stack direction={"row"} key={subIndex} sx={{ mb: 2 }}>
                                                        <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                            <TextField
                                                                key={subIndex}
                                                                label={`Retrieval Name`}
                                                                value={ret?.retrievalName || ''}
                                                                InputProps={{
                                                                    readOnly: true,
                                                                }}
                                                                variant="outlined"
                                                                fullWidth
                                                                margin="normal"
                                                            />
                                                        </Box>
                                                        <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                            <TextField
                                                                key={subIndex}
                                                                label={`Retrieval Group`}
                                                                value={ret?.retrievalGroup || ''}
                                                                InputProps={{
                                                                    readOnly: true,
                                                                }}
                                                                variant="outlined"
                                                                fullWidth
                                                                margin="normal"
                                                            />
                                                        </Box>
                                                        <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                            <TextField
                                                                key={subIndex}
                                                                label={`Platform`}
                                                                value={ret?.retrievalItemInstruction.retrievalPlatform || ''}
                                                                InputProps={{
                                                                    readOnly: true,
                                                                }}
                                                                variant="outlined"
                                                                fullWidth
                                                                margin="normal"
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 0, ml: 2 }}>
                                                            <Button fullWidth variant="contained" onClick={(event)=>handleRemoveRetrievalFromWorkflow(event, ret)}>Remove</Button>
                                                        </Box>
                                                    </Stack>
                                                ))}
                                            </Box>
                                            <Box flexGrow={1} sx={{ mb: 0, mt: 2, ml: 2 }}>
                                                <Button  variant="contained" onClick={() => addRetrievalStageView()} >{addRetrievalView ? 'Done Adding': 'Add Retrieval Stages'}</Button>
                                            </Box>
                                        </Collapse>

                                        <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                <Divider/>
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 2, mb:2}}>
                                                <Stack direction={"row"} >
                                                    <IconButton
                                                        onClick={() => setOpenAnalysis(!openAnalysis)}
                                                        aria-expanded={openAnalysis}
                                                        aria-label="show more"
                                                    >
                                                        {openAnalysis ? <ExpandLess /> : <ExpandMore />}
                                                    </IconButton>
                                                    <Typography gutterBottom variant="h5" component="div">
                                                        Analysis
                                                    </Typography>
                                                </Stack>
                                            </Box>
                                            <Collapse in={openAnalysis} timeout="auto" unmountOnExit>
                                                <Box flexGrow={2} sx={{mt: 2}}>
                                                    {!loading && analysisStages && analysisStages.map((task, subIndex) => (
                                                        <Stack direction={"row"} key={subIndex} sx={{ mb: 2 }}>
                                                            <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                                <TextField
                                                                    key={subIndex}
                                                                    label={`Analysis Name`}
                                                                    value={task.taskName}
                                                                    InputProps={{
                                                                        readOnly: true,
                                                                    }}
                                                                    variant="outlined"
                                                                    fullWidth
                                                                    margin="normal"
                                                                />
                                                            </Box>
                                                            <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                                <TextField
                                                                    key={subIndex + task.taskGroup}
                                                                    label={`Analysis Group`}
                                                                    value={task.taskGroup}
                                                                    InputProps={{
                                                                        readOnly: true,
                                                                    }}
                                                                    variant="outlined"
                                                                    fullWidth
                                                                    margin="normal"
                                                                />
                                                            </Box>
                                                            <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                                <TextField
                                                                    key={subIndex+task.model}
                                                                    label={`Analysis Model`}
                                                                    value={task.model}
                                                                    InputProps={{
                                                                        readOnly: true,
                                                                    }}
                                                                    variant="outlined"
                                                                    fullWidth
                                                                    margin="normal"
                                                                />
                                                            </Box>
                                                            <Box flexGrow={1} sx={{ mb: 0, mt: -1, ml:2 }}>
                                                                <TextField
                                                                    type="number"
                                                                    label="Analysis Cycle Count"
                                                                    variant="outlined"
                                                                    value={(task?.taskID !== undefined && taskMap[task.taskID]?.cycleCount > 0) ? taskMap[task.taskID].cycleCount : 1}
                                                                    inputProps={{ min: 1 }}  // Set minimum value to 0
                                                                    onChange={(event) => handleTaskCycleCountChange(parseInt(event.target.value, 10), task)}
                                                                    fullWidth
                                                                />
                                                            </Box>
                                                            <Box flexGrow={1} sx={{ mb: 0, ml: 2 }}>
                                                                <Button fullWidth variant="contained" onClick={(event)=>handleRemoveAnalysisFromWorkflow(event, task)}>Remove</Button>
                                                            </Box>
                                                        </Stack>
                                                    ))}
                                                </Box>
                                                <Box flexGrow={1} sx={{ mb: 0, mt: 2, ml: 2 }}>
                                                    <Button  variant="contained" onClick={() => addAnalysisStageView()} >{addAnalysisView ? 'Done Adding': 'Add Analysis Stages'}</Button>
                                                </Box>
                                            <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                <Divider/>
                                            </Box>
                                            { retrievalStages && analysisStages && retrievalStages.length > 0 && analysisStages.length > 0 &&
                                                <div>
                                                    <Box sx={{ mt: 2, ml: 2 }} >
                                                        <Typography variant="h6" color="text.secondary">
                                                            Add Retrieval Stages to Analysis
                                                        </Typography>
                                                    </Box>
                                                    { workflowAnalysisRetrievalsMap &&
                                                        <Box sx={{ mt:2,  ml: 2, mr: 2 }} >
                                                            <Box >
                                                                {Object.entries(workflowAnalysisRetrievalsMap).map(([key, value], index) => {
                                                                    const retrievalNameForKey= retrievalsMap[(Number(key))]?.retrievalName || '';
                                                                    if (!retrievalNameForKey || retrievalNameForKey.length <= 0) {
                                                                        return null;
                                                                    }
                                                                    return Object.entries(value).map(([subKey, subValue], subIndex) => {
                                                                        const subKeyNumber = Number(subKey);
                                                                        const subTask = taskMap[(Number(subKeyNumber))]
                                                                        const subTaskName = subTask?.taskName || '';
                                                                        if (!subValue || subKey.length <= 0) {
                                                                            return null;
                                                                        }
                                                                        if (subTaskName.length <= 0) {
                                                                            return null;
                                                                        }
                                                                        return (
                                                                            <Stack direction={"row"} key={`${key}-${subKey}`}>
                                                                                <React.Fragment key={subIndex}>
                                                                                    <TextField
                                                                                        label={`Retrieval`}
                                                                                        value={retrievalNameForKey || ''}
                                                                                        InputProps={{ readOnly: true }}
                                                                                        variant="outlined"
                                                                                        fullWidth
                                                                                        margin="normal"
                                                                                    />
                                                                                    <Box flexGrow={1} sx={{ mt: 4, ml: 2, mr: 2 }}>
                                                                                        <ArrowForwardIcon />
                                                                                    </Box>
                                                                                    <TextField
                                                                                        label={`Analysis`}
                                                                                        value={subTaskName || ''}
                                                                                        InputProps={{ readOnly: true }}
                                                                                        variant="outlined"
                                                                                        fullWidth
                                                                                        margin="normal"
                                                                                    />
                                                                                    <Box flexGrow={1} sx={{mt: 3, ml: 2}}>
                                                                                        <Button variant="contained" onClick={(event) => handleRemoveRetrievalRelationshipFromWorkflow(event, key, subKeyNumber)}>Remove</Button>
                                                                                    </Box>
                                                                                </React.Fragment>
                                                                            </Stack>
                                                                        );
                                                                    });
                                                                })}
                                                            </Box>
                                                        </Box>
                                                    }
                                                    <Stack direction={"row"} sx={{ mt: 2 }}>
                                                        <Box flexGrow={3} sx={{ml: 2, mt: 2}}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id={`retrieval-stage-select-label-${1}`}>Retrieval</InputLabel>
                                                                <Select
                                                                    labelId={`retrieval-stage-select-label-${1}`}
                                                                    id={`retrieval-stage-select-${1}`}
                                                                    value={selectedRetrievalForAnalysis} // Use the state for the selected value
                                                                    label="Retrieval Source"
                                                                    onChange={(event) => setSelectedRetrievalForAnalysis(event.target.value)} // Update the state on change
                                                                >
                                                                    {retrieval && retrievalStages.map((ret, subIndex) => (
                                                                        <MenuItem key={subIndex} value={ret.retrievalID || 0}>{ret.retrievalName}</MenuItem>
                                                                    ))}
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                        <Box flexGrow={3} sx={{mt: 2, ml: 2, mr: 2}}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id={`ret-analysis-stage-select-label-${1}`}>Analysis</InputLabel>
                                                                <Select
                                                                    labelId={`ret-analysis-stage-select-label-${1}`}
                                                                    id={`ret-analysis-stage-select-${1}`}
                                                                    value={selectedAnalysisStageForRetrieval} // Use the state for the selected value
                                                                    label="Analysis Source"
                                                                    onChange={(event) => setSelectedAnalysisStageForRetrieval(event.target.value)} // Update the state on change
                                                                >
                                                                    {analysisStages && analysisStages.map((stage, subIndex) => (
                                                                        <MenuItem key={subIndex} value={stage.taskID}>{stage.taskName}</MenuItem>
                                                                    ))}
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                    </Stack>
                                                    <Box flexGrow={3} sx={{mt: 2, ml: 2, mr: 2 }}>
                                                        <Button variant="contained" onClick={handleAddRetrievalToAnalysis}>Add Source</Button>
                                                    </Box>
                                                </div>
                                            }
                                            </Collapse>
                                            { analysisStages && analysisStages.length > 0 &&
                                                <div>
                                                    <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                        <Divider/>
                                                    </Box>
                                                    <Box flexGrow={2} sx={{mt:2 , mb: 2}}>
                                                        <Stack direction={"row"}>
                                                            <IconButton
                                                                onClick={() => setOpenAggregation(!openAggregation)}
                                                                aria-expanded={openAggregation}
                                                                aria-label="show more"
                                                            >
                                                                {openAggregation ? <ExpandLess /> : <ExpandMore />}
                                                            </IconButton>
                                                            <Typography gutterBottom variant="h5" component="div">
                                                                Aggregation
                                                            </Typography>
                                                        </Stack>
                                                    </Box>
                                                    <Collapse in={openAggregation} timeout="auto" unmountOnExit>
                                                        <Box flexGrow={2} sx={{ml: 2 , mb: 0}}>
                                                            <Typography gutterBottom variant="body2" component="div">
                                                                One aggregation cycle is equal to the longest of any dependent analysis cycles.
                                                                If you have an analysis stage that occurs every 2 time cycles, and set the aggregation cycle count to 2,
                                                                it will run on time cycle 4 after the analysis stage completes.
                                                            </Typography>
                                                        </Box>
                                                    <Box flexGrow={2} sx={{mt: 4}}>
                                                        <Stack direction={"column"} key={0}>
                                                        {aggregationStages && aggregationStages.map((task, subIndex) => (
                                                                <Stack direction={"row"} key={subIndex} sx={{ mb: 1 }}>
                                                                    <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                                        <TextField
                                                                            key={subIndex}
                                                                            label={`Aggregation Name`}
                                                                            value={task.taskName}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                            variant="outlined"
                                                                            fullWidth
                                                                            margin="normal"
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                                        <TextField
                                                                            key={subIndex}
                                                                            label={`Aggregation Group`}
                                                                            value={task.taskGroup}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                            variant="outlined"
                                                                            fullWidth
                                                                            margin="normal"
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                                        <TextField
                                                                            key={subIndex+task.model}
                                                                            label={`Aggregation Model`}
                                                                            value={task.model}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                            variant="outlined"
                                                                            fullWidth
                                                                            margin="normal"
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={2} sx={{ mb: 0, mt: -1, ml:2 }}>
                                                                        <TextField
                                                                            type="number"
                                                                            label="Aggregation Cycle Count"
                                                                            variant="outlined"
                                                                            value={task?.taskID !== undefined && taskMap[task.taskID]?.cycleCount > 0 ? taskMap[task.taskID].cycleCount : 1}
                                                                            inputProps={{ min: 1 }}  // Set minimum value to 0
                                                                            onChange={(event) => handleTaskCycleCountChange(parseInt(event.target.value, 10),task)}
                                                                            fullWidth
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ mb: 4, ml: 2 }}>
                                                                        <Button fullWidth variant="contained" onClick={(event)=> handleRemoveAggregationFromWorkflow(event,task)}>Remove</Button>
                                                                    </Box>
                                                                </Stack>
                                                                ))}
                                                                <Box flexGrow={1} sx={{ mb: 4, mt: -1, ml: 2 }}>
                                                                    <Button variant="contained" onClick={() => addAggregationStageView()} >{addAggregateView ? 'Done Adding' : 'Add Aggregation Stages' }</Button>
                                                                </Box>
                                                                <Divider />
                                                            {aggregationStages.length > 0 && analysisStages.length > 0 &&
                                                                <div>
                                                                    <Box sx={{ mt: 2, ml: 2 }} >
                                                                        <Typography variant="h6" color="text.secondary">
                                                                            Analysis {'->'} Aggregation Dependencies
                                                                        </Typography>
                                                                    </Box>
                                                                    { workflowBuilderTaskMap &&
                                                                        <Box sx={{ mt:2,  ml: 2, mr: 2 }} >
                                                                            <Box >
                                                                                {Object.entries(workflowBuilderTaskMap).map(([key, value], index) => {
                                                                                    const taskNameForKey = taskMap[Number(key)]?.taskName || '';
                                                                                    if (!taskNameForKey) {
                                                                                        return null;
                                                                                    }
                                                                                    return Object.entries(value).map(([subKey, subValue], subIndex) => {
                                                                                        if (!subValue) {
                                                                                            return null;
                                                                                        }
                                                                                        const subKeyNumber = Number(subKey);
                                                                                        const subTaskName = taskMap[subKeyNumber]?.taskName || '';
                                                                                        if (!subTaskName) {
                                                                                            return null;
                                                                                        }
                                                                                        return (
                                                                                            <Stack direction={"row"} key={`${key}-${subKey}`}>
                                                                                                <React.Fragment key={subIndex}>
                                                                                                    <TextField
                                                                                                        label={`Analysis`}
                                                                                                        value={subTaskName || ''}
                                                                                                        InputProps={{ readOnly: true }}
                                                                                                        variant="outlined"
                                                                                                        fullWidth
                                                                                                        margin="normal"
                                                                                                    />
                                                                                                    <Box flexGrow={1} sx={{ mt: 4, ml: 2, mr: 2 }}>
                                                                                                        <ArrowForwardIcon />
                                                                                                    </Box>
                                                                                                    <TextField
                                                                                                        label={`Aggregate`}
                                                                                                        value={taskNameForKey || ''}
                                                                                                        InputProps={{ readOnly: true }}
                                                                                                        variant="outlined"
                                                                                                        fullWidth
                                                                                                        margin="normal"
                                                                                                    />
                                                                                                    <Box flexGrow={1} sx={{mt: 3, ml: 2}}>
                                                                                                        <Button variant="contained" onClick={(event) => handleRemoveTaskRelationshipFromWorkflow(event, key, subKeyNumber)}>Remove</Button>
                                                                                                    </Box>
                                                                                                </React.Fragment>
                                                                                            </Stack>
                                                                                        );
                                                                                    });
                                                                                })}
                                                                            </Box>
                                                                        </Box>
                                                                    }

                                                                    { analysisStages &&
                                                                        <div>
                                                                            <Box sx={{ mt: 4 }} >
                                                                                <Divider />
                                                                            </Box>
                                                                            <Box sx={{ mt: 2, ml: 2 }} >
                                                                                <Typography variant="h6" color="text.secondary">
                                                                                    Add Analysis Stages to Aggregates
                                                                                </Typography>
                                                                            </Box>
                                                                        </div>
                                                                    }
                                                                    <Stack sx={{ mt: 6, ml: 0 }} direction={"row"} key={1}>
                                                                        { analysisStages &&
                                                                        <Box flexGrow={3} sx={{ mt: -3, ml: 2 }}>
                                                                            <FormControl fullWidth>
                                                                                <InputLabel id={`analysis-stage-select-label-${1}`}>Analysis</InputLabel>
                                                                                <Select
                                                                                    labelId={`analysis-stage-select-label-${1}`}
                                                                                    id={`analysis-stage-select-${1}`}
                                                                                    value={selectedAnalysisStageForAggregation} // Use the state for the selected value
                                                                                    label="Analysis Source"
                                                                                    onChange={(event) => setSelectedAnalysisStageForAggregation(event.target.value)} // Update the state on change
                                                                                >
                                                                                    {analysisStages && analysisStages.map((stage, subIndex) => (
                                                                                        <MenuItem key={subIndex} value={stage.taskID}>{stage.taskName}</MenuItem>
                                                                                    ))}
                                                                                </Select>
                                                                            </FormControl>
                                                                        </Box>
                                                                        }
                                                                        <Box flexGrow={1} sx={{ mt: -1.2, ml: 6, mr: -1 }}>
                                                                            <ArrowForwardIcon />
                                                                        </Box>
                                                                        { aggregationStages &&
                                                                        <Box flexGrow={3} sx={{ mt: -3, ml: 0 }}>
                                                                            <FormControl fullWidth>
                                                                                <InputLabel id={`agg-stage-select-label-${2}`}>Aggregate</InputLabel>
                                                                                <Select
                                                                                    labelId={`agg-stage-select-label-${2}`}
                                                                                    id={`agg-stage-select-${2}`}
                                                                                    value={selectedAggregationStageForAnalysis} // This should correspond to the selected stage for each task
                                                                                    label="Aggregate Source"
                                                                                    onChange={(event) => setSelectedAggregationStageForAnalysis(event.target.value)} // Update the state on change)
                                                                                >
                                                                                    {aggregationStages && aggregationStages.map((stage: any, subIndex: number) => (
                                                                                        <MenuItem key={stage.taskID} value={stage.taskID}>{stage.taskName}</MenuItem>
                                                                                    ))}
                                                                                </Select>
                                                                            </FormControl>
                                                                        </Box>
                                                                        }
                                                                        <Box flexGrow={3} sx={{mt: -2, ml: 2 }}>
                                                                            <Button variant="contained" onClick={handleAddSubTaskToAggregate}>Add Source</Button>
                                                                        </Box>
                                                                    </Stack>
                                                                </div>
                                                            }
                                                            </Stack>
                                                    </Box>
                                                    </Collapse>
                                                </div>
                                            }
                                        </CardContent>
                                        <Box flexGrow={1} sx={{ mt: 0, mb: 0}}>
                                            <Divider/>
                                        </Box>
                                        <Box flexGrow={2} sx={{mt: 2, ml: 2}}>
                                            <Stack direction={"row"}>
                                                <IconButton
                                                    onClick={() => setOpenEvals(!openEvals)}
                                                    aria-expanded={openEvals}
                                                    aria-label="show more"
                                                >
                                                    {openEvals ? <ExpandLess /> : <ExpandMore />}
                                                </IconButton>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Automated Evals
                                                </Typography>
                                            </Stack>
                                        </Box>
                                        <Collapse in={openEvals} timeout="auto" unmountOnExit>
                                            <Box flexGrow={2} sx={{ml: 4 , mr:2, mb: 0}}>
                                                <Typography gutterBottom variant="body2" component="div">
                                                    One eval cycle is equal N * the attached task cycle.
                                                    If you have an analysis stage that occurs every 2 time cycles, and set the eval cycle count to 2,
                                                    it will run on time cycle 4 after the analysis stage completes for the second time.
                                                </Typography>
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 2}}>
                                                {addedEvalFns && addedEvalFns.map((ef: EvalFn, subIndex) => (
                                                    <Stack direction={"row"} key={subIndex} sx={{ mb: 2 }}>
                                                        <Box flexGrow={2} sx={{ mt: -1, ml: 4 }}>
                                                            <TextField
                                                                key={subIndex}
                                                                label={`Eval Name`}
                                                                value={ef.evalName}
                                                                InputProps={{
                                                                    readOnly: true,
                                                                }}
                                                                variant="outlined"
                                                                fullWidth
                                                                margin="normal"
                                                            />
                                                        </Box>
                                                        <Box flexGrow={2} sx={{ mt: 1, ml: 2 }}>
                                                            <TextField
                                                                type="number"
                                                                label="Eval Cycle Count"
                                                                variant="outlined"
                                                                value={
                                                                    ef.evalID !== undefined && evalMap[ef.evalID] !== undefined &&
                                                                    evalMap[ef.evalID].evalCycleCount? evalMap[ef.evalID].evalCycleCount : 1
                                                                }
                                                                inputProps={{ min: 1 }}  // Set minimum value to 0
                                                                onChange={(event) => handleEvalCycleCountChange(parseInt(event.target.value, 10), ef)}
                                                                fullWidth
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mt: 2, mb: 0, ml: 2, mr: 4 }}>
                                                            <Button fullWidth variant="contained" onClick={(event)=>handleRemoveEvalFnFromWorkflow(event, ef)}>Remove Eval</Button>
                                                        </Box>
                                                    </Stack>
                                                ))}
                                            </Box>
                                            <Box flexGrow={1} sx={{ mb: 0, mt: 2, ml: 4 }}>
                                                <Button  variant="contained" onClick={() => addEvalsStageView()} >{addEvalsView ? 'Done Adding': 'Add Eval Stages'}</Button>
                                            </Box>
                                            { evalFnStages && ((analysisStages && analysisStages.length > 0) || (aggregationStages && aggregationStages.length > 0)) && evalFnStages.length > 0 &&
                                                <div>
                                                    <Box flexGrow={1} sx={{ mt: 4, mb: 0}}>
                                                        <Divider/>
                                                    </Box>
                                                    <Box sx={{ mt:2, ml: 4 }} >
                                                        <Typography variant="h6" color="text.secondary">
                                                            Analysis/Aggregation {'->'} Eval Dependencies
                                                        </Typography>
                                                    </Box>
                                                    { workflowBuilderEvalsTaskMap &&
                                                        <Box sx={{ mt:2,  ml: 2, mr: 2 }} >
                                                            <Box >
                                                                {Object.entries(workflowBuilderEvalsTaskMap).map(([key, value], index) => {
                                                                    // these are the tasks
                                                                    if (key === undefined){
                                                                        return null;
                                                                    }
                                                                    const taskNameForKey= taskMap[(Number(key))]?.taskName || '';
                                                                    if (!taskNameForKey || taskNameForKey.length <= 0) {
                                                                        return null;
                                                                    }
                                                                    return Object.entries(value).map(([subKey, subValue], subIndex) => {
                                                                        // these are evals
                                                                        if (subKey === undefined){
                                                                            return null;
                                                                        }
                                                                        const evalID = Number(subKey);
                                                                        const evalFn = evalMap[(Number(evalID))]
                                                                        const subTaskName = evalFn?.evalName || '';

                                                                        if (!subValue || subKey.length <= 0) {
                                                                            return null;
                                                                        }
                                                                        if (subTaskName.length <= 0) {
                                                                            return null;
                                                                        }
                                                                        return (
                                                                            <Stack direction={"row"} key={`${key}-${subKey}`}>
                                                                                <React.Fragment key={subIndex}>
                                                                                    <Box flexGrow={2} sx={{ mt: 0, ml: 2, mr: 2 }}>
                                                                                        <TextField
                                                                                            label={`Task`}
                                                                                            value={taskNameForKey || ''}
                                                                                            InputProps={{ readOnly: true }}
                                                                                            variant="outlined"
                                                                                            fullWidth
                                                                                            margin="normal"
                                                                                        />
                                                                                    </Box>
                                                                                    <Box flexGrow={1} sx={{ mt: 4, ml: 4, mr: -2 }}>
                                                                                        <ArrowForwardIcon />
                                                                                    </Box>
                                                                                    <Box flexGrow={2} sx={{ mt: 0, ml: 0, mr: 2 }}>
                                                                                        <TextField
                                                                                            label={`Eval`}
                                                                                            value={subTaskName || ''}
                                                                                            InputProps={{ readOnly: true }}
                                                                                            variant="outlined"
                                                                                            fullWidth
                                                                                            margin="normal"
                                                                                        />
                                                                                    </Box>
                                                                                    <Box flexGrow={1} sx={{mt: 3, ml: 2}}>
                                                                                        <Button variant="contained" onClick={(event) => handleRemoveEvalFnRelationshipFromWorkflow(event, evalID, Number(key))}>Remove</Button>
                                                                                    </Box>
                                                                                </React.Fragment>
                                                                            </Stack>
                                                                        );
                                                                    });
                                                                })}
                                                            </Box>
                                                        </Box>}
                                                    <Box flexGrow={1} sx={{ mt: 4, mb: 0}}>
                                                        <Divider/>
                                                    </Box>
                                                    <Box sx={{ mt:2, ml: 4 }} >
                                                        <Typography variant="h6" color="text.secondary">
                                                            Connect (Aggregation/Analysis) Task Outputs  {'->'} Eval Stages
                                                        </Typography>
                                                    </Box>
                                                        <Stack sx={{ mt: 6, ml: 0 }} direction={"row"} key={1}>
                                                            { analysisStages && !toggleEvalToTaskType &&
                                                                <Box flexGrow={3} sx={{ mt: -3, ml: 4, mr: 2 }}>
                                                                    <FormControl fullWidth>
                                                                        <InputLabel id={`analysis-stage-select-label-${1}`}>Analysis</InputLabel>
                                                                        <Select
                                                                            labelId={`analysis-stage-select-label-${1}`}
                                                                            id={`analysis-stage-select-${1}`}
                                                                            value={selectedAnalysisStageForEval} // Use the state for the selected value
                                                                            label="Analysis Source"
                                                                            onChange={(event) => setSelectedAnalysisStageForEval(event.target.value)} // Update the state on change
                                                                        >
                                                                            {analysisStages && analysisStages.map((stage, subIndex) => (
                                                                                <MenuItem key={subIndex} value={stage.taskID}>{stage.taskName}</MenuItem>
                                                                            ))}
                                                                        </Select>
                                                                    </FormControl>
                                                                </Box>
                                                            }
                                                            { aggregationStages && toggleEvalToTaskType &&
                                                                <Box flexGrow={3} sx={{ mt: -3, ml: 4, mr: 2 }}>
                                                                    <FormControl fullWidth>
                                                                        <InputLabel id={`agg-stage-select-label-${2}`}>Aggregate</InputLabel>
                                                                        <Select
                                                                            labelId={`agg-stage-select-label-${2}`}
                                                                            id={`agg-stage-select-${2}`}
                                                                            value={selectedAggregationStageForEval} // This should correspond to the selected stage for each task
                                                                            label="Aggregate Source"
                                                                            onChange={(event) => setSelectedAggregationStageForEval(event.target.value)} // Update the state on change)
                                                                        >
                                                                            {aggregationStages && aggregationStages.map((stage: any, subIndex: number) => (
                                                                                <MenuItem key={stage.taskID} value={stage.taskID}>{stage.taskName}</MenuItem>
                                                                            ))}
                                                                        </Select>
                                                                    </FormControl>
                                                                </Box>
                                                            }
                                                            <Box flexGrow={1} sx={{ mt: -1.2, ml: 5, mr: -4 }}>
                                                                <ArrowForwardIcon />
                                                            </Box>
                                                            { evalFnStages &&
                                                                <Box flexGrow={3} sx={{ mt: -3, ml: 4 }}>
                                                                    <FormControl fullWidth>
                                                                        <InputLabel id={`eval-stage-select-label-${1}`}>Evals</InputLabel>
                                                                        <Select
                                                                            labelId={`eval-stage-select-label-${1}`}
                                                                            id={`eval-stage-select-${1}`}
                                                                            value={selectedEvalStage} // Use the state for the selected value
                                                                            label="Eval Source"
                                                                            onChange={(event) => setSelectedEvalStage(event.target.value)} // Update the state on change
                                                                        >
                                                                            {evalFnStages && evalFnStages.map((stage, subIndex) => (
                                                                                <MenuItem key={subIndex} value={stage.evalID ? stage.evalID : ''}>{stage.evalName}</MenuItem>
                                                                            ))}
                                                                        </Select>
                                                                    </FormControl>
                                                                </Box>
                                                            }
                                                            <Box  sx={{mt: -2, ml: 2, mr: 0 }}>
                                                                <Button variant="contained" onClick={(e) => handleAddEvalToSubTask(e)}>Attach Eval</Button>
                                                            </Box>
                                                            <Box sx={{mt: -2, ml: 2, mr: 2 }}>
                                                                <Button variant="contained" onClick={setToggleEvalTaskType}>Toggle {toggleEvalToTaskType ? 'Analysis':'Aggregation'} Evals</Button>
                                                            </Box>
                                                        </Stack>
                                                        </div>
                                            }

                                        </Collapse>
                                        <Box flexGrow={1} sx={{ mt: 4, mb: 0}}>
                                            <Divider/>
                                        </Box>
                                        <CardContent>
                                            <Box sx={{ ml: 2, mr: 2 }} >
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Fundamental Time Period
                                                </Typography>
                                            </Box>
                                            <Box sx={{ ml: 2, mr: 2 }} >
                                                <Typography variant="body2" color="text.secondary">
                                                    This is the time period that each cycle is referenced against for determining its next execution time. The workflow will run every time period, and will run all analysis and aggregation stages that are due to run if
                                                    any during that discrete time step. If an analysis cycle is set to 1 and the fundamental time period is 5 minutes, it will run every 5 minutes.
                                                </Typography>
                                            </Box>
                                            <Stack direction="row" spacing={2} sx={{ ml: 2, mr: 2, mt: 4, mb: 2 }}>
                                                <Box sx={{ width: '33%' }}> {/* Adjusted Box for TextField */}
                                                    <TextField
                                                        type="number"
                                                        label="Time Step Size"
                                                        variant="outlined"
                                                        inputProps={{ min: 1 }}  // Set minimum value to 1
                                                        value={stepSize}
                                                        onChange={handleTimeStepChange}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box sx={{ width: '33%' }}> {/* Adjusted Box for FormControl */}
                                                    <FormControl fullWidth>
                                                        <InputLabel id="time-unit-label">Time Unit</InputLabel>
                                                        <Select
                                                            labelId="time-unit-label"
                                                            id="time-unit-select"
                                                            value={stepSizeUnit}
                                                            label="Time Unit"
                                                            onChange={handleUpdateStepSizeUnit}
                                                        >
                                                            <MenuItem value="minutes">Minutes</MenuItem>
                                                            <MenuItem value="hours">Hours</MenuItem>
                                                            <MenuItem value="days">Days</MenuItem>
                                                            <MenuItem value="weeks">Weeks</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>
                                            </Stack>
                                        </CardContent>
                                        <CardActions>
                                            <Box flexGrow={1} sx={{ mb: -7, mt: -8, ml: 2, mr: 2}}>
                                                <Button fullWidth variant="contained" onClick={createOrUpdateWorkflow} >Save Workflow</Button>
                                            </Box>
                                        </CardActions>
                                            {requestStatus !== '' && (
                                                <Container sx={{  mt: 2}}>
                                                    <Typography variant="h6" color={requestStatusError}>
                                                        {requestStatus}
                                                    </Typography>
                                                </Container>
                                            )}
                                </div>
                                }
                                <CardContent>
                                    {!addAnalysisView && !addAggregateView && !addRetrievalView && (selectedMainTabBuilder === 1 || (taskType === 'analysis' && addSchemasView && selectedMainTabBuilder === 7)) && !addEvalsView &&
                                        <div>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Analysis Instructions
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Token overflow strategy will determine how the AI will handle requests that are projected to exceed the maximum token length for the model you select, or has returned a result with that error.
                                                Deduce will chunk your analysis into smaller pieces and aggregate them into a final analysis result. Truncate will simply truncate the request
                                                to the maximum token length it can support. If you set the max tokens field greater than 0, it becomes the maximum number of tokens to spend per task request.
                                            </Typography>
                                            <Stack direction="row" >
                                                <Box flexGrow={3} sx={{ width: '50%', mb: 0, mt: 2, mr: 1 }}>
                                                    <TextField
                                                        label={`Analysis Name`}
                                                        variant="outlined"
                                                        value={editAnalysisTask.taskName}
                                                        onChange={(event) => dispatch(setEditAnalysisTask({ ...editAnalysisTask, taskName: event.target.value }))}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box flexGrow={3} sx={{ width: '50%', mb: 0, mt: 2, ml: 1 }}>
                                                    <TextField
                                                        label={`Analysis Group`}
                                                        variant="outlined"
                                                        value={editAnalysisTask.taskGroup}
                                                        onChange={(event) => dispatch(setEditAnalysisTask({ ...editAnalysisTask, taskGroup: event.target.value }))}
                                                        fullWidth
                                                    />
                                                </Box>
                                            </Stack>
                                            <Stack direction="row" >
                                                { editAggregateTask.responseFormat === 'json' ?
                                                    <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                        <FormControl fullWidth>
                                                            <InputLabel id="analysis-model-label">Analysis Model</InputLabel>
                                                            <Select
                                                                labelId="analysis-model-label"
                                                                id="analysis-model-select"
                                                                value={editAnalysisTask.model}
                                                                label="Analysis Model"
                                                                onChange={(event) => dispatch(setEditAnalysisTask({ ...editAnalysisTask, model: event.target.value }))}
                                                            >
                                                                <MenuItem value="gpt-3.5-turbo-1106">gpt-3.5-turbo-1106</MenuItem>
                                                                <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
                                                            </Select>
                                                        </FormControl>
                                                    </Box>
                                                    :
                                                    <Box flexGrow={3} sx={{ mb: 4, mt: 4 }}>
                                                        <FormControl fullWidth>
                                                            <InputLabel id="analysis-model-label">Analysis Model</InputLabel>
                                                            <Select
                                                                labelId="analysis-model-label"
                                                                id="analysis-model-select"
                                                                value={editAnalysisTask.model}
                                                                label="Analysis Model"
                                                                onChange={(event) => dispatch(setEditAnalysisTask({ ...editAnalysisTask, model: event.target.value }))}
                                                            >
                                                                <MenuItem value="gpt-3.5-turbo-instruct">gpt-3.5-turbo-instruct</MenuItem>
                                                                <MenuItem value="gpt-3.5-turbo-1106">gpt-3.5-turbo-1106</MenuItem>
                                                                <MenuItem value="gpt-3.5-turbo">gpt-3.5-turbo</MenuItem>
                                                                <MenuItem value="gpt-4">gpt-4</MenuItem>
                                                                <MenuItem value="gpt-4-32k">gpt-4-32k</MenuItem>
                                                                <MenuItem value="gpt-4-32k-0613">gpt-4-32k-0613</MenuItem>
                                                                <MenuItem value="gpt-4-0613">gpt-4-0613</MenuItem>
                                                                <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
                                                                <MenuItem value="babbage-002">babbage-002</MenuItem>
                                                                <MenuItem value="davinci-002">davinci-002</MenuItem>
                                                            </Select>
                                                        </FormControl>
                                                    </Box>
                                                }
                                                <Box flexGrow={2} sx={{ mb: 4, mt: 4, ml:2 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="token-overflow-analysis-label">Token Overflow Strategy</InputLabel>
                                                        <Select
                                                            labelId="token-overflow-analysis-label"
                                                            id="analysis-overflow-analysis-select"
                                                            value={editAnalysisTask.tokenOverflowStrategy}
                                                            label="Token Overflow Strategy"
                                                            onChange={(event) => dispatch(setEditAnalysisTask({ ...editAnalysisTask, tokenOverflowStrategy: event.target.value }))} // Dispatch action directly
                                                        >
                                                            <MenuItem value="deduce">deduce</MenuItem>
                                                            <MenuItem value="truncate">truncate</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>
                                            </Stack>
                                            { editAnalysisTask.responseFormat === 'json' || editAnalysisTask.responseFormat === 'social-media-engagement' ?
                                                <div>
                                                    <Box  sx={{ mb: 4, mt: 0, ml: -1}}>
                                                        <Button variant="contained" color="secondary" onClick={addSchemasViewToggle} style={{marginLeft: '10px'}}>
                                                            { addSchemasView ? 'Done Adding':'Add Schemas' }
                                                        </Button>
                                                    </Box>
                                                    { editAnalysisTask.schemas && editAnalysisTask.schemas.length > 0 &&
                                                        editAnalysisTask.schemas.map((schema: JsonSchemaDefinition, index: number) => (
                                                            <div>
                                                                <Stack direction="row" key={index}>
                                                                    <Box sx={{ mb: 2, mt: 2, width: '50%' }}>
                                                                        <TextField
                                                                            key={`schema-name-${index}`}
                                                                            fullWidth
                                                                            id={`schema-${index}`}
                                                                            label={`Schema-Name-${index}`}
                                                                            variant="outlined"
                                                                            value={schema.schemaName}
                                                                            InputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Box sx={{ mb: 2, mt: 2, ml: 2, width: '50%' }}>
                                                                        <TextField
                                                                            key={`schema-group-${index}`}
                                                                            fullWidth
                                                                            id={`schema-group-${index}`}
                                                                            label={`Schema-Group-${index}`}
                                                                            variant="outlined"
                                                                            value={schema.schemaGroup}
                                                                            InputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Box sx={{ ml: 2, mb: 2, mt: 3 }}>
                                                                        <Button
                                                                            variant="contained"
                                                                            id={`sm-button-${index}`}
                                                                            color="primary"
                                                                            key={`sm-button-${index}`}
                                                                            onClick={(e) => removeSchemasViewToggle(e, index)}
                                                                        >
                                                                            Remove
                                                                        </Button>
                                                                    </Box>
                                                                </Stack>
                                                            </div>

                                                        ))
                                                    }
                                                </div>
                                                :
                                                <Box  sx={{ mb: 2, mt: -2 }}>
                                                    <TextareaAutosize
                                                        minRows={18}
                                                        value={editAnalysisTask.prompt}
                                                        onChange={(event) => dispatch(setEditAnalysisTask({ ...editAnalysisTask, prompt: event.target.value }))} // Dispatch action directly
                                                        style={{ resize: "both", width: "100%" }}
                                                    />
                                                </Box>
                                            }
                                        </div>
                                    }
                                    { !addAggregateView && !addAnalysisView && !addRetrievalView && !addEvalsView && selectedMainTabBuilder === 2 &&
                                        <div>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Aggregation Instructions
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Use this to tell the AI how to aggregate the results of your analysis chunks into a rolling aggregation window. If aggregating on a single analysis, the aggregation cycle count sets how many
                                                base analysis cycles to aggregate on. If aggregating on multiple analysis, it will aggregate whenever the the underlying analysis is run.
                                            </Typography>
                                            <Stack direction="row" >
                                                <Box flexGrow={3} sx={{ width: '50%', mb: 0, mt: 2, mr: 1 }}>
                                                    <TextField
                                                        label={`Aggregation Name`}
                                                        variant="outlined"
                                                        value={editAggregateTask.taskName}
                                                        onChange={(event) => dispatch(setEditAggregateTask({ ...editAggregateTask, taskName: event.target.value }))}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box flexGrow={3} sx={{ width: '50%', mb: 0, mt: 2, ml: 1 }}>
                                                    <TextField
                                                        label={`Aggregation Group`}
                                                        variant="outlined"
                                                        value={editAggregateTask.taskGroup}
                                                        onChange={(event) => dispatch(setEditAggregateTask({ ...editAggregateTask, taskGroup: event.target.value }))}
                                                        fullWidth
                                                    />
                                                </Box>
                                            </Stack>
                                            <Stack direction="row" >
                                                { editAggregateTask.responseFormat === 'json' ?
                                                    <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="model-label">Aggregation Model</InputLabel>
                                                        <Select
                                                            labelId="model-label"
                                                            id="model-select"
                                                            value={editAggregateTask.model}
                                                            label="Aggregation Model"
                                                            onChange={(event) => handleEditAggregateTaskModel(event)}
                                                        >
                                                            <MenuItem value="gpt-3.5-turbo-1106">gpt-3.5-turbo-1106</MenuItem>
                                                            <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>
                                                    :
                                                    <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                        <FormControl fullWidth>
                                                            <InputLabel id="model-label">Aggregation Model</InputLabel>
                                                            <Select
                                                                labelId="model-label"
                                                                id="model-select"
                                                                value={editAggregateTask.model}
                                                                label="Aggregation Model"
                                                                onChange={(event) => handleEditAggregateTaskModel(event)}
                                                            >
                                                                <MenuItem value="gpt-3.5-turbo-instruct">gpt-3.5-turbo-instruct</MenuItem>
                                                                <MenuItem value="gpt-3.5-turbo-1106">gpt-3.5-turbo-1106</MenuItem>
                                                                <MenuItem value="gpt-3.5-turbo">gpt-3.5-turbo</MenuItem>
                                                                <MenuItem value="gpt-4">gpt-4</MenuItem>
                                                                <MenuItem value="gpt-4-32k">gpt-4-32k</MenuItem>
                                                                <MenuItem value="gpt-4-32k-0613">gpt-4-32k-0613</MenuItem>
                                                                <MenuItem value="gpt-4-0613">gpt-4-0613</MenuItem>
                                                                <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
                                                                <MenuItem value="babbage-002">babbage-002</MenuItem>
                                                                <MenuItem value="davinci-002">davinci-002</MenuItem>

                                                            </Select>
                                                        </FormControl>
                                                    </Box>
                                                }
                                                <Box flexGrow={2} sx={{ mb: 2, mt: 4, ml:2 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="aggregation-token-overflow-analysis-label">Token Overflow Strategy</InputLabel>
                                                        <Select
                                                            labelId="aggregation-token-overflow-analysis-label"
                                                            id="aggregation-token-overflow-analysis-select"
                                                            value={editAggregateTask.tokenOverflowStrategy}
                                                            label="Token Overflow Strategy"
                                                            onChange={(event) => dispatch(setEditAggregateTask({ ...editAggregateTask, tokenOverflowStrategy: event.target.value }))}
                                                        >
                                                            <MenuItem value="deduce">deduce</MenuItem>
                                                            <MenuItem value="truncate">truncate</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>

                                            </Stack>
                                            { editAggregateTask.responseFormat === 'json' || editAggregateTask.responseFormat === 'social-media-engagement'?
                                                <div>
                                                    <Box  sx={{ mb: 4, mt: 2, ml: -1 }}>
                                                        <Button variant="contained" color="secondary" onClick={addSchemasViewToggle} style={{marginLeft: '10px'}}>
                                                            { addSchemasView ? 'Done Adding':'Add Schemas' }
                                                        </Button>
                                                    </Box>
                                                        { editAggregateTask.schemas && editAggregateTask.schemas.length > 0 &&
                                                            editAggregateTask.schemas.map((schema: JsonSchemaDefinition, index: number) => (
                                                                <div>
                                                                    <Stack direction="row" key={index}>
                                                                        <Box sx={{ mb: 2, mt: 2, width: '50%' }}>
                                                                            <TextField
                                                                                key={`schema-name-${index}`}
                                                                                fullWidth
                                                                                id={`schema-${index}`}
                                                                                label={`Schema-Name-${index}`}
                                                                                variant="outlined"
                                                                                value={schema.schemaName}
                                                                                InputProps={{ readOnly: true }}
                                                                            />
                                                                        </Box>
                                                                        <Box sx={{ mb: 2, mt: 2, ml: 2, width: '50%' }}>
                                                                            <TextField
                                                                                key={`schema-group-${index}`}
                                                                                fullWidth
                                                                                id={`schema-group-${index}`}
                                                                                label={`Schema-Group-${index}`}
                                                                                variant="outlined"
                                                                                value={schema.schemaGroup}
                                                                                InputProps={{ readOnly: true }}
                                                                            />
                                                                        </Box>
                                                                        <Box sx={{ ml: 2, mb: 2, mt: 3 }}>
                                                                            <Button
                                                                                variant="contained"
                                                                                id={`sm-button-${index}`}
                                                                                color="primary"
                                                                                key={`sm-button-${index}`}
                                                                                onClick={(e) => removeSchemasViewToggle(e, index)}
                                                                            >
                                                                                Remove
                                                                            </Button>
                                                                        </Box>
                                                                    </Stack>
                                                                </div>

                                                            ))
                                                        }
                                                </div>
                                                :
                                                <Box  sx={{ mb: 2, mt: 2 }}>
                                                    <TextareaAutosize
                                                        minRows={18}
                                                        value={editAggregateTask.prompt}
                                                        onChange={(event) => dispatch(setEditAggregateTask({ ...editAggregateTask, prompt: event.target.value }))}
                                                        style={{ resize: "both", width: "100%" }}
                                                    />
                                                </Box>
                                            }
                                        </div>
                                    }
                                    {
                                        selectedMainTabBuilder == 3 && !addRetrievalView && !loading && !addAnalysisView && !addAggregateView && !addEvalsView &&
                                    <CardContent>
                                        <div>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Retrieval Procedures
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                This allows you to modularize your retrieval procedures and reuse them across multiple workflows.
                                            </Typography>
                                            <Stack direction="column" spacing={2} sx={{ mt: 4, mb: 0 }}>
                                                <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="retrieval-name"
                                                                label="Retrieval Name"
                                                                variant="outlined"
                                                                value={retrieval.retrievalName}
                                                                onChange={(event) => dispatch(setRetrieval({ ...retrieval, retrievalName: event.target.value }))}
                                                            />
                                                        </Box>
                                                    <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                        <TextField
                                                            fullWidth
                                                            id="retrieval-group"
                                                            label="Retrieval Group"
                                                            variant="outlined"
                                                            value={retrieval.retrievalGroup}
                                                            onChange={(event) => dispatch(setRetrieval({ ...retrieval, retrievalGroup: event.target.value }))}
                                                        />
                                                    </Box>
                                                    </Stack>
                                                    <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="platform-label">Platform</InputLabel>
                                                        <Select
                                                            labelId="platform-label"
                                                            id="platforms-input"
                                                            value={retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.retrievalPlatform ? retrieval.retrievalItemInstruction.retrievalPlatform : ''}
                                                            label="Platform"
                                                            onChange={(e) => {
                                                                const updatedRetrieval = {
                                                                    ...retrieval,
                                                                    retrievalItemInstruction: {
                                                                        ...retrieval.retrievalItemInstruction,
                                                                        retrievalPlatform: e.target.value
                                                                    }
                                                                };
                                                                dispatch(setRetrieval(updatedRetrieval));
                                                            }}
                                                        >
                                                            <MenuItem value="web">Web</MenuItem>
                                                            <MenuItem value="reddit">Reddit</MenuItem>
                                                            <MenuItem value="twitter">Twitter</MenuItem>
                                                            <MenuItem value="discord">Discord</MenuItem>
                                                            <MenuItem value="telegram">Telegram</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                    </Box>
                                                { retrieval.retrievalItemInstruction !== undefined && retrieval.retrievalItemInstruction.retrievalPlatform !== 'web' &&
                                                    <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                        <TextField
                                                            fullWidth
                                                            id="group-input"
                                                            label={"Platform Groups"}
                                                            variant="outlined"
                                                            value={retrieval.retrievalItemInstruction.retrievalPlatformGroups || ''}
                                                            onChange={(e) => {
                                                                const updatedRetrieval = {
                                                                    ...retrieval,
                                                                    retrievalItemInstruction: {
                                                                        ...retrieval.retrievalItemInstruction,
                                                                        retrievalPlatformGroups: e.target.value
                                                                    }
                                                                };
                                                                dispatch(setRetrieval(updatedRetrieval));
                                                            }}
                                                        />
                                                    </Box>
                                                }
                                                { retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.retrievalPlatform === 'web' &&
                                                    <div>
                                                        <Typography variant="h6" color="text.secondary">
                                                            Use a Load Balancer group for web data retrieval.
                                                        </Typography>
                                                        <FormControl sx={{ mt: 3 }} fullWidth variant="outlined">
                                                            <InputLabel key={`groupNameLabel`} id={`groupName`}>
                                                                Routing Group
                                                            </InputLabel>
                                                            <Select
                                                                labelId={`groupNameLabel`}
                                                                id={`groupName`}
                                                                name="groupName"
                                                                value={retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.webFilters && retrieval.retrievalItemInstruction.webFilters.routingGroup ? retrieval.retrievalItemInstruction.webFilters.routingGroup : ''}
                                                                onChange={(e) => {
                                                                    const updatedRetrieval = {
                                                                        ...retrieval,
                                                                        retrievalItemInstruction: {
                                                                            ...retrieval.retrievalItemInstruction,
                                                                            webFilters: {
                                                                                ...retrieval.retrievalItemInstruction.webFilters,
                                                                                routingGroup: e.target.value, // Correctly update the routingGroup field
                                                                            }
                                                                        }
                                                                    };
                                                                    dispatch(setRetrieval(updatedRetrieval));
                                                                }}
                                                                label="Routing Group"
                                                            >
                                                                {Object.keys(groups).map((name) => <MenuItem key={name} value={name}>{name}</MenuItem>)}
                                                            </Select>
                                                        </FormControl>

                                                        <FormControl sx={{ mt: 3 }} fullWidth variant="outlined">
                                                            <InputLabel key={`groupNameLabel`} id={`groupName`}>
                                                                Load Balancing
                                                            </InputLabel>
                                                            <Select
                                                                labelId={`lbStrategy`}
                                                                id={`lbStrategy`}
                                                                name="lbStrategy"
                                                                value={retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.webFilters && retrieval.retrievalItemInstruction.webFilters.lbStrategy ? retrieval.retrievalItemInstruction.webFilters.lbStrategy : 'round-robin'}
                                                                onChange={(e) => {
                                                                    const updatedRetrieval = {
                                                                        ...retrieval,
                                                                        retrievalItemInstruction: {
                                                                            ...retrieval.retrievalItemInstruction,
                                                                            webFilters: {
                                                                                ...retrieval.retrievalItemInstruction.webFilters,
                                                                                lbStrategy: e.target.value, // Correctly update the routingGroup field
                                                                            }
                                                                        }
                                                                    };
                                                                    dispatch(setRetrieval(updatedRetrieval));
                                                                }}
                                                                label="Load Balancing"
                                                            >
                                                                <MenuItem value="round-robin">Round Robin</MenuItem>
                                                                <MenuItem value="poll-table">Poll Table</MenuItem>
                                                            </Select>
                                                        </FormControl>

                                                    </div>
                                                }
                                                    { retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.retrievalPlatform === 'discord' &&
                                                        <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="category-name-input"
                                                                label="Discord Category Name"
                                                                variant="outlined"
                                                                value={retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.discordFilters ? retrieval.retrievalItemInstruction.discordFilters?.categoryName : ''}
                                                                onChange={(e) => {
                                                                    const updatedRetrieval = {
                                                                        ...retrieval,
                                                                        retrievalItemInstruction: {
                                                                            ...retrieval.retrievalItemInstruction,
                                                                            discordFilters: {
                                                                                ...retrieval.retrievalItemInstruction.discordFilters,
                                                                                categoryName: e.target.value, // Correctly update the routingGroup field
                                                                            }
                                                                        }
                                                                    };
                                                                    dispatch(setRetrieval(updatedRetrieval));
                                                                }}
                                                            />
                                                    </Box>
                                                    }
                                                {/*{ retrieval.retrievalPlatform !== 'web' &&*/}
                                                {/*    <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>*/}
                                                {/*        <TextField*/}
                                                {/*            fullWidth*/}
                                                {/*            id="usernames-input"*/}
                                                {/*            label="Usernames"*/}
                                                {/*            variant="outlined"*/}
                                                {/*            value={retrieval.retrievalUsernames}*/}
                                                {/*            onChange={(e) => dispatch(setRetrievalUsernames(e.target.value))}*/}
                                                {/*        />*/}
                                                {/*    </Box>*/}
                                                {/*}*/}

                                                {  retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.retrievalPlatform === 'web' &&
                                                    <Typography variant="h5" color="text.secondary">
                                                        Describe how the AI should extract data from the website address.
                                                    </Typography>
                                                }
                                                {/*{ retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.retrievalPlatform !== 'web' &&*/}
                                                {/*    <Typography variant="h5" color="text.secondary">*/}
                                                {/*        Describe what you're looking for, and the AI will generate a list of keywords to search for after it runs for the first time.*/}
                                                {/*        You can preview, edit, or give the AI more information to refine the search.*/}
                                                {/*    </Typography>*/}
                                                {/*}*/}
                                                {/*    <Box  sx={{ mb: 2, mt: 2 }}>*/}
                                                {/*        <TextareaAutosize*/}
                                                {/*            minRows={18}*/}
                                                {/*            value={retrieval.retrievalItemInstruction?.retrievalPrompt || ''}*/}
                                                {/*            onChange={(e) => {*/}
                                                {/*                const updatedRetrieval = {*/}
                                                {/*                    ...retrieval,*/}
                                                {/*                    retrievalItemInstruction: {*/}
                                                {/*                        ...retrieval.retrievalItemInstruction,*/}
                                                {/*                        retrievalPrompt: e.target.value*/}
                                                {/*                    }*/}
                                                {/*                };*/}
                                                {/*                dispatch(setRetrieval(updatedRetrieval));*/}
                                                {/*            }}*/}
                                                {/*            style={{ resize: "both", width: "100%" }}*/}
                                                {/*        />*/}
                                                {/*    </Box>*/}
                                                    <Typography variant="h5" color="text.secondary">
                                                        Add search keywords using comma separated values below.
                                                    </Typography>
                                                    <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                        <TextField
                                                            fullWidth
                                                            id="keywords-input"
                                                            label="Keywords"
                                                            variant="outlined"
                                                            value={retrieval.retrievalItemInstruction && retrieval.retrievalItemInstruction.retrievalKeywords ? retrieval.retrievalItemInstruction.retrievalKeywords : ''}
                                                            onChange={(e) => {
                                                                const updatedRetrieval = {
                                                                    ...retrieval,
                                                                    retrievalItemInstruction: {
                                                                        ...retrieval.retrievalItemInstruction,
                                                                        retrievalKeywords: e.target.value
                                                                    }
                                                                };
                                                                dispatch(setRetrieval(updatedRetrieval));
                                                            }}
                                                        />
                                                    </Box>
                                                    {requestRetrievalStatus != '' && (
                                                        <Container sx={{ mb: 2, mt: -2}}>
                                                            <Typography variant="h6" color={requestRetrievalStatusError}>
                                                                {requestRetrievalStatus}
                                                            </Typography>
                                                        </Container>
                                                    )}
                                                    <Box flexGrow={1} sx={{ mb: 0 }}>
                                                        <Button fullWidth variant="contained" onClick={createOrUpdateRetrieval} >Save Retrieval</Button>
                                                    </Box>
                                                    </Stack>
                                                </div>
                                            </CardContent>
                                    }
                                    {
                                        selectedMainTabBuilder === 5 && !addRetrievalView && !loading && !addAnalysisView && !addAggregateView && !addEvalsView && !addAssistantsView && !addTriggersToEvalFnView &&
                                        <CardContent>
                                            <div>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Trigger Action Procedures
                                                </Typography>
                                                <Typography variant="body2" color="text.secondary">
                                                    This allows you to setup event driven rules for AI systems. Email and text actions require a verified email and phone number.
                                                    Contact us to use these outputs for now, as it currently requires semi-manual work from us to integrate you,
                                                    and we require you to use a human-in-the loop to approve all AI orchestrated actions that communicate with the outside world.
                                                    You'll need to generate metrics via the adaptive load balancer to use the metrics based actions. PromQL triggers are only available to enterprise users for now.
                                                    You can add a trigger to any analysis or aggregation stage output. If you add multiple metrics to a trigger, it will only trigger if all metric thresholds are met.
                                                    Post-trigger you can add an operator to the metric score to adjust the score before it is used in the trigger. Eg. you can multiply the score, or perform a general math
                                                    operation on it.
                                                </Typography>
                                                <Stack direction="column" spacing={2} sx={{ mt: 4, mb: 0 }}>
                                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="trigger-name"
                                                                label="Trigger Name"
                                                                variant="outlined"
                                                                value={action.triggerName}
                                                                onChange={(e) => dispatch(setTriggerAction({
                                                                    ...action, // Spread the existing action properties
                                                                    triggerName: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="trigger-group"
                                                                label="Trigger Group"
                                                                variant="outlined"
                                                                value={action.triggerGroup}
                                                                onChange={(e) => dispatch(setTriggerAction({
                                                                    ...action, // Spread the existing action properties
                                                                    triggerGroup: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                    </Stack>
                                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                        <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id="trigger-source-label">Trigger Source</InputLabel>
                                                                <Select
                                                                    labelId="trigger-source--label"
                                                                    id="trigger-source-input"
                                                                    value={'eval'}
                                                                    label="Trigger Source"
                                                                    // onChange={(e) => dispatch(setTriggerAction({
                                                                    //     ...action,
                                                                    // }))}
                                                                >
                                                                    <MenuItem value="eval">Eval</MenuItem>
                                                                    {/*<MenuItem value="metrics">Metrics</MenuItem>*/}
                                                                    {/*<MenuItem value="email">Email</MenuItem>*/}
                                                                    {/*<MenuItem value="text">Text</MenuItem>*/}
                                                                    {/*<MenuItem value="reddit">Reddit</MenuItem>*/}
                                                                    {/*<MenuItem value="twitter">Twitter</MenuItem>*/}
                                                                    {/*<MenuItem value="discord">Discord</MenuItem>*/}
                                                                    {/*<MenuItem value="telegram">Telegram</MenuItem>*/}
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                        <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id="trigger-env-label">Trigger Action</InputLabel>
                                                                <Select
                                                                    labelId="trigger-env-label"
                                                                    id="trigger-env-input"
                                                                    value={action.triggerAction}
                                                                    label="Trigger Action"
                                                                    onChange={(e) => dispatch(setTriggerAction({
                                                                        ...action,
                                                                        triggerAction: e.target.value
                                                                    }))}
                                                                >
                                                                    <MenuItem value="social-media-engagement">Social Media Engagement Approval</MenuItem>
                                                                    {/*<MenuItem value="email">Email</MenuItem>*/}
                                                                    {/*<MenuItem value="text">Text</MenuItem>*/}
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                    </Stack>
                                                    { !loading &&
                                                    <Stack direction="row" >
                                                        <Box flexGrow={1} sx={{ mb: 0,ml: 0, mr:2  }}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id="eval-state-trigger">Eval State</InputLabel>
                                                                <Select
                                                                    id="eval-state-trigger"
                                                                    label="Eval State Trigger"
                                                                    value={action && action.evalTriggerAction.evalTriggerState}
                                                                    onChange={(e) => dispatch(setTriggerAction({
                                                                        ...action, // Spread the existing action properties
                                                                        evalTriggerAction: { ...action.evalTriggerAction, evalTriggerState: e.target.value }
                                                                    }))}
                                                                >
                                                                    <MenuItem value="info">{'info'}</MenuItem>
                                                                    {/*<MenuItem value="optional">{'optional'}</MenuItem>*/}
                                                                    <MenuItem value="filter">{'filter'}</MenuItem>
                                                                    {/*<MenuItem value="warning">{'warning'}</MenuItem>*/}
                                                                    {/*<MenuItem value="critical">{'critical'}</MenuItem>*/}
                                                                    <MenuItem value="error">{'error'}</MenuItem>
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 0,ml: 0, mr:2  }}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id="eval-completion-trigger">Eval Completion</InputLabel>
                                                                <Select
                                                                    id="eval-completion-trigger"
                                                                    label="Eval Completion Trigger"
                                                                    value={action && action.evalTriggerAction.evalResultsTriggerOn || ''}
                                                                    onChange={(e) => dispatch(setTriggerAction({
                                                                        ...action,
                                                                        evalTriggerAction: { ...action.evalTriggerAction, evalResultsTriggerOn: e.target.value }
                                                                    }))}
                                                                >
                                                                    <MenuItem value="all-pass">{'all-pass'}</MenuItem>
                                                                    <MenuItem value="any-pass">{'any-pass'}</MenuItem>
                                                                    <MenuItem value="all-fail">{'all-fail'}</MenuItem>
                                                                    <MenuItem value="any-fail">{'any-fail'}</MenuItem>
                                                                    <MenuItem value="mixed-status">{'mixed-status'}</MenuItem>
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                    </Stack>
                                                    }
                                                    { !loading && action.evalResultsTriggerOn == 'metrics' &&
                                                        <Stack direction="row" >
                                                            {/*<Box flexGrow={1} sx={{ mb: 0,ml: 0, mr:2  }}>*/}
                                                            {/*    <TextField*/}
                                                            {/*        fullWidth*/}
                                                            {/*        id="metric-name"*/}
                                                            {/*        label="Metric Name"*/}
                                                            {/*        variant="outlined"*/}
                                                            {/*        value={actionMetric.metricName}*/}
                                                            {/*        onChange={(e) => dispatch(setActionMetric({*/}
                                                            {/*            ...actionMetric, // Spread the existing action properties*/}
                                                            {/*            metricName: e.target.value // Update the actionName*/}
                                                            {/*        }))}*/}
                                                            {/*    />*/}
                                                            {/*</Box>*/}
                                                            {/*<Box flexGrow={1} sx={{ mb: 0,ml: 0, mr:2  }}>*/}
                                                            {/*    <TextField*/}
                                                            {/*        fullWidth*/}
                                                            {/*        type={"number"}*/}
                                                            {/*        id="metric-score-threshold"*/}
                                                            {/*        label="Metric Score Threshold"*/}
                                                            {/*        variant="outlined"*/}
                                                            {/*        value={actionMetric.metricScoreThreshold}*/}
                                                            {/*        onChange={(e) => dispatch(setActionMetric({*/}
                                                            {/*            ...actionMetric, // Spread the existing action properties*/}
                                                            {/*            metricScoreThreshold: e.target.value // Update the actionName*/}
                                                            {/*        }))}*/}
                                                            {/*    />*/}
                                                            {/*</Box>*/}
                                                            {/*<Box flexGrow={7} >*/}
                                                            {/*    <FormControl fullWidth >*/}
                                                            {/*        <InputLabel id="metric-action-operator">Operator</InputLabel>*/}
                                                            {/*        <Select*/}
                                                            {/*            labelId="metric-action-operator-label"*/}
                                                            {/*            id="metric-action-operator-label"*/}
                                                            {/*            value={actionMetric.metricOperator}*/}
                                                            {/*            label="Metric Action Operator"*/}
                                                            {/*            fullWidth*/}
                                                            {/*            onChange={(e) => dispatch(setActionMetric({*/}
                                                            {/*                ...actionMetric, // Spread the existing action properties*/}
                                                            {/*                metricOperator: e.target.value // Update the actionName*/}
                                                            {/*            }))}*/}
                                                            {/*        >*/}
                                                            {/*            <MenuItem value="add">Add</MenuItem>*/}
                                                            {/*            <MenuItem value="subtract">Subtract</MenuItem>*/}
                                                            {/*            <MenuItem value="multiply">Multiply</MenuItem>*/}
                                                            {/*            <MenuItem value="modulus">Modulus</MenuItem>*/}
                                                            {/*            <MenuItem value="assign">Assign</MenuItem>*/}
                                                            {/*        </Select>*/}
                                                            {/*    </FormControl>*/}
                                                            {/*</Box>*/}
                                                            {/*<Box flexGrow={1} sx={{ mb: 0,ml: 2, mr:0  }}>*/}
                                                            {/*    <TextField*/}
                                                            {/*        fullWidth*/}
                                                            {/*        type={"number"}*/}
                                                            {/*        id="metric-action-number"*/}
                                                            {/*        label="Post-Trigger Operator Value"*/}
                                                            {/*        variant="outlined"*/}
                                                            {/*        value={actionMetric.metricPostActionMultiplier}*/}
                                                            {/*        onChange={(e) => dispatch(setActionMetric({*/}
                                                            {/*            ...actionMetric, // Spread the existing action properties*/}
                                                            {/*            metricPostActionMultiplier: e.target.value // Update the actionName*/}
                                                            {/*        }))}*/}
                                                            {/*    />*/}
                                                            {/*</Box>*/}
                                                            {/*<Box flexGrow={2} sx={{ mt:1, mb: 0,ml: 2, mr:0  }}>*/}
                                                            {/*    <Button fullWidth variant={"contained"} onClick={addActionMetricRow}>Add</Button>*/}
                                                            {/*</Box>*/}
                                                        </Stack>
                                                    }
                                                    {/*{*/}
                                                    {/*    !loading && action && action.actionMetrics && action.actionMetrics.map((metric: ActionMetric, index: number) => (*/}
                                                    {/*        <Stack key={index} direction="row" alignItems="center" spacing={2} sx={{ mt: 4, mb: 4 }}>*/}
                                                    {/*            /!* Metric Name *!/*/}
                                                    {/*            <Box flexGrow={1} sx={{ ml: 4, mr: 4 }}>*/}
                                                    {/*                <TextField*/}
                                                    {/*                    fullWidth*/}
                                                    {/*                    id={`metric-name-${index}`}*/}
                                                    {/*                    label="Metric Name"*/}
                                                    {/*                    variant="outlined"*/}
                                                    {/*                    value={metric.metricName}*/}
                                                    {/*                    inputProps={{ readOnly: true }}*/}
                                                    {/*                />*/}
                                                    {/*            </Box>*/}
                                                    {/*            /!* Metric Score Threshold *!/*/}
                                                    {/*            <Box flexGrow={1} sx={{ mr: 4 }}>*/}
                                                    {/*                <TextField*/}
                                                    {/*                    fullWidth*/}
                                                    {/*                    type="number"*/}
                                                    {/*                    id={`metric-score-threshold-${index}`}*/}
                                                    {/*                    label="Metric Score Threshold"*/}
                                                    {/*                    variant="outlined"*/}
                                                    {/*                    value={metric.metricScoreThreshold}*/}
                                                    {/*                    inputProps={{ readOnly: true }}*/}
                                                    {/*                />*/}
                                                    {/*            </Box>*/}
                                                    {/*            /!* Metric Action Multiplier *!/*/}
                                                    {/*            <Box flexGrow={1} sx={{ mr: 4 }}>*/}
                                                    {/*                <TextField*/}
                                                    {/*                    fullWidth*/}
                                                    {/*                    type="number"*/}
                                                    {/*                    id={`metric-action-multiplier-${index}`}*/}
                                                    {/*                    label="Metric Action Multiplier"*/}
                                                    {/*                    variant="outlined"*/}
                                                    {/*                    value={metric.metricPostActionMultiplier}*/}
                                                    {/*                    inputProps={{ readOnly: true }}*/}
                                                    {/*                />*/}
                                                    {/*            </Box>*/}
                                                    {/*            /!* Remove Button *!/*/}
                                                    {/*            <Box sx={{ mr: 4 }}>*/}
                                                    {/*                <Button onClick={() => removeActionMetricRow(index)}>Remove</Button>*/}
                                                    {/*            </Box>*/}
                                                    {/*        </Stack>*/}
                                                    {/*    ))*/}
                                                    {/*}*/}
                                                    {requestActionStatus != '' && (
                                                        <Container sx={{ mb: 2, mt: -2}}>
                                                            <Typography variant="h6" color={requestActionStatusError}>
                                                                {requestActionStatus}
                                                            </Typography>
                                                        </Container>
                                                    )}
                                                    <Box flexGrow={1} sx={{ mb: 0 }}>
                                                        <Button fullWidth variant="contained" onClick={createOrUpdateAction}>Save Action</Button>
                                                    </Box>
                                                </Stack>
                                            </div>
                                        </CardContent>
                                    }
                                    {
                                        (selectedMainTabBuilder === 4 || addTriggersToEvalFnView) && !addRetrievalView && !loading && !addAnalysisView && !addAggregateView && !addEvalsView && !addAssistantsView &&
                                        <CardContent>
                                            <div>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Eval Functions
                                                </Typography>
                                                <Typography variant="body2" color="text.secondary">
                                                    This allows you to setup scoring rules and triggers for AI system outputs that set metrics for the AI to use in its decision making process. For metric
                                                    array types, the comparison value returns the true/false if every array element item passes the comparison eval test.
                                                </Typography>
                                                <Stack direction="column" spacing={2} sx={{ mt: 2, mb: 0 }}>
                                                    <Stack direction="row" spacing={2} sx={{ mt: 0, mb: 4 }}>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="eval-name"
                                                                label="Eval Name"
                                                                variant="outlined"
                                                                value={evalFn &&  evalFn.evalName}
                                                                onChange={(e) => dispatch(setEval({
                                                                    ...evalFn, // Spread the existing action properties
                                                                    evalName: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="eval-group"
                                                                label="Eval Group"
                                                                variant="outlined"
                                                                value={evalFn && evalFn.evalGroupName}
                                                                onChange={(e) => dispatch(setEval({
                                                                    ...evalFn, // Spread the existing action properties
                                                                    evalGroupName: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                        <Box flexGrow={7} >
                                                            <FormControl fullWidth >
                                                                <InputLabel id="eval-type-operator">Type</InputLabel>
                                                                <Select
                                                                    labelId="eval-type-label"
                                                                    id="eval-type-label"
                                                                    value={evalFn && evalFn.evalType}
                                                                    label="Eval Type"
                                                                    fullWidth
                                                                    onChange={(e) => dispatch(setEval({
                                                                        ...evalFn, // Spread the existing action properties
                                                                        evalType: e.target.value // Update the actionName
                                                                    }))}
                                                                >
                                                                    <MenuItem value="model">Model</MenuItem>
                                                                    {/*<MenuItem value="api">API</MenuItem>*/}
                                                                    {/*<MenuItem value="adaptive">Adaptive</MenuItem>*/}
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                    </Stack>
                                                    { evalFn && evalFn.evalType == 'model' &&
                                                    <Stack direction={"row"} >
                                                        <Box flexGrow={3} sx={{ mr: 2}}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id="eval-model-label">Eval Model</InputLabel>
                                                                <Select
                                                                    labelId="eval-model-label"
                                                                    id="eval-model-select"
                                                                    value={evalFn.evalModel}
                                                                    label="Eval Model"
                                                                    onChange={(e) => dispatch(setEval({
                                                                        ...evalFn, // Spread the existing action properties
                                                                        evalModel: e.target.value // Update the actionName
                                                                    }))}
                                                                >
                                                                    <MenuItem value="gpt-3.5-turbo-1106">gpt-3.5-turbo-1106</MenuItem>
                                                                    <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                        <Box flexGrow={3} sx={{ }}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id="eval-format-label">Eval Format</InputLabel>
                                                                <Select
                                                                    labelId="eval-format-label"
                                                                    id="eval-format-select"
                                                                    value={evalFn.evalFormat}
                                                                    label="Eval Model"
                                                                    onChange={(e) => dispatch(setEval({
                                                                        ...evalFn, // Spread the existing action properties
                                                                        evalFormat: e.target.value // Update the actionName
                                                                    }))}
                                                                >
                                                                    {/*<MenuItem value="code">code</MenuItem>*/}
                                                                    <MenuItem value="json">json</MenuItem>
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                    </Stack>
                                                    }
                                                    <Stack direction="column" >
                                                        <Stack direction="row" >
                                                            <Box flexGrow={7} sx={{ mb: 2,ml: 0, mr:2  }}>
                                                                <TextField
                                                                    fullWidth
                                                                    id="metric-name"
                                                                    label="Metric Name"
                                                                    variant="outlined"
                                                                    value={evalMetric.evalMetricName}
                                                                    onChange={(e) => dispatch(setEvalMetric({
                                                                        ...evalMetric, // Spread the existing action properties
                                                                        evalMetricName: e.target.value // Update the actionName
                                                                    }))}
                                                                />
                                                            </Box>
                                                            <Box flexGrow={7} sx={{ mb: 2,ml: 0, mr:2  }}>
                                                                <FormControl fullWidth >
                                                                    <InputLabel id="metric-state-operator">State</InputLabel>
                                                                    <Select
                                                                        labelId="metric-state-operator-label"
                                                                        id="metric-state-operator-label"
                                                                        value={evalMetric.evalState}
                                                                        label="Metric State Operator"
                                                                        fullWidth
                                                                        onChange={(e) => dispatch(setEvalMetric({
                                                                            ...evalMetric, // Spread the existing action properties
                                                                            evalState: e.target.value // Update the actionName
                                                                        }))}
                                                                    >
                                                                        <MenuItem value="info">{'info'}</MenuItem>
                                                                        <MenuItem value="filter">{'filter'}</MenuItem>
                                                                        {/*<MenuItem value="optional">{'optional'}</MenuItem>*/}
                                                                        {/*<MenuItem value="warning">{'warning'}</MenuItem>*/}
                                                                        {/*<MenuItem value="critical">{'critical'}</MenuItem>*/}
                                                                        <MenuItem value="error">{'error'}</MenuItem>
                                                                    </Select>
                                                                </FormControl>
                                                            </Box>
                                                            <Box flexGrow={7} sx={{ mb: 2,ml: 0, mr:2  }}>
                                                                <FormControl fullWidth >
                                                                    <InputLabel id="eval-metric-type">Metric Data Type</InputLabel>
                                                                    <Select
                                                                        labelId="eval-metric-type-label"
                                                                        id="eval-metric-type-label"
                                                                        value={evalMetric.evalMetricDataType}
                                                                        label="Eval Metric Type"
                                                                        fullWidth
                                                                        onChange={(e) => dispatch(setEvalMetric({
                                                                            ...evalMetric, // Spread the existing action properties
                                                                            evalMetricDataType: e.target.value // Update the actionName
                                                                        }))}
                                                                    >
                                                                        <MenuItem value="number">{'number'}</MenuItem>
                                                                        <MenuItem value="string">{'string'}</MenuItem>
                                                                        <MenuItem value="boolean">{'boolean'}</MenuItem>
                                                                        <MenuItem value="array[boolean]">{'array[boolean]'}</MenuItem>
                                                                        <MenuItem value="array[number]">{'array[number]'}</MenuItem>
                                                                        <MenuItem value="array[string]">{'array[string]'}</MenuItem>
                                                                    </Select>
                                                                </FormControl>
                                                            </Box>
                                                            {( evalMetric.evalMetricDataType === 'boolean' || evalMetric.evalMetricDataType === 'array[boolean]') &&
                                                                <Box flexGrow={7} sx={{ mb: 2,ml: 0, mr:2  }}>
                                                                    <FormControl fullWidth >
                                                                        <InputLabel id="metric-action-operator">Operator</InputLabel>
                                                                        <Select
                                                                            labelId="metric-action-operator-label"
                                                                            id="metric-action-operator-label"
                                                                            value={evalMetric.evalComparisonBoolean}
                                                                            label="Metric Action Operator"
                                                                            fullWidth
                                                                            onChange={(e) => dispatch(setEvalMetric({
                                                                                ...evalMetric, // Spread the existing action properties
                                                                                evalComparisonBoolean: e.target.value === 'true'// Update the actionName
                                                                            }))}
                                                                        >
                                                                            <MenuItem value="true">{'true'}</MenuItem>
                                                                            <MenuItem value="false">{'false'}</MenuItem>
                                                                        </Select>
                                                                    </FormControl>
                                                                </Box>
                                                            }
                                                            { (evalMetric.evalMetricDataType === 'number' || evalMetric.evalMetricDataType === 'array[number]') &&
                                                            <Box flexGrow={7} >
                                                                <FormControl fullWidth >
                                                                    <InputLabel id="metric-action-operator">Operator</InputLabel>
                                                                    <Select
                                                                        labelId="metric-action-operator-label"
                                                                        id="metric-action-operator-label"
                                                                        value={evalMetric.evalOperator}
                                                                        label="Metric Action Operator"
                                                                        fullWidth
                                                                        onChange={(e) => dispatch(setEvalMetric({
                                                                            ...evalMetric, // Spread the existing action properties
                                                                            evalOperator: e.target.value // Update the actionName
                                                                        }))}
                                                                    >
                                                                        <MenuItem value="gt">{'>'}</MenuItem>
                                                                        <MenuItem value="gte">{'>='}</MenuItem>
                                                                        <MenuItem value="lt">{'<'}</MenuItem>
                                                                        <MenuItem value="lte">{'<='}</MenuItem>
                                                                        <MenuItem value="eq">{'=='}</MenuItem>
                                                                    </Select>
                                                                </FormControl>
                                                            </Box>
                                                            }
                                                            { evalMetric && (evalMetric.evalMetricDataType === 'string' || evalMetric.evalMetricDataType === 'array[string]') &&
                                                                <Box flexGrow={7} sx={{ mb: 2,ml: 0, mr:2  }}>
                                                                    <FormControl fullWidth >
                                                                        <InputLabel id="metric-action-operator">Operator</InputLabel>
                                                                        <Select
                                                                            labelId="metric-action-operator-label"
                                                                            id="metric-action-operator-label"
                                                                            value={evalMetric.evalOperator}
                                                                            label="Metric Action Operator"
                                                                            fullWidth
                                                                            onChange={(e) => dispatch(setEvalMetric({
                                                                                ...evalMetric, // Spread the existing action properties
                                                                                evalOperator: e.target.value // Update the actionName
                                                                            }))}
                                                                        >
                                                                            {evalMetric.evalMetricDataType === 'array[string]' &&
                                                                                <MenuItem value="all-unique-words">{'all-unique-words'}</MenuItem>
                                                                            }
                                                                            <MenuItem value="contains">{'contains'}</MenuItem>
                                                                            <MenuItem value="has-prefix">{'has-prefix'}</MenuItem>
                                                                            <MenuItem value="has-suffix">{'has-suffix'}</MenuItem>
                                                                            <MenuItem value="does-not-start-with-any">{'does-not-start-with'}</MenuItem>
                                                                            <MenuItem value="does-not-include">{'does-not-include'}</MenuItem>
                                                                            <MenuItem value="equals">{'equals'}</MenuItem>
                                                                            <MenuItem value="length-less-than">{'length-less-than'}</MenuItem>
                                                                            <MenuItem value="length-less-than-eq">{'length-less-than-eq'}</MenuItem>
                                                                            <MenuItem value="length-greater-than">{'length-greater-than'}</MenuItem>
                                                                            <MenuItem value="length-greater-than-eq">{'length-greater-than-eq'}</MenuItem>
                                                                            <MenuItem value="length-eq">{'length-eq'}</MenuItem>
                                                                        </Select>
                                                                    </FormControl>
                                                                </Box>
                                                            }
                                                            { (evalMetric.evalMetricDataType === 'number' || evalMetric.evalMetricDataType === 'array[number]'
                                                                || (evalMetric.evalOperator == 'unique-words') || (evalMetric.evalOperator == 'length-eq')
                                                                ) &&
                                                            <Box flexGrow={1} sx={{ mb: 0,ml: 2, mr:2  }}>
                                                                <TextField
                                                                    fullWidth
                                                                    id="eval-comparison-value"
                                                                    label="Comparison Value"
                                                                    variant="outlined"
                                                                    type={"number"}
                                                                    value={evalMetric && evalMetric.evalComparisonNumber || 0}
                                                                    onChange={(e) => dispatch(setEvalMetric({
                                                                        ...evalMetric, // Spread the existing action properties
                                                                        evalComparisonNumber: Number(e.target.value) // Update the actionName
                                                                    }))}
                                                                />
                                                            </Box>
                                                            }
                                                            { (evalMetric.evalMetricDataType === 'string'|| evalMetric.evalMetricDataType === 'array[string]') &&
                                                                (evalMetric.evalOperator !== 'all-unique-words') &&  (evalMetric.evalOperator != 'length-eq') &&
                                                                <Box flexGrow={1} sx={{ mb: 0,ml: 0, mr:2  }}>
                                                                    <TextField
                                                                        fullWidth
                                                                        id="eval-comparison-string"
                                                                        label="Comparison String"
                                                                        variant="outlined"
                                                                        value={evalMetric && evalMetric.evalComparisonString || ''}
                                                                        onChange={(e) => dispatch(setEvalMetric({
                                                                            ...evalMetric, // Spread the existing action properties
                                                                            evalComparisonString: e.target.value // Update the actionName
                                                                        }))}
                                                                    />
                                                                </Box>
                                                            }
                                                            <Box flexGrow={3} sx={{ mb: 0,ml: 0, mr:2  }}>
                                                                <FormControl fullWidth>
                                                                    <InputLabel id="eval-result">Result</InputLabel>
                                                                    <Select
                                                                        labelId="eval-result-label"
                                                                        id="eval-result-label"
                                                                        value={evalMetric.evalMetricResult}
                                                                        label="Result"
                                                                        fullWidth
                                                                        onChange={(e) => dispatch(setEvalMetric({
                                                                            ...evalMetric, // Spread the existing action properties
                                                                            evalMetricResult: e.target.value // Update the metricOperator
                                                                        }))}
                                                                    >
                                                                        <MenuItem value="pass">Pass</MenuItem>
                                                                        <MenuItem value="fail">Fail</MenuItem>
                                                                    </Select>
                                                                </FormControl>
                                                            </Box>
                                                            <Box flexGrow={2} sx={{ mt:1, mb: 0,ml: 0, mr:0  }}>
                                                                <Button fullWidth variant={"contained"} onClick={addEvalMetricRow}>Add</Button>
                                                            </Box>
                                                            <Box flexGrow={2} sx={{ mt:1, mb: 0,ml: 2, mr:0  }}>
                                                                <Button fullWidth variant={"contained"} onClick={clearEvalMetricRow}>Clear</Button>
                                                            </Box>
                                                            </Stack>
                                                        { evalFn.evalType == 'model' &&
                                                            <Box  sx={{ mb: 2, mt: 2 }}>
                                                                <TextareaAutosize
                                                                    minRows={18}
                                                                    value={evalMetric.evalModelPrompt}
                                                                    onChange={(e) => dispatch(setEvalMetric({
                                                                        ...evalMetric, // Spread the existing action properties
                                                                        evalModelPrompt: e.target.value // Update the metricOperator
                                                                    }))}
                                                                    style={{ resize: "both", width: "100%" }}
                                                                />
                                                            </Box>
                                                        }
                                                        </Stack>
                                                    {
                                                        !loading && action && evalFn.evalMetrics && evalFn.evalMetrics.map((metric: EvalMetric, index: number) => (
                                                            <Stack key={index} direction="column" sx={{ mt: 4, mb: 4, mr: 0 }}>
                                                                <Stack key={index} direction="row" alignItems="center" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                                    {/* Metric Name */}
                                                                    <Box flexGrow={1} sx={{ ml: 4, mr: 4 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id={`metric-name-${index}`}
                                                                            label="Metric Name"
                                                                            variant="outlined"
                                                                            value={metric.evalMetricName}
                                                                            inputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ ml: 4, mr: 4 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id={`metric-state-${index}`}
                                                                            label="Metric State"
                                                                            variant="outlined"
                                                                            value={metric.evalState}
                                                                            inputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ ml: 4, mr: 4 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id={`metric-state-${index}`}
                                                                            label="Data Type"
                                                                            variant="outlined"
                                                                            value={metric.evalMetricDataType}
                                                                            inputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ ml: 4, mr: 4 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id={`metric-op-${index}`}
                                                                            label="Operator"
                                                                            variant="outlined"
                                                                            value={metric.evalOperator}
                                                                            inputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ ml: 4, mr: 4 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id={`metric-comp-${index}`}
                                                                            label="Comparison Value"
                                                                            variant="outlined"
                                                                            value={GetValue(metric)}
                                                                            inputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ ml: 4, mr: 4 }}>
                                                                        <TextField
                                                                            fullWidth
                                                                            id={`metric-result-${index}`}
                                                                            label="Result"
                                                                            variant="outlined"
                                                                            value={metric.evalMetricResult}
                                                                            inputProps={{ readOnly: true }}
                                                                        />
                                                                    </Box>
                                                                    <Stack direction="row" spacing={2} sx={{ ml: 4, mr: 4 }}>
                                                                        <Box sx={{ mr: 4 }}>
                                                                            <Button variant={"contained"} onClick={() => editEvalMetricRow(index)}>Edit</Button>
                                                                        </Box>
                                                                        <Box sx={{ mr: 4 }}>
                                                                            <Button variant={"contained"} onClick={() => removeEvalMetricRow(index)}>Remove</Button>
                                                                        </Box>
                                                                    </Stack>
                                                                </Stack>
                                                                <Box flexGrow={1} sx={{ ml: 0, mr: 12 }}>
                                                                    <TextField
                                                                        fullWidth
                                                                        id={`eval-model-scoring-instruction-${index}`}
                                                                        label="Model Scoring Instructions"
                                                                        variant="outlined"
                                                                        value={metric.evalModelPrompt}
                                                                        inputProps={{ readOnly: true }}
                                                                    />
                                                                </Box>
                                                            </Stack>
                                                        ))
                                                    }
                                                    <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                        <Divider/>
                                                    </Box>
                                                    <Typography variant="h6" color="text.secondary">
                                                       Eval Triggers
                                                    </Typography>
                                                    {/*<Collapse in={openRetrievals} timeout="auto" unmountOnExit>*/}
                                                        <Box flexGrow={2} sx={{mt: 2}}>
                                                            {evalFn && evalFn.triggerFunctions && evalFn.triggerFunctions.map((trigger: TriggerAction, subIndex: React.Key | null | undefined) => (
                                                                <Stack direction={"row"} key={subIndex} sx={{ mb: 2 }}>
                                                                    <Box flexGrow={2} sx={{ mt: 0, ml: 0 }}>
                                                                        <TextField
                                                                            key={subIndex}
                                                                            label={`Trigger Name`}
                                                                            value={trigger && trigger.triggerName || ''}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                            variant="outlined"
                                                                            fullWidth
                                                                            margin="normal"
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={2} sx={{ mt: 0, ml: 2 }}>
                                                                        <TextField
                                                                            key={subIndex}
                                                                            label={`Trigger Group`}
                                                                            value={trigger && trigger.triggerGroup || ''}
                                                                            InputProps={{
                                                                                readOnly: true,
                                                                            }}
                                                                            variant="outlined"
                                                                            fullWidth
                                                                            margin="normal"
                                                                        />
                                                                    </Box>
                                                                    <Box flexGrow={1} sx={{ mb: 0, ml: 2, mt: 3 }}>
                                                                        <Button fullWidth variant="contained" onClick={(event)=>handleRemoveTriggerFromEvalFn(event, trigger)}>Remove</Button>
                                                                    </Box>
                                                                </Stack>
                                                            ))}
                                                        </Box>
                                                    <Box flexGrow={1} sx={{ mb: 0 }}>
                                                        <Button variant="contained" onClick={addTriggersToEvalFn} >{addTriggersToEvalFnView ? 'Done Adding' : 'Add Triggers'}</Button>
                                                    </Box>
                                                    <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                        <Divider/>
                                                    </Box>
                                                    {requestEvalCreateOrUpdateStatus != '' && (
                                                        <Container sx={{ mb: 2, mt: -2}}>
                                                            <Typography variant="h6" color={requestEvalCreateOrUpdateStatusError}>
                                                                {requestEvalCreateOrUpdateStatus}
                                                            </Typography>
                                                        </Container>
                                                    )}

                                                    <Box flexGrow={1} sx={{ mb: 0 }}>
                                                        <Button fullWidth variant="contained" onClick={createOrUpdateEval} >Save Eval</Button>
                                                    </Box>
                                                </Stack>
                                            </div>
                                        </CardContent>
                                    }
                                    {
                                        (selectedMainTabBuilder === 6) && !addRetrievalView && !loading && !addAnalysisView && !addAggregateView && !addEvalsView && !addTriggerActionsView &&
                                        <Assistants assistant={assistant}
                                                    createOrUpdateAssistant={createOrUpdateAssistant}
                                                    requestStatusAssistant={requestStatusAssistant} requestStatusAssistantError={requestStatusAssistantError}/>
                                    }
                                    { !addAnalysisView && !addAggregateView && selectedMainTabBuilder === 1 && !loading && !addRetrievalView && !addEvalsView && !addAssistantsView && !addTriggerActionsView &&
                                        <div>
                                            <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                <Box flexGrow={2} sx={{ mb: 2, mt: 4, ml:2 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="response-format-label">Response Format</InputLabel>
                                                        <Select
                                                            labelId="response-format-label"
                                                            id="response-format-label"
                                                            value={editAnalysisTask.responseFormat}
                                                            label="Response Format"
                                                            onChange={(e) => handleEditAnalysisTaskResponseFormat(e)}

                                                        >
                                                            <MenuItem value="text">text</MenuItem>
                                                            <MenuItem value="social-media-content-writer">social-media-content-writer</MenuItem>
                                                            <MenuItem value="social-media-engagement">social-media-engagement</MenuItem>
                                                            <MenuItem value="json">json</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>
                                                <Box sx={{ width: '50%', mb: 4, mt: 4 }}>
                                                    <TextField
                                                        type="number"
                                                        label={`Max Tokens Analysis Model`}
                                                        variant="outlined"
                                                        value={editAnalysisTask && editAnalysisTask.maxTokens || 0}
                                                        inputProps={{ min: 0 }}
                                                        onChange={(e) => dispatch(setEditAnalysisTask({
                                                            ...editAnalysisTask, // Spread the existing action properties
                                                            maxTokens: Number(e.target.value) // Update the actionName
                                                        }))}
                                                        fullWidth
                                                    />
                                                </Box>
                                            </Stack>
                                            {requestAnalysisStatus != '' && (
                                                <Container sx={{ mb: 2, mt: -2}}>
                                                    <Typography variant="h6" color={requestAnalysisStatusError}>
                                                        {requestAnalysisStatus}
                                                    </Typography>
                                                </Container>
                                            )}
                                            <Box flexGrow={1} sx={{ mb: 2 }}>
                                                <Button fullWidth variant="contained" onClick={() =>  createOrUpdateTask('analysis')} >Save Analysis</Button>
                                            </Box>
                                        </div>
                                    }
                                    {  !addAnalysisView && !addAggregateView && selectedMainTabBuilder === 2 && !loading && !addRetrievalView && !addEvalsView && !addTriggerActionsView && !addAssistantsView &&
                                        <div>
                                            <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                            <Box flexGrow={2} sx={{ mb: 2, mt: 4, ml:2 }}>
                                                <FormControl fullWidth>
                                                    <InputLabel id="response-format-label">Response Format</InputLabel>
                                                    <Select
                                                        labelId="response-format-label"
                                                        id="response-format-label"
                                                        value={editAggregateTask.responseFormat}
                                                        label="Response Format"
                                                        onChange={(e) => handleEditAggTaskResponseFormat(e)}
                                                    >
                                                        <MenuItem value="text">text</MenuItem>
                                                        <MenuItem value="social-media-content-writer">social-media-content-writer</MenuItem>
                                                        <MenuItem value="social-media-engagement">social-media-engagement</MenuItem>
                                                        <MenuItem value="json">json</MenuItem>
                                                    </Select>
                                                </FormControl>
                                            </Box>
                                            <Box sx={{ width: '50%', mb: 4, mt: 4 }}>
                                                <TextField
                                                    type="number"
                                                    label={`Max Aggregation Token Usage`}
                                                    variant="outlined"
                                                    value={editAggregateTask && editAggregateTask.maxTokens || 0}
                                                    onChange={(e) => dispatch(setEditAggregateTask({
                                                        ...editAggregateTask, // Spread the existing action properties
                                                        maxTokens: Number(e.target.value)
                                                    }))}
                                                    inputProps={{ min: 0 }}
                                                    fullWidth
                                                />
                                            </Box>
                                            </Stack>
                                            {requestAggStatus != '' && (
                                                <Container sx={{ mb: 2, mt: -2}}>
                                                    <Typography variant="h6" color={requestAggStatusError}>
                                                        {requestAggStatus}
                                                    </Typography>
                                                </Container>
                                            )}
                                            <Box flexGrow={1} sx={{ mb: 2 }}>
                                                <Button fullWidth variant="contained" onClick={() => createOrUpdateTask('aggregation')} >Save Aggregation</Button>
                                            </Box>
                                        </div>
                                    }
                                    { !addAnalysisView && !addAggregateView &&
                                        (selectedMainTabBuilder === 7 ||
                                            ((selectedMainTabBuilder === 1 || selectedMainTabBuilder === 2) && addSchemasView) &&
                                            !loading && !addRetrievalView && !addEvalsView && !addTriggerActionsView && !addAssistantsView) &&
                                        <div>
                                            <Schemas schemaField={schemaField}
                                                     schema={schema}
                                                     removeSchemaField={removeSchemaField}
                                                     addJsonSchemaFieldRow={addJsonSchemaFieldRow}
                                                     createOrUpdateSchema={createOrUpdateSchema}/>
                                            {requestStatusSchema !== '' && (
                                                <Container sx={{ mt: 2 }}>
                                                    <Typography variant="h6" color={requestStatusSchemaError}>
                                                        {requestStatusSchema}
                                                    </Typography>
                                                </Container>
                                            )}
                                        </div>
                                    }
                                </CardContent>
                            </Card>
                        </Stack>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                            <Tabs value={selectedMainTabBuilder} onChange={handleMainTabChange} aria-label="basic tabs">
                                <Tab className="onboarding-card-highlight-all-workflows" label="Workflows"  />
                                <Tab className="onboarding-card-highlight-all-analysis" label="Analysis" />
                                <Tab className="onboarding-card-highlight-all-aggregation" label="Aggregations" />
                                <Tab className="onboarding-card-highlight-all-retrieval" label="Retrievals" />
                                <Tab className="onboarding-card-highlight-all-evals" label="Evals" />
                                <Tab className="onboarding-card-highlight-all-actions" label="Actions" />
                                <Tab className="onboarding-card-highlight-all-assistants" label="Assistants" />
                                <Tab className="onboarding-card-highlight-all-schemas" label="Schemas" />
                            </Tabs>
                        </Box>
                    </Container>
                    { selectedMainTabBuilder === 0 && !addAnalysisView && !addAggregateView && !addRetrievalView && selectedWorkflows && selectedWorkflows.length > 0 && !addEvalsView && !addAssistantsView &&
                        !addTriggerActionsView &&
                        <div>
                                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                    <Box sx={{ mb: 2 }}>
                                        <span>({selectedWorkflows.length} Selected Workflows)</span>
                                        <Button variant="outlined" color="secondary" onClick={(event) => handleDeleteWorkflows(event)} style={{marginLeft: '10px'}}>
                                            Delete { selectedWorkflows.length === 1 ? 'Workflow' : 'Workflows' }
                                        </Button>
                                    </Box>
                                </Container>

                        </div>
                    }
                    { selectedMainTabBuilder === 0 &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <WorkflowTable selected={selected} setSelected={setSelected}/>
                        </Container>
                    }
                    { (selectedMainTabBuilder === 1 || selectedMainTabBuilder === 2) && (addAggregateView || addAnalysisView) &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <Box sx={{ mb: 2 }}>
                                <span>({Object.values(selected).filter(value => value).length} Selected Tasks)</span>
                                <Button variant="outlined" color="secondary" onClick={handleAddTasksToWorkflow} style={{marginLeft: '10px'}}>
                                    Add {addAnalysisView ? 'Analysis' : 'Aggregation'} Stages
                                </Button>
                            </Box>
                        </Container>
                    }
                    { (selectedMainTabBuilder === 4) && (addEvalsView) &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <Box sx={{ mb: 2 }}>
                                <span>(Selected {Object.values(selected).filter(value => value).length} {Object.values(selected).filter(value => value).length > 1 ? 'Evals': 'Eval'})</span>
                                <Button variant="outlined" color="secondary" onClick={handleAddTasksToWorkflow} style={{marginLeft: '10px'}}>
                                    Add Eval Stages
                                </Button>
                            </Box>
                        </Container>
                    }
                    { (selectedMainTabBuilder === 1 || selectedMainTabBuilder === 2) && !addSchemasView &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <TasksTable tasks={tasks} selected={selected} handleClick={handleClick} handleSelectAllClick={handleSelectAllClick} />
                            </Container>
                        </div>
                    }
                    { (selectedMainTabBuilder === 3) && addRetrievalView &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4,  ml: 2 }}>
                            <Box sx={{ mb: 2 }}>
                                <span>({Object.values(selected).filter(value => value).length} Selected Retrievals)</span>
                                <Button variant="outlined" color="secondary" onClick={handleAddTasksToWorkflow} style={{marginLeft: '10px'}}>
                                    Add Retrievals
                                </Button>
                            </Box>
                        </Container>
                    }
                    { (selectedMainTabBuilder === 3) &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <RetrievalsTable retrievals={retrievals} selected={selected} handleSelectAllClick={handleSelectAllClick} handleClick={handleClick} />
                            </Container>
                        </div>
                    }
                    { addTriggersToEvalFnView &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <Box sx={{ mb: 2 }}>
                                    <span>({Object.values(selected).filter(value => value).length} Selected Triggers)</span>
                                    <Button variant="outlined" color="secondary" onClick={(event) => handleAddTriggersToEvalFn(event)} style={{marginLeft: '10px'}}>
                                        Add {Object.values(selected).filter(value => value).length === 1 ? 'Trigger' : 'Triggers'}
                                    </Button>
                                </Box>
                            </Container>

                        </div>
                    }
                    { addSchemasView &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <Box sx={{ mb: 2 }}>
                                    <span>({Object.values(selected).filter(value => value).length} Selected Schemas)</span>
                                    <Button variant="outlined" color="secondary" onClick={(event) => addSchemaToTask(event)} style={{marginLeft: '10px'}}>
                                        Add {Object.values(selected).filter(value => value).length === 1 ? 'Schema' : 'Schemas'}
                                    </Button>
                                </Box>
                            </Container>

                        </div>
                    }
                    { (selectedMainTabBuilder === 4) && !addTriggersToEvalFnView &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <EvalsTable evalFns={evalFns} selected={selected} handleSelectAllClick={handleSelectAllClick} handleClick={handleClick} />
                            </Container>
                        </div>
                    }
                    { (selectedMainTabBuilder === 5 || addTriggersToEvalFnView) &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <ActionsTable actions={actions} selected={selected} handleSelectAllClick={handleSelectAllClick} handleClick={handleClick} />
                            </Container>
                        </div>
                    }
                    { (selectedMainTabBuilder === 6) &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <AssistantsTable assistants={assistants} selected={selected} handleSelectAllClick={handleSelectAllClick} handleClick={handleClick} />
                            </Container>
                        </div>
                    }

                    { (selectedMainTabBuilder === 7 || addSchemasView) &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <SchemasTable schemas={schemas}
                                              requestStatusSchema={requestStatusSchema}
                                              requestStatusSchemaError={requestStatusSchemaError}
                                              selected={selected} handleSelectAllClick={handleSelectAllClick} handleClick={handleClick} />
                            </Container>
                        </div>
                    }
                    <ZeusCopyright sx={{ pt: 4 }} />
                </Box>
            </Box>
        </ThemeProvider>
    );
}
type ValuePiece = Date | string | null;

type Value = ValuePiece | [ValuePiece, ValuePiece];
export default function AiWorkflowsEngineBuilderDashboard() {
    return <WorkflowEngineBuilder />;
}

function GetValue(evm: EvalMetric) {
    if (evm.evalMetricDataType === 'string') {
        return evm.evalComparisonString
    }
    if (evm.evalMetricDataType === 'number') {
        return evm.evalComparisonNumber
    }
    if (evm.evalMetricDataType === 'boolean') {
        return evm.evalComparisonBoolean
    }
    if (evm.evalMetricDataType === 'array[string]') {
        return evm.evalComparisonString
    }
    if (evm.evalMetricDataType === 'array[number]') {
        return evm.evalComparisonNumber
    }
    if (evm.evalMetricDataType === 'array[boolean]') {
        return evm.evalComparisonBoolean
    }
    return ''
}
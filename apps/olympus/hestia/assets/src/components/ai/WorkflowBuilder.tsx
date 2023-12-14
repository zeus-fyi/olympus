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
    setAction,
    setActionMetric,
    setActionPlatformAccount,
    setAddAggregateTasks,
    setAddAggregationView,
    setAddAnalysisTasks,
    setAddAnalysisView,
    setAddRetrievalTasks,
    setAddRetrievalView,
    setAggregationWorkflowInstructions,
    setAnalysisRetrievalsMap,
    setAnalysisWorkflowInstructions,
    setDiscordOptionsCategoryName,
    setRetrievalGroup,
    setRetrievalKeywords,
    setRetrievalName,
    setRetrievalPlatform,
    setRetrievalPlatformGroups,
    setRetrievalPrompt,
    setRetrievals,
    setRetrievalUsernames,
    setSelectedMainTabBuilder,
    setSelectedWorkflows,
    setTaskMap,
    setWebRoutingGroup,
    setWorkflowBuilderTaskMap,
    setWorkflowGroupName,
    setWorkflowName,
} from "../../redux/ai/ai.reducer";
import {aiApiGateway} from "../../gateway/ai";
import {
    DeleteWorkflowsActionRequest,
    PostWorkflowsRequest,
    Retrieval,
    TaskModelInstructions
} from "../../redux/ai/ai.types";
import {TasksTable} from "./TasksTable";
import {isValidLabel} from "../clusters/wizard/builder/AddComponentBases";
import {RetrievalsTable} from "./RetrievalsTable";
import {loadBalancingApiGateway} from "../../gateway/loadbalancing";
import {setEndpoints, setGroupEndpoints} from "../../redux/loadbalancing/loadbalancing.reducer";

const mdTheme = createTheme();

function WorkflowEngineBuilder(props: any) {
    const [open, setOpen] = useState(true);
    const groups = useSelector((state: RootState) => state.loadBalancing.groups);
    const [loading, setIsLoading] = useState(false);
    const selectedWorkflows = useSelector((state: any) => state.ai.selectedWorkflows);
    const selectedMainTabBuilder = useSelector((state: any) => state.ai.selectedMainTabBuilder);
    const [selected, setSelected] = useState<{ [key: number]: boolean }>({});
    const analysisWorkflowInstructions = useSelector((state: RootState) => state.ai.analysisWorkflowInstructions);
    const aggregationWorkflowInstructions = useSelector((state: RootState) => state.ai.aggregationWorkflowInstructions);
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
    const aggregationStages = useSelector((state: RootState) => state.ai.addedAggregateTasks);
    const [tasks, setTasks] = useState(allTasks && allTasks.filter((task: TaskModelInstructions) => task.taskType === taskType));
    const retrievals = useSelector((state: RootState) => state.ai.retrievals);
    const workflowBuilderTaskMap = useSelector((state: RootState) => state.ai.workflowBuilderTaskMap);
    const taskMap = useSelector((state: RootState) => state.ai.taskMap);
    const retrievalsMap = useSelector((state: RootState) => state.ai.retrievalsMap);
    const workflowAnalysisRetrievalsMap = useSelector((state: RootState) => state.ai.workflowAnalysisRetrievalsMap);
    const retrieval = useSelector((state: RootState) => state.ai.retrieval);
    const workflowName = useSelector((state: RootState) => state.ai.workflowName);
    const workflowGroupName = useSelector((state: RootState) => state.ai.workflowGroupName);
    const workflows = useSelector((state: any) => state.ai.workflows);
    const action = useSelector((state: any) => state.ai.action);
    const actionMetric = useSelector((state: any) => state.ai.actionMetric);
    const actionPlatformAccount = useSelector((state: any) => state.ai.actionPlatformAccount);

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

    }, [addAggregateView, addAnalysisView,addRetrievalView, selectedMainTabBuilder, analysisStages, aggregationStages,retrievals, retrievalStages, workflowBuilderTaskMap, workflowAnalysisRetrievalsMap, taskMap]);
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
        }
        setIsLoading(false)
    }

    const addRetrievalStageView = async () => {
        const toggle = !addRetrievalView;
        dispatch(setAddAnalysisView(false));
        dispatch(setAddAggregationView(false));
        dispatch(setAddRetrievalView(toggle));
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
    const [analysisModel, setAnalysisModel] = React.useState('gpt-3.5-turbo-1106');
    const handleUpdateAnalysisModel = (event: any) => {
        setAnalysisModel(event.target.value);
    };
    const [analysisModelMaxTokens, setAnalysisModelMaxTokens] = React.useState(0);
    const handleUpdateAnalysisModelMaxTokens = (val: number) => {
        setAnalysisModelMaxTokens(val);
    };
    const [analysisModelTokenOverflowStrategy, setAnalysisModelTokenOverflowStrategy] = React.useState('deduce');
    const handleUpdateAnalysisModelTokenOverflowStrategy = (event: any) => {
        setAnalysisModelTokenOverflowStrategy(event.target.value);
    };
    const [aggregationModel, setAggregationModel] = React.useState('gpt-4');
    const handleUpdateAggregationModel = (event: any) => {
        setAggregationModel(event.target.value);
    };

    const [analysisName, setAnalysisName] = React.useState('');
    const handleUpdateAnalysisName = (event: any) => {
        setAnalysisName(event.target.value);
    };
    const [analysisGroupName, setAnalysisGroupName] = React.useState('default');
    const handleUpdateAnalysisGroupName = (event: any) => {
        setAnalysisGroupName(event.target.value);
    };
    const [aggregationName, setAggregationName] = React.useState('');
    const handleUpdateAggregationName = (event: any) => {
        setAggregationName(event.target.value);
    };
    const [aggregationGroupName, setAggregationGroupName] = React.useState('default');
    const handleUpdateAggregationGroupName = (event: any) => {
        setAggregationGroupName(event.target.value);
    };
    const [aggregationModelTokenOverflowStrategy, setAggregationModelTokenOverflowStrategy] = React.useState('deduce');
    const handleUpdateAggregationModelTokenOverflowStrategy = (event: any) => {
        setAggregationModelTokenOverflowStrategy(event.target.value);
    };
    const [aggregationModelMaxTokens, setAggregationModelMaxTokens] = React.useState(0);
    const handleUpdateAggregationModelMaxTokens = (val: number) => {
        setAggregationModelMaxTokens(val);
    };
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const handleUpdateAnalysisWorkflowInstructions =(value: string) => {
        dispatch(setAnalysisWorkflowInstructions(value));
    };
    const handleUpdateAggregationWorkflowInstructions =(value: string) => {
        dispatch(setAggregationWorkflowInstructions(value));
    };
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
            if (Object.keys(workflowAnalysisRetrievalsMap) && retrievalStages.length <= 0) {
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
                analysisRetrievalsMap: workflowAnalysisRetrievalsMap
            }
            setIsLoading(true)
            const response = await aiApiGateway.createAiWorkflowRequest(payload);
            const statusCode = response.status
            if (statusCode < 400) {
                const data = response.data;
                setSelected({});
                dispatch(setAddAnalysisView(false));
                dispatch(setAddAggregationView(false));
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

    const createOrUpdateRetrieval= async () => {
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

            if ((retrieval.retrievalKeywords.length <= 0 && retrieval.retrievalPrompt.length <= 0 && retrieval.retrievalGroup.length <= 0 ))  {
                setRequestRetrievalStatus('At least one of retrieval keywords or prompt or group must be set')
                setRequestRetrievalStatusError('error')
                return;
            }

            if (retrieval.retrievalPlatform.length <= 0) {
                setRequestRetrievalStatus('Retrieval platform must be set')
                setRequestRetrievalStatusError('error')
                return;
            }
            const response = await aiApiGateway.createOrUpdateRetrieval(retrieval);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data as Retrieval;
                dispatch(setRetrievals([...retrievals, data]))
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
            const tn = (taskType === 'analysis' ? analysisName : aggregationName);

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

            const taskGn = (taskType === 'analysis' ? analysisGroupName : aggregationGroupName);
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

            const prompt = (taskType === 'analysis' ? analysisWorkflowInstructions : aggregationWorkflowInstructions);
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
                taskType: taskType,
                taskGroup:taskGn,
                taskName: tn,
                model: (taskType === 'analysis' ? analysisModel : aggregationModel),
                group: (taskType === 'analysis' ? analysisGroupName : aggregationGroupName),
                prompt: (taskType === 'analysis' ? analysisWorkflowInstructions : aggregationWorkflowInstructions),
                maxTokens:  (taskType === 'analysis' ? analysisModelMaxTokens : aggregationModelMaxTokens),
                cycleCount: (taskType === 'analysis' ? 1 : 1),
                tokenOverflowStrategy: (taskType === 'analysis' ? analysisModelTokenOverflowStrategy : aggregationModelTokenOverflowStrategy),
            };
            const response = await aiApiGateway.createOrUpdateTaskRequest(task);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                if (taskType === 'analysis') {
                    setRequestAnalysisStatus('Task created successfully')
                    setRequestAnalysisStatusError('success')
                } else if (taskType === 'aggregation') {
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

    const createOrUpdateAction= async () => {
        try {
            setIsLoading(true)

            if (!isValidLabel(action.actionName)) {
                setRequestActionStatus('Action name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestActionStatusError('error')
                return;
            }

            if (!isValidLabel(action.actionGroupName)) {
                setRequestActionStatus('Action group name is invalid. It must be must be 63 characters or less and begin and end with an alphanumeric character and can contain contain dashes (-), underscores (_), dots (.), and alphanumerics between')
                setRequestActionStatusError('error')
                return;
            }

            const response = await aiApiGateway.createOrUpdateAction(retrieval);
            const statusCode = response.status;
            if (statusCode < 400) {
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
        if (addRetrievalView || selectedMainTabBuilder === 3) {
            const newSelection = retrievals.reduce((acc: { [key: number]: boolean }, task: any, index: number) => {
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
                            Mockingbird — A Time Series RAG-LLM Workflow Engine
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
                            <Card sx={{ minWidth: 500, maxWidth: 900 }}>
                                {( selectedMainTabBuilder === 0 || addAnalysisView || addAggregateView || addRetrievalView) &&
                                    <div>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Workflow Generation
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                            This allows you to write natural language instructions to chain to your search queries. Add a name
                                            for your workflow, and then write instructions for the AI to follow, and it will save the workflow for you.
                                            </Typography>

                                            <Stack direction={"row"} >
                                                <Box sx={{ width: '50%', mb: 0, mt: 2 }}>
                                                    <TextField
                                                        label={`Workflow Name`}
                                                        variant="outlined"
                                                        value={workflowName}
                                                        onChange={(event) => dispatch(setWorkflowName(event.target.value))}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box sx={{ width: '50%', mb: 0, mt: 2, ml: 2 }}>
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
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Retrieval Stages
                                                </Typography>
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 4}}>
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
                                                                value={ret?.retrievalPlatform || ''}
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
                                            <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                <Divider/>
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 2}}>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Analysis Stages
                                                </Typography>
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 4}}>
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
                                                    <Box sx={{ mt: 2 }} >
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
                                            { analysisStages && analysisStages.length > 0 &&
                                                <div>
                                                    <Box flexGrow={1} sx={{ mt: 4, mb: 2}}>
                                                        <Divider/>
                                                    </Box>
                                                    <Box flexGrow={2} sx={{mt:2 , mb: 4}}>
                                                        <Typography gutterBottom variant="h5" component="div">
                                                            Aggregation Stages
                                                        </Typography>
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
                                                                    <Box sx={{ mt: 2 }} >
                                                                        <Typography variant="h6" color="text.secondary">
                                                                            Analysis {'->'} Aggregation Dependencies
                                                                        </Typography>
                                                                    </Box>
                                                                    { workflowBuilderTaskMap &&
                                                                        <Box sx={{ mt:2,  ml: 2, mr: 2 }} >
                                                                            <Box >
                                                                                {Object.entries(workflowBuilderTaskMap).map(([key, value], index) => {
                                                                                    const taskNameForKey = taskMap[(Number(key))]?.taskName || '';
                                                                                    if (!taskNameForKey || taskNameForKey.length <= 0) {
                                                                                        return null;
                                                                                    }
                                                                                    return Object.entries(value).map(([subKey, subValue], subIndex) => {
                                                                                        if (!subValue || subKey.length <= 0) {
                                                                                            return null;
                                                                                        }
                                                                                        const subKeyNumber = Number(subKey);
                                                                                        const subTaskName = taskMap[(subKeyNumber)]?.taskName || '';
                                                                                        if (subTaskName.length <= 0) {
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
                                                                            <Box sx={{ mt: 2 }} >
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
                                                                        <Box flexGrow={1} sx={{ mt: -1.2, ml: 4, mr: -1 }}>
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
                                                </div>
                                            }
                                        </CardContent>
                                        <Box flexGrow={1} sx={{ mt: 4, mb: 0}}>
                                            <Divider/>
                                        </Box>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Fundamental Time Period
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                This is the time period that each cycle is referenced against for determining its next execution time. The workflow will run every time period, and will run all analysis and aggregation stages that are due to run if
                                                any during that discrete time step. If an analysis cycle is set to 1 and the fundamental time period is 5 minutes, it will run every 5 minutes.
                                            </Typography>
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
                                            {requestStatus != '' && (
                                                <Container sx={{  mt: 2}}>
                                                    <Typography variant="h6" color={requestStatusError}>
                                                        {requestStatus}
                                                    </Typography>
                                                </Container>
                                            )}
                                </div>
                                }
                                <CardContent>
                                    {!addAnalysisView && !addAggregateView && !addRetrievalView && selectedMainTabBuilder == 1 &&
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
                                                        value={analysisName}
                                                        onChange={handleUpdateAnalysisName}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box flexGrow={3} sx={{ width: '50%', mb: 0, mt: 2, ml: 1 }}>
                                                    <TextField
                                                        label={`Analysis Group`}
                                                        variant="outlined"
                                                        value={analysisGroupName}
                                                        onChange={handleUpdateAnalysisGroupName}
                                                        fullWidth
                                                    />
                                                </Box>
                                            </Stack>
                                            <Stack direction="row" >
                                                <Box flexGrow={3} sx={{ mb: 4, mt: 4 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="analysis-model-label">Analysis Model</InputLabel>
                                                        <Select
                                                            labelId="analysis-model-label"
                                                            id="analysis-model-select"
                                                            value={analysisModel}
                                                            label="Analysis Model"
                                                            onChange={handleUpdateAnalysisModel}
                                                        >
                                                            <MenuItem value="gpt-3.5-turbo-instruct">gpt-3.5-turbo-instruct</MenuItem>
                                                            <MenuItem value="gpt-3.5-turbo-1106">gpt-3.5-turbo-1106</MenuItem>
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
                                                <Box flexGrow={2} sx={{ mb: 4, mt: 4, ml:2 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="token-overflow-analysis-label">Token Overflow Strategy</InputLabel>
                                                        <Select
                                                            labelId="token-overflow-analysis-label"
                                                            id="analysis-overflow-analysis-select"
                                                            value={analysisModelTokenOverflowStrategy}
                                                            label="Token Overflow Strategy"
                                                            onChange={handleUpdateAnalysisModelTokenOverflowStrategy}
                                                        >
                                                            <MenuItem value="deduce">deduce</MenuItem>
                                                            <MenuItem value="truncate">truncate</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>
                                            </Stack>
                                            <Box  sx={{ mb: 2, mt: -2 }}>
                                                <TextareaAutosize
                                                    minRows={18}
                                                    value={analysisWorkflowInstructions}
                                                    onChange={(e) => handleUpdateAnalysisWorkflowInstructions(e.target.value)}
                                                    style={{ resize: "both", width: "100%" }}
                                                />
                                            </Box>
                                        </div>
                                    }
                                    { !addAggregateView && !addAnalysisView && !addRetrievalView && selectedMainTabBuilder == 2 &&
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
                                                        value={aggregationName}
                                                        onChange={handleUpdateAggregationName}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box flexGrow={3} sx={{ width: '50%', mb: 0, mt: 2, ml: 1 }}>
                                                    <TextField
                                                        label={`Aggregation Group`}
                                                        variant="outlined"
                                                        value={aggregationGroupName}
                                                        onChange={handleUpdateAggregationGroupName}
                                                        fullWidth
                                                    />
                                                </Box>
                                            </Stack>
                                            <Stack direction="row" >
                                                <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="platform-label">Aggregation Model</InputLabel>
                                                        <Select
                                                            labelId="platform-label"
                                                            id="platform-select"
                                                            value={aggregationModel}
                                                            label="Aggregation Model"
                                                            onChange={handleUpdateAggregationModel}
                                                        >
                                                            <MenuItem value="gpt-3.5-turbo-instruct">gpt-3.5-turbo-instruct</MenuItem>
                                                            <MenuItem value="gpt-3.5-turbo-1106">gpt-3.5-turbo-1106</MenuItem>
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
                                                <Box flexGrow={2} sx={{ mb: 2, mt: 4, ml:2 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="aggregation-token-overflow-analysis-label">Token Overflow Strategy</InputLabel>
                                                        <Select
                                                            labelId="aggregation-token-overflow-analysis-label"
                                                            id="aggregation-token-overflow-analysis-select"
                                                            value={aggregationModelTokenOverflowStrategy}
                                                            label="Token Overflow Strategy"
                                                            onChange={handleUpdateAggregationModelTokenOverflowStrategy}
                                                        >
                                                            <MenuItem value="deduce">deduce</MenuItem>
                                                            <MenuItem value="truncate">truncate</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>

                                            </Stack>
                                            <Box  sx={{ mb: 2, mt: 2 }}>
                                                <TextareaAutosize
                                                    minRows={18}
                                                    value={aggregationWorkflowInstructions}
                                                    onChange={(e) => handleUpdateAggregationWorkflowInstructions(e.target.value)}
                                                    style={{ resize: "both", width: "100%" }}
                                                />
                                            </Box>
                                        </div>
                                    }
                                    {
                                        selectedMainTabBuilder == 3 && !addRetrievalView && !loading &&
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
                                                                onChange={(e) => dispatch(setRetrievalName(e.target.value))}
                                                            />
                                                        </Box>
                                                    <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                        <TextField
                                                            fullWidth
                                                            id="retrieval-group"
                                                            label="Retrieval Group"
                                                            variant="outlined"
                                                            value={retrieval.retrievalGroup}
                                                            onChange={(e) => dispatch(setRetrievalGroup(e.target.value))}
                                                        />
                                                    </Box>
                                                    </Stack>
                                                    <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                    <FormControl fullWidth>
                                                        <InputLabel id="platform-label">Platform</InputLabel>
                                                        <Select
                                                            labelId="platform-label"
                                                            id="platforms-input"
                                                            value={retrieval.retrievalPlatform}
                                                            label="Platform"
                                                            onChange={(e) => dispatch(setRetrievalPlatform(e.target.value))}
                                                        >
                                                            <MenuItem value="web">Web</MenuItem>
                                                            <MenuItem value="reddit">Reddit</MenuItem>
                                                            <MenuItem value="twitter">Twitter</MenuItem>
                                                            <MenuItem value="discord">Discord</MenuItem>
                                                            <MenuItem value="telegram">Telegram</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                    </Box>
                                                { retrieval.retrievalPlatform !== 'web' &&
                                                    <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                        <TextField
                                                            fullWidth
                                                            id="group-input"
                                                            label={"Platform Groups"}
                                                            variant="outlined"
                                                            value={retrieval.retrievalPlatformGroups}
                                                            onChange={(e) => dispatch(setRetrievalPlatformGroups(e.target.value))}
                                                        />
                                                    </Box>
                                                }
                                                { retrieval.retrievalPlatform === 'web' &&
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
                                                                value={retrieval.webFilters?.routingGroup || ''}
                                                                onChange={(e) => dispatch(setWebRoutingGroup(e.target.value))}
                                                                label="Routing Group"
                                                            >
                                                                {Object.keys(groups).map((name) => <MenuItem key={name} value={name}>{name}</MenuItem>)}
                                                            </Select>
                                                        </FormControl>
                                                    </div>
                                                }
                                                    { retrieval.retrievalPlatform === 'discord' &&
                                                        <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="category-name-input"
                                                                label="Discord Category Name"
                                                                variant="outlined"
                                                                value={retrieval.discordFilters?.categoryName || ''}
                                                                onChange={(e) => dispatch(setDiscordOptionsCategoryName(e.target.value))}
                                                            />
                                                    </Box>
                                                    }
                                                { retrieval.retrievalPlatform !== 'web' &&
                                                    <Box flexGrow={1} sx={{ mb: 2, ml: 4, mr:4  }}>
                                                        <TextField
                                                            fullWidth
                                                            id="usernames-input"
                                                            label="Usernames"
                                                            variant="outlined"
                                                            value={retrieval.retrievalUsernames}
                                                            onChange={(e) => dispatch(setRetrievalUsernames(e.target.value))}
                                                        />
                                                    </Box>
                                                }

                                                { retrieval.retrievalPlatform === 'web' &&
                                                    <Typography variant="h5" color="text.secondary">
                                                        Describe how the AI should extract data from the website address.
                                                    </Typography>
                                                }
                                                { retrieval.retrievalPlatform !== 'web' &&
                                                    <Typography variant="h5" color="text.secondary">
                                                        Describe what you're looking for, and the AI will generate a list of keywords to search for after it runs for the first time.
                                                        You can preview, edit, or give the AI more information to refine the search.
                                                    </Typography>
                                                }
                                                    <Box  sx={{ mb: 2, mt: 2 }}>
                                                        <TextareaAutosize
                                                            minRows={18}
                                                            value={retrieval.retrievalPrompt}
                                                            onChange={(e) => dispatch(setRetrievalPrompt(e.target.value))}
                                                            style={{ resize: "both", width: "100%" }}
                                                        />
                                                    </Box>
                                                    <Typography variant="h5" color="text.secondary">
                                                        Add search keywords using comma separated values below.
                                                    </Typography>
                                                    <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                        <TextField
                                                            fullWidth
                                                            id="keywords-input"
                                                            label="Keywords"
                                                            variant="outlined"
                                                            value={retrieval.retrievalKeywords}
                                                            onChange={(e) => dispatch(setRetrievalKeywords(e.target.value))}
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
                                        selectedMainTabBuilder == 4 && !addRetrievalView && !loading && !addAnalysisView && !addAggregateView &&
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
                                                    You can add a trigger to any analysis or aggregation stage output.
                                                </Typography>
                                                <Stack direction="column" spacing={2} sx={{ mt: 4, mb: 0 }}>
                                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="action-name"
                                                                label="Action Name"
                                                                variant="outlined"
                                                                value={action.actionName}
                                                                onChange={(e) => dispatch(setAction({
                                                                    ...action, // Spread the existing action properties
                                                                    actionName: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="action-group"
                                                                label="Action Group"
                                                                variant="outlined"
                                                                value={action.actionGroupName}
                                                                onChange={(e) => dispatch(setAction({
                                                                    ...action, // Spread the existing action properties
                                                                    actionGroupName: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                    </Stack>
                                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                        <Box flexGrow={2} sx={{ mb: 2, mt: 4 }}>
                                                            <FormControl fullWidth>
                                                                <InputLabel id="platform-label">Platform</InputLabel>
                                                                <Select
                                                                    labelId="platform-label"
                                                                    id="platforms-input"
                                                                    value={actionPlatformAccount.actionPlatformName}
                                                                    label="Platform"
                                                                    onChange={(e) => dispatch(setActionPlatformAccount({
                                                                        ...actionPlatformAccount,
                                                                        actionPlatformName: e.target.value
                                                                    }))}
                                                                >
                                                                    <MenuItem value="email">Email</MenuItem>
                                                                    <MenuItem value="text">Text</MenuItem>
                                                                    <MenuItem value="reddit">Reddit</MenuItem>
                                                                    <MenuItem value="twitter">Twitter</MenuItem>
                                                                    <MenuItem value="discord">Discord</MenuItem>
                                                                    <MenuItem value="telegram">Telegram</MenuItem>
                                                                </Select>
                                                            </FormControl>
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="platform-account"
                                                                label="Platform Account"
                                                                variant="outlined"
                                                                value={actionMetric.metricName}
                                                                onChange={(e) => dispatch(setActionMetric({
                                                                    ...actionMetric, // Spread the existing action properties
                                                                    metricName: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                    </Stack>
                                                    <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                id="metric-name"
                                                                label="Metric Name"
                                                                variant="outlined"
                                                                value={actionMetric.metricName}
                                                                onChange={(e) => dispatch(setActionMetric({
                                                                    ...actionMetric, // Spread the existing action properties
                                                                    metricName: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                type={"number"}
                                                                id="metric-score-threshold"
                                                                label="Metric Score Threshold"
                                                                variant="outlined"
                                                                value={actionMetric.metricScoreThreshold}
                                                                onChange={(e) => dispatch(setActionMetric({
                                                                    ...actionMetric, // Spread the existing action properties
                                                                    metricScoreThreshold: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 2,ml: 4, mr:4  }}>
                                                            <TextField
                                                                fullWidth
                                                                type={"number"}
                                                                id="metric-action-multiplier"
                                                                label="Metric Action Multiplier"
                                                                variant="outlined"
                                                                value={actionMetric.metricPostActionMultiplier}
                                                                onChange={(e) => dispatch(setAction({
                                                                    ...actionMetric, // Spread the existing action properties
                                                                    metricPostActionMultiplier: e.target.value // Update the actionName
                                                                }))}
                                                            />
                                                        </Box>
                                                    </Stack>
                                                    {requestActionStatus != '' && (
                                                        <Container sx={{ mb: 2, mt: -2}}>
                                                            <Typography variant="h6" color={requestActionStatusError}>
                                                                {requestActionStatus}
                                                            </Typography>
                                                        </Container>
                                                    )}
                                                    <Box flexGrow={1} sx={{ mb: 0 }}>
                                                        <Button fullWidth variant="contained" onClick={createOrUpdateAction} >Save Action</Button>
                                                    </Box>
                                                </Stack>
                                            </div>
                                        </CardContent>
                                    }
                                    {/*<Typography gutterBottom variant="h5" component="div">*/}
                                    {/*    Time Intervals*/}
                                    {/*</Typography>*/}
                                    {/*<Typography variant="body2" color="text.secondary">*/}
                                    {/*    You can run an analysis on demand or use this to define an analysis chunk interval as part of an aggregate analysis.*/}
                                    {/*</Typography>*/}
                                    {/*<Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>*/}
                                    {/*    <Box sx={{ width: '33%' }}> /!* Adjusted Box for TextField *!/*/}
                                    {/*        <TextField*/}
                                    {/*            type="number"*/}
                                    {/*            label="Time Step Size"*/}
                                    {/*            variant="outlined"*/}
                                    {/*            inputProps={{ min: 1 }}  // Set minimum value to 1*/}
                                    {/*            value={stepSize}*/}
                                    {/*            onChange={handleTimeStepChange}*/}
                                    {/*            fullWidth*/}
                                    {/*        />*/}
                                    {/*    </Box>*/}
                                    {/*    <Box sx={{ width: '33%' }}> /!* Adjusted Box for FormControl *!/*/}
                                    {/*        <FormControl fullWidth>*/}
                                    {/*            <InputLabel id="time-unit-label">Time Unit</InputLabel>*/}
                                    {/*            <Select*/}
                                    {/*                labelId="time-unit-label"*/}
                                    {/*                id="time-unit-select"*/}
                                    {/*                value={stepSizeUnit}*/}
                                    {/*                label="Time Unit"*/}
                                    {/*                onChange={handleUpdateStepSizeUnit}*/}
                                    {/*            >*/}
                                    {/*                <MenuItem value="seconds">Seconds</MenuItem>*/}
                                    {/*                <MenuItem value="minutes">Minutes</MenuItem>*/}
                                    {/*                <MenuItem value="hours">Hours</MenuItem>*/}
                                    {/*                <MenuItem value="days">Days</MenuItem>*/}
                                    {/*                <MenuItem value="weeks">Weeks</MenuItem>*/}
                                    {/*            </Select>*/}
                                    {/*        </FormControl>*/}
                                    {/*    </Box>*/}
                                    {/*    <Box sx={{ width: '33%', mb: 4 }}>*/}
                                    {/*        <TextField*/}
                                    {/*            label={`Total Time (${stepSizeUnit})`} // Label now reflects the selected unit*/}
                                    {/*            variant="outlined"*/}
                                    {/*            value={stepSize* cycleCount}*/}
                                    {/*            InputProps={{*/}
                                    {/*                readOnly: true,*/}
                                    {/*            }}*/}
                                    {/*            fullWidth*/}
                                    {/*        />*/}
                                    {/*    </Box>*/}
                                    {/*</Stack>*/}
                                    {/*<Typography gutterBottom variant="h5" component="div">*/}
                                    {/*    Token Usage Limits*/}
                                    {/*</Typography>*/}
                                    {/*<Typography variant="body2" color="text.secondary">*/}
                                    {/*    Use this to limit how many tokens you want to use for your LLM stages. Set to 0 for unlimited.*/}
                                    {/*</Typography>*/}
                                    { !addAnalysisView && !addAggregateView && selectedMainTabBuilder == 1 &&
                                        <div>
                                            <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                                <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                                    <TextField
                                                        label={`Analysis Model`}
                                                        variant="outlined"
                                                        value={analysisModel}
                                                        InputProps={{
                                                            readOnly: true,
                                                        }}
                                                        fullWidth
                                                    />
                                                </Box>
                                                <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                                    <TextField
                                                        type="number"
                                                        label={`Max Tokens Analysis Model`}
                                                        variant="outlined"
                                                        value={analysisModelMaxTokens}
                                                        inputProps={{ min: 0 }}
                                                        onChange={(event) => handleUpdateAnalysisModelMaxTokens(parseInt(event.target.value, 10))}
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
                                    {  !addAnalysisView && !addAggregateView && selectedMainTabBuilder == 2 &&
                                        <div>
                                            <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 4 }}>
                                            <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                                <TextField
                                                    label={`Aggregation Model`}
                                                    variant="outlined"
                                                    value={aggregationModel}
                                                    InputProps={{
                                                        readOnly: true,
                                                    }}
                                                    fullWidth
                                                />
                                            </Box>
                                            <Box sx={{ width: '100%', mb: 4, mt: 4 }}>
                                                <TextField
                                                    type="number"
                                                    label={`Max Aggregation Token Usage`}
                                                    variant="outlined"
                                                    value={aggregationModelMaxTokens}
                                                    onChange={(event) => handleUpdateAggregationModelMaxTokens(parseInt(event.target.value, 10))}
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
                                    {/*{(selectedMainTab === 4) &&*/}
                                    {/*    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>*/}
                                    {/*        <Actions />*/}
                                    {/*    </Container>*/}
                                    {/*}*/}
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
                                <Tab className="onboarding-card-highlight-all-actions" label="Actions" />
                            </Tabs>
                        </Box>
                    </Container>

                    { selectedMainTabBuilder === 0 && !addAnalysisView && !addAggregateView && !addRetrievalView && selectedWorkflows && selectedWorkflows.length > 0  &&
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
                    { selectedMainTabBuilder == 0 &&
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
                    { (selectedMainTabBuilder === 1 || selectedMainTabBuilder === 2) &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <TasksTable tasks={tasks} selected={selected} handleClick={handleClick} handleSelectAllClick={handleSelectAllClick} />
                            </Container>
                        </div>
                    }
                    { (selectedMainTabBuilder === 3) && addRetrievalView &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <Box sx={{ mb: 2 }}>
                                <span>({Object.values(selected).filter(value => value).length} Selected Retrievals)</span>
                                <Button variant="outlined" color="secondary" onClick={handleAddTasksToWorkflow} style={{marginLeft: '10px'}}>
                                    Add Retrieval Stages
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
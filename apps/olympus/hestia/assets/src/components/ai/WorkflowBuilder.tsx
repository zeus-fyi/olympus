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
    setAddAggregateTasks,
    setAddAggregationView,
    setAddAnalysisTasks,
    setAddAnalysisView,
    setAggregationWorkflowInstructions,
    setAnalysisWorkflowInstructions,
} from "../../redux/ai/ai.reducer";
import {aiApiGateway} from "../../gateway/ai";
import {PostWorkflowsRequest, TaskModelInstructions, WorkflowModelInstructions} from "../../redux/ai/ai.types";
import {TasksTable} from "./TasksTable";

const mdTheme = createTheme();

function WorkflowEngineBuilder(props: any) {
    const [open, setOpen] = useState(true);
    const [loading, setIsLoading] = useState(false);
    const [selectedMainTab, setSelectedMainTab] = useState(0);
    const [selected, setSelected] = useState<{ [key: number]: boolean }>({});
    const analysisWorkflowInstructions = useSelector((state: RootState) => state.ai.analysisWorkflowInstructions);
    const aggregationWorkflowInstructions = useSelector((state: RootState) => state.ai.aggregationWorkflowInstructions);
    const addAnalysisView = useSelector((state: RootState) => state.ai.addAnalysisView);
    const addAggregateView = useSelector((state: RootState) => state.ai.addAggregationView);
    const allTasks = useSelector((state: any) => state.ai.tasks);
    const [taskType, setTaskType] = useState('analysis');
    const analysisStages = useSelector((state: RootState) => state.ai.addedAnalysisTasks);
    const aggregationStages = useSelector((state: RootState) => state.ai.addedAggregateTasks);
    const [tasks, setTasks] = useState(allTasks.filter((task: TaskModelInstructions) => task.taskType === taskType));
    useEffect(() => {
    }, [addAggregateView, addAnalysisView, selectedMainTab, analysisStages, aggregationStages]);
    const dispatch = useDispatch();
    const [analysisCycleCount, setAnalysisCycleCount] = useState(0);
    const handleAnalysisCycleCountChange = (val: number) => {
        if (val < aggregationCycleCount) {
            setAggregationCycleCount(val);
        }
        setAnalysisCycleCount(val);
    };

    const handleAddTasksToWorkflow = async (event: any) => {
        setIsLoading(true)
        const selectedTasks: TaskModelInstructions[] = Object.keys(selected)
            .filter(key => selected[Number(key)])
            .map(key => tasks[Number(key)]);
        if (addAnalysisView){
            dispatch(setAddAnalysisTasks(selectedTasks));
        } else if (addAggregateView){
            dispatch(setAddAggregateTasks(selectedTasks));
        }
        setIsLoading(false)
    }

    const handleRemoveTasksFromWorkflow = async (event: any) => {
        console.log('selectedKeys', selected);
    }

    const addAnalysisStageView = async () => {
        const toggle = !addAnalysisView;
        dispatch(setAddAnalysisView(toggle));
        dispatch(setAddAggregationView(false));
        if (toggle) {
            setSelectedMainTab(1)
            setSelected({});
            setTaskType('analysis');
            setTasks(allTasks.filter((task: any) => task.taskType === 'analysis'));
        } else {
            setSelectedMainTab(0)
        }
    }
    const addAggregationStageView = async () => {
        const toggle = !addAggregateView;
        dispatch(setAddAnalysisView(false));
        dispatch(setAddAggregationView(toggle));
        if (toggle) {
            setSelected({});
            setTaskType('aggregation');
            setTasks(allTasks.filter((task: any) => task.taskType === 'aggregation'));
            setSelectedMainTab(2)
        } else {
            setSelectedMainTab(0)
        }
    }
    const [aggregationCycleCount, setAggregationCycleCount] = useState(0);
    const handleAggregationCycleCountChange = (val: number) => {
        setAggregationCycleCount(val);
    };
    const [stepSize, setStepSize] = useState(5);
    const handleTimeStepChange = (event: any) => {
        setStepSize(event.target.value);
    };
    const [stepSizeUnit, setStepSizeUnit] = React.useState('minutes');
    const handleUpdateStepSizeUnit = (event: any) => {
        setStepSizeUnit(event.target.value);
    };
    const [analysisModel, setAnalysisModel] = React.useState('gpt-3.5-turbo');
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
    const [aggregationModel, setAggregationModel] = React.useState('gpt-4-1106-preview');
    const handleUpdateAggregationModel = (event: any) => {
        setAggregationModel(event.target.value);
    };
    const [workflowName, setWorkflowName] = React.useState('');
    const handleUpdateWorkflowName = (event: any) => {
        setWorkflowName(event.target.value);
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

    const createOrUpdateWorkflow = async (timeRange: '1 hour'| '24 hours' | '7 days'| '30 days' | 'window' | 'all') => {
        try {
            setIsLoading(true)
            const models: WorkflowModelInstructions[] = [
                {
                    instructionType: "analysis",
                    model: analysisModel,
                    prompt: analysisWorkflowInstructions,
                    maxTokens: analysisModelMaxTokens, // Optional
                    tokenOverflowStrategy: analysisModelTokenOverflowStrategy, // Optional
                    cycleCount: analysisCycleCount
                },
                {
                    instructionType: "aggregation",
                    model: aggregationModel,
                    prompt: aggregationWorkflowInstructions,
                    maxTokens: aggregationModelMaxTokens, // Optional
                    tokenOverflowStrategy: aggregationModelTokenOverflowStrategy, // Optional
                    cycleCount: aggregationCycleCount
                }
            ];
            const payload: PostWorkflowsRequest = {
                workflowName: workflowName,
                stepSize: stepSize,
                stepSizeUnit: stepSizeUnit,
                models: models,
            }
            const response = await aiApiGateway.createAiWorkflowRequest(payload);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
            } else {
                console.log('Failed to search', response);
            }
        } catch (e) {
        } finally {
            setIsLoading(false);
        }
    }

    const createOrUpdateTask = async (taskType: string) => {
        try {
            setIsLoading(true)
            const task: TaskModelInstructions = {
                taskType: taskType,
                taskGroup: (taskType === 'analysis' ? analysisGroupName : aggregationGroupName),
                taskName: (taskType === 'analysis' ? analysisName : aggregationName),
                model: (taskType === 'analysis' ? analysisModel : aggregationModel),
                group: (taskType === 'analysis' ? analysisGroupName : aggregationGroupName),
                prompt: (taskType === 'analysis' ? analysisWorkflowInstructions : aggregationWorkflowInstructions),
                maxTokens:  (taskType === 'analysis' ? analysisModelMaxTokens : aggregationModelMaxTokens),
                tokenOverflowStrategy: (taskType === 'analysis' ? analysisModelTokenOverflowStrategy : aggregationModelTokenOverflowStrategy),
            };
            const response = await aiApiGateway.createOrUpdateTaskRequest(task);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
            } else {
                console.log('Failed to search', response);
            }
        } catch (e) {
        } finally {
            setIsLoading(false);
        }
    }

    if (loading) {
        return <div>Loading...</div>;
    }

    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        if (newValue === 1) {
            setSelected({});
            setTaskType('analysis');
            setTasks(allTasks.filter((task: any) => task.taskType === 'analysis'));
        } else if (newValue === 2) {
            setSelected({});
            setTaskType('aggregation');
            setTasks(allTasks.filter((task: any) => task.taskType === 'aggregation'));
        }
        dispatch(setAddAnalysisView(false));
        dispatch(setAddAnalysisView(false));
        setSelectedMainTab(newValue);
    };
    const handleClick = (index: number) => {
        setSelected((prevSelected) => ({
            ...prevSelected,
            [index]: !prevSelected[index]
        }));
    }

    const handleSelectAllClick = (event: React.ChangeEvent<HTMLInputElement>) => {
        const isChecked = event.target.checked;
        const newSelection = tasks.reduce((acc: { [key: number]: boolean }, task: any, index: number) => {
            acc[index] = isChecked;
            return acc;
        }, {});
        setSelected(newSelection);
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
                            LLM Workflow Engine
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
                                {( selectedMainTab === 0 || addAnalysisView || addAggregateView) &&
                                    <div>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Workflow Generation
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                            This allows you to write natural language instructions to chain to your search queries. Add a name
                                            for your workflow, and then write instructions for the AI to follow, and it will save the workflow for you.
                                            </Typography>
                                            <Box sx={{ width: '100%', mb: 0, mt: 2 }}>
                                                <TextField
                                                    label={`Workflow Name`}
                                                    variant="outlined"
                                                    value={workflowName}
                                                    onChange={handleUpdateWorkflowName}
                                                    fullWidth
                                                />
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 2}}>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Analysis Stages
                                                </Typography>
                                                {/*<Typography variant="body2" color="text.secondary">*/}
                                                {/*    Add Analysis Stages*/}
                                                {/*</Typography>*/}
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 4}}>
                                                {analysisStages && analysisStages.map((task, index) => (
                                                    <Stack direction={"row"} key={index}>
                                                        <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                            <TextField
                                                                key={index}
                                                                // label={`Task ${index + 1}`}
                                                                value={task.taskName}
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
                                                                label="Analysis Cycle Count"
                                                                variant="outlined"
                                                                value={analysisCycleCount}
                                                                inputProps={{ min: 0 }}  // Set minimum value to 0
                                                                onChange={(event) => handleAnalysisCycleCountChange(parseInt(event.target.value, 10))}
                                                                fullWidth
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 0, ml: 2 }}>
                                                            <Button fullWidth variant="contained" >Remove</Button>
                                                        </Box>
                                                    </Stack>
                                                ))}
                                            </Box>
                                            <Box flexGrow={1} sx={{ mb: 0, mt: 2 }}>
                                                <Button fullWidth variant="contained" onClick={() => addAnalysisStageView()} >{addAnalysisView ? 'Done Adding': 'Add Analysis Stages'}</Button>
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 2, mb: 2}}>
                                                <Typography gutterBottom variant="h5" component="div">
                                                    Aggregation Stages
                                                </Typography>
                                                {/*<Typography variant="body2" color="text.secondary">*/}
                                                {/*    Add Aggregation Stages*/}
                                                {/*</Typography>*/}
                                            </Box>
                                            <Box flexGrow={2} sx={{mt: 4}}>
                                                {aggregationStages && aggregationStages.map((task, index) => (
                                                    <Stack direction={"row"} key={index}>
                                                        <Box flexGrow={2} sx={{ mt: -3, ml: 2 }}>
                                                            <TextField
                                                                key={index}
                                                                // label={`Task ${index + 1}`}
                                                                value={task.taskName}
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
                                                                value={aggregationCycleCount}
                                                                inputProps={{ min: 0 }}  // Set minimum value to 0
                                                                onChange={(event) => handleAggregationCycleCountChange(parseInt(event.target.value, 10))}
                                                                fullWidth
                                                            />
                                                        </Box>
                                                        <Box flexGrow={1} sx={{ mb: 0, ml: 2 }}>
                                                            <Button fullWidth variant="contained" >Remove</Button>
                                                        </Box>
                                                    </Stack>
                                                ))}
                                            </Box>
                                            <Box flexGrow={1} sx={{ mb: 0, mt: 2 }}>
                                                <Button fullWidth variant="contained" onClick={() => addAggregationStageView()} >{addAggregateView ? 'Done Adding' : 'Add Aggregation Stages' }</Button>
                                            </Box>
                                        </CardContent>
                                        <CardContent>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Time Intervals
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                You can run an analysis on demand or use this to define an analysis chunk interval as part of an aggregate analysis.
                                            </Typography>
                                            <Stack direction="row" spacing={2} sx={{ mt: 4, mb: 2 }}>
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
                                                            <MenuItem value="seconds">Seconds</MenuItem>
                                                            <MenuItem value="minutes">Minutes</MenuItem>
                                                            <MenuItem value="hours">Hours</MenuItem>
                                                            <MenuItem value="days">Days</MenuItem>
                                                            <MenuItem value="weeks">Weeks</MenuItem>
                                                        </Select>
                                                    </FormControl>
                                                </Box>
                                                <Box sx={{ width: '33%' }}>
                                                    <TextField
                                                        label={`Total Time (${stepSizeUnit})`} // Label now reflects the selected unit
                                                        variant="outlined"
                                                        value={stepSize* analysisCycleCount}
                                                        InputProps={{
                                                            readOnly: true,
                                                        }}
                                                        fullWidth
                                                    />
                                                </Box>
                                            </Stack>
                                        </CardContent>
                                        <CardActions>
                                            <Box flexGrow={1} sx={{ mb: -6, mt: -4 }}>
                                                <Button fullWidth variant="contained" onClick={() => createOrUpdateWorkflow('all')} >Save Workflow</Button>
                                            </Box>
                                        </CardActions>
                                </div>
                                }
                                <CardContent>
                                    {!addAnalysisView && !addAggregateView && selectedMainTab == 1 &&
                                        <div>
                                            <Typography gutterBottom variant="h5" component="div">
                                                Analysis Instructions
                                            </Typography>
                                            <Typography variant="body2" color="text.secondary">
                                                Token overflow strategy will determine how the AI will handle requests that are projected to exceed the maximum token length for the model you select, or has returned a result with that error.
                                                Deduce will chunk your analysis into smaller pieces and aggregate them into a final analysis result. Truncate will simply truncate the request
                                                to the maximum token length it can support.
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
                                                            <MenuItem value="gpt-3.5-turbo">gpt-3.5-turbo</MenuItem>
                                                            <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
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
                                    { !addAggregateView && !addAnalysisView && selectedMainTab == 2 &&
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
                                                        <InputLabel id="aggregation-model-label">Aggregation Model</InputLabel>
                                                        <Select
                                                            labelId="aggregation-model-label"
                                                            id="aggregation-model-select"
                                                            value={aggregationModel}
                                                            label="Aggregation Model"
                                                            onChange={handleUpdateAggregationModel}
                                                        >
                                                            <MenuItem value="gpt-3.5-turbo">gpt-3.5-turbo</MenuItem>
                                                            <MenuItem value="gpt-4-1106-preview">gpt-4-1106-preview</MenuItem>
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
                                    { !addAnalysisView && !addAggregateView && selectedMainTab == 1 &&
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
                                            <Box flexGrow={1} sx={{ mb: 2 }}>
                                                <Button fullWidth variant="contained" onClick={() =>  createOrUpdateTask('analysis')} >Save Analysis</Button>
                                            </Box>
                                        </div>
                                    }
                                    {  !addAnalysisView && !addAggregateView && selectedMainTab == 2 &&
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
                                            <Box flexGrow={1} sx={{ mb: 2 }}>
                                                <Button fullWidth variant="contained" onClick={() => createOrUpdateTask('aggregation')} >Save Aggregation</Button>
                                            </Box>
                                        </div>
                                    }
                                </CardContent>
                            </Card>
                        </Stack>
                    </Container>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                            <Tabs value={selectedMainTab} onChange={handleMainTabChange} aria-label="basic tabs">
                                <Tab className="onboarding-card-highlight-all-workflows" label="Workflows"  />
                                <Tab className="onboarding-card-highlight-all-analysis" label="Analysis" />
                                <Tab className="onboarding-card-highlight-all-aggregation" label="Aggregations" />
                            </Tabs>
                        </Box>
                    </Container>
                    { selectedMainTab == 0 &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <WorkflowTable />
                        </Container>
                    }
                    { (selectedMainTab === 1 || selectedMainTab === 2) && (addAggregateView || addAnalysisView) &&
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            <Box sx={{ mb: 2 }}>
                                <span>({Object.values(selected).filter(value => value).length} Selected Tasks)</span>
                                <Button variant="outlined" color="secondary" onClick={handleAddTasksToWorkflow} style={{marginLeft: '10px'}}>
                                    Add {addAnalysisView ? 'Analysis' : 'Aggregation'} Stages
                                </Button>
                            </Box>
                        </Container>
                    }
                    { (selectedMainTab === 1 || selectedMainTab === 2) &&
                        <div>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                <TasksTable tasks={tasks} selected={selected} handleClick={handleClick} handleSelectAllClick={handleSelectAllClick} />
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
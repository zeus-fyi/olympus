import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {Card, MenuItem, Select, Stack} from "@mui/material";
import Box from "@mui/material/Box";
import {setContactsFlowMap, setContactsOverrideFlowMap, setPromptFlowMap} from "../../redux/flows/flows.reducer";
import {Task, WorkflowTemplate} from "../../redux/ai/ai.types";

export function ContactsTextFieldRows(props: any) {
    const stages = useSelector((state: RootState) => state.flows.stages);
    const headers = useSelector((state: RootState) => state.flows.csvHeaders);
    const dispatch = useDispatch();
    const stageContactsOverrideMap = useSelector((state: RootState) => state.flows.stageContactsOverrideMap);
    const stageContactsMap = useSelector((state: RootState) => state.flows.stageContactsMap);
    // const selected = useSelector((state: any) => state.ai.selectedWorkflows);
    let workflows = useSelector((state: any) => state.ai.workflows);
    if (workflows) {
        workflows = workflows.filter((workflow: WorkflowTemplate) => {
            return workflow.tasks.some((task: Task) => task.responseFormat === 'csv');
        })
    }

    const flowNames = workflows.map((workflow: WorkflowTemplate) => {
        return workflow.workflowName
    })

    const handleSelectChange = async (promptName: string, wfTaskName: string) => {
        const payload = {
            key: promptName, // taskID
            subKey: wfTaskName, // retrievalID
        };
        dispatch(setContactsFlowMap(payload));
    }
    const handleSelectChangeContactsOverrideFlowMap = async (promptName: string, wfTaskName: string) => {
        const payload = {
            key: promptName, // taskID
            subKey: wfTaskName, // retrievalID
        };
        dispatch(setContactsOverrideFlowMap(payload));
    }

    return (
        <div>
            <Card sx={{ mb: 2, mt: 4 }}>
                {headers.map((header, index) => (
                    <Stack direction={"row"} key={index} sx={{ mb: 2, mt: 2 }}>
                        <Box flexGrow={1} sx={{ mt: 0, ml: 2 }}>
                            <TextField
                                fullWidth
                                label={header}
                                value={header}
                                inputProps={{ readOnly: true }}
                                variant="outlined"
                            />
                        </Box>
                        <Box flexGrow={1} sx={{ mt:0, ml: 2, mr: 2 }}>
                            <TextField
                                fullWidth
                                label={header}
                                variant="outlined"
                                value={header && stageContactsOverrideMap[header]}
                                onChange={(event) => handleSelectChangeContactsOverrideFlowMap(header, event.target.value as string)}
                            />
                        </Box>
                        <Box flexGrow={1} sx={{ mt: 0, ml: 2, mr: 2 }}>
                            <Select
                                fullWidth
                                value={stageContactsMap[header] || 'Default'}
                                label={header}
                                onChange={(event) => handleSelectChange(header, event.target.value as string)}
                                displayEmpty
                                variant="outlined"
                            >
                                <MenuItem value="Default">
                                    <em>Default</em>
                                </MenuItem>
                                {stages && Object.keys(stages).map((flow, flowIndex) => (
                                    <MenuItem key={flowIndex} value={flow}>{flow}</MenuItem>
                                ))}
                                {flowNames && flowNames.map((name: string, nameIndex: number) => (
                                    <MenuItem key={`wf-${nameIndex}`} value={name}>{name}</MenuItem>
                                ))}
                            </Select>
                        </Box>
                    </Stack>
                ))}
            </Card>
        </div>
    );
}

export function PromptsTextFieldRows(props: any) {
    const stages = useSelector((state: RootState) => state.flows.stages);
    // const flowList = useSelector((state: RootState) => state.flows.flowList);
    const stagePromptMap = useSelector((state: RootState) => state.flows.stagePromptMap);
    const dispatch = useDispatch();
    const headers = useSelector((state: RootState) => state.flows.promptHeaders);
    let workflows = useSelector((state: any) => state.ai.workflows);
    if (workflows) {
        workflows = workflows.filter((workflow: WorkflowTemplate) => {
            return workflow.tasks.some((task: Task) => task.responseFormat === 'csv');
        })
    }

    const flowNames = workflows.map((workflow: WorkflowTemplate) => {
        return workflow.workflowName
    })
    const handleSelectChange = async (promptName: string, wfTaskName: string) => {
        const payload = {
            key: promptName, // taskID
            subKey: wfTaskName, // retrievalID
        };
        dispatch(setPromptFlowMap(payload));
    }

    return (
        <div>
            <Card sx={{ mb: 2, mt: 4 }}>
                {headers.map((header, index) => (
                    <Stack direction={"row"} key={index} sx={{ mb: 2, mt: 2 }}>
                        <Box flexGrow={1} sx={{ mt: 0, ml: 2 }}>
                            <TextField
                                fullWidth
                                label={header}
                                value={header}
                                inputProps={{ readOnly: true }}
                                variant="outlined"
                            />
                        </Box>
                        {/*<Box flexGrow={1} sx={{ mt:0, ml: 2, mr: 2 }}>*/}
                        {/*    <TextField*/}
                        {/*        fullWidth*/}
                        {/*        label={header}*/}
                        {/*        variant="outlined"*/}
                        {/*    />*/}
                        {/*</Box>*/}
                        <Box flexGrow={1} sx={{ mt: 0, ml: 2, mr: 2 }}>
                            <Select
                                fullWidth
                                value={stagePromptMap[header] || 'Default'}
                                label={header}
                                onChange={(event) => handleSelectChange(header, event.target.value as string)}
                                displayEmpty
                                variant="outlined"
                            >
                                <MenuItem value="Default">
                                    <em>Default</em>
                                </MenuItem>
                                {stages && Object.keys(stages).map((flow, flowIndex) => (
                                    <MenuItem key={flowIndex} value={flow}>{flow}</MenuItem>
                                ))}
                                {flowNames && flowNames.map((name: string, nameIndex: number) => (
                                    <MenuItem key={`wf-${nameIndex}`} value={name}>{name}</MenuItem>
                                ))}
                            </Select>
                        </Box>
                    </Stack>
                ))}
            </Card>
        </div>
    );
}
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {Card, MenuItem, Select, Stack} from "@mui/material";
import Box from "@mui/material/Box";
import {setColFlowMap, setPromptFlowMap} from "../../redux/flows/flows.reducer";

export function ContactsTextFieldRows(props: any) {
    const stages = useSelector((state: RootState) => state.flows.stages);
    const headers = useSelector((state: RootState) => state.flows.csvHeaders);
    const dispatch = useDispatch();
    const contactsColMap = useSelector((state: RootState) => state.flows.stageColMap);

    const handleSelectChange = async (promptName: string, wfTaskName: string) => {
        const payload = {
            key: promptName, // taskID
            subKey: wfTaskName, // retrievalID
        };
        dispatch(setColFlowMap(payload));
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
                                value={contactsColMap[header] || 'Default'}
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
                            </Select>
                        </Box>
                    </Stack>
                ))}
            </Card>
        </div>
    );
}
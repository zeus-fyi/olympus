import * as React from "react";
import {useState} from "react";
import {Card, CircularProgress, Tab, Tabs} from "@mui/material";
import {MbTaskCmdPrompt} from "./CommandPrompt";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import {SetupCard} from "./Setup";
import {aiApiGateway} from "../../gateway/ai";
import {useDispatch, useSelector} from "react-redux";
import {setCommandPrompt} from "../../redux/flows/flows.reducer";

export function Commands(props: any) {
    const bodyPrompts = useSelector((state: any) => state.flows.uploadContentTasks);
    const contacts = useSelector((state: any) => state.flows.uploadContentContacts);
    const cmds = useSelector((state: any) => state.flows.commandPrompts);
    const [checked, setChecked] = React.useState(false);
    const [gs, setGsChecked] = React.useState(false);
    const [selectedMainTab, setSelectedMainTab] = useState(0);
    const dispatch = useDispatch();
    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedMainTab(newValue);
    }
    const handleChangeGs = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setGsChecked(event.target.checked);
    };
    const handleChange = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setChecked(event.target.checked);
    };

    const handleChangeGoogleSearchPrompt = (event: string) => {
        // Construct the new commandPrompts object
        const newCommandPrompts = {
            ...cmds,
            googleSearch: event
        };
        // Dispatch an action to update the state with the new commandPrompts object
        dispatch(setCommandPrompt(newCommandPrompts));
    };

    const handleChangeLinkedInPrompt = (event: string) => {
        // Construct the new commandPrompts object
        const newCommandPrompts = {
            ...cmds,
            linkedIn: event
        };
        console.log(newCommandPrompts)
        // Dispatch an action to update the state with the new commandPrompts object
        dispatch(setCommandPrompt(newCommandPrompts));
    };


    let buttonLabelCreate;
    let buttonDisabledCreate;
    let statusMessageCreate;
    const [flowsRequestStatus, setFlowsRequestStatus] = useState('');
    switch (flowsRequestStatus) {
        case 'pending':
            buttonLabelCreate = <CircularProgress size={20}/>;
            buttonDisabledCreate = true;
            break;
        case 'success':
            buttonLabelCreate = 'Send';
            buttonDisabledCreate = false;
            statusMessageCreate = 'Request Sent Successfully!';
            break;
        case 'insufficientTokenBalance':
            buttonLabelCreate = 'Send';
            buttonDisabledCreate = true;
            statusMessageCreate = 'Insufficient Token Balance. Email alex@zeus.fyi to request more tokens.'
            break;
        case 'error':
            buttonLabelCreate = 'Send';
            buttonDisabledCreate = false;
            statusMessageCreate = ''
            break;
        default:
            buttonLabelCreate = 'Send';
            buttonDisabledCreate = false;
            break;
    }
    const onClickSubmit = async () => {
        try {
            setFlowsRequestStatus('pending');
            const fa = {
                contentContactsCsv: contacts,
                contentContactsFieldMaps: {},
                promptsCsv: bodyPrompts,
                stages: {
                    linkedIn: checked,
                    googleSearch: gs
                },
               commandPrompts: cmds
            }
            let res: any = await aiApiGateway.flowsRequest(fa)
            const statusCode = res.status;
            if (statusCode === 200 || statusCode === 204) {
                setFlowsRequestStatus('success');
            } else if (statusCode === 412) {
                setFlowsRequestStatus('insufficientTokenBalance');
            } else {
                setFlowsRequestStatus('error');
            }
        } catch (e: any) {
            if (e.response && e.response.status === 412) {
                setFlowsRequestStatus('insufficientTokenBalance');
            } else {
                setFlowsRequestStatus('error');
            }
        }
    }
    const getTabName = (selectedTab: number): string => {
        if (selectedTab === 0) {
            return 'Google Search: ';
        } else if (selectedTab === 1) {
            return 'LinkedIn: ';
        } else {
            return '';
        }
    }

    return (
        <div>
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'left', justifyContent: 'center', mb: 2 }}>
                <SetupCard checked={checked} gs={gs} handleChangeGs={handleChangeGs} handleChange={handleChange} />
            </Box>
            <Card sx={{ maxWidth: 1200, justifyContent: 'center' }}>
            <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large', fontWeight: 'thin', marginRight: '15px', color: '#151C2F' }}>
                <span style={{ fontSize: 'large', fontWeight: 'thin', color: '#151C2F' }}>{getTabName(selectedMainTab)}</span> Agent Tasking Commands
            </Typography>
                {selectedMainTab === 0 && (
                <MbTaskCmdPrompt language={"plaintext"} code={cmds.googleSearch} onChange={handleChangeGoogleSearchPrompt} height={"200px"} width={"1200px"}/>
            )}
            {selectedMainTab === 1 && (
                <MbTaskCmdPrompt language={"plaintext"} code={cmds.linkedIn} onChange={handleChangeLinkedInPrompt} height={"200px"} width={"1200px"}/>
            )}
            <Box mt={2} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'right' }}>
                <Button
                    variant="contained"
                    onClick={() => onClickSubmit()}
                    sx={{ backgroundColor: '#00C48C', '&:hover': { backgroundColor: '#00A678' }}}
                >
                    Submit
                </Button>
                {statusMessageCreate && (
                    <Typography variant="body2" color={flowsRequestStatus === 'error' ? 'error' : 'success'}>
                        {statusMessageCreate}
                    </Typography>
                )}
            </Box>
            </Card>
            <Box sx={{ mb: 2, mt: 2, ml: 0, mr:0  }}>
                <Tabs value={selectedMainTab} onChange={handleMainTabChange} aria-label="basic tabs">
                    <Tab label="Google Search"/>
                    <Tab label="LinkedIn" />
                </Tabs>
            </Box>
        </div>
    );
}

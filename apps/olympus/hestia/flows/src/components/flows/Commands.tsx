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
    const bodyPrompts = useSelector((state: any) => state.flows.promptsCsvContent);
    const contacts = useSelector((state: any) => state.flows.uploadContentContacts);
    const cmds = useSelector((state: any) => state.flows.commandPrompts);
    const [checked, setChecked] = React.useState(false);
    const [checkedLi, setCheckedLi] = React.useState(false);
    const [multiPromptOn, setMultiPromptOn] = React.useState(false);
    const [gs, setGsChecked] = React.useState(false);
    const [webChecked, setWebChecked] = React.useState(false);
    const [vesChecked, setVesChecked] = React.useState(false);
    const [selectedMainTab, setSelectedMainTab] = useState(0);
    const dispatch = useDispatch();
    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setSelectedMainTab(newValue);
    }
    const handleChangeWebChecked = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setWebChecked(event.target.checked);
    }
    const handleChangeVesChecked = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setVesChecked(event.target.checked);
    }
    const handleChangeLi = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setCheckedLi(event.target.checked);
    }
    const handleChangeGs = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setGsChecked(event.target.checked);
    };
    const handleChange = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setChecked(event.target.checked);
    };
    const handleChangeMultiPromptOn = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setMultiPromptOn(event.target.checked);
    }
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
        // Dispatch an action to update the state with the new commandPrompts object
        dispatch(setCommandPrompt(newCommandPrompts));
    };

    const handleChangeWebScrapePrompt = (event: string) => {
        // Construct the new commandPrompts object
        const newCommandPrompts = {
            ...cmds,
            websiteScrape: event
        };
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
                contentContactsCsv: [] as [],
                contentContactsCsvStr: objectArrayToCsv(contacts),
                contentContactsFieldMaps: {},
                promptsCsv: [] as [],
                promptsCsvStr: objectArrayToCsv(bodyPrompts),
                stages: {
                    linkedIn: checked,
                    linkedInBiz: checkedLi,
                    googleSearch: gs,
                    validateEmails: vesChecked,
                    websiteScrape: webChecked
                },
               commandPrompts: cmds
            }
            // console.log(fa, 'sffsf')
            let res: any = await aiApiGateway.flowsRequest(fa)
            const statusCode = res.status;
            if (statusCode >= 200 && statusCode < 300) {
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
            return 'Google Search ';
        } else if (selectedTab === 1) {
            return 'LinkedIn ';
        } else if (selectedTab === 2) {
            return 'Website ';
        } else {
            return '';
        }
    }

    return (
        <div>
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'left', justifyContent: 'center', mb: 2 }}>
                <SetupCard
                    webChecked={webChecked} handleChangeWebChecked={handleChangeWebChecked}
                    vesChecked={vesChecked} handleChangeVesChecked={handleChangeVesChecked}
                    multiPromptOn={multiPromptOn} handleChangeMultiPromptOn={handleChangeMultiPromptOn}
                    checkedLi={checkedLi} handleChangeLi={handleChangeLi}
                    checked={checked} gs={gs} handleChangeGs={handleChangeGs}
                    handleChange={handleChange} />
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
            {selectedMainTab === 2 && (
                <MbTaskCmdPrompt language={"plaintext"} code={cmds.websiteScrape} onChange={handleChangeWebScrapePrompt} height={"200px"} width={"1200px"}/>
            )}
            <Box mt={2} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'right' }}>
                <Button
                    variant="contained"
                    disabled={buttonDisabledCreate}
                    onClick={() => onClickSubmit()}
                    sx={{ backgroundColor: '#00C48C', '&:hover': { backgroundColor: '#00A678' }}}
                >
                    {buttonLabelCreate}
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
                    <Tab label="Website" />
                </Tabs>
            </Box>
        </div>
    );
}
const objectArrayToCsv = <T extends Record<string, unknown>>(data: T[]): string => {
    if (data.length === 0) {
        return '';
    }

    // Extract headers
    const headers = Object.keys(data[0]).join(',');

    // Extract rows
    const rows = data.map(obj =>
        Object.values(obj).map(val => {
            // Handle values that contain commas, double-quotes, or newlines by enclosing in double quotes
            if (typeof val === 'string') {
                return `"${val.replace(/"/g, '""')}"`; // Escape double quotes
            }
            return String(val);
        }).join(',')
    );

    // Combine headers and rows
    return [headers, ...rows].join('\r\n');
};
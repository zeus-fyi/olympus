import * as React from "react";
import {useState} from "react";
import {Card, CardActionArea, CircularProgress} from "@mui/material";
import {MbTaskCmdPrompt} from "./CommandPrompt";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import {SetupCard} from "./Setup";
import {aiApiGateway} from "../../gateway/ai";
import {useSelector} from "react-redux";

export function Commands(props: any) {
    const bodyPrompts = useSelector((state: any) => state.flows.uploadContentTasks);
    const contacts = useSelector((state: any) => state.flows.uploadContentContacts);

    const [checked, setChecked] = React.useState(false);
    const [gs, setGsChecked] = React.useState(false);
    const handleChangeGs = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setGsChecked(event.target.checked);
    };
    const handleChange = (event: { target: { checked: boolean | ((prevState: boolean) => boolean); }; }) => {
        setChecked(event.target.checked);
    };
    const [code, setCode] = React.useState("");
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
                promptsCsv: bodyPrompts,
                stages: {
                    linkedin: checked,
                    googleSearch: gs
                },
                commandPrompts: {
                    linkedin: '',
                    googleSearch: code

                }
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
    return (
        <div>
            <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'left', justifyContent: 'center', mb: 2 }}>
                <SetupCard checked={checked} gs={gs} handleChangeGs={handleChangeGs} handleChange={handleChange} />
            </Box>
            <Card sx={{ maxWidth: 1200, justifyContent: 'center' }}>
            <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                Agent Tasking Commands
            </Typography>
                <CardActionArea>
                    <MbTaskCmdPrompt language={"plaintext"} code={code} onChange={setCode} height={"200px"} width={"1200px"}/>
                    <Box mt={2} sx={{ display: 'flex', flexDirection: 'column', alignItems: 'right' }}>
                        <Button
                            variant="contained"
                            onClick={onClickSubmit}
                            // disabled={buttonDisabledCreate}
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
                </CardActionArea>
            </Card>

        </div>
    );
}
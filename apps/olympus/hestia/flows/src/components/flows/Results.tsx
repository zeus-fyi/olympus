import * as React from "react";
import {useState} from "react";
import {Card, CardContent, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import {useDispatch, useSelector} from "react-redux";
import {WorkflowAnalysisTable} from "./WorkflowAnalysisTable";
import Container from "@mui/material/Container";
import Button from "@mui/material/Button";
import {PostRunsActionRequest} from "../../redux/ai/ai.types";
import {aiApiGateway} from "../../gateway/ai";
import {setSelectedRuns} from "../../redux/ai/ai.reducer";

export function Results(props: any) {
    const selectedRuns = useSelector((state: any) => state.ai.selectedRuns);
    const [requestRunsStatus, setRequestRunsStatus] = useState('');
    const [requestRunsStatusError, setRequestRunsStatusError] = useState('');
    const [loading, setIsLoading] = React.useState(false);
    const dispatch = useDispatch();
    const runs = useSelector((state: any) => state.ai.runs);
    const handleRunsActionRequest = async (event: any, action: string) => {
        const params: PostRunsActionRequest = {
            action: action,
            runs: selectedRuns.map((index: number) => {
                return runs[index].orchestration
            })
        }
        if (params.runs.length === 0) {
            return
        }
        try {
            setIsLoading(true)
            const response = await aiApiGateway.execRunsActionRequest(params);
            const statusCode = response.status;
            if (statusCode < 400) {
                const data = response.data;
                dispatch(setSelectedRuns([]));
                setRequestRunsStatus('Run cancellation submitted successfully')
                setRequestRunsStatusError('success')
            }
        } catch (error: any) {
            const status: number = await error?.response?.status || 500;
            if (status === 412) {
                setRequestRunsStatus('Billing setup required. Please configure your billing information to continue using this service.');
                setRequestRunsStatusError('error')
            }
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <div>
            <Card sx={{ maxWidth: 320 }}>
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large', fontWeight: 'thin', marginRight: '15px', color: '#151C2F' }}>
                        Results
                    </Typography>
                </CardContent>
            </Card>
            <Box sx={{ mt: 4 }}>
                {requestRunsStatus != '' && (
                    <Container sx={{  mt: 2}}>
                        <Typography variant="h6" color={requestRunsStatusError}>
                            {requestRunsStatus}
                        </Typography>
                    </Container>
                )}
                { selectedRuns && selectedRuns.length > 0  &&
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Stack direction="row" spacing={2}>
                            <Box sx={{ mb: 2 }}>
                                <span>({selectedRuns.length} Selected Runs)</span>
                                <Button variant="outlined" color="secondary" onClick={(event) => handleRunsActionRequest(event, 'stop')} style={{marginLeft: '10px'}}>
                                    Stop { selectedRuns.length === 1 ? 'Run' : 'Runs' }
                                </Button>
                            </Box>
                            <Box sx={{ ml: -2, mb: 2 }}>
                                <Button variant="outlined" color="secondary" onClick={(event) => handleRunsActionRequest(event, 'archive')} style={{marginLeft: '10px'}}>
                                    Archive { selectedRuns.length === 1 ? 'Run' : 'Runs' }
                                </Button>
                            </Box>
                        </Stack>
                    </Container>
                }
                <WorkflowAnalysisTable csvExport={true} />
            </Box>
        </div>
    );
}
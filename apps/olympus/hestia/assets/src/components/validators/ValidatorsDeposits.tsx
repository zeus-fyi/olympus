import {Card, CardActions, CardContent, CircularProgress, Container, Stack} from "@mui/material";
import * as React from "react";
import {useState} from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {Network} from "./ZeusServiceRequest";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {
    createValidatorsDepositsDataJSON,
    createValidatorsDepositServiceRequest,
    ValidatorDepositDataRxJSON,
    validatorsApiGateway
} from "../../gateway/validators";
import {ValidatorsUploadActionAreaCard} from "./ValidatorsUpload";
import {setDepositData} from "../../redux/aws_wizard/aws.wizard.reducer";

export function ValidatorsDepositRequestAreaCardWrapper(props: any) {
    const { activeStep, onValidatorsDepositsUpload, authorizedNetworks} = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ValidatorsUploadActionAreaCard onValidatorsDepositsUpload={onValidatorsDepositsUpload} authorizedNetworks={authorizedNetworks}/>,
            <ValidatorsDepositRequestAreaCard authorizedNetworks={authorizedNetworks} />
        </Stack>
    );
}

export function ValidatorsDepositRequestAreaCard(props: any) {
    const {authorizedNetworks } = props;
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <SubmitValidators authorizedNetworks={authorizedNetworks}/>
            </Container >
        </div>
    );
}

export function SubmitValidators(props: any) {
    const {authorizedNetworks } = props;
    const depositData = useSelector((state: RootState) => state.awsCredentials.depositData);
    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    const dispatch = useDispatch();
    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const [requestStatus, setRequestStatus] = useState('');

    switch (requestStatus) {
        case 'pending':
            buttonLabel = <CircularProgress size={20} />;
            buttonDisabled = true;
            break;
        case 'success':
            buttonLabel = 'Submitted successfully';
            buttonDisabled = true;
            statusMessage = 'Validator deposits submitted successfully!';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while sending validator deposits.';
            break;
        default:
            buttonLabel = 'Submit Deposits is Paused';
            buttonDisabled = true;
            break;
    }
    const onClickSendValidatorsDeposits = async () => {
        try {
            setRequestStatus('pending');
            const depositParams = depositData.map((dd: any) => {
                return createValidatorsDepositsDataJSON(dd.pubkey, dd.withdrawal_credentials, dd.signature, dd.deposit_data_root,dd.amount,dd.deposit_message_root,dd.fork_version);
            });
            const reqParams = createValidatorsDepositServiceRequest(network, depositParams)
            const response = await validatorsApiGateway.depositValidatorsServiceRequest(reqParams);
            if (response.status === 202) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
            }
            const sd :[ValidatorDepositDataRxJSON] = response.data;
            const hm: { [key: string]: string } = {};
            sd.forEach((d: ValidatorDepositDataRxJSON) => {
                hm[d.pubkey] = d.rx;
            });
            const rxDepositData =  depositData.map((obj: any) => {
                    if (obj.hasOwnProperty('rx')) {
                        return {
                            ...obj,
                            ['rx']: hm[obj.pubkey],
                        };
                    }
                });
            dispatch(setDepositData(rxDepositData));
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }}
    return (
        <Card sx={{ maxWidth: 500 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Send Validator Deposits to Network
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    If you didn't generate your validator deposits from the previous step, you can upload your
                    deposit data JSON file to the left. We'll automatically pay for the 32Eth deposit fee & submit your
                    deposits when you use the Ephemery network.
                </Typography>
            </CardContent>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <Network network={network} authorizedNetworks={authorizedNetworks}/>
            </Container>
            <CardActions>
                <Button size="small" onClick={onClickSendValidatorsDeposits} disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}

import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import * as React from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {setFeeRecipient, setKeyGroupName, setNetworkName} from "../../redux/validators/ethereum.validators.reducer";
import {AgeEncryptionKeySecretName} from "./AwsSecrets";
import {ExternalAccessSecretName} from "./AwsExtUserAndLambdaVerify";

export function ZeusServiceRequestAreaCardWrapper(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ZeusServiceRequestAreaCard />
        </Stack>
    );
}

export function ZeusServiceRequestAreaCard() {
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <ZeusServiceRequestParams />
            </Container >
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <ZeusServiceRequest />
            </Container >
        </div>
    );
}

// TODO
export function ZeusServiceRequest() {
    const handleZeusServiceRequest = async () => {
        try {
            // TODO, get external accesss key and secret key from redux store
            //const response = await awsApiGateway.verifyLambdaKeySigning();
        } catch (error) {
            console.log("error", error);
        }};

    return (
        <Card sx={{ maxWidth: 500 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Create Zeus Validators Service Request
                </Typography>
            </CardContent>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <AgeEncryptionKeySecretName />
                <ExternalAccessSecretName />
            </Container>
            <CardActions>
                <Button onClick={handleZeusServiceRequest} size="small">Submit</Button>
            </CardActions>
        </Card>
    );
}

export function ZeusServiceRequestParams() {
    return (
        <div style={{ display: 'flex' }}>

        <Card sx={{ maxWidth: 500 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Set Zeus Validators Service Params
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Sets Zeus Validators Service Params
                </Typography>
            </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Network />
                    <KeyGroupName />
                    <FeeRecipient />
                </Container>
        </Card>
        </div>
    );
}

export function KeyGroupName() {
    const dispatch = useDispatch();
    const keyGroupName = useSelector((state: RootState) => state.validatorSecrets.keyGroupName);
    const onAccessKeyGroupName = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newKeyGroupName = event.target.value;
        dispatch(setKeyGroupName(newKeyGroupName));
    };
    return (
        <TextField
            fullWidth
            id="keyGroupName"
            label="Key Group Name"
            variant="outlined"
            value={keyGroupName}
            onChange={onAccessKeyGroupName}
            sx={{ width: '100%' }}
        />
    );
}

export function FeeRecipient() {
    const dispatch = useDispatch();
    const feeRecipient = useSelector((state: RootState) => state.validatorSecrets.feeRecipient);
    const onAccessFeeRecipientChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newFeeRecipient = event.target.value;
        dispatch(setFeeRecipient(newFeeRecipient));
    };
    return (
        <TextField
            fullWidth
            id="feeRecipient"
            label="Fee Recipient"
            variant="outlined"
            value={feeRecipient}
            onChange={onAccessFeeRecipientChange}
            sx={{ width: '100%' }}
        />
    );
}

export function Network(props: any) {
    const dispatch = useDispatch();
    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    const onAccessSetNetwork = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newNetworkName = event.target.value;
        dispatch(setNetworkName(newNetworkName));
    };
    return (
        <TextField
            fullWidth
            id="network"
            label="Network"
            variant="outlined"
            value={network}
            onChange={onAccessSetNetwork}
            sx={{ width: '100%' }}
        />
    );
}
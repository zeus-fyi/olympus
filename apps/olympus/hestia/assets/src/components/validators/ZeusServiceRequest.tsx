import {
    Card,
    CardActions,
    CardContent,
    Container,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    Stack
} from "@mui/material";
import * as React from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {setFeeRecipient, setKeyGroupName, setNetworkName} from "../../redux/validators/ethereum.validators.reducer";
import {AgeEncryptionKeySecretName} from "./AwsSecrets";
import {ExternalAccessSecretName} from "./AwsExtUserAndLambdaVerify";
import {
    createAuthAwsLambda,
    createValidatorOrgGroup,
    createValidatorServiceRequest,
    validatorsApiGateway
} from "../../gateway/validators";
import {getNetworkId} from "./Validators";
import {awsApiGateway} from "../../gateway/aws";

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

export function ZeusServiceRequest() {
    const feeRecipient = useSelector((state: RootState) => state.validatorSecrets.feeRecipient);
    const depositData = useSelector((state: RootState) => state.awsCredentials.depositData);
    const keyGroupName = useSelector((state: RootState) => state.validatorSecrets.keyGroupName);
    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const externalAccessUserName = useSelector((state: RootState) => state.awsCredentials.externalAccessUserName);
    const externalAccessSecretName = useSelector((state: RootState) => state.awsCredentials.externalAccessSecretName);
    const ageSecretName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);
    const blsSignerFunctionName = useSelector((state: RootState) => state.awsCredentials.blsSignerFunctionName);

    const handleZeusServiceRequest = async () => {
        try {
            const validatorServiceRequestSlice = depositData.map((dd: any) => {
                return createValidatorOrgGroup(dd.pubkey, feeRecipient)
                })
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            const signerUrl = await awsApiGateway.getLambdaFunctionURL(creds, blsSignerFunctionName);
            const url = await signerUrl.data
            const getExtCreds = await awsApiGateway.createOrFetchExternalLambdaUserAccessKeys(creds,externalAccessUserName, externalAccessSecretName);
            const extCreds = {accessKeyId: getExtCreds.data.accessKey, secretAccessKey: getExtCreds.data.secretKey};
            const serviceAuth = createAuthAwsLambda(url, ageSecretName,extCreds);
            const protocolID = getNetworkId(network);
            const sr = createValidatorServiceRequest(keyGroupName,protocolID,serviceAuth,validatorServiceRequestSlice)
            const response = await validatorsApiGateway.createValidatorsServiceRequest(sr);
            console.log("response", response)
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
        <div>
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
    const onAccessSetNetwork = (selectedNetwork: string) => {
        console.log('Selected network:', selectedNetwork);
        dispatch(setNetworkName(selectedNetwork));
    };

    return (
        <FormControl variant="outlined" style={{ minWidth: '100%' }}>
            <InputLabel id="network-label">Network</InputLabel>
            <Select
                labelId="network-label"
                id="network"
                value={network}
                label="Network"
                onChange={(event) => onAccessSetNetwork(event.target.value as string)}
                sx={{ width: '100%' }}
            >
                <MenuItem value="Ephemery">Ephemery</MenuItem>
            </Select>
        </FormControl>
    );
}


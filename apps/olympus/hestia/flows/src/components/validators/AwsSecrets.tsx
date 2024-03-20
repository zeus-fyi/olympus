import * as React from "react";
import {useState} from "react";
import {Box, Card, CardActions, CardContent, CircularProgress, Container, Stack} from "@mui/material";
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import CryptoJS from 'crypto-js';
import {
    LambdaFunctionGenEncZipFileCreation,
    LambdaFunctionGenValidatorDepositsCreation,
    LambdaFunctionSecretsCreation
} from "./AwsLambdaCreation";
import Button from "@mui/material/Button";
import {awsLambdaApiGateway} from "../../gateway/aws.lambda";
import Typography from "@mui/material/Typography";
import {setAgeSecretName, setValidatorSecretsName} from "../../redux/aws_wizard/aws.wizard.reducer";
import {Network} from "./ZeusServiceRequest";

export const charsets = {
    NUMBERS: '0123456789',
    LOWERCASE: 'abcdefghijklmnopqrstuvwxyz',
    UPPERCASE: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
    SYMBOLS: '!#$%&()*+,-/<=>?@[]^`{|}~',
};

export const generatePassword = (length: number, charset: string): string => {
    const x = CryptoJS.lib.WordArray.random(charset.length)
    let result = '';
    for (let i = 0; i < length; i++) {
        const index = x.words[i] % charset.length;
        result += charset.charAt(index);
    }
    return result;
}

export function CreateAwsInternalLambdasActionAreaCardWrapper(props: any) {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <LambdaFunctionSecretsCreation />
            <LambdaFunctionGenEncZipFileCreation />
            <LambdaFunctionGenValidatorDepositsCreation />
        </Stack>
    );
}

export function CreateAwsSecretsActionAreaCardWrapper(props: any) {
    const {authorizedNetworks, pageView,onGenerateValidatorDepositsAndZip } = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <CreateAwsSecretNamesAreaCard pageView={pageView} authorizedNetworks={authorizedNetworks} onGenerateValidatorDepositsAndZip={onGenerateValidatorDepositsAndZip}/>
        </Stack>
    );
}

export function CreateAwsSecretNamesAreaCard(props: any) {
    const {authorizedNetworks, pageView,onGenerateValidatorDepositsAndZip } = props;
    const sgLambdaURL = useSelector((state: RootState) => state.awsCredentials.secretGenLambdaFnUrl);
    const ak = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const sk = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const awsValidatorSecretName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const awsAgeEncryptionKeyName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);
    const network = useSelector((state: RootState) => state.validatorSecrets.network);

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
            buttonLabel = 'Created successfully';
            buttonDisabled = true;
            statusMessage = 'Secrets generated successfully!';
            break;
        case 'error':
            buttonLabel = 'Error creating secrets';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating the secrets.';
            break;
        case 'errorAuth':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'Update your AWS credentials on step 1 and try again.';
            break;
        default:
            buttonLabel = 'Create';
            buttonDisabled = false;
            break;
    }

    const validatorSecretName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const dispatch = useDispatch();
    const onCreateNewValidatorSecrets = async () => {
        try {
            if (!ak || !sk) {
                setRequestStatus('errorAuth');
                return;
            }
            const creds = {accessKeyId: ak, secretAccessKey: sk};
            const response = await awsLambdaApiGateway.invokeValidatorSecretsGeneration(sgLambdaURL, creds, awsValidatorSecretName, awsAgeEncryptionKeyName);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};
    return (
        <Card sx={{ maxWidth: 500 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Generate New Secrets Using Lambda
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Securely generates a new mnemonic, hd wallet password, and age encryption key and saves them in your secret manager with the below key names. If
                    secrets already exist with the same key name, it will not overwrite them.
                </Typography>
            </CardContent>
                <Stack direction="column" alignItems="center" spacing={2}>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        <Network authorizedNetworks={authorizedNetworks}/>
                    </Box>
                    <Box mt={2}>
                        <ValidatorSecretName validatorSecretName={validatorSecretName}/>
                    </Box>
                    <Box mt={2}>
                        <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName+network}/>
                    </Box>
                </Container>
                <CardActions>
                    <Button onClick={onCreateNewValidatorSecrets} size="small" disabled={buttonDisabled}>{buttonLabel}</Button>
                </CardActions>
                    {statusMessage && (
                        <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                            {statusMessage}
                        </Typography>
                    )}
                </Stack>
        </Card>
    );
}

export function ValidatorSecretName(props: any) {
    const dispatch = useDispatch();
    const validatorSecretName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const onValidatorSecretNameNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newValidatorSecretsName = event.target.value;
        dispatch(setValidatorSecretsName(newValidatorSecretsName));
    };
    return (
        <TextField
            fullWidth
            id="validatorSecretName"
            label="AWS Validator Key Secret Name"
            variant="outlined"
            value={validatorSecretName}
            onChange={onValidatorSecretNameNameChange}
            sx={{ width: '100%' }}
        />
    );
}

export function Mnemonic() {
    const mnemonic = useSelector((state: RootState) => state.validatorSecrets.mnemonic);
    const dispatch = useDispatch();
    const onMnemonicChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newMnemonicValue = event.target.value;
    };
    return (
        <TextField
            fullWidth
            id="mnemonic"
            label="24 Word Mnemonic"
            variant="outlined"
            value={mnemonic}
            onChange={onMnemonicChange}
            sx={{ width: '100%' }}
        />
    );
}

export function HDWalletPassword() {
    const dispatch = useDispatch();
    const hdWalletPw = useSelector((state: RootState) => state.validatorSecrets.hdWalletPw);
    const onHdWalletPwChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newHdWalletPw = event.target.value;
    };
    return (
        <TextField
            fullWidth
            id="hdWalletPassword"
            label="HD Wallet Password"
            variant="outlined"
            value={hdWalletPw}
            onChange={onHdWalletPwChange}
            sx={{ width: '100%' }}
        />
    );
}

export function AgeEncryptionKeySecretName(props: any) {
    const dispatch = useDispatch();
    const awsAgeEncryptionKeyName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);
    const onAccessAwsAgeEncryptionKeyName = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newAgeSecretName = event.target.value;
        dispatch(setAgeSecretName(newAgeSecretName));
    };
    return (
        <TextField
            fullWidth
            id="ageEncryptionKeySecretName"
            label="AWS Age Encryption Key Secret Name"
            variant="outlined"
            value={awsAgeEncryptionKeyName}
            onChange={onAccessAwsAgeEncryptionKeyName}
            sx={{ width: '100%' }}
        />
    );
}

import * as React from "react";
import {useState} from "react";
import {Box, Card, CardActions, CardContent, CircularProgress} from "@mui/material";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {awsApiGateway} from "../../gateway/aws";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {
    setDepositsGenLambdaFnUrl,
    setEncKeystoresZipLambdaFnUrl,
    setSecretGenLambdaFnUrl
} from "../../redux/aws_wizard/aws.wizard.reducer";
import TextField from "@mui/material/TextField";

export function LambdaFunctionSecretsCreation() {
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const sgLambdaURL = useSelector((state: RootState) => state.awsCredentials.secretGenLambdaFnUrl);
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
            statusMessage = 'Lambda function created successfully!';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating the lambda function.';
            break;
        default:
            buttonLabel = 'Create | Restore';
            buttonDisabled = false;
            break;
    }
    const dispatch = useDispatch();
    const onCreateLambdaSecretsFn = async () => {
        try {
            setRequestStatus('pending');
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            const response = await awsApiGateway.createValidatorSecretsLambda(creds);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setSecretGenLambdaFnUrl(response.data));
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Creation for Trustless Secrets Generation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a lambda function in AWS that securely generates a mnemonic, hdWalletPassword, and Age Encryption key and stores them in your secrets manager.
                </Typography>
            </CardContent>
            <Box ml={2} mr={2}>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="secretGenLambdaFnUrl"
                    label="SecretGenLambdaFnUrl"
                    name="secretGenLambdaFnUrl"
                    value={sgLambdaURL}
                    autoFocus
                />
            </Box>
            <CardActions>
                <Button onClick={onCreateLambdaSecretsFn} size="small" disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}

export function LambdaFunctionGenValidatorDepositsCreation() {
    const accKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const depositsGenLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.depositsGenLambdaFnUrl);

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
            statusMessage = 'Lambda function created successfully!';
            break;
        case 'error':
            buttonLabel = 'Error creating lambda function';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating the lambda function.';
            break;
        default:
            buttonLabel = 'Create | Restore';
            buttonDisabled = false;
            break;
    }
    const dispatch = useDispatch();
    const onCreateLambdaValidatorDepositsFn = async () => {
        try {
            setRequestStatus('pending');
            const creds = {accessKeyId: accKey, secretAccessKey: secKey};
            const response = await awsApiGateway.createValidatorsDepositDataLambda(creds);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setDepositsGenLambdaFnUrl(response.data));
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function For Secure Validator Deposits Generation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a lambda function in AWS that securely generates validator deposit messages using your mnemonic from secrets manager.
                </Typography>
            </CardContent>
            <Box ml={2} mr={2}>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="depositsGenLambdaFnUrl"
                    label="DepositsGenLambdaFnUrl"
                    name="depositsGenLambdaFnUrl"
                    value={depositsGenLambdaFnUrl}
                    autoFocus
                />
            </Box>
            <CardActions>
                <Button onClick={onCreateLambdaValidatorDepositsFn} size="small" disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}

export function LambdaFunctionGenEncZipFileCreation() {
    const ak = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const sk = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const encKeystoresZipLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.encKeystoresZipLambdaFnUrl);

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
            buttonLabel = 'Created successfully';
            buttonDisabled = true;
            statusMessage = 'Lambda function created successfully!';
            break;
        case 'error':
            buttonLabel = 'Error creating lambda function';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating the lambda function.';
            break;
        default:
            buttonLabel = 'Create | Restore';
            buttonDisabled = false;
            break;
    }
    const onCreateLambdaEncryptedKeystoresZipFn = async () => {
        try {
            setRequestStatus('pending');
            const creds = {accessKeyId: ak, secretAccessKey: sk};
            const response = await awsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(creds);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
            }
            dispatch(setEncKeystoresZipLambdaFnUrl(response.data));
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Encrypted Keystores Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a lambda function in AWS that securely generates an encrypted zip file with validator signing keys
                    using your mnemonic from secret manager.
                </Typography>
            </CardContent>
            <Box ml={2} mr={2}>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="encKeystoresZipLambdaFnUrl"
                    label="EncryptedKeystoresZipLambdaFnUrl"
                    name="encKeystoresZipLambdaFnUrl"
                    value={encKeystoresZipLambdaFnUrl}
                    autoFocus
                />
            </Box>
            <CardActions>
                <Button onClick={onCreateLambdaEncryptedKeystoresZipFn} size="small" disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}


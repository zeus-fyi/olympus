import * as React from "react";
import {Card, CardActions, CardContent} from "@mui/material";
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

    const dispatch = useDispatch();
    const onCreateLambdaSecretsFn = async () => {
        try {
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            const response = await awsApiGateway.createValidatorSecretsLambda(creds);
            dispatch(setSecretGenLambdaFnUrl(response.data));
        } catch (error) {
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
            <CardActions>
                <Button onClick={onCreateLambdaSecretsFn} size="small">Create | Restore</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionGenValidatorDepositsCreation() {
    const accKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const depositsGenLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.depositsGenLambdaFnUrl);

    const dispatch = useDispatch();
    const onCreateLambdaValidatorDepositsFn = async () => {
        try {
            const creds = {accessKeyId: accKey, secretAccessKey: secKey};
            const response = await awsApiGateway.createValidatorsDepositDataLambda(creds);
            dispatch(setDepositsGenLambdaFnUrl(response.data));
        } catch (error) {
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
            <CardActions>
                <Button onClick={onCreateLambdaValidatorDepositsFn} size="small">Create | Restore</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionGenEncZipFileCreation() {
    const ak = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const sk = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const encKeystoresZipLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.encKeystoresZipLambdaFnUrl);

    const dispatch = useDispatch();

    const onCreateLambdaEncryptedKeystoresZipFn = async () => {
        try {
            const creds = {accessKeyId: ak, secretAccessKey: sk};
            const response = await awsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(creds);
            dispatch(setEncKeystoresZipLambdaFnUrl(response.data));
        } catch (error) {
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
            <CardActions>
                <Button onClick={onCreateLambdaEncryptedKeystoresZipFn} size="small">Create | Restore</Button>
            </CardActions>
        </Card>
    );
}


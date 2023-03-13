import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {EncryptedKeystoresZipUploadActionAreaCard} from "./AwsExtUserAndLambdaVerify";
import {awsApiGateway} from "../../gateway/aws";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {
    setBlsSignerLambdaFnUrl,
    setDepositsGenLambdaFnUrl,
    setEncKeystoresZipLambdaFnUrl,
    setSecretGenLambdaFnUrl
} from "../../redux/aws_wizard/aws.wizard.reducer";

export function CreateAwsLambdaFunctionActionAreaCardWrapper(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <EncryptedKeystoresZipUploadActionAreaCard />
            <CreateAwsLambdaFunctionAreaCard />
        </Stack>
    );
}

export function CreateAwsLambdaFunctionAreaCard() {
    return (
            <div style={{ display: 'flex' }}>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <LambdaFunctionKeystoresLayerCreation />
                </Container >
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <LambdaFunctionCreation />
                </Container >
            </div>
    );
}

export function LambdaFunctionKeystoresLayerCreation() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Keystores Layer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates your encrypted keystores layer for usage in your AWS lambda function using your generated encrypted zip keystores file.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionCreation() {
    const acKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const seKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const dispatch = useDispatch();

    const onCreateLambdaSignerFn = async () => {
        try {
            console.log('onCreateLambdaSignerFn')
            const response = await awsApiGateway.createLambdaFunction(acKey, seKey);
            console.log(response.data)
            dispatch(setBlsSignerLambdaFnUrl(response.data));

        } catch (error) {
            console.log("error", error);
        }};
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Signer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a BLS signer lambda function in AWS that decrypts your keystores with your Age Encryption key to sign messages.
                </Typography>
            </CardContent>
            <CardActions>
                <Button onClick={onCreateLambdaSignerFn} size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionSecretsCreation() {
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const dispatch = useDispatch();

    const onCreateLambdaSecretsFn = async () => {
        try {
            console.log('onCreateLambdaSecretsFn')
            const response = await awsApiGateway.createValidatorSecretsLambda(accessKey, secretKey);
            console.log(response.data)
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
            <CardActions>
                <Button onClick={onCreateLambdaSecretsFn} size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionGenValidatorDepositsCreation() {
    const accKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const dispatch = useDispatch();

    const onCreateLambdaValidatorDepositsFn = async () => {
        try {
            console.log('onCreateLambdaValidatorDepositsFn')
            const response = await awsApiGateway.createValidatorsDepositDataLambda(accKey, secKey);
            console.log(response.data)
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
            <CardActions>
                <Button onClick={onCreateLambdaValidatorDepositsFn} size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionGenEncZipFileCreation() {
    const ak = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const sk = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const dispatch = useDispatch();

    const onCreateLambdaEncryptedKeystoresZipFn = async () => {
        try {
            console.log('onCreateLambdaEncryptedKeystoresZipFn')
            const response = await awsApiGateway.createValidatorsDepositDataLambda(ak, sk);
            console.log(response.data)
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
                    Creates a lambda function in AWS that generates an encrypted zip file with validator signing keys.
                </Typography>
            </CardContent>
            <CardActions>
                <Button onClick={onCreateLambdaEncryptedKeystoresZipFn} size="small">Create</Button>
            </CardActions>
        </Card>
    );
}


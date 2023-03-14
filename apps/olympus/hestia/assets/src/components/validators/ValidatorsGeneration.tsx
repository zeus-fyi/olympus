import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import TextField from "@mui/material/TextField";
import {AgeEncryptionKeySecretName, ValidatorSecretName} from "./AwsSecrets";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {awsApiGateway} from "../../gateway/aws";
import {setDepositsGenLambdaFnUrl, setEncKeystoresZipLambdaFnUrl} from "../../redux/aws_wizard/aws.wizard.reducer";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";

export function GenerateValidatorKeysAndDepositsAreaCardWrapper(props: any) {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <GenerateValidatorsParams />
            <GenValidatorDepositsCreationActionsCard />
            <GenerateZipValidatorActionsCard />
        </Stack>
    );
}

export function GenerateValidatorsParams() {
    const awsValidatorsNetwork = useSelector((state: RootState) => state.validatorSecrets.network);
    const validatorCount = useSelector((state: RootState) => state.validatorSecrets.validatorCount);
    const offset = useSelector((state: RootState) => state.validatorSecrets.hdOffset);
    const awsValidatorSecretName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const awsAgeEncryptionKeyName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                    <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName}/>
                    <ValidatorsNetwork awsValidatorsNetwork={awsValidatorsNetwork}/>
                    <ValidatorCount validatorCount={validatorCount}/>
                    <ValidatorOffsetHD offset={offset}/>
                </Container>
            </div>
        </Card>

    );
}

export function GenerateZipValidatorActionsCard() {
    const ak = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const sk = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const encKeystoresZipLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.encKeystoresZipLambdaFnUrl);

    const dispatch = useDispatch();

    const onCreateLambdaEncryptedKeystoresZipFn = async () => {
        try {
            const response = await awsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(ak, sk);
            dispatch(setEncKeystoresZipLambdaFnUrl(response.data));
        } catch (error) {
            console.log("error", error);
        }};
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Generate Encrypted Keystores.zip
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Uses your lambda function to generate an encrypted keystores zip file for creating your lambda function signing layer in the next step using
                    the age encryption key name you've set on the left.
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
                <Button onClick={onCreateLambdaEncryptedKeystoresZipFn} size="small">Generate</Button>
            </CardActions>
        </Card>
    );
}

export function GenValidatorDepositsCreationActionsCard() {
    const accKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const depositsGenLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.depositsGenLambdaFnUrl);

    const dispatch = useDispatch();
    const onCreateLambdaValidatorDepositsFn = async () => {
        try {
            const response = await awsApiGateway.createValidatorsDepositDataLambda(accKey, secKey);
            dispatch(setDepositsGenLambdaFnUrl(response.data));
        } catch (error) {
            console.log("error", error);
        }};
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Generate Validator Deposits
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Uses your lambda function in AWS to securely generates validator deposit messages using your mnemonic
                    from secrets manager with the secret key name you've set on the left.
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
                <Button onClick={onCreateLambdaValidatorDepositsFn} size="small">Generate</Button>
            </CardActions>
        </Card>
    );
}

// TODO set these in redux
export function ValidatorsNetwork(props: any) {
    const { awsValidatorsNetwork, setAwsValidatorsNetwork } = props;

    return (
        <TextField
            fullWidth
            id="validatorsNetwork"
            label="Network Name"
            variant="outlined"
            value={awsValidatorsNetwork}
            onChange={setAwsValidatorsNetwork}
            sx={{ width: '100%' }}
        />
    );
}
// TODO set these in redux
export function ValidatorCount(props: any) {
    const { validatorCount, onValidatorCountChange } = props;
    return (
        <TextField
            fullWidth
            id="validatorCount"
            label="Validator Count"
            variant="outlined"
            value={validatorCount}
            onChange={onValidatorCountChange}
            sx={{ width: '100%' }}
        />
    );
}
// TODO set these in redux
export function ValidatorOffsetHD(props: any) {
    const { offset, setOffset } = props;
    return (
        <TextField
            fullWidth
            id="validatorOffsetHD"
            label="Validator HD Offset"
            variant="outlined"
            value={offset}
            onChange={setOffset}
            sx={{ width: '100%' }}
        />
    );
}
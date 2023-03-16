import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import TextField from "@mui/material/TextField";
import {AgeEncryptionKeySecretName, ValidatorSecretName} from "./AwsSecrets";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {setHdOffset, setNetworkName, setValidatorCount} from "../../redux/validators/ethereum.validators.reducer";

export function GenerateValidatorKeysAndDepositsAreaCardWrapper(props: any) {
    const { activeStep, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip } = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <GenerateValidatorsParams />
            <GenValidatorDepositsCreationActionsCard onGenerateValidatorDeposits={onGenerateValidatorDeposits}/>
            <GenerateZipValidatorActionsCard onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip} />
        </Stack>
    );
}

export function GenerateValidatorsParams() {
    const awsValidatorSecretName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const awsAgeEncryptionKeyName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                    <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName}/>
                    <ValidatorsNetwork />
                    <ValidatorCount />
                    <ValidatorOffsetHD />
                </Container>
            </div>
        </Card>
    );
}

export function GenerateZipValidatorActionsCard(props: any) {
    const { activeStep, onGenerateValidatorEncryptedKeystoresZip } = props;
    const encKeystoresZipLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.encKeystoresZipLambdaFnUrl);

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
                <Button onClick={onGenerateValidatorEncryptedKeystoresZip} size="small">Generate</Button>
            </CardActions>
        </Card>
    );
}

export function GenValidatorDepositsCreationActionsCard(props: any) {
    const { activeStep, onGenerateValidatorDeposits } = props;

    const depositsGenLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.depositsGenLambdaFnUrl);

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
                <Button onClick={onGenerateValidatorDeposits} size="small">Generate</Button>
            </CardActions>
        </Card>
    );
}

export function ValidatorsNetwork() {
    const dispatch = useDispatch();
    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    const onValidatorsNetworkChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        try {
            const network = event.target.value;
            dispatch(setNetworkName(network));
        } catch (error) {
            console.log("error", error);
        }};
    return (
        <TextField
            fullWidth
            id="validatorsNetwork"
            label="Network Name"
            variant="outlined"
            value={network}
            onChange={onValidatorsNetworkChange}
            sx={{ width: '100%' }}
        />
    );
}
export function ValidatorCount() {
    const dispatch = useDispatch();
    const validatorCount = useSelector((state: RootState) => state.validatorSecrets.validatorCount);
    const onValidatorCountChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        try {
            const validatorCount = parseInt(event.target.value);
            dispatch(setValidatorCount(validatorCount));
        } catch (error) {
            console.log("error", error);
        }};
    return (
        <TextField
            fullWidth
            id="validatorCount"
            label="Validator Count"
            variant="outlined"
            type="number"
            value={validatorCount}
            onChange={onValidatorCountChange}
            sx={{ width: '100%' }}
        />
    );
}
export function ValidatorOffsetHD() {
    const dispatch = useDispatch();
    const hdOffset = useSelector((state: RootState) => state.validatorSecrets.hdOffset);
    const onHdOffsetChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        try {
            const hdOffset = parseInt(event.target.value);
            dispatch(setHdOffset(hdOffset));
        } catch (error) {
            console.log("error", error);
        }};
    return (
        <TextField
            fullWidth
            id="validatorOffsetHD"
            label="Validator HD Offset"
            variant="outlined"
            type="number"
            value={hdOffset}
            onChange={onHdOffsetChange}
            sx={{ width: '100%' }}
        />
    );
}
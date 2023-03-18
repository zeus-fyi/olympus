import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import TextField from "@mui/material/TextField";
import {AgeEncryptionKeySecretName, ValidatorSecretName} from "./AwsSecrets";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {setHdOffset, setValidatorCount} from "../../redux/validators/ethereum.validators.reducer";
import {Network} from "./ZeusServiceRequest";

export function GenerateValidatorKeysAndDepositsAreaCardWrapper(props: any) {
    const { activeStep, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip,
        zipGenButtonLabel, zipGenButtonEnabled, zipGenStatus,
        buttonLabelVd, buttonDisabledVd, statusMessageVd,} = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <GenerateValidatorsParams />
            <GenValidatorDepositsCreationActionsCard onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                                                     buttonLabelVd={buttonLabelVd}
                                                     buttonDisabledVd={buttonDisabledVd}
                                                     statusMessageVd={statusMessageVd}
            />
            <GenerateZipValidatorActionsCard onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
                                             zipGenButtonLabel={zipGenButtonLabel}
                                             zipGenButtonEnabled={zipGenButtonEnabled}
                                             zipGenStatus={zipGenStatus}
            />
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
                    <Network />
                    <ValidatorCount />
                    <ValidatorOffsetHD />
                </Container>
            </div>
        </Card>
    );
}

export function GenerateZipValidatorActionsCard(props: any) {
    const { activeStep, onGenerateValidatorEncryptedKeystoresZip, zipGenButtonLabel, zipGenButtonEnabled, zipGenStatus } = props;
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
                <Button onClick={onGenerateValidatorEncryptedKeystoresZip} size="small" disabled={zipGenButtonEnabled}>{zipGenButtonLabel}</Button>
            </CardActions>
            {zipGenStatus && (
                <Typography variant="body2" color={zipGenStatus === 'error' ? 'error' : 'success'}>
                    {zipGenStatus}
                </Typography>
            )}
        </Card>
    );
}

export function GenValidatorDepositsCreationActionsCard(props: any) {
    const { activeStep, onGenerateValidatorDeposits, buttonLabelVd, buttonDisabledVd, statusMessageVd } = props;

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
                <Button onClick={onGenerateValidatorDeposits} size="small" disabled={buttonDisabledVd}>{buttonLabelVd}</Button>
            </CardActions>
            {statusMessageVd && (
                <Typography variant="body2" color={statusMessageVd === 'error' ? 'error' : 'success'}>
                    {statusMessageVd}
                </Typography>
            )}
        </Card>
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
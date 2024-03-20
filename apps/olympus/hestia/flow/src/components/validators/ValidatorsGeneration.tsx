import * as React from "react";
import {Box, Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import TextField from "@mui/material/TextField";
import {AgeEncryptionKeySecretName, ValidatorSecretName} from "./AwsSecrets";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {
    setHdOffset,
    setValidatorCount,
    setWithdrawalCredentials
} from "../../redux/validators/ethereum.validators.reducer";
import {Network} from "./ZeusServiceRequest";

export function GenerateValidatorKeysAndDepositsAreaCardWrapper(props: any) {
    const { activeStep, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip,
        zipGenButtonLabel, zipGenButtonEnabled, zipGenStatus, requestStatusZipGen,
        buttonLabelVd, buttonDisabledVd, statusMessageVd,authorizedNetworks, pageView,onGenerateValidatorDepositsAndZip} = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <GenerateValidatorsParams
                authorizedNetworks={authorizedNetworks}
                pageView={pageView}
                onGenerateValidatorDepositsAndZip={onGenerateValidatorDepositsAndZip}
                buttonLabelVd={buttonLabelVd}
                buttonDisabledVd={buttonDisabledVd}
                statusMessageVd={statusMessageVd}
            />
            {pageView ?
                <React.Fragment>
                    <GenValidatorDepositsCreationActionsCard
                        onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                        buttonLabelVd={buttonLabelVd}
                        buttonDisabledVd={buttonDisabledVd}
                        statusMessageVd={statusMessageVd}
                    />
                    <GenerateZipValidatorActionsCard
                        onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
                        zipGenButtonLabel={zipGenButtonLabel}
                        zipGenButtonEnabled={zipGenButtonEnabled}
                        zipGenStatus={zipGenStatus}
                        requestStatusZipGen={requestStatusZipGen}
                    />
                </React.Fragment> :
                <div></div>
                }
        </Stack>
    );
}

export function GenerateValidatorsParams(props: any) {
    const {authorizedNetworks, pageView, onGenerateValidatorDepositsAndZip, buttonLabelVd, buttonDisabledVd, statusMessageVd, requestStatusVd } = props;
    const awsValidatorSecretName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const awsAgeEncryptionKeyName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);

    return (
        <div>
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                    </Box>
                    <Box mt={2}>
                        <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName}/>
                    </Box>
                    <Box mt={2}>
                        <Network authorizedNetworks={authorizedNetworks}/>
                    </Box>
                    <Box mt={2}>
                        <ValidatorCount />
                    </Box>
                    <Box mt={2}>
                        <Typography variant="body2" color="text.secondary">
                            The Validator HD Offset and Withdrawal Credentials fields are optional. If you don't set the withdrawal credentials, it will generate one for you from your mnemonic.
                            The offset is used to set the offset validator keys from your mnemonic. Eg. if you set the offset to 1, it will generate the
                            validator keys starting from the second validator key from your mnemonic onwards.
                        </Typography>
                    </Box>
                    <Box mt={2}>
                        <ValidatorOffsetHD />
                    </Box>
                    <Box mt={2}>
                        <WithdrawalCredentials />
                    </Box>
                </Container>
            </div>
            {pageView ? <div></div>
                : (
                    <div>
                        <CardActions sx={{justifyContent: 'center'}}>
                            <Button onClick={onGenerateValidatorDepositsAndZip} size="small"
                                    disabled={buttonDisabledVd}>
                                {buttonLabelVd}
                            </Button>
                        </CardActions>
                        {statusMessageVd && (
                            <Typography variant="body2" color={requestStatusVd === 'error' ? 'error' : 'success'}>
                                {statusMessageVd}
                            </Typography>
                        )}
                    </div>)
            }
        </Card>
        </div>
    );
}

export function GenerateZipValidatorActionsCard(props: any) {
    const { activeStep, onGenerateValidatorEncryptedKeystoresZip, zipGenButtonLabel, zipGenButtonEnabled, zipGenStatus, requestStatusZipGen } = props;
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
                <Button onClick={onGenerateValidatorEncryptedKeystoresZip} size="small" disabled={zipGenButtonEnabled}>{zipGenButtonLabel}</Button>
            </CardActions>
            {zipGenStatus && (
                <Typography variant="body2" color={requestStatusZipGen === 'error' ? 'error' : 'success'}>
                    {zipGenStatus}
                </Typography>
            )}
        </Card>
    );
}

export function GenValidatorDepositsCreationActionsCard(props: any) {
    const { activeStep, onGenerateValidatorDeposits, buttonLabelVd, buttonDisabledVd, statusMessageVd, requestStatusVd } = props;
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
                <Button onClick={onGenerateValidatorDeposits} size="small" disabled={buttonDisabledVd}>{buttonLabelVd}</Button>
            </CardActions>
            {statusMessageVd && (
                <Typography variant="body2" color={requestStatusVd === 'error' ? 'error' : 'success'}>
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

export function WithdrawalCredentials() {
    const dispatch = useDispatch();
    const withdrawalCredentials = useSelector((state: RootState) => state.validatorSecrets.withdrawalCredentials);
    const onWithdrawalCredentialsChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        try {
            const wc = event.target.value;
            dispatch(setWithdrawalCredentials(wc));
        } catch (error) {
            console.log("error", error);
        }};
    return (
        <TextField
            fullWidth
            id="WithdrawalCredentials"
            label="WithdrawalCredentials"
            variant="outlined"
            value={withdrawalCredentials}
            onChange={onWithdrawalCredentialsChange}
            sx={{ width: '100%' }}
        />
    );
}
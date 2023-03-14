import * as React from "react";
import {useState} from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {setHdWalletPw, setMnemonic} from "../../redux/validators/ethereum.validators.reducer";
import CryptoJS from 'crypto-js';
import {
    LambdaFunctionGenEncZipFileCreation,
    LambdaFunctionGenValidatorDepositsCreation,
    LambdaFunctionSecretsCreation
} from "./AwsLambdaCreation";
import Button from "@mui/material/Button";
import {awsLambdaApiGateway} from "../../gateway/aws.lambda";
import Typography from "@mui/material/Typography";

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
    const { activeStep, onGenerate, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <LambdaFunctionSecretsCreation />
            <LambdaFunctionGenEncZipFileCreation />
            <LambdaFunctionGenValidatorDepositsCreation />
        </Stack>
    );
}

export function CreateAwsSecretsActionAreaCardWrapper(props: any) {
    const { activeStep, onGenerate, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <CreateAwsSecretNamesAreaCard />
        </Stack>
    );
}

export function CreateAwsSecretNamesAreaCard() {
    const [awsAgeEncryptionKeyName, setAwsAgeEncryptionKeyName] = useState('ageEncryptionKeyEphemery');
    const [awsValidatorSecretName, setAwsValidatorSecretName] = useState('mnemonicAndHDWalletEphemery');
    const sgLambdaURL = useSelector((state: RootState) => state.awsCredentials.secretGenLambdaFnUrl);
    const ak = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const sk = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const dispatch = useDispatch();
    const onCreateNewValidatorSecrets = async () => {
        try {
            const response = await awsLambdaApiGateway.invokeValidatorSecretsGeneration(sgLambdaURL, ak, sk, awsValidatorSecretName, awsAgeEncryptionKeyName);
            console.log(response.data)
        } catch (error) {
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
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                    <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName}/>
                </Container>
                <CardActions>
                    <Button onClick={onCreateNewValidatorSecrets} size="small">Create</Button>
                </CardActions>
                </Stack>
        </Card>
    );
}

export function ValidatorSecretName(props: any) {
    const { validatorSecretName, onValidatorSecretNameNameChange } = props;
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
        dispatch(setMnemonic(newMnemonicValue));
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
        dispatch(setHdWalletPw(newHdWalletPw));
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
    const { awsAgeEncryptionKeyName, onAccessAwsAgeEncryptionKeyName} = props;
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

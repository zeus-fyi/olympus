import * as React from "react";
import {useState} from "react";
import {Card, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {setHdWalletPw, setMnemonic} from "../../redux/validators/ethereum.validators.reducer";
import CryptoJS from 'crypto-js';
import {LambdaFunctionSecretsCreation} from "./AwsLambdaCreation";

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

export function CreateAwsSecretsActionAreaCardWrapper(props: any) {
    const { activeStep, onGenerate, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard activeStep={activeStep} onGenerate={onGenerate} onGenerateValidatorDeposits={onGenerateValidatorDeposits} onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}/>
            <LambdaFunctionSecretsCreation />
            <CreateAwsSecretNamesAreaCard />
        </Stack>
    );
}

export function CreateAwsSecretNamesAreaCard() {
    const [awsAgeEncryptionKeyName, setAwsAgeEncryptionKeyName] = useState('ageEncryptionKeyEphemery');
    const [awsValidatorSecretName, setAwsValidatorSecretName] = useState('mnemonicAndHDWalletEphemery');
    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName}/>
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                </Container>
            </div>
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

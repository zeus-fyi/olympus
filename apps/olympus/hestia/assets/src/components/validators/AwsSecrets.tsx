import * as React from "react";
import {useState} from "react";
import {Card, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {ethers} from "ethers";
import {setHdWalletPw, setMnemonic} from "../../redux/validators/ethereum.validators.reducer";
import {awsApiGateway} from "../../gateway/aws";
import {setAgePrivKey, setAgePubKey} from "../../redux/aws_wizard/aws.wizard.reducer";
import CryptoJS from 'crypto-js';

export const charsets = {
    NUMBERS: '0123456789',
    LOWERCASE: 'abcdefghijklmnopqrstuvwxyz',
    UPPERCASE: 'ABCDEFGHIJKLMNOPQRSTUVWXYZ',
    SYMBOLS: '!#$%&()*+,-./<=>?@[]^`{|}~',
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

export function CreateAwsSecretsActionAreaCardWrapper() {
    const dispatch = useDispatch();

    const onGenerate = async () => {
        try {
            const response = await awsApiGateway.getGeneratedAgeKey();
            const ageKeyGenData: any = response.data;
            dispatch(setAgePrivKey(ageKeyGenData.agePrivateKey));
            dispatch(setAgePubKey(ageKeyGenData.agePublicKey));
            const entropyBytes = ethers.randomBytes(32); // 16 bytes = 128 bits of entropy
            let phrase = ethers.Mnemonic.fromEntropy(entropyBytes).phrase;
            dispatch(setMnemonic(phrase));
            const password = generatePassword(20, charsets.NUMBERS + charsets.LOWERCASE + charsets.UPPERCASE + charsets.SYMBOLS);
            dispatch(setHdWalletPw(password));

        } catch (error) {
            console.log("error", error);
        }};

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard onGenerate={onGenerate}/>
            <CreateAwsSecretsValidatorSecretsActionAreaCard />
            <CreateAwsSecretsAgeEncryptionActionAreaCard />
        </Stack>
    );
}

export function CreateAwsSecretsAgeEncryptionActionAreaCard() {
    const [awsAgeEncryptionKeyName, setAwsAgeEncryptionKeyName] = useState('ageEncryptionKeyEphemery');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName}/>
                    <AgeCredentialsPublicKey />
                    <AgeCredentialsPrivateKey />
                </Container>
            </div>
        </Card>
    );
}

export function CreateAwsSecretsValidatorSecretsActionAreaCard(props: any) {
    const [awsValidatorSecretName, setAwsValidatorSecretName] = useState('mnemonicAndHDWalletEphemery');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                    <HDWalletPassword />
                    <Mnemonic />
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


export function AgeCredentialsPublicKey(props: any) {
    const dispatch = useDispatch();
    const onAgePubKeyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newPubKeyValue = event.target.value;
        dispatch(setAgePubKey(newPubKeyValue));
    };
    const agePubKey = useSelector((state: RootState) => state.awsCredentials.agePubKey);
    return (
        <TextField
            fullWidth
            id="AgePubKey"
            label="Age Encryption Public Key"
            variant="outlined"
            value={agePubKey}
            onChange={onAgePubKeyChange}
            sx={{ width: '100%' }}
        />
    );
}

export function AgeCredentialsPrivateKey(props: any) {
    const dispatch = useDispatch();

    const onAgePrivKeyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newPrivKeyValue = event.target.value;
        dispatch(setAgePrivKey(newPrivKeyValue));
    };
    const agePrivKey = useSelector((state: RootState) => state.awsCredentials.agePrivKey);
    return (
        <TextField
            fullWidth
            id="AgePrivKey"
            label="Age Encryption Secret Key"
            variant="outlined"
            value={agePrivKey}
            onChange={onAgePrivKeyChange}
            sx={{ width: '100%' }}
        />
    );
}

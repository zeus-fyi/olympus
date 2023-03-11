import * as React from "react";
import {useState} from "react";
import {Card, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import TextField from "@mui/material/TextField";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {ethers} from "ethers";
import {setMnemonic} from "../../redux/validators/ethereum.validators.reducer";

export function CreateAwsSecretsActionAreaCardWrapper() {
    const dispatch = useDispatch();

    const onGenerate = async () => {
        try {
            const entropyBytes = ethers.randomBytes(32); // 16 bytes = 128 bits of entropy
            let phrase = ethers.Mnemonic.fromEntropy(entropyBytes).phrase;
            console.log("mnemonic: ", phrase)
            dispatch(setMnemonic(phrase));
            return ethers.Mnemonic.fromEntropy(entropyBytes).phrase;
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
    const [agePubKey, setAgePubKey] = useState('');
    const [agePrivKey, setAgePrivKey] = useState('');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <AgeEncryptionKeySecretName awsAgeEncryptionKeyName={awsAgeEncryptionKeyName}/>
                    <AgeCredentialsPublicKey agePubKey={agePubKey}/>
                    <AgeCredentialsPrivateKey agePrivKey={agePrivKey} />
                </Container>
            </div>
        </Card>
    );
}

export function CreateAwsSecretsValidatorSecretsActionAreaCard(props: any) {
    const [awsValidatorSecretName, setAwsValidatorSecretName] = useState('mnemonicAndHDWalletEphemery');
    const [hdWalletPw, setHDWalletPw] = useState('');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                    <HDWalletPassword hdWalletPw={hdWalletPw}/>
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

export function HDWalletPassword(props: any) {
    const { hdWalletPw, onHDWalletPwChange } = props;
    return (
        <TextField
            fullWidth
            id="hdWalletPassword"
            label="HD Wallet Password"
            variant="outlined"
            value={hdWalletPw}
            onChange={onHDWalletPwChange}
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
    const { awsAgeSecretName, onAwsAgeSecretName } = props;
    return (
        <TextField
            fullWidth
            id="AgePubKey"
            label="Age Encryption Public Key"
            variant="outlined"
            value={awsAgeSecretName}
            onChange={onAwsAgeSecretName}
            sx={{ width: '100%' }}
        />
    );
}

export function AgeCredentialsPrivateKey(props: any) {
    const { agePrivKey, onAgePrivKeyChange } = props;
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

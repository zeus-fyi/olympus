import * as React from "react";
import {useState} from "react";
import {Card, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import TextField from "@mui/material/TextField";

export function CreateAwsSecretsActionAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard />
            <CreateAwsSecretsValidatorSecretsActionAreaCard />
            <CreateAwsSecretsAgeEncryptionActionAreaCard />
        </Stack>
    );
}

export function CreateAwsSecretsAgeEncryptionActionAreaCard() {
    const [awsAgeEncryptionKeyName, setAwsAgeEncryptionKeyName] = useState('ageEncryptionKeyEphemery');
    const [agePubKey, setAgePubKey] = useState('');
    const [agePrivKey, setAgePrivKey] = useState('');

    // const onAwsAgeEncryptionKeyNameChange = (event: { target: { value: React.SetStateAction<string>; }; }) => {
    //     setAwsAgeEncryptionKeyName(event.target.value);
    // };

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

export function CreateAwsSecretsValidatorSecretsActionAreaCard() {
    const [awsValidatorSecretName, setAwsValidatorSecretName] = useState('mnemonicAndHDWalletEphemery');
    const [mnemonic, setMnemonic] = useState('');
    const [hdWalletPw, setHDWalletPw] = useState('');

    const handleAccessKeyChange = (event: { target: { value: React.SetStateAction<string>; }; }) => {
        setMnemonic(event.target.value);
    };
    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretName}/>
                    <HDWalletPassword hdWalletPw={hdWalletPw}/>
                    <Mnemonic mnemonic={mnemonic}/>
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

export function Mnemonic(props: any) {
    const { accessKey, onAccessKeyChange } = props;
    return (
        <TextField
            fullWidth
            id="mnemonic"
            label="24 Word Mnemonic"
            variant="outlined"
            value={accessKey}
            onChange={onAccessKeyChange}
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
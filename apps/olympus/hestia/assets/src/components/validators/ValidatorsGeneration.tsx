import * as React from "react";
import {useState} from "react";
import {Card, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import TextField from "@mui/material/TextField";

export function GenerateValidatorKeysAndDepositsAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard />
            <GenerateValidatorKeysAndDeposits />
            <GenerateValidatorsParams />
        </Stack>
    );
}

export function GenerateValidatorKeysAndDeposits() {
    const [awsValidatorSecretName, setAwsValidatorSecretName] = useState('');
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
    const { accessKey, onAccessKeyChange } = props;
    return (
        <TextField
            fullWidth
            id="validatorSecretName"
            label="AWS Validator Key Secret Name"
            variant="outlined"
            value={accessKey}
            onChange={onAccessKeyChange}
            sx={{ width: '100%' }}
        />
    );
}

export function AgeEncryptionKeySecretName(props: any) {
    const { awsValidatorSecretName, onAccessValidatorSecretNameChange } = props;
    return (
        <TextField
            fullWidth
            id="ageEncryptionKeySecretName"
            label="AWS Age Encryption Key Secret Name"
            variant="outlined"
            value={awsValidatorSecretName}
            onChange={onAccessValidatorSecretNameChange}
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

export function GenerateValidatorsParams() {
    const [awsValidatorSecretName, setAwsValidatorSecretName] = useState('');
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
                    <ValidatorsNetwork validatorSecretName={awsValidatorSecretName}/>
                    <ValidatorCount hdWalletPw={hdWalletPw}/>
                    <ValidatorOffsetHD mnemonic={mnemonic}/>
                </Container>
            </div>
        </Card>

    );
}

export function ValidatorsNetwork(props: any) {
    const { accessKey, onAccessKeyChange } = props;
    return (
        <TextField
            fullWidth
            id="validatorsNetwork"
            label="Network Name"
            variant="outlined"
            value={accessKey}
            onChange={onAccessKeyChange}
            sx={{ width: '100%' }}
        />
    );
}

export function ValidatorCount(props: any) {
    const { accessKey, onAccessKeyChange } = props;
    return (
        <TextField
            fullWidth
            id="validatorCount"
            label="Validator Count"
            variant="outlined"
            value={accessKey}
            onChange={onAccessKeyChange}
            sx={{ width: '100%' }}
        />
    );
}

export function ValidatorOffsetHD(props: any) {
    const { accessKey, onAccessKeyChange } = props;
    return (
        <TextField
            fullWidth
            id="validatorOffsetHD"
            label="Validator HD Offset"
            variant="outlined"
            value={accessKey}
            onChange={onAccessKeyChange}
            sx={{ width: '100%' }}
        />
    );
}
import * as React from "react";
import {useState} from "react";
import {Card, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import TextField from "@mui/material/TextField";
import {ValidatorSecretName} from "./AwsSecrets";

export function GenerateValidatorKeysAndDepositsAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard />
            <GenerateValidatorsParams />
        </Stack>
    );
}

export function GenerateValidatorKeysAndDeposits() {
    const [awsValidatorSecretNameDeposits, awsValidatorSecretNameDepositsName] = useState('mnemonicAndHDWalletEphemery');
    const [mnemonic, setMnemonic] = useState('');
    const [hdWalletPw, setHDWalletPw] = useState('');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorSecretName validatorSecretName={awsValidatorSecretNameDeposits}/>
                    <HDWalletPassword hdWalletPw={hdWalletPw}/>
                    <Mnemonic mnemonic={mnemonic}/>
                </Container>
            </div>
        </Card>

    );
}

export function Mnemonic(props: any) {
    const { accessKey, onAccessKeyChange } = props;

    const onAccessMnemonicSecretChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        onAccessKeyChange(event.target.value);
    };

    return (
        <TextField
            fullWidth
            id="mnemonic"
            label="24 Word Mnemonic"
            variant="outlined"
            value={accessKey}
            onChange={onAccessMnemonicSecretChange}
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
    const [awsValidatorsNetwork, setAwsValidatorsNetwork] = useState('Ephemery');
    const [validatorCount, onValidatorCountChange ] = useState('1');
    const [offset, setOffset] = useState('0');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorsNetwork awsValidatorsNetwork={awsValidatorsNetwork}/>
                    <ValidatorCount validatorCount={validatorCount}/>
                    <ValidatorOffsetHD offset={offset}/>
                </Container>
            </div>
        </Card>

    );
}

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
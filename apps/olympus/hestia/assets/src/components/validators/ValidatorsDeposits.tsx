import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import * as React from "react";
import {useState} from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import {Network} from "./ZeusServiceRequest";

export function ValidatorsDepositRequestAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ValidatorsDepositRequestAreaCard />
        </Stack>
    );
}

export function ValidatorsDepositRequestAreaCard() {
    const [mnemonic, setMnemonic] = useState('');
    const [hdWalletPw, setHDWalletPw] = useState('');
    const [agePubKey, setAgePubKey] = useState('');
    const [agePrivKey, setAgePrivKey] = useState('');

    const handleAccessKeyChange = (event: { target: { value: React.SetStateAction<string>; }; }) => {
        setAgePubKey(event.target.value);
    };
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <ValidatorsDepositsSubmitWrapper />
            </Container >
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <SubmitValidators />
            </Container >
        </div>

    );
}

export function SubmitValidators() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Send Validator Deposits to Network
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Sends Validator Deposits to the Network
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Send</Button>
            </CardActions>
        </Card>
    );
}


export function ValidatorsDepositsSubmitWrapper() {
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
                    <Eth1WalletPrivateKey validatorSecretName={awsValidatorSecretName}/>
                    <Network mnemonic={mnemonic}/>
                </Container>
            </div>
        </Card>

    );
}

export function Eth1WalletPrivateKey(props: any) {
    const { accessKey, onAccessKeyChange } = props;
    return (
        <TextField
            fullWidth
            id="eth1WalletPrivateKey"
            label="Eth1 Wallet Private Key"
            variant="outlined"
            value={accessKey}
            onChange={onAccessKeyChange}
            sx={{ width: '100%' }}
        />
    );
}
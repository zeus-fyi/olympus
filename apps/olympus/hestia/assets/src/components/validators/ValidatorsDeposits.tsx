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
    const [network, setNetwork] = useState('Ephemery');
    const [eth1Pk, setEth1Pk] = useState('');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Network network={network}/>
                    <Eth1WalletPrivateKey eth1Pk={eth1Pk}/>
                </Container>
            </div>
        </Card>

    );
}

export function Eth1WalletPrivateKey(props: any) {
    const { eth1Pk, onAccessEth1PkChange } = props;
    return (
        <TextField
            fullWidth
            id="eth1WalletPrivateKey"
            label="Eth1 Wallet Private Key"
            variant="outlined"
            value={eth1Pk}
            onChange={onAccessEth1PkChange}
            sx={{ width: '100%' }}
        />
    );
}
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import * as React from "react";
import {useState} from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";

export function ZeusServiceRequestAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ZeusServiceRequestAreaCard />
        </Stack>
    );
}

export function ZeusServiceRequestAreaCard() {
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
                <ZeusServiceRequestWrapper />
            </Container >
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <ZeusServiceRequest />
            </Container >
        </div>

    );
}

export function ZeusServiceRequest() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Create Zeus Validators Service Request
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates Zeus Validators Service Request
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Send</Button>
            </CardActions>
        </Card>
    );
}

export function ZeusServiceRequestWrapper() {
    const [keyGroupName, setKeyGroupName] = useState('EphemeryDemoGroup');
    const [feeRecipient, setMnemonic] = useState('');
    const [network, setNetwork] = useState('Ephemery');

    const handleAccessKeyChange = (event: { target: { value: React.SetStateAction<string>; }; }) => {
        setMnemonic(event.target.value);
    };
    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <KeyGroupName keyGroupName={keyGroupName}/>
                    <Network network={network}/>
                    <FeeRecipient feeRecipient={feeRecipient}/>
                </Container>
            </div>
        </Card>

    );
}

export function KeyGroupName(props: any) {
    const { keyGroupName, onAccessKeyGroupName } = props;
    return (
        <TextField
            fullWidth
            id="keyGroupName"
            label="Key Group Name"
            variant="outlined"
            value={keyGroupName}
            onChange={onAccessKeyGroupName}
            sx={{ width: '100%' }}
        />
    );
}

export function FeeRecipient(props: any) {
    const { feeRecipient, onAccessFeeRecipientChange } = props;
    return (
        <TextField
            fullWidth
            id="feeRecipient"
            label="Fee Recipient"
            variant="outlined"
            value={feeRecipient}
            onChange={onAccessFeeRecipientChange}
            sx={{ width: '100%' }}
        />
    );
}

export function Network(props: any) {
    const { network, onAccessSetNetwork } = props;
    return (
        <TextField
            fullWidth
            id="network"
            label="Network"
            variant="outlined"
            value={network}
            onChange={onAccessSetNetwork}
            sx={{ width: '100%' }}
        />
    );
}
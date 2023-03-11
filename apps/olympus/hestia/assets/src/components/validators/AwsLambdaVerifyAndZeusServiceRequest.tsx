import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import * as React from "react";
import {useState} from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {ValidatorsUploadActionAreaCard} from "./ValidatorsUpload";
import {ZeusServiceRequest} from "./ZeusServiceRequest";

export function LambdaVerifyAndZeusServiceRequest() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ValidatorsUploadActionAreaCard />,
            <AwsLambdaFunctionVerifyAreaCard />
        </Stack>
    );
}

export function AwsLambdaFunctionVerifyAreaCard() {
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
                <LambdaVerifyCard />
            </Container >
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <ZeusServiceRequest />
            </Container >
        </div>

    );
}

export function LambdaVerifyCard() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Verify Lambda Key Signing
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Sends random hex string payloads to your AWS lambda function and verifies the returned signatures match the public keys.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Send Request</Button>
            </CardActions>
        </Card>
    );
}


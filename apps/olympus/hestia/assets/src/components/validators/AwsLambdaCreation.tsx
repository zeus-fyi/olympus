import * as React from "react";
import {useState} from "react";
import {Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";

export function CreateAwsLambdaFunctionActionAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard />
            <CreateAwsLambdaFunctionAreaCard />
        </Stack>
    );
}

export function CreateAwsLambdaFunctionAreaCard() {
    const [mnemonic, setMnemonic] = useState('');
    const [hdWalletPw, setHDWalletPw] = useState('');
    const [agePubKey, setAgePubKey] = useState('');
    const [agePrivKey, setAgePrivKey] = useState('');

    const handleAccessKeyChange = (event: { target: { value: React.SetStateAction<string>; }; }) => {
        setAgePubKey(event.target.value);
    };
    return (
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                </Container >
            </div>

    );
}

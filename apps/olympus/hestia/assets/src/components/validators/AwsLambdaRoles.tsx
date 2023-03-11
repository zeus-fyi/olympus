import * as React from "react";
import {useState} from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";

export function CreateInternalAwsLambdaUserRolesActionAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard />
            <CreateAwsLambdaUserRolesActionAreaCard />
        </Stack>
    );
}

export function CreateAwsLambdaUserRolesActionAreaCard() {
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
                    <InternalLambdaUserRolePolicySetup />
                </Container >
            </div>

    );
}
export function InternalLambdaUserRolePolicySetup() {
    return (
        <Card sx={{ maxWidth: 345 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Internal
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Internal Lambda User & RolePolicy Setup
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Share</Button>
                <Button size="small">Learn More</Button>
            </CardActions>
        </Card>
    );
}

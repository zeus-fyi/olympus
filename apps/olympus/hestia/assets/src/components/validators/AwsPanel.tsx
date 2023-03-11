import {Card, CardActionArea, CardMedia, Container, Stack} from "@mui/material";
import * as React from "react";
import {useState} from "react";
import TextField from "@mui/material/TextField";

export function AwsUploadActionAreaCard() {
    const [accessKey, setAccessKey] = useState('');
    const [secretKey, setSecretKey] = useState('');

    const handleAccessKeyChange = (event: { target: { value: React.SetStateAction<string>; }; }) => {
        setAccessKey(event.target.value);
    };
    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
            <CardActionArea>
                <CardMedia
                    component="img"
                    height="230"
                    image={require("../../static/aws.jpg")}
                    alt="aws"
                />
            </CardActionArea>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <AwsCredentialsAccessKey accessKey={accessKey}/>
                    <AwsCredentialsSecret secretKey={secretKey} />
                </Container>
            </div>
        </Card>

);
}
export function AwsCredentialsAccessKey(props: any) {
    const { accessKey, onAccessKeyChange } = props;
    return (
            <TextField
                fullWidth
                id="AccessKey"
                label="Access Key"
                variant="outlined"
                value={accessKey}
                onChange={onAccessKeyChange}
                sx={{ width: '100%' }}
            />
    );
}

export function AwsCredentialsSecret(props: any) {
    const { secretKey, onSecretKeyChange } = props;
    return (
            <TextField
                fullWidth
                id="SecretKey"
                label="Secret Key"
                variant="outlined"
                value={secretKey}
                onChange={onSecretKeyChange}
                sx={{ width: '100%' }}

            />
    );
}
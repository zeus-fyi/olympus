import {Card, CardActionArea, CardMedia} from "@mui/material";
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
        <Card sx={{ maxWidth: 320 }}>
            <CardActionArea>
                <CardMedia
                    component="img"
                    height="230"
                    image={require("../../static/aws.jpg")}
                    alt="aws"
                />
            </CardActionArea>
            <AwsCredentialsAccessKey accessKey={accessKey}/>
            <AwsCredentialsSecret secretKey={secretKey} />
        </Card>
    );
}
export function AwsCredentialsAccessKey(props: any) {
    const { accessKey, onAccessKeyChange } = props;
    return (
            <TextField
                fullWidth
                id="AccessKey"
                label="AccessKey"
                variant="outlined"
                value={accessKey}
                onChange={onAccessKeyChange}
            />
    );
}

export function AwsCredentialsSecret(props: any) {
    const { secretKey, onSecretKeyChange } = props;
    return (
            <TextField
                fullWidth
                id="SecretKey"
                label="SecretKey"
                variant="outlined"
                value={secretKey}
                onChange={onSecretKeyChange}
            />
    );
}
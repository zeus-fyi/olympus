import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";

export function CreateAwsLambdaFunctionActionAreaCardWrapper(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard activeStep={activeStep}/>
            <CreateAwsLambdaFunctionAreaCard />
        </Stack>
    );
}

export function CreateAwsLambdaFunctionAreaCard() {
    return (
            <div style={{ display: 'flex' }}>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <LambdaFunctionKeystoresLayerCreation />
                </Container >
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <LambdaFunctionCreation />
                </Container >
            </div>

    );
}

export function LambdaFunctionKeystoresLayerCreation() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Keystores Layer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates your encrypted keystores layer for usage in your AWS lambda function using your generated encrypted zip keystores file.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}


export function LambdaFunctionCreation() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Signer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a BLS signer lambda function in AWS that decrypts your keystores with your Age Encryption key to sign messages.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionSecretsCreation() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Creation for Trustless Secrets Generation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a lambda function in AWS that securely generates a mnemonic, hdWalletPassword, and Age Encryption key and stores them in your secrets manager.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionGenValidatorDepositsCreation() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function For Secure Validator Deposits Generation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a lambda function in AWS that securely generates validator deposit messages using your mnemonic from secrets manager.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

export function LambdaFunctionGenEncZipFileCreation() {
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Encrypted Keystores Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a lambda function in AWS that generates an encrypted zip file with validator signing keys.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}


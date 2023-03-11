import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";

export function CreateAwsLambdaFunctionActionAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard />
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
                    Lambda Encrypted Keystores Layer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates your encrypted keystores layer for usage in your AWS lambda function.
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
                    Lambda Function Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a BLS signer lambda function in AWS.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}

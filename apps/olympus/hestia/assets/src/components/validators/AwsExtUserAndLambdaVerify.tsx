import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import * as React from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {ValidatorsUploadActionAreaCard} from "./ValidatorsUpload";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {awsApiGateway} from "../../gateway/aws";
import TextField from "@mui/material/TextField";
import {validatorsApiGateway} from "../../gateway/validators";

export function LambdaExtUserVerify(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ValidatorsUploadActionAreaCard />,
            <CreateAwsExternalLambdaUser />
            <AwsLambdaFunctionVerifyAreaCard />
        </Stack>
    );
}

export function CreateAwsExternalLambdaUser() {
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const externalAccessUserName = useSelector((state: RootState) => state.awsCredentials.externalAccessUserName);
    const externalAccessSecretName = useSelector((state: RootState) => state.awsCredentials.externalAccessSecretName);

    const handleCreateUser = async () => {
        try {
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            const response = await awsApiGateway.createExternalLambdaUser(creds);
            console.log("response", response);
        } catch (error) {
            console.log("error", error);
        }
        try {
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            await awsApiGateway.createOrFetchExternalLambdaUserAccessKeys(creds,externalAccessUserName, externalAccessSecretName);
        } catch (error) {
            console.log("error", error);
        }
    };

    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Setup External Access
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates an AWS Lambda User and RolePolicy for external function calls. This is the user
                    that we will use to send authorized messages to your validators.
                </Typography>
            </CardContent>
            <TextField
                margin="normal"
                required
                fullWidth
                id="externalAccessUserName"
                label="ExternalAccessUserName"
                name="externalAccessUserName"
                value={externalAccessUserName}
                autoFocus
            />
            <TextField
                margin="normal"
                required
                fullWidth
                id="externalAccessSecretName"
                label="ExternalAccessSecretName"
                name="externalAccessSecretName"
                value={externalAccessSecretName}
                autoFocus
            />
            <CardActions>
                <Button size="small" onClick={handleCreateUser}>Create User & Auth Keys</Button>
            </CardActions>
        </Card>
    );
}

export function AwsLambdaFunctionVerifyAreaCard() {
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <LambdaVerifyCard />
            </Container >
        </div>
    );
}

export function LambdaVerifyCard() {
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const depositsData = useSelector((state: RootState) => state.awsCredentials.depositData);
    const blsSignerLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.blsSignerLambdaFnUrl);
    const blsSignerFunctionName = useSelector((state: RootState) => state.awsCredentials.blsSignerFunctionName);
    const externalAccessUserName = useSelector((state: RootState) => state.awsCredentials.externalAccessUserName);
    const externalAccessSecretName = useSelector((state: RootState) => state.awsCredentials.externalAccessSecretName);

    const handleVerifySigners = async () => {
        try {
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            const r = await awsApiGateway.createOrFetchExternalLambdaUserAccessKeys(creds, externalAccessUserName, externalAccessSecretName);
            const extCreds = {accessKeyId: r.data.accessKey, secretAccessKey: r.data.secretKey};
            const response = await validatorsApiGateway.verifyValidators(extCreds,blsSignerLambdaFnUrl, depositsData);
            // TODO set dispatch/update the table
            console.log("r", response);
        } catch (error) {
            console.log("error", error);
        }};

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
            <TextField
                margin="normal"
                required
                fullWidth
                id="blsSignerLambdaFunctionName"
                label="BlsSignerLambdaFunctionName"
                name="blsSignerLambdaFunctionName"
                value={blsSignerFunctionName}
                autoFocus
            />
            <TextField
                margin="normal"
                required
                fullWidth
                id="blsSignerLambdaFnUrl"
                label="BlsSignerLambdaFnUrl"
                name="blsSignerLambdaFnUrl"
                value={blsSignerLambdaFnUrl}
                autoFocus
            />
            <CardActions>
                <Button size="small" onClick={handleVerifySigners}>Send Request</Button>
            </CardActions>
        </Card>
    );
}

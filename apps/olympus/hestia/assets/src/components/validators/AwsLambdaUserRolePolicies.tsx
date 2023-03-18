import * as React from "react";
import {useState} from "react";
import {Card, CardActions, CardContent, CircularProgress, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {awsApiGateway} from "../../gateway/aws";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import TextField from "@mui/material/TextField";

export function CreateInternalAwsLambdaUserRolesActionAreaCardWrapper(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard activeStep={activeStep}/>
            <CreateAwsLambdaUserRolesActionAreaCard />
        </Stack>
    );
}

export function CreateAwsLambdaUserRolesActionAreaCard() {
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
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);

    const handleCreateUser = async () => {
        try {
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            if (!accessKey || !secretKey) {
                setRequestStatus('errorAuth');
                return;
            }
            setRequestStatus('pending');

            const response = await awsApiGateway.createInternalLambdaUser(creds);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
            }
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};

    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const [requestStatus, setRequestStatus] = useState('');

    switch (requestStatus) {
        case 'pending':
            buttonLabel = <CircularProgress size={20} />;
            buttonDisabled = true;
            break;
        case 'success':
            buttonLabel = 'Created successfully';
            buttonDisabled = true;
            statusMessage = 'User and role policy created successfully!';
            break;
        case 'error':
            buttonLabel = 'Error creating user';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating the user.';
            break;
        case 'errorAuth':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'Update your AWS credentials on step 1 and try again.';
            break;
        default:
            buttonLabel = 'Create';
            buttonDisabled = false;
            break;
    }

    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Internal User & RolePolicy Setup
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a new user and role policy for your own internal usage, e.g. for running, testing, development, etc.
                </Typography>
            </CardContent>
            <TextField
                margin="normal"
                required
                fullWidth
                id="internalLambdaUser"
                label="internalLambdaUser"
                name="internalLambdaUser"
                value={"internalLambdaUser"}
                autoFocus
            />
            <CardActions>
                <Button size="small" onClick={handleCreateUser} disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}
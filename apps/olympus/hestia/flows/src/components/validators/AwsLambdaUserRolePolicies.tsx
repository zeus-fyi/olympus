import * as React from "react";
import {useState} from "react";
import {Box, Card, CardActions, CardContent, CircularProgress, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {awsApiGateway} from "../../gateway/aws";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import TextField from "@mui/material/TextField";
import {setEncKeystoresZipLambdaFnUrl, setSecretGenLambdaFnUrl} from "../../redux/aws_wizard/aws.wizard.reducer";

export function CreateInternalAwsLambdaUserRolesActionAreaCardWrapper(props: any) {
    const { activeStep, pageView } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard activeStep={activeStep}/>
            <CreateAwsLambdaUserRolesActionAreaCard pageView={pageView}/>
        </Stack>
    );
}

export function CreateAwsLambdaUserRolesActionAreaCard(props: any) {
    const {pageView} = props;
    return (
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <InternalLambdaUserRolePolicySetup pageView={pageView}/>
                </Container >
            </div>
    );
}

export function InternalLambdaUserRolePolicySetup(props: any) {
    const {pageView} = props;
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const dispatch = useDispatch();
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
                return
            }
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};

    const handleSetupAllInternal = async () => {
        try {
            const creds = {accessKeyId: accessKey, secretAccessKey: secretKey};
            if (!accessKey || !secretKey) {
                setRequestStatus('errorAuth');
                return;
            }
            setRequestStatus('pending');

            const response = await awsApiGateway.createInternalLambdaUser(creds);
            if (response.status === 200) {
            } else {
                setRequestStatus('error');
                return
            }
            const secretsLambdaResponse = await awsApiGateway.createValidatorSecretsLambda(creds);
            if (secretsLambdaResponse.status === 200) {
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setSecretGenLambdaFnUrl(secretsLambdaResponse.data));
            const responseDepositFnSetup = await awsApiGateway.createValidatorsDepositDataLambda(creds);
            if (responseDepositFnSetup.status === 200) {
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(responseDepositFnSetup(responseDepositFnSetup.data));
            const encZipFnSetup = await awsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(creds);
            if (encZipFnSetup.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setEncKeystoresZipLambdaFnUrl(encZipFnSetup.data));
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }
    };

    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const [requestStatus, setRequestStatus] = useState('');
    const buttonOnClick = pageView ? handleSetupAllInternal : handleCreateUser;

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
        <Card sx={{ maxWidth: 700 }}>
            {pageView ? (
                <div>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Internal User & RolePolicy Setup
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Creates a new user and role policy for your own internal usage, e.g. for running, testing, development, etc.
                        </Typography>
                    </CardContent>
                    <Box ml={2} mr={2}>
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
                    </Box>
                </div>
                ) :
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                       Setup Internal Management Components
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        One Time Setup for Internal User, Roles, Policy, and Functions
                    </Typography>
                </CardContent>
            }
            <CardActions>
                <Button size="small" onClick={buttonOnClick} disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {awsApiGateway} from "../../gateway/aws";
import {
    setBlsSignerLambdaFnUrl,
    setKeystoreLayerName,
    setKeystoreLayerNumber,
    setSignerFunctionName
} from "../../redux/aws_wizard/aws.wizard.reducer";
import {
    Box,
    Card,
    CardActionArea,
    CardActions,
    CardContent,
    CardMedia,
    CircularProgress,
    Container,
    Stack
} from "@mui/material";
import Typography from "@mui/material/Typography";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import * as React from "react";
import {useState} from "react";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";

export function CreateAwsLambdaFunctionActionAreaCardWrapper(props: any) {
    const { activeStep, onEncZipFileUpload, zipBlob, pageView, onHandleVerifySigners} = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            {pageView ? (
                <React.Fragment>
                    <EncryptedKeystoresZipUploadActionAreaCard onEncZipFileUpload={onEncZipFileUpload} />
                    <CreateAwsLambdaFunctionAreaCard zipBlob={zipBlob} />
                </React.Fragment>
                ) : (
                <LambdaFunctionCreationBundled zipBlob={zipBlob} onHandleVerifySigners={onHandleVerifySigners}/>
            )}
        </Stack>
    );
}

export function CreateAwsLambdaFunctionAreaCard(props: any) {
    const { zipBlob } = props;
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <LambdaFunctionKeystoresLayerCreation zipBlob={zipBlob}/>
            </Container >
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <LambdaFunctionCreation />
            </Container >
        </div>
    );
}


export function LambdaFunctionCreationBundled(props: any) {
    const { zipBlob, onHandleVerifySigners} = props;
    const acKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const seKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const signerName = useSelector((state: RootState) => state.awsCredentials.blsSignerFunctionName);
    const signerLayerName = useSelector((state: RootState) => state.awsCredentials.blsSignerKeystoresLayerName);

    const dispatch = useDispatch();

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
            buttonLabel = 'Successfully updated, or created function';
            buttonDisabled = true;
            statusMessage = 'Lambda function created successfully!';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating or updating the lambda function.';
            break;
        case 'errorAuth':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating or updating the lambda function.';
            break;
        default:
            buttonLabel = 'Create | Update';
            buttonDisabled = false;
            break;
    }

    const onCreateLambdaSignerFn = async () => {
        try {
            if (zipBlob.size === 0) {
                console.log("ZipBlob is empty. Skipping creation of Lambda Keystores Layer.");
                setRequestStatus('errorNoZip');
                return;
            }
            const creds = {accessKeyId: acKey, secretAccessKey: seKey};
            if (!creds.accessKeyId || !creds.secretAccessKey) {
                setRequestStatus('errorAuth');
                return;
            }
            setRequestStatus('pending')
            const responseKeystores = await awsApiGateway.createLambdaFunctionKeystoresLayer(creds, signerLayerName, zipBlob);
            if (responseKeystores.status === 200) {
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setKeystoreLayerNumber(responseKeystores.data));
            const response = await awsApiGateway.createLambdaSignerFunction(creds, signerName, signerLayerName);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setBlsSignerLambdaFnUrl(response.data));
            await onHandleVerifySigners();
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};


    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Signer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a BLS signer lambda function in AWS that decrypts your keystores layer
                    with your Age Encryption key, and will sign messages for your validators.
                </Typography>
            </CardContent>
            <CardActions>
                <Button onClick={onCreateLambdaSignerFn} size="small" disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}

export function LambdaFunctionCreation() {
    const acKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const seKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const signerName = useSelector((state: RootState) => state.awsCredentials.blsSignerFunctionName);
    const signerLayerName = useSelector((state: RootState) => state.awsCredentials.blsSignerKeystoresLayerName);
    const signerUrl = useSelector((state: RootState) => state.awsCredentials.blsSignerLambdaFnUrl);

    const dispatch = useDispatch();

    const blsSignerFunctionName = useSelector((state: RootState) => state.awsCredentials.blsSignerFunctionName);
    const onBlsSignerFunctionName = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newBlsSignerFunctionName = event.target.value;
        dispatch(setSignerFunctionName(newBlsSignerFunctionName));
    };
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
            buttonLabel = 'Successfully updated, or created function';
            buttonDisabled = true;
            statusMessage = 'Lambda function created successfully!';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating or updating the lambda function.';
            break;
        case 'errorAuth':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating or updating the lambda function.';
            break;
        default:
            buttonLabel = 'Create | Update';
            buttonDisabled = false;
            break;
    }
    const onCreateLambdaSignerFn = async () => {
        try {
            const creds = {accessKeyId: acKey, secretAccessKey: seKey};
            if (!acKey || !seKey) {
                setRequestStatus('errorAuth');
                return;
            }
            setRequestStatus('pending')
            const response = await awsApiGateway.createLambdaSignerFunction(creds, signerName, signerLayerName);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setBlsSignerLambdaFnUrl(response.data));
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};

    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Function Signer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a BLS signer lambda function in AWS that decrypts your keystores layer (with the name you supplied on the left panel)
                    with your Age Encryption key, and will sign messages for your validators. You only need to share the age
                    key name reference, not the actual public or private key.
                </Typography>
            </CardContent>
            <Box ml={2} mr={2}>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="blsSignerLambdaFunctionName"
                    label="BlsSignerLambdaFunctionName"
                    name="blsSignerLambdaFunctionName"
                    value={blsSignerFunctionName}
                    onChange={onBlsSignerFunctionName}
                    autoFocus
                />
            </Box>
            <Box ml={2} mr={2}>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="blsSignerLambdaFnUrl"
                    label="BlsSignerLambdaFnUrl"
                    name="blsSignerLambdaFnUrl"
                    value={signerUrl}
                    autoFocus
                />
            </Box>
            <CardActions>
                <Button onClick={onCreateLambdaSignerFn} size="small" disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}

export function LambdaFunctionKeystoresLayerCreation(props: any) {
    const { zipBlob } = props;
    const dispatch = useDispatch();
    const acKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const seKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const signerLayerName = useSelector((state: RootState) => state.awsCredentials.blsSignerKeystoresLayerName);

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
            statusMessage = 'Keystore layer created successfully!';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while creating the keystore layer.';
            break;
        case 'errorNoZip':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'No zip file detected, please regenerate or manually upload one.';
            break;
        case 'errorAuth':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'Update your AWS credentials on step 1 and try again.';
            break;
        default:
            buttonLabel = 'Create | Update';
            buttonDisabled = false;
            break;
    }
    const onCreateLambdaKeystoresLayer = async () => {
        try {
            if (zipBlob.size === 0) {
                console.log("ZipBlob is empty. Skipping creation of Lambda Keystores Layer.");
                setRequestStatus('errorNoZip');
                return;
            }
            const creds = {accessKeyId: acKey, secretAccessKey: seKey};
            if (!creds.accessKeyId || !creds.secretAccessKey) {
                setRequestStatus('errorAuth');
                return;
            }
            setRequestStatus('pending')
            const response = await awsApiGateway.createLambdaFunctionKeystoresLayer(creds, signerLayerName, zipBlob);
            if (response.status === 200) {
                setRequestStatus('success');
            } else {
                setRequestStatus('error');
                return
            }
            dispatch(setKeystoreLayerNumber(response.data));
        } catch (error) {
            setRequestStatus('error');
            console.log("error", error);
        }};

    const blsSignerKeystoresLayerName = useSelector((state: RootState) => state.awsCredentials.blsSignerKeystoresLayerName);
    const onBlsSignerKeystoresLayerNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newBlsSignerKeystoresLayerName = event.target.value;
        dispatch(setKeystoreLayerName(newBlsSignerKeystoresLayerName));
    };
    const keystoresLayerNumber = useSelector((state: RootState) => state.awsCredentials.keystoreLayerNumber);

    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Lambda Keystores Layer Creation
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a new layer from your keystores.zip for usage in your AWS lambda signing function.
                    If you did not create your zip file in the previous step you'll need to manually upload it on the left.
                    If you create a new layer, then you can re-run the create lambda function on the right to update
                    the signer function with the new layer.
                </Typography>
            </CardContent>
            <Box ml={2} mr={2}>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="blsSignerKeystoresLayerName"
                    label="BlsSignerKeystoresLayerName"
                    name="blsSignerKeystoresLayerName"
                    value={blsSignerKeystoresLayerName}
                    onChange={onBlsSignerKeystoresLayerNameChange}
                    autoFocus
                />
            </Box>
            <Box ml={2} mr={2}>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    id="keystoresLayerNumber"
                    label="KeystoresLayerNumber"
                    name="keystoresLayerNumber"
                    type={"number"}
                    value={keystoresLayerNumber}
                    autoFocus
                />
            </Box>
            <CardActions>
                <Button size="small" onClick={onCreateLambdaKeystoresLayer} disabled={buttonDisabled}>{buttonLabel}</Button>
            </CardActions>
            {statusMessage && (
                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                    {statusMessage}
                </Typography>
            )}
        </Card>
    );
}

export function EncryptedKeystoresZipUploadActionAreaCard(props: any) {
    const { activeStep, onEncZipFileUpload } = props;

    return (
        <Card sx={{ maxWidth: 320 }}>
            <CardActionArea>
                <CardMedia
                    component="img"
                    height="230"
                    image={require("../../static/ethereum-logo.png")}
                    alt="ethereum"
                />
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', backgroundColor: '#8991B0'}}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                        Upload Keystores.zip
                    </Typography>
                    <UploadKeystoresZipButton onEncZipFileUpload={onEncZipFileUpload}/>
                </CardContent>
            </CardActionArea>
        </Card>
    );
}

export function UploadKeystoresZipButton(props: any) {
    const { activeStep, onEncZipFileUpload } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <Button variant="contained" component="label" style={{ backgroundColor: '#8991B0', color: '#151C2F' }}>
                <CloudUploadIcon />
                <input hidden accept="application/zip" type="file" onChange={onEncZipFileUpload}/>
            </Button>
        </Stack>
    );
}

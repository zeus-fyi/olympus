import {Card, CardActionArea, CardActions, CardContent, CardMedia, Container, Stack} from "@mui/material";
import * as React from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {ValidatorsUploadActionAreaCard} from "./ValidatorsUpload";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {awsApiGateway} from "../../gateway/aws";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";

export function LambdaExtUserVerify(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ValidatorsUploadActionAreaCard />,
            <AwsLambdaFunctionVerifyAreaCard />
        </Stack>
    );
}

export function EncryptedKeystoresZipUploadActionAreaCard() {
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
                    <UploadValidatorsButton />
                </CardContent>
            </CardActionArea>
        </Card>
    );
}
export function UploadValidatorsButton() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <Button variant="contained" component="label" style={{ backgroundColor: '#8991B0', color: '#151C2F' }}>
                <CloudUploadIcon />
                <input hidden accept="image/*" multiple type="file" />
            </Button>
        </Stack>
    );
}

export function AwsLambdaFunctionVerifyAreaCard() {
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <CreateAwsExternalLambdaUser />
            </Container >
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <LambdaVerifyCard />
            </Container >
        </div>
    );
}

export function LambdaVerifyCard() {
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
            <CardActions>
                <Button size="small">Send Request</Button>
            </CardActions>
        </Card>
    );
}

export function CreateAwsExternalLambdaUser() {
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);

    const handleCreateUser = async () => {
        try {
            const response = await awsApiGateway.createExternalLambdaUser(accessKey,secretKey);
            console.log("response", response);
        } catch (error) {
            console.log("error", error);
        }};

    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Create AWS External Lambda User
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates AWS External Lambda User
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small" onClick={handleCreateUser}>Create</Button>
            </CardActions>
        </Card>
    );
}

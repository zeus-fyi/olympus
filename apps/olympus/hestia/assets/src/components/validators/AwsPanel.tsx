import {Card, CardActionArea, CardMedia, Container, Stack} from "@mui/material";
import * as React from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {useDispatch, useSelector} from 'react-redux';
import {setAccessKey, setSecretKey} from '../../redux/aws_wizard/aws.wizard.reducer';
import {RootState} from "../../redux/store";

export function AwsUploadActionAreaCard(props: any) {
    const { activeStep, onGenerate, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip } = props;
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const dispatch = useDispatch();

    const onAccessKeyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newAccessKey = event.target.value;
        dispatch(setAccessKey(newAccessKey));
    };
    const onSecretKeyChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newSecretKey = event.target.value;
        dispatch(setSecretKey(newSecretKey));
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
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    {activeStep === 0 ? (
                        <div>
                        <AwsCredentialsAccessKey accessKey={accessKey} onAccessKeyChange={onAccessKeyChange}/>
                        <AwsCredentialsSecret secretKey={secretKey} onSecretKeyChange={onSecretKeyChange}/>
                        </div>
                    ): (  <div>   </div> )}
                    <AwsCredentialsButtons activeStep={activeStep} onGenerate={onGenerate} onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip} onGenerateValidatorDeposits={onGenerateValidatorDeposits}/>
                </Container>
            </div>
        </Card>
);
}

export function AwsCredentialsButtons(props: any) {
    const { activeStep, onGenerate, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip, onSave } = props;
    return (
        <Stack direction="row" spacing={2} sx={{ mt: 2 }}>
            {activeStep === 3 ? (
                <div>
                    <Button variant="contained" onClick={onGenerateValidatorDeposits}>
                        Generate Deposit Data
                    </Button>
                    <Button variant="outlined" onClick={onGenerateValidatorEncryptedKeystoresZip}>
                        Generate Keystores Zip
                    </Button>
                </div>
            ) : (
                <div>
                    <Button variant="contained" onClick={onSave}>
                        Create
                    </Button>
                    <Button variant="outlined" onClick={onGenerate}>
                        Generate
                    </Button>
                </div>
            )}
        </Stack>
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
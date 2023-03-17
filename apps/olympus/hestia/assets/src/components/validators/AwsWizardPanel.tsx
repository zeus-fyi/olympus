import * as React from 'react';
import {ChangeEvent, useState} from 'react';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepButton from '@mui/material/StepButton';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Container from "@mui/material/Container";
import {CreateAwsInternalLambdasActionAreaCardWrapper, CreateAwsSecretsActionAreaCardWrapper,} from "./AwsSecrets";
import {CreateInternalAwsLambdaUserRolesActionAreaCardWrapper} from "./AwsLambdaUserRolePolicies";
import {LambdaExtUserVerify} from "./AwsExtUserAndLambdaVerify";
import {GenerateValidatorKeysAndDepositsAreaCardWrapper} from "./ValidatorsGeneration";
import {ZeusServiceRequestAreaCardWrapper} from "./ZeusServiceRequest";
import {ValidatorsDepositRequestAreaCardWrapper} from "./ValidatorsDeposits";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {
    setDepositData,
    setDepositsGenLambdaFnUrl,
    setEncKeystoresZipLambdaFnUrl
} from "../../redux/aws_wizard/aws.wizard.reducer";
import {awsApiGateway} from "../../gateway/aws";
import {awsLambdaApiGateway} from "../../gateway/aws.lambda";
import {CreateAwsLambdaFunctionActionAreaCardWrapper} from './AwsLambdaKeystoreSigners';
import {ValidatorsDepositsTable} from "./ValidatorsDepositsTable";
import {ValidatorDepositDataJSON} from "../../gateway/validators";

const steps = [
    'AWS Auth & Internal User Roles',
    'Create Internal Lambdas',
    'Generate Secrets',
    'Generate Validator Keys/Deposits',
    'Create/Update External Lambda Function',
    'Verify Lambda Function',
    'Request Zeus Service',
    'Submit Deposits',
];

function stepComponents(activeStep: number, onGenerateValidatorDeposits: any, onGenerateValidatorEncryptedKeystoresZip: any, onEncZipFileUpload: any, zipBlob: Blob, onHandleVerifySigners: any, onValidatorsDepositsUpload: any) {
    const steps = [
        <CreateInternalAwsLambdaUserRolesActionAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <CreateAwsInternalLambdasActionAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <CreateAwsSecretsActionAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <GenerateValidatorKeysAndDepositsAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <CreateAwsLambdaFunctionActionAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
            onEncZipFileUpload={onEncZipFileUpload}
            zipBlob={zipBlob}
        />,
        <LambdaExtUserVerify
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
            onHandleVerifySigners={onHandleVerifySigners}
        />,
        <ZeusServiceRequestAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <ValidatorsDepositRequestAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
            onValidatorsDepositsUpload={onValidatorsDepositsUpload}
        />]
    return steps[activeStep]
}

export default function AwsWizardPanel() {
    const [activeStep, setActiveStep] = React.useState(0);
    const [completed, setCompleted] = React.useState<{
        [k: number]: boolean;
    }>({});

    const totalSteps = () => {
        return steps.length;
    };

    const completedSteps = () => {
        return Object.keys(completed).length;
    };

    const isLastStep = () => {
        return activeStep === totalSteps() - 1;
    };

    const allStepsCompleted = () => {
        return completedSteps() === totalSteps();
    };

    const handleNext = () => {
        const newActiveStep =
            isLastStep() && !allStepsCompleted()
                ? // It's the last step, but not all steps have been completed,
                  // find the first step that has been completed
                steps.findIndex((step, i) => !(i in completed))
                : activeStep + 1;
        setActiveStep(newActiveStep);
    };

    const handleBack = () => {
        setActiveStep((prevActiveStep) => prevActiveStep - 1);
    };

    const handleStep = (step: number) => () => {
        setActiveStep(step);
    };

    const handleComplete = () => {
        const newCompleted = completed;
        newCompleted[activeStep] = true;
        setCompleted(newCompleted);
        handleNext();
    };

    const handleReset = () => {
        setActiveStep(0);
        setCompleted({});
    };

    const [encZipFile, setEncZipFile] = useState<Blob>(new Blob([], {type: 'application/zip'}));

    const onEncZipFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files && event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = (e) => {
                const blob = new Blob([e.target!.result as ArrayBuffer], {type: 'application/zip'});
                setEncZipFile(blob);
            };
            reader.readAsArrayBuffer(file);
        }
    };
    const dispatch = useDispatch();
    const akey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const skey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const validatorSecretsName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const ageSecretName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);
    const validatorCount = useSelector((state: RootState) => state.validatorSecrets.validatorCount);
    const hdOffset = useSelector((state: RootState) => state.validatorSecrets.hdOffset);
    const onGenerateValidatorEncryptedKeystoresZip = async () => {
        try {
            const creds = { accessKeyId: akey, secretAccessKey: skey };
            const res = await awsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(creds);
            const updatedEncKeystoresZipLambdaFnUrl = res.data;
            dispatch(setEncKeystoresZipLambdaFnUrl(updatedEncKeystoresZipLambdaFnUrl));
            const zip = await awsLambdaApiGateway.invokeEncryptedKeystoresZipGeneration(updatedEncKeystoresZipLambdaFnUrl, creds, ageSecretName, validatorSecretsName, validatorCount, hdOffset);
            const zipBlob = await zip.blob();
            const blob = new Blob([zipBlob], { type: 'application/octet-stream' });
            download(blob, "keystores.zip");
            setEncZipFile(blob);
        } catch (error) {
            console.log("error", error);
        } finally {
        }
    };
    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    let depositsGenLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.depositsGenLambdaFnUrl);
    const depositData = useSelector((state: RootState) => state.awsCredentials.depositData);

    const onGenerateValidatorDeposits = async () => {
        console.log("onGenerateValidatorDeposits");
        const creds = {accessKeyId: akey, secretAccessKey: skey};
        try {
            const response = await awsApiGateway.createValidatorsDepositDataLambda(creds);
            dispatch(setDepositsGenLambdaFnUrl(response.data));
            depositsGenLambdaFnUrl = response.data;
        } catch (error) {
            console.log("error", error);
        }
        try {
            console.log("invokeValidatorDepositsGeneration");
            const dpSlice = await awsLambdaApiGateway.invokeValidatorDepositsGeneration(depositsGenLambdaFnUrl,creds,network,validatorSecretsName,validatorCount,hdOffset);
            const body = await dpSlice.json();
            body.forEach((item: any) => {
                item.verified = false;
            });
            dispatch(setDepositData(body));
        } catch (error) {
            console.log("error", error);
        }
    };
    const onValidatorsDepositsUpload = (event: ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                const jsonData = JSON.parse(e.target?.result as string) as ValidatorDepositDataJSON[]
                console.log(jsonData);
                dispatch(setDepositData(jsonData));
            } catch (error) {
                console.error("Error parsing JSON file:", error);
            }
        };
        reader.readAsText(file);
    }

    const externalAccessUserName = useSelector((state: RootState) => state.awsCredentials.externalAccessUserName);
    const externalAccessSecretName = useSelector((state: RootState) => state.awsCredentials.externalAccessSecretName);
    const blsSignerFunctionName = useSelector((state: RootState) => state.awsCredentials.blsSignerFunctionName);

    const onHandleVerifySigners = async () => {
        const creds = {accessKeyId: akey, secretAccessKey: skey};
        try {
            const r = await awsApiGateway.createOrFetchExternalLambdaUserAccessKeys(creds, externalAccessUserName, externalAccessSecretName);
            const url = await awsApiGateway.getLambdaFunctionURL(creds, blsSignerFunctionName);
            const extCreds = {accessKeyId: r.data.accessKey, secretAccessKey: r.data.secretKey};
            const response = await awsApiGateway.verifyLambdaFunctionSigner(extCreds,ageSecretName,url.data, depositData);
            const verifiedKeys = response.data
            let hm = createHashMap(verifiedKeys);
            const verifiedDepositData = depositData.map((obj: any) => {
                if (obj.hasOwnProperty('verified')) {
                    return {
                        ...obj,
                        ['verified']: hm[obj.pubkey],
                    };
                }});
            dispatch(setDepositData(verifiedDepositData));
        } catch (error) {
            console.log("error", error);
        }}
    return (
        <div>
        <Box sx={{ width: '100%' }}>
            <Stepper nonLinear activeStep={activeStep}>
                {steps.map((label, index) => (
                    <Step key={label} completed={completed[index]}>
                        <StepButton color="inherit" onClick={handleStep(index)}>
                            {label}
                        </StepButton>
                    </Step>
                ))}
            </Stepper>
            <div>
                {allStepsCompleted() ? (
                    <React.Fragment>
                        <Typography sx={{ mt: 2, mb: 1 }}>
                            All steps completed - you&apos;re finished
                        </Typography>
                        <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
                            <Box sx={{ flex: '1 1 auto' }} />
                            <Button onClick={handleReset}>Reset</Button>
                        </Box>
                    </React.Fragment>
                ) : (
                    <React.Fragment>
                        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                            {stepComponents(activeStep, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip, onEncZipFileUpload, encZipFile, onHandleVerifySigners, onValidatorsDepositsUpload)}
                        </Container>
                        <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
                            <Button
                                color="inherit"
                                disabled={activeStep === 0}
                                onClick={handleBack}
                                sx={{ mr: 1 }}
                            >
                                Back
                            </Button>
                            <Box sx={{ flex: '1 1 auto' }} />
                            <Button onClick={handleNext} sx={{ mr: 1 }}>
                                Next
                            </Button>
                            {activeStep !== steps.length &&
                                (completed[activeStep] ? (
                                    <Typography variant="caption" sx={{ display: 'inline-block' }}>
                                    </Typography>
                                ) : (
                                    <Button onClick={handleComplete}>
                                        {completedSteps() === totalSteps() - 1
                                            ? 'Finish'
                                            : 'Complete Step'}
                                    </Button>
                                ))}
                        </Box>
                    </React.Fragment>
                )}
            </div>
        </Box>
            <ValidatorsDepositsTable depositData={depositData} activeStep={activeStep}/>
</div>
);
}

export function download(blob: any, filename: string) {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.setAttribute('download', `${filename}`);
    // the filename you want
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
}

function createHashMap(keys: string[]): { [key: string]: boolean } {
    return keys.reduce((hashMap: { [key: string]: boolean }, key: string) => {
        hashMap[key] = true;
        return hashMap;
    }, {});
}

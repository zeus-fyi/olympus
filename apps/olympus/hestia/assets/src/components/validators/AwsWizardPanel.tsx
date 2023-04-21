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
    initialState,
    setAgeSecretName,
    setDepositData,
    setDepositsGenLambdaFnUrl,
    setEncKeystoresZipLambdaFnUrl,
    setKeystoreLayerName,
    setKeystoreLayerNumber,
    setSignerFunctionName,
    setValidatorSecretsName
} from "../../redux/aws_wizard/aws.wizard.reducer";
import {PageToggleView} from './SimplifiedView';
import {awsApiGateway} from "../../gateway/aws";
import {awsLambdaApiGateway} from "../../gateway/aws.lambda";
import {CreateAwsLambdaFunctionActionAreaCardWrapper} from './AwsLambdaKeystoreSigners';
import {ValidatorsDepositsTable} from "./ValidatorsDepositsTable";
import {ValidatorDepositDataJSON, validatorsApiGateway} from "../../gateway/validators";
import {
    setAuthorizedNetworks,
    setKeyGroupName,
    setNetworkAppended
} from "../../redux/validators/ethereum.validators.reducer";
import {CircularProgress, Stack} from "@mui/material";

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

const stepsSimplified = [
    'AWS Setup',
    'Generate Validator Keys/Deposits',
    'Request Zeus Service',
    'Submit Deposits',
];

function stepComponents(activeStep: number,
                        onGenerateValidatorDeposits: any, 
                        onGenerateValidatorEncryptedKeystoresZip: any,
                        onEncZipFileUpload: any,
                        zipBlob: Blob,
                        onHandleVerifySigners: any,
                        onValidatorsDepositsUpload: any,
                        zipGenButtonLabel: any, zipGenButtonEnabled: any, zipGenStatus: any, requestStatusZipGen: any,
                        buttonLabelVd: any, buttonDisabledVd: any,statusMessageVd: any, requestStatusVd: any,
                        buttonLabelVerify: any, buttonDisabledVerify: any, statusMessageVerify: any, statusVerify: any,
                        pageView: any, onGenerateValidatorDepositsAndZip: any
) {

    const steps = pageView ? ([
        <CreateInternalAwsLambdaUserRolesActionAreaCardWrapper
            pageView={pageView}
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
            onGenerateValidatorDepositsAndZip={onGenerateValidatorDepositsAndZip}
        />,
        <GenerateValidatorKeysAndDepositsAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
            buttonLabelVd={buttonLabelVd}
            buttonDisabledVd={buttonDisabledVd}
            statusMessageVd={statusMessageVd}
            zipGenButtonLabel={zipGenButtonLabel}
            zipGenButtonEnabled={zipGenButtonEnabled}
            zipGenStatus={zipGenStatus}
            requestStatusZipGen={requestStatusZipGen}
            pageView={pageView}
            onGenerateValidatorDepositsAndZip={onGenerateValidatorDepositsAndZip}
        />,
        <CreateAwsLambdaFunctionActionAreaCardWrapper
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
            onEncZipFileUpload={onEncZipFileUpload}
            zipBlob={zipBlob}
            pageView={pageView}
            onHandleVerifySigners={onHandleVerifySigners}
        />,
        <LambdaExtUserVerify
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
            onHandleVerifySigners={onHandleVerifySigners}
            buttonLabelVerify={buttonLabelVerify}
            buttonDisabledVerify={buttonDisabledVerify}
            statusMessageVerify={statusMessageVerify}
            statusVerify={statusVerify}
            requestStatusVd={requestStatusVd}
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
        />]) :
        ([
            <CreateInternalAwsLambdaUserRolesActionAreaCardWrapper
                pageView={pageView}
                activeStep={activeStep}
                onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
            />,
            <div>
                <Stack direction="row" spacing={2}>
                    <CreateAwsSecretsActionAreaCardWrapper
                        activeStep={activeStep}
                        onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                        onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
                        pageView={pageView}
                        onGenerateValidatorDepositsAndZip={onGenerateValidatorDepositsAndZip}
                    />
                    <GenerateValidatorKeysAndDepositsAreaCardWrapper
                        activeStep={activeStep}
                        onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                        onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
                        buttonLabelVd={buttonLabelVd}
                        buttonDisabledVd={buttonDisabledVd}
                        statusMessageVd={statusMessageVd}
                        zipGenButtonLabel={zipGenButtonLabel}
                        zipGenButtonEnabled={zipGenButtonEnabled}
                        zipGenStatus={zipGenStatus}
                        requestStatusZipGen={requestStatusZipGen}
                        pageView={pageView}
                        onGenerateValidatorDepositsAndZip={onGenerateValidatorDepositsAndZip}
                    />
                </Stack>
            </div>,
            <Stack direction="row" spacing={2}>
                <CreateAwsLambdaFunctionActionAreaCardWrapper
                    activeStep={activeStep}
                    onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                    onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
                    onEncZipFileUpload={onEncZipFileUpload}
                    zipBlob={zipBlob}
                    pageView={pageView}
                    onHandleVerifySigners={onHandleVerifySigners}
                />
                <ZeusServiceRequestAreaCardWrapper
                    activeStep={activeStep}
                    onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                    onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
                />,
            </Stack>,
            <ValidatorsDepositRequestAreaCardWrapper
                activeStep={activeStep}
                onGenerateValidatorDeposits={onGenerateValidatorDeposits}
                onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
                onValidatorsDepositsUpload={onValidatorsDepositsUpload}
            />
        ]);
    return steps[activeStep]
}

export default function AwsWizardPanel(props: any) {
    const {pageView, setPageView} = props;
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

    let zipGenButtonLabel;
    let zipGenButtonEnabled;
    let zipGenStatus;
    const [requestStatusZipGen, setRequestStatusZipGen] = useState('');
    
    switch (requestStatusZipGen) {
        case 'pending':
            zipGenButtonLabel = <CircularProgress size={20} />;
            zipGenButtonEnabled = true;
            break;
        case 'success':
            zipGenButtonLabel = 'Generate';
            zipGenButtonEnabled = true;
            zipGenStatus = 'Encrypted keystores request completed successfully!';
            break;
        case 'error':
            zipGenButtonLabel = 'Retry';
            zipGenButtonEnabled = false;
            zipGenStatus = 'An error occurred while sending the keystores zip generation request.';
            break;
        case 'errorAuth':
            zipGenButtonLabel = 'Retry';
            zipGenButtonEnabled = false;
            zipGenStatus = 'Update your AWS credentials on step 1 and try again.';
            break;
        default:
            zipGenButtonLabel = 'Generate';
            zipGenButtonEnabled = false;
            break;
    }

    const dispatch = useDispatch();
    const akey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const skey = useSelector((state: RootState) => state.awsCredentials.secretKey);
    const validatorSecretsName = useSelector((state: RootState) => state.awsCredentials.validatorSecretsName);
    const ageSecretName = useSelector((state: RootState) => state.awsCredentials.ageSecretName);
    const validatorCount = useSelector((state: RootState) => state.validatorSecrets.validatorCount);
    const hdOffset = useSelector((state: RootState) => state.validatorSecrets.hdOffset);
    const onGenerateValidatorEncryptedKeystoresZip = async () => {
        try {
            setRequestStatusZipGen('pending');
            const creds = { accessKeyId: akey, secretAccessKey: skey };
            if (!akey || !skey) {
                setRequestStatusZipGen('errorAuth');
                return;
            }
            const res = await awsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(creds);
            if (res.status !== 200) {
                setRequestStatusZipGen('error');
                return
            }
            const updatedEncKeystoresZipLambdaFnUrl = res.data;
            dispatch(setEncKeystoresZipLambdaFnUrl(updatedEncKeystoresZipLambdaFnUrl));
            const zip = await awsLambdaApiGateway.invokeEncryptedKeystoresZipGeneration(updatedEncKeystoresZipLambdaFnUrl, creds, ageSecretName, validatorSecretsName, validatorCount, hdOffset);
            if (zip.status !== 200) {
                setRequestStatusZipGen('error');
                return
            }
            const zipBlob = await zip.blob();
            const blob = new Blob([zipBlob], { type: 'application/octet-stream' });
            download(blob, "keystores.zip");
            setEncZipFile(blob);
            setRequestStatusZipGen('success');
        } catch (error) {
            setRequestStatusZipGen('error');
            console.log("error", error);
        } finally {
        }
    };

    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    let depositsGenLambdaFnUrl = useSelector((state: RootState) => state.awsCredentials.depositsGenLambdaFnUrl);
    const depositData = useSelector((state: RootState) => state.awsCredentials.depositData);
    const withdrawalCredentials = useSelector((state: RootState) => state.validatorSecrets.withdrawalCredentials);
    let buttonLabelVd;
    let buttonDisabledVd;
    let statusMessageVd;
    const [requestStatusVd, setRequestStatusVd] = useState('');

    switch (requestStatusVd) {
        case 'pending':
            buttonLabelVd = <CircularProgress size={20} />;
            buttonDisabledVd = true;
            break;
        case 'success':
            buttonLabelVd = 'Created successfully';
            buttonDisabledVd = false;
            statusMessageVd = 'Validator deposits created successfully!';
            break;
        case 'error':
            buttonLabelVd = 'Error creating validator deposits';
            buttonDisabledVd = false;
            statusMessageVd = 'An error occurred while creating the validator deposits.';
            break;
        case 'errorAuth':
            buttonLabelVd = 'Retry';
            buttonDisabledVd = false;
            statusMessageVd = 'Update your AWS credentials on step 1 and try again.';
            break;
        default:
            buttonLabelVd = 'Generate';
            buttonDisabledVd = false;
            break;
    }

    const onGenerateValidatorDeposits = async () => {
        setRequestStatusVd('pending');
        const creds = {accessKeyId: akey, secretAccessKey: skey};
        if (!akey || !skey) {
            setRequestStatusVd('errorAuth');
            return;
        }
        try {
            const response = await awsApiGateway.createValidatorsDepositDataLambda(creds);
            if (response.status !== 200) {
                setRequestStatusVd('error');
                return
            }
            dispatch(setDepositsGenLambdaFnUrl(response.data));
            depositsGenLambdaFnUrl = response.data;
        } catch (error) {
            setRequestStatusVd('error');
            console.log("error", error);
            return
        }
        try {
            const dpSlice = await awsLambdaApiGateway.invokeValidatorDepositsGeneration(depositsGenLambdaFnUrl,creds,network,validatorSecretsName,validatorCount,hdOffset,withdrawalCredentials);
            if (dpSlice.status === 200) {
                setRequestStatusVd('success');
            } else {
                setRequestStatusVd('error');
                return
            }
            const body = await dpSlice.json();
            const jsonString = JSON.stringify(body, null, 2); // Convert JSON object to string with 2-space indentation
            const blob = new Blob([jsonString], { type: 'application/json' }); // Create a blob from the JSON string
            download(blob, "deposit_data.json");

            body.forEach((item: any) => {
                item.verified = false;
                item.rx = '';
            });
            dispatch(setDepositData(body));
        } catch (error) {
            setRequestStatusVd('error');
            console.log("error", error);
            return
        }
    };

    const onGenerateValidatorDepositsAndZip = async () => {
        setRequestStatusVd('pending');
        const creds = {accessKeyId: akey, secretAccessKey: skey};
        if (!akey || !skey) {
            setRequestStatusVd('errorAuth');
            return;
        }
        try {
            const creds = { accessKeyId: akey, secretAccessKey: skey };
            if (!akey || !skey) {
                setRequestStatusVd('errorAuth');
                return;
            }
            const res = await awsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(creds);
            if (res.status !== 200) {
                setRequestStatusVd('error');
                return
            }
            const updatedEncKeystoresZipLambdaFnUrl = res.data;
            dispatch(setEncKeystoresZipLambdaFnUrl(updatedEncKeystoresZipLambdaFnUrl));
            const zip = await awsLambdaApiGateway.invokeEncryptedKeystoresZipGeneration(updatedEncKeystoresZipLambdaFnUrl, creds, ageSecretName, validatorSecretsName, validatorCount, hdOffset);
            if (zip.status !== 200) {
                setRequestStatusVd('error');
                return
            }
            const zipBlob = await zip.blob();
            const blob = new Blob([zipBlob], { type: 'application/octet-stream' });
            download(blob, "keystores.zip");
            setEncZipFile(blob);

            const response = await awsApiGateway.createValidatorsDepositDataLambda(creds);
            if (response.status !== 200) {
                setRequestStatusVd('error');
                return
            }
            dispatch(setDepositsGenLambdaFnUrl(response.data));
            depositsGenLambdaFnUrl = response.data;
        } catch (error) {
            setRequestStatusVd('error');
            console.log("error", error);
            return
        }
        try {
            const dpSlice = await awsLambdaApiGateway.invokeValidatorDepositsGeneration(depositsGenLambdaFnUrl,creds,network,validatorSecretsName,validatorCount,hdOffset,withdrawalCredentials);
            if (dpSlice.status === 200) {
                setRequestStatusVd('success');
            } else {
                setRequestStatusVd('error');
                return
            }
            const body = await dpSlice.json();
            const jsonString = JSON.stringify(body, null, 2); // Convert JSON object to string with 2-space indentation
            const blob = new Blob([jsonString], { type: 'application/json' }); // Create a blob from the JSON string
            download(blob, "deposit_data.json");

            body.forEach((item: any) => {
                item.verified = false;
                item.rx = '';
            });
            dispatch(setDepositData(body));
        } catch (error) {
            setRequestStatusVd('error');
            console.log("error", error);
            return
        }
    };

    const onValidatorsDepositsUpload = (event: ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = (e) => {
            try {
                const jsonData = JSON.parse(e.target?.result as string) as ValidatorDepositDataJSON[]
                jsonData.forEach((item: any) => {
                    item.verified = false;
                    item.rx = '';
                });
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

    let buttonLabelVerify;
    let buttonDisabledVerify;
    let statusMessageVerify;
    const [statusVerify, setStatusVerify] = useState('');

    switch (statusVerify) {
        case 'pending':
            buttonLabelVerify = <CircularProgress size={20} />;
            buttonDisabledVerify = true;
            break;
        case 'success':
            buttonLabelVerify = 'Verify request completed successfully';
            buttonDisabledVerify = true;
            statusMessageVerify = 'Verify request completed successfully!';
            break;
        case 'error':
            buttonLabelVerify = 'Error sending verify request';
            buttonDisabledVerify = false;
            statusMessageVerify = 'An error occurred while sending the verify request.';
            break;
        case 'errorAuth':
            buttonLabelVerify = 'Retry';
            buttonDisabledVerify = false;
            statusMessageVerify = 'Update your AWS credentials on step 1 and try again.';
            break;
        default:
            buttonLabelVerify = 'Send Verify Request';
            buttonDisabledVerify = false;
            break;
    }
    const onHandleVerifySigners = async () => {
        const creds = {accessKeyId: akey, secretAccessKey: skey};
        if (!akey || !skey) {
            setStatusVerify('errorAuth');
            return;
        }
        try {
            setStatusVerify('pending');
            const r = await awsApiGateway.createOrFetchExternalLambdaUserAccessKeys(creds, externalAccessUserName, externalAccessSecretName);
            const url = await awsApiGateway.getLambdaFunctionURL(creds, blsSignerFunctionName);
            const extCreds = {accessKeyId: r.data.accessKey, secretAccessKey: r.data.secretKey};
            const response = await awsApiGateway.verifyLambdaFunctionSigner(extCreds,ageSecretName,url.data, depositData);
            if (response.status !== 200) {
                setStatusVerify('error');
                return
            } else {
                setStatusVerify('success');
            }
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
            setStatusVerify('error');
            console.log("error", error);
        }}

    const networkAppended = useSelector((state: RootState) => state.validatorSecrets.networkAppended);
    const handleNetworkAppend = () => {
        if (!networkAppended) {
            dispatch(setKeystoreLayerNumber(initialState.keystoreLayerNumber));
            const newValidatorSecretsName = initialState.validatorSecretsName + network;
            dispatch(setValidatorSecretsName(newValidatorSecretsName));
            const newAgeSecretName = initialState.ageSecretName + network;
            dispatch(setAgeSecretName(newAgeSecretName));
            const newBlsSignerFunctionName  = initialState.blsSignerFunctionName + network;
            dispatch(setSignerFunctionName(newBlsSignerFunctionName));
            const newKeystoresLayerName = initialState.blsSignerKeystoresLayerName + network;
            dispatch(setKeystoreLayerName(newKeystoresLayerName));
            const newKeyGroupName = 'DefaultKeyGroup' + network;
            dispatch(setKeyGroupName(newKeyGroupName));
            dispatch(setNetworkAppended(true));
        }
    };

    const fetchData = async () => {
        const response = await validatorsApiGateway.getAuthedValidatorsServiceRequest();
        const verifiedNetworks: [string] = response.data;
        dispatch(setAuthorizedNetworks(verifiedNetworks));
    };

    React.useEffect(() => {
        fetchData().catch((e) => {
            console.log("error", e);
        })
        handleNetworkAppend();
    }, [network]);
    return (
        <div>
        <Box sx={{ width: '100%' }}>
            {pageView ? (
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Stepper nonLinear activeStep={activeStep}>
                        {steps.map((label, index) => (
                            <Step key={label} completed={completed[index]}>
                                <StepButton color="inherit" onClick={handleStep(index)}>
                                    {label}
                                </StepButton>
                            </Step>
                        ))}
                    </Stepper>
                </Container>
            ) :
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Stepper nonLinear activeStep={activeStep}>
                        {stepsSimplified.map((label, index) => (
                            <Step key={label} completed={completed[index]}>
                                <StepButton color="inherit" onClick={handleStep(index)}>
                                    {label}
                                </StepButton>
                            </Step>
                        ))}
                    </Stepper>
                </Container>
            }
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
                            {stepComponents(activeStep,
                                onGenerateValidatorDeposits,
                                onGenerateValidatorEncryptedKeystoresZip,
                                onEncZipFileUpload,
                                encZipFile,
                                onHandleVerifySigners,
                                onValidatorsDepositsUpload,
                                zipGenButtonLabel,
                                zipGenButtonEnabled,
                                zipGenStatus,
                                requestStatusZipGen,
                                buttonLabelVd,
                                buttonDisabledVd,
                                statusMessageVd,
                                requestStatusVd,
                                buttonLabelVerify,
                                buttonDisabledVerify,
                                statusMessageVerify,
                                statusVerify,
                                pageView,
                                onGenerateValidatorDepositsAndZip
                            )}
                        </Container>

                        <Box sx={{mb: 2}}>
                            <PageToggleView pageView={pageView} setPageView={setPageView}/>
                        </Box>
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

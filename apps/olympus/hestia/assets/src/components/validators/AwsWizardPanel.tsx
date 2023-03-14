import * as React from 'react';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepButton from '@mui/material/StepButton';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Container from "@mui/material/Container";
import {CreateAwsInternalLambdasActionAreaCardWrapper, CreateAwsSecretsActionAreaCardWrapper,} from "./AwsSecrets";
import {CreateInternalAwsLambdaUserRolesActionAreaCardWrapper} from "./AwsLambdaUserRolePolicies";
import {CreateAwsLambdaFunctionActionAreaCardWrapper} from "./AwsLambdaCreation";
import {LambdaExtUserVerify} from "./AwsExtUserAndLambdaVerify";
import {GenerateValidatorKeysAndDepositsAreaCardWrapper} from "./ValidatorsGeneration";
import {ZeusServiceRequestAreaCardWrapper} from "./ZeusServiceRequest";
import {ValidatorsDepositRequestAreaCardWrapper} from "./ValidatorsDeposits";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";

const steps = [
    'AWS Auth & Internal User Roles',
    'Create Internal Lambdas',
    'Generate Secrets',
    'Generate Validator Keys/Deposits',
    'Create/Update External Lambda Function',
    'Create Zeus Service Request',
    'Submit Deposits',
];

function stepComponents(activeStep: number, onGenerateValidatorDeposits: any, onGenerateValidatorEncryptedKeystoresZip: any) {
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
        />,
        <LambdaExtUserVerify
            activeStep={activeStep}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
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


    const count = useSelector((state: RootState) => state.validatorSecrets.validatorCount);
    const hdOffset = useSelector((state: RootState) => state.validatorSecrets.hdOffset);

    const onGenerateValidatorDeposits = async () => {
        try {
            // TODO this is a stub
            console.log('onGenerateValidatorDeposits')
            //const response = await validatorsApiGateway.generateValidatorsDepositDataLambda(mnemonic, hdWalletPw,count,hdOffset);
            //console.log(response.data)
        } catch (error) {
            console.log("error", error);
        }};

    const onGenerateValidatorEncryptedKeystoresZip = async () => {
        try {
            // TODO this is a stub
            console.log('onGenerateValidatorEncryptedKeystoresZip')
            //const response = await validatorsApiGateway.createValidatorsAgeEncryptedKeystoresZipLambda(agePubKey, agePrivKey,mnemonic, hdWalletPw, count, hdOffset);
            //console.log(response.data)
        } catch (error) {
            console.log("error", error);
        }};
    return (
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
                            {stepComponents(activeStep, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip)}
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
    );
}
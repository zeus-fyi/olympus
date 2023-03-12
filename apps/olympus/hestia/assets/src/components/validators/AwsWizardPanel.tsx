import * as React from 'react';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepButton from '@mui/material/StepButton';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Container from "@mui/material/Container";
import {charsets, CreateAwsSecretsActionAreaCardWrapper, generatePassword} from "./AwsSecrets";
import {CreateInternalAwsLambdaUserRolesActionAreaCardWrapper} from "./AwsLambdaUserRolePolicies";
import {CreateAwsLambdaFunctionActionAreaCardWrapper} from "./AwsLambdaCreation";
import {LambdaExtUserVerify} from "./AwsExtUserAndLambdaVerify";
import {GenerateValidatorKeysAndDepositsAreaCardWrapper} from "./ValidatorsGeneration";
import {ZeusServiceRequestAreaCardWrapper} from "./ZeusServiceRequest";
import {ValidatorsDepositRequestAreaCardWrapper} from "./ValidatorsDeposits";
import {awsApiGateway} from "../../gateway/aws";
import {setAgePrivKey, setAgePubKey} from "../../redux/aws_wizard/aws.wizard.reducer";
import {ethers} from "ethers";
import {setHdWalletPw, setMnemonic} from "../../redux/validators/ethereum.validators.reducer";
import {validatorsApiGateway} from "../../gateway/validators";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";

const steps = [
    'Create AWS Secrets',
    'Generate Validator Deposits',
    'Create Lambda User Roles',
    'Create or Update Lambda Function',
    'Verify Lambda Function',
    'Create Zeus Service Request',
    'Submit Deposits',
];

function stepComponents(activeStep: number, onGenerate: any, onGenerateValidatorDeposits: any, onGenerateValidatorEncryptedKeystoresZip: any) {
    const steps = [
        <CreateAwsSecretsActionAreaCardWrapper
            activeStep={activeStep}
            onGenerate={onGenerate}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <GenerateValidatorKeysAndDepositsAreaCardWrapper
            activeStep={activeStep}
            onGenerate={onGenerate}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <CreateInternalAwsLambdaUserRolesActionAreaCardWrapper
            activeStep={activeStep}
            onGenerate={onGenerate}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <CreateAwsLambdaFunctionActionAreaCardWrapper
            activeStep={activeStep}
            onGenerate={onGenerate}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <LambdaExtUserVerify
            activeStep={activeStep}
            onGenerate={onGenerate}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <ZeusServiceRequestAreaCardWrapper
            activeStep={activeStep}
            onGenerate={onGenerate}
            onGenerateValidatorDeposits={onGenerateValidatorDeposits}
            onGenerateValidatorEncryptedKeystoresZip={onGenerateValidatorEncryptedKeystoresZip}
        />,
        <ValidatorsDepositRequestAreaCardWrapper
            activeStep={activeStep}
            onGenerate={onGenerate}
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

    const dispatch = useDispatch();
    const onGenerate = async () => {
        try {
            console.log('onGenerate')
            const response = await awsApiGateway.getGeneratedAgeKey();
            const ageKeyGenData: any = response.data;
            dispatch(setAgePrivKey(ageKeyGenData.agePrivateKey));
            dispatch(setAgePubKey(ageKeyGenData.agePublicKey));
            const entropyBytes = ethers.randomBytes(32); // 16 bytes = 128 bits of entropy
            let phrase = ethers.Mnemonic.fromEntropy(entropyBytes).phrase;
            dispatch(setMnemonic(phrase));
            const password = generatePassword(20, charsets.NUMBERS + charsets.LOWERCASE + charsets.UPPERCASE + charsets.SYMBOLS);
            dispatch(setHdWalletPw(password));
        } catch (error) {
            console.log("error", error);
        }};

    const mnemonic = useSelector((state: RootState) => state.validatorSecrets.mnemonic);
    const hdWalletPw = useSelector((state: RootState) => state.validatorSecrets.hdWalletPw);

    const agePubKey = useSelector((state: RootState) => state.awsCredentials.agePubKey);
    const agePrivKey = useSelector((state: RootState) => state.awsCredentials.agePrivKey);
    const count = useSelector((state: RootState) => state.validatorSecrets.validatorCount);
    const hdOffset = useSelector((state: RootState) => state.validatorSecrets.hdOffset);

    const onGenerateValidatorDeposits = async () => {
        try {
            // TODO this is a stub
            console.log('onGenerateValidatorDeposits')
            const response = await validatorsApiGateway.generateValidatorsDepositData(mnemonic, hdWalletPw,count,hdOffset);
            console.log(response.data)
        } catch (error) {
            console.log("error", error);
        }};

    const onGenerateValidatorEncryptedKeystoresZip = async () => {
        try {
            // TODO this is a stub
            console.log('onGenerateValidatorEncryptedKeystoresZip')
            const response = await validatorsApiGateway.generateValidatorsAgeEncryptedKeystoresZip(agePubKey, agePrivKey,mnemonic, hdWalletPw, count, hdOffset);
            console.log(response.data)
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
                            {stepComponents(activeStep, onGenerate, onGenerateValidatorDeposits, onGenerateValidatorEncryptedKeystoresZip)}
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
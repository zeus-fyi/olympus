import * as React from "react";
import {clusterBuilderSteps} from "../clusters/wizard/builder/WizardSteps";
import Box from "@mui/material/Box";
import Stepper from "@mui/material/Stepper";
import Step from "@mui/material/Step";
import StepButton from "@mui/material/StepButton";
import Button from "@mui/material/Button";
import Container from "@mui/material/Container";
import {AppPage} from "./AppPage";
import {DeployPage} from "./DeployPage";

export const appBuildToggleSteps = [
    'Deploy App',
    'Configs',
    // 'Resources'
];

export function appPageStepComponents(activeStep: number, app: string, cloudProvider: any, setCloudProvider: any, region: any, setRegion:any) {
    const steps = [
        <DeployPage app={app} cloudProvider={cloudProvider} setCloudProvider={setCloudProvider} region={region} setRegion={setRegion}/>,
        <AppPage />,
        // <AppResourceNodesResourcesTable />
    ];
    return (steps[activeStep])
}

export default function DeployConfigToggle(props: any) {
    const {app, region, setRegion, cloudProvider, setCloudProvider} = props
    const [activeStep, setActiveStep] = React.useState(0);
    const [completed, setCompleted] = React.useState<{
        [k: number]: boolean;
    }>({});

    const totalSteps = () => {
        return appBuildToggleSteps.length;
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
                clusterBuilderSteps.findIndex((step, i) => !(i in completed))
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

    React.useEffect(() => {

    }, []);
    return (
        <div>
            <Box sx={{ width: '100%' }}>
                <Stepper nonLinear activeStep={activeStep}>
                    {appBuildToggleSteps.map((label, index) => (
                        <Step key={label} completed={completed[index]}>
                            <StepButton color="inherit" onClick={handleStep(index)}>
                                {label}
                            </StepButton>
                        </Step>
                    ))}
                </Stepper>
                <div>
                        <React.Fragment>
                            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                                {appPageStepComponents(activeStep, app, cloudProvider, setCloudProvider, region, setRegion)}
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
                            </Box>
                        </React.Fragment>
                </div>
            </Box>
        </div>
    );
}
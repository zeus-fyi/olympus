import * as React from 'react';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepButton from '@mui/material/StepButton';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import Container from "@mui/material/Container";
import {CsvUploadActionAreaCard} from "./Upload";
import {AnalyzeActionAreaCard} from "./Analyze";
import {Commands} from "./Commands";
import {Results} from "./Results";

const steps = [
    'AWS Auth & Internal User Roles',
    'Create Internal Lambdas',
];

const stepsSimplified = [
    'Contacts',
    'Prompts',
    'Analysis',
    'Results',
];

function stepComponents(activeStep: number,
                        pageView: any,
) {
    const steps = [
        <CsvUploadActionAreaCard />,
        <AnalyzeActionAreaCard />,
        <Commands />,
        <Results />,
    ];
    return steps[activeStep]
}

export default function WizardPanel(props: any) {
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
                                {stepComponents(activeStep, pageView,)}
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
            {/*<UploadTable activeStep={activeStep}/>*/}
        </div>
    );
}
export const PageToggleView = (props: any) => {
    const { pageView, setPageView } = props;
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setPageView(event.target.checked);
    };

    return (
        <div>
            {/*<Stack direction={"row"} spacing={2} alignItems={"center"}>*/}
            {/*    {pageView ? (*/}
            {/*        <p>Advanced View</p>*/}
            {/*    ) : (*/}
            {/*        <p>Simplified View</p>*/}
            {/*    )}*/}
            {/*    <Switch*/}
            {/*        checked={pageView}*/}
            {/*        onChange={handleChange}*/}
            {/*        color="primary"*/}
            {/*        name="pageView"*/}
            {/*        inputProps={{ 'aria-label': 'toggle page view' }}*/}
            {/*    />*/}
            {/*</Stack>*/}
        </div>
    );
};
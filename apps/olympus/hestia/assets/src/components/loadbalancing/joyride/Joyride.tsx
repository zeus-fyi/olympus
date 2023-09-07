import React from 'react';
import Joyride, {CallBackProps, STATUS, Step} from 'react-joyride';
import {useSetState} from 'react-use';

interface State {
    run: boolean;
    steps: Step[];
}

export default function JoyrideTutorialBegin(props: any) {
    const { runTutorial, setSelectedMainTab, handleChangeGroup } = props;
    const [{ run, steps }, setState] = useSetState<State>({
        run: runTutorial,
        steps: [
            {
                content: <h2>Let's get started!</h2>,
                locale: { skip: <strong aria-label="skip"></strong> },
                placement: 'center',
                target: 'body',
            },
            {
                content: 'This view shows all the routes you have registered for use with the Load Balancer.',
                placement: 'bottom',
                target: '.onboarding-card-highlight-all-routes', // css class we'll add to the Card for targeting
                title: 'All Routes',
            },
            {
                content: 'This view shows all registered routing procedures you have access to.',
                placement: 'bottom',
                target: '.onboarding-card-highlight-all-procedures', // css class we'll add to the Card for targeting
                title: 'All Procedures',
            },
        ],
    });
    const handleJoyrideCallback = (data: CallBackProps) => {
        const { status, index } = data;
        const finishedStatuses: string[] = [STATUS.FINISHED, STATUS.SKIPPED];

        if (status === STATUS.RUNNING && index === 1) {
            setSelectedMainTab(1);
            // Just before the last step starts, we call the handleChangeGroup function
            //handleChangeGroup('ethereum-mainnet');
        }
        if (finishedStatuses.includes(status)) {
            setState({ run: false });
        }
    };

    return (
            <Joyride
                callback={handleJoyrideCallback}
                continuous
                hideCloseButton
                run={run}
                scrollToFirstStep
                showProgress
                showSkipButton
                steps={steps}
                styles={{
                    options: {
                        zIndex: 10000,
                    },
                }}
            />
    );
}

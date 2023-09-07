import React from 'react';
import Joyride, {Step} from 'react-joyride';
import {useSetState} from 'react-use';

export interface State {
    run: boolean;
    steps: Step[];
}

export default function JoyrideTutorialBegin(props: any) {
    const { runTutorial, setSelectedMainTab, handleChangeGroup, setSelectedTab,
        setTableRoutes, groups, handleJoyrideCallback } = props;
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
            {
                content: 'This view your generated routing table.',
                placement: 'bottom',
                target: '.onboarding-card-highlight-qn-routing-table', // css class we'll add to the Card for targeting
                title: 'QuickNode Generated Routing Table',
            },
        ],
    });


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

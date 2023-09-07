import React from 'react';
import Joyride, {Step} from 'react-joyride';
import {useSetState} from 'react-use';

export interface State {
    run: boolean;
    steps: Step[];
}

export default function JoyrideTutorialBegin(props: any) {
    const { runTutorial, setSelectedMainTab, handleChangeGroup, setSelectedTab,
        setTableRoutes, groups, groupName, handleJoyrideCallback } = props;

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
                content: 'This view shows your generated routing table.',
                placement: 'bottom',
                target: '.onboarding-card-highlight-qn-routing-table',
                title: 'QuickNode Generated Routing Table',
            },
            {
                content: 'When you\'re on the table view, you\'ll be able to submit a sample request with the procedure. Send one now, so you can see the metrics chart!',
                placement: 'bottom',
                target: '.onboarding-card-highlight-procedures',
                title: 'Procedures',
            },
            {
                content: 'This view shows your priority score routes table & scale factors. You can adjust the scale factors to change the priority scoring weights if desired.',
                placement: 'bottom',
                target: '.onboarding-card-highlight-priority-scores',
                title: 'Priority Scores',
            },
            {
                content: 'This view shows your route latency metrics.',
                placement: 'bottom',
                target: '.onboarding-card-highlight-metrics',
                title: 'Metrics',
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

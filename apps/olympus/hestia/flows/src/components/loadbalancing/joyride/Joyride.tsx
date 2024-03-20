import React from 'react';
import Joyride, {Step} from 'react-joyride';

export interface State {
    run: boolean;
    steps: Step[];
}

export default function JoyrideTutorialBegin(props: any) {
    const { run, steps, handleJoyrideCallback } = props;
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

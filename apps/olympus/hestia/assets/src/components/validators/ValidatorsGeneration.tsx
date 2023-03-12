import * as React from "react";
import {useState} from "react";
import {Card, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import TextField from "@mui/material/TextField";

export function GenerateValidatorKeysAndDepositsAreaCardWrapper(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard activeStep={activeStep}/>
            <GenerateValidatorsParams />
        </Stack>
    );
}

export function GenerateValidatorsParams() {
    const [awsValidatorsNetwork, setAwsValidatorsNetwork] = useState('Ephemery');
    const [validatorCount, onValidatorCountChange ] = useState('1');
    const [offset, setOffset] = useState('0');

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <ValidatorsNetwork awsValidatorsNetwork={awsValidatorsNetwork}/>
                    <ValidatorCount validatorCount={validatorCount}/>
                    <ValidatorOffsetHD offset={offset}/>
                </Container>
            </div>
        </Card>

    );
}

export function ValidatorsNetwork(props: any) {
    const { awsValidatorsNetwork, setAwsValidatorsNetwork } = props;

    return (
        <TextField
            fullWidth
            id="validatorsNetwork"
            label="Network Name"
            variant="outlined"
            value={awsValidatorsNetwork}
            onChange={setAwsValidatorsNetwork}
            sx={{ width: '100%' }}
        />
    );
}

export function ValidatorCount(props: any) {
    const { validatorCount, onValidatorCountChange } = props;
    return (
        <TextField
            fullWidth
            id="validatorCount"
            label="Validator Count"
            variant="outlined"
            value={validatorCount}
            onChange={onValidatorCountChange}
            sx={{ width: '100%' }}
        />
    );
}

export function ValidatorOffsetHD(props: any) {
    const { offset, setOffset } = props;
    return (
        <TextField
            fullWidth
            id="validatorOffsetHD"
            label="Validator HD Offset"
            variant="outlined"
            value={offset}
            onChange={setOffset}
            sx={{ width: '100%' }}
        />
    );
}
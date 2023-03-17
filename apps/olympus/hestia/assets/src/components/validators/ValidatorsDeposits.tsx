import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import * as React from "react";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {Network} from "./ZeusServiceRequest";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {
    createValidatorsDepositsDataJSON,
    createValidatorsDepositServiceRequest,
    validatorsApiGateway
} from "../../gateway/validators";

export function ValidatorsDepositRequestAreaCardWrapper(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ValidatorsDepositRequestAreaCard />
        </Stack>
    );
}

export function ValidatorsDepositRequestAreaCard() {
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <ValidatorsDepositsSubmitWrapper />
            </Container >
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <SubmitValidators />
            </Container >
        </div>

    );
}

export function SubmitValidators() {
    const depositData = useSelector((state: RootState) => state.awsCredentials.depositData);
    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    const onClickSendValidatorsDeposits = async () => {
        try {
            const depositParams = depositData.map((dd: any) => {
                return createValidatorsDepositsDataJSON(dd.pubkey, dd.withdrawal_credentials, dd.signature, dd.deposit_data_root,dd.amount,dd.deposit_message_root,dd.fork_version);
            });
            console.log("depositParams", depositParams)
            const reqParams = createValidatorsDepositServiceRequest(network, depositParams)
            console.log("reqParams", reqParams)
            const response = await validatorsApiGateway.depositValidatorsServiceRequest(reqParams);
            console.log("response", response);
        } catch (error) {
            console.log("error", error);
        }}
    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Send Validator Deposits to Network
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Sends Validator Deposits to the Network
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small" onClick={onClickSendValidatorsDeposits}>Send</Button>
            </CardActions>
        </Card>
    );
}


export function ValidatorsDepositsSubmitWrapper() {
    const network = useSelector((state: RootState) => state.validatorSecrets.network);

    return (
        <Card sx={{ maxWidth: 500 }}>
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Network network={network}/>
                </Container>
            </div>
        </Card>

    );
}

// export function Eth1WalletPrivateKey(props: any) {
//     const { eth1Pk, onAccessEth1PkChange } = props;
//     return (
//         <TextField
//             fullWidth
//             id="eth1WalletPrivateKey"
//             label="Eth1 Wallet Private Key"
//             variant="outlined"
//             value={eth1Pk}
//             onChange={onAccessEth1PkChange}
//             sx={{ width: '100%' }}
//         />
//     );
// }
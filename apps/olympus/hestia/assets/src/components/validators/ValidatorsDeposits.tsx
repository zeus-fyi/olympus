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
import {ValidatorsUploadActionAreaCard} from "./ValidatorsUpload";

export function ValidatorsDepositRequestAreaCardWrapper(props: any) {
    const { activeStep, onValidatorsDepositsUpload } = props;

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <ValidatorsUploadActionAreaCard onValidatorsDepositsUpload={onValidatorsDepositsUpload}/>,
            <ValidatorsDepositRequestAreaCard />
        </Stack>
    );
}

export function ValidatorsDepositRequestAreaCard() {
    return (
        <div style={{ display: 'flex' }}>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <SubmitValidators />
            </Container >
        </div>

    );
}

export function SubmitValidators(props: any) {
    const depositData = useSelector((state: RootState) => state.awsCredentials.depositData);
    const network = useSelector((state: RootState) => state.validatorSecrets.network);
    const onClickSendValidatorsDeposits = async () => {
        try {
            const depositParams = depositData.map((dd: any) => {
                return createValidatorsDepositsDataJSON(dd.pubkey, dd.withdrawal_credentials, dd.signature, dd.deposit_data_root,dd.amount,dd.deposit_message_root,dd.fork_version);
            });
            const reqParams = createValidatorsDepositServiceRequest(network, depositParams)
            const response = await validatorsApiGateway.depositValidatorsServiceRequest(reqParams);
        } catch (error) {
            console.log("error", error);
        }}
    return (
        <Card sx={{ maxWidth: 500 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Send Validator Deposits to Network
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    If you didn't generate your validator deposits from the previous step, you can upload your
                    deposit data JSON file to the left.
                </Typography>
            </CardContent>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <Network network={network}/>
            </Container>
            <CardActions>
                <Button size="small" onClick={onClickSendValidatorsDeposits}>Send</Button>
            </CardActions>
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
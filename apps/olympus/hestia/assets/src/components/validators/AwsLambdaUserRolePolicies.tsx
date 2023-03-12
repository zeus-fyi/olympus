import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import {awsApiGateway} from "../../gateway/aws";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";

export function CreateInternalAwsLambdaUserRolesActionAreaCardWrapper(props: any) {
    const { activeStep } = props;
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard activeStep={activeStep}/>
            <CreateAwsLambdaUserRolesActionAreaCard />
        </Stack>
    );
}

export function CreateAwsLambdaUserRolesActionAreaCard() {
    return (
            <div style={{ display: 'flex' }}>
                <Stack direction="column" alignItems="center" spacing={2}>
                </Stack>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <InternalLambdaUserRolePolicySetup />
                </Container >
            </div>
    );
}

export function InternalLambdaUserRolePolicySetup() {
    const accessKey = useSelector((state: RootState) => state.awsCredentials.accessKey);
    const secretKey = useSelector((state: RootState) => state.awsCredentials.secretKey);

    const handleCreateUser = async () => {
        try {
            const response = await awsApiGateway.createInternalLambdaUser(accessKey,secretKey);
            console.log("response", response);
        } catch (error) {
            console.log("error", error);
        }};

    return (
        <Card sx={{ maxWidth: 400 }}>
            <CardContent>
                <Typography gutterBottom variant="h5" component="div">
                    Internal User & RolePolicy Setup
                </Typography>
                <Typography variant="body2" color="text.secondary">
                    Creates a new user and role policy for your own internal usage, e.g. for running, testing, development, etc.
                </Typography>
            </CardContent>
            <CardActions>
                <Button size="small" onClick={handleCreateUser}>Create</Button>
            </CardActions>
        </Card>
    );
}
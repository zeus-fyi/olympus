import * as React from "react";
import {Card, CardActions, CardContent, Container, Stack} from "@mui/material";
import {AwsUploadActionAreaCard} from "./AwsPanel";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";

export function CreateInternalAwsLambdaUserRolesActionAreaCardWrapper() {
    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <AwsUploadActionAreaCard />
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
                <Button size="small">Create</Button>
            </CardActions>
        </Card>
    );
}
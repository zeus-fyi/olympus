import {Box, Card, CardContent, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import * as React from "react";
import {SelectedSkeletonBaseName} from "./AddSkeletonBaseDockerConfigs";

export function WorkloadPreviewAndSubmitPage(props: any) {
    const {} = props;
    return (
        <div>
            <Stack direction="row" spacing={2}>
                <Card sx={{ maxWidth: 500 }}>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Workload Config
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Sets Infra and App Configs
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <Box mt={2}>
                            <SelectedComponentBaseName />
                        </Box>
                        <Box mt={2}>
                            <SelectedSkeletonBaseName />
                        </Box>
                    </Container>
                </Card>
            </Stack>
        </div>
    );
}
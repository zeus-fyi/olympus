import * as React from "react";
import {Box, Button, Card, CardContent, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import {AddSkeletonBaseDockerConfigs, SelectedSkeletonBaseName} from "./AddSkeletonBaseDockerConfigs";

export function WorkloadConfigPage(props: any) {
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
                        <Box mt={2}>
                            <Button variant="contained">
                                Add Deployment
                            </Button>
                        </Box>
                        <Box mt={2}>
                            <Button variant="contained">
                                Add StatefulSet
                            </Button>
                        </Box>
                        <Box mt={2}>
                            <Button variant="contained">
                                Add Service
                            </Button>
                        </Box>
                        <Box mt={2}>
                            <Button variant="contained">
                                Add Ingress
                            </Button>
                        </Box>
                        <Box mt={2}>
                            <Button variant="contained">
                                Add ServiceMonitor
                            </Button>
                        </Box>
                    </Container>
                </Card>
                <AddSkeletonBaseDockerConfigs />
            </Stack>
        </div>
    );
}
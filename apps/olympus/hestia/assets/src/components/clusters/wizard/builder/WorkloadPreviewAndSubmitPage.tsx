import {Box, Button, Card, CardContent, Container, Stack} from "@mui/material";
import Typography from "@mui/material/Typography";
import {SelectedComponentBaseName} from "./DefineComponentBases";
import * as React from "react";
import {useState} from "react";
import {SelectedSkeletonBaseName} from "./AddSkeletonBaseDockerConfigs";

export function WorkloadPreviewAndSubmitPage(props: any) {
    const {} = props;
    const [viewField, setViewField] = useState('');
    const onChangeComponentOrSkeletonBase = () => {
        setViewField('')
    }
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
                            <SelectedComponentBaseName onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                        <Box mt={2}>
                            <SelectedSkeletonBaseName onChangeComponentOrSkeletonBase={onChangeComponentOrSkeletonBase}/>
                        </Box>
                        <Box mt={2}>
                            <Button variant="contained">
                                Generate Preview
                            </Button>
                        </Box>
                        <Box mt={2}>
                            <Button variant="contained">
                                Create Cluster
                            </Button>
                        </Box>
                    </Container>
                </Card>
            </Stack>
        </div>
    );
}
import {Box, Card, CardContent, Container} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import TextField from "@mui/material/TextField";
import {AddPortsInputFields} from "./DefinePorts";
import {setSelectedDockerImageName} from "../../../../redux/clusters/clusters.builder.reducer";

export function DefineDockerParams(props: any) {
    const {} = props;
    return (
        <div>
            <Card sx={{ maxWidth: 500 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Set Container Configs
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Sets Docker Image Default
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <DockerImageName />
                    <DockerImageCmdArgs />
                    <Box mt={2}>
                        <AddPortsInputFields />
                    </Box>
                </Container>
            </Card>
        </div>
    );
}

export function DockerImageName() {
    const dispatch = useDispatch();
    const cluster  = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const dockerImageName = useSelector((state: RootState) => state.clusterBuilder.selectedDockerImageName);
    const onDockerImageNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const newDockerImageName = event.target.value;
        dispatch(setSelectedDockerImageName(newDockerImageName));
    };
    return (
        <div>
            <Box mt={2}>
                <TextField
                    fullWidth
                    id="containerName"
                    label="Container Name"
                    variant="outlined"
                    value={dockerImageName}
                    onChange={onDockerImageNameChange}
                    sx={{ width: '100%' }}
                />
            </Box>
            <Box mt={2}>
                <TextField
                    fullWidth
                    id="dockerImage"
                    label="Docker Image Name"
                    variant="outlined"
                    value={dockerImageName}
                    onChange={onDockerImageNameChange}
                    sx={{ width: '100%' }}
                />
            </Box>
        </div>
    );
}

export function DockerImageCmdArgs() {
    const dispatch = useDispatch();

    const cmd = ''
    const args = ''
    return (
        <div>
            <Box mt={2}>
                <TextField
                    fullWidth
                    id="dockerImageCmd"
                    label="Docker Cmd"
                    variant="outlined"
                    value={cmd}
                    sx={{ width: '100%' }}
                />
            </Box>
            <Box mt={2}>
                <TextField
                    fullWidth
                    id="dockerImageArgs"
                    label="Docker Args"
                    variant="outlined"
                    value={args}
                    sx={{ width: '100%' }}
                />
            </Box>
        </div>
    );
}
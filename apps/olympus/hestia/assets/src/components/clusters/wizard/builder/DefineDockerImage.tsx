import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import TextField from "@mui/material/TextField";
import {AddPortsInputFields} from "./DefinePorts";
import {setSelectedContainerName} from "../../../../redux/clusters/clusters.builder.reducer";

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
                    <ContainerConfig />
                    <DockerConfig />
                </Container>
            </Card>
        </div>
    );
}

export function ContainerConfig() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);

    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    if (cluster.componentBases === undefined) {
        return <div></div>
    }
    let show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }

    const skeletonBaseContainerNames = skeletonBaseKeys[selectedSkeletonBaseName];
    show = skeletonBaseContainerNames !== undefined && Object.keys(skeletonBaseContainerNames.containers).length > 0;
    if (!show) {
        return <div></div>
    }
    const onContainerNameChange = (newContainerName: string) => {
        dispatch(setSelectedContainerName(newContainerName));
    };
    return (
        <div>
            {show &&
                <FormControl variant="outlined" style={{ minWidth: '100%' }}>
                    <InputLabel id="network-label">Containers</InputLabel>
                    <Select
                        labelId="containerName-label"
                        id="containerName"
                        value={selectedContainerName}
                        label="Container Name"
                        onChange={(event) => onContainerNameChange(event.target.value as string)}
                        sx={{ width: '100%' }}
                    >
                        {Object.keys(skeletonBaseContainerNames.containers).map((key: any, i: number) => (
                            <MenuItem key={i} value={key}>
                                {key}
                            </MenuItem>))
                        }
                    </Select>
                </FormControl>
            }
        </div>);
}


export function DockerConfig() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    if (cluster.componentBases === undefined) {
        return <div></div>
    }
    let show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }

    const skeletonBaseContainerNames = skeletonBaseKeys[selectedSkeletonBaseName];
    show = skeletonBaseContainerNames !== undefined && Object.keys(skeletonBaseContainerNames.containers).length > 0;
    if (!show) {
        return <div></div>
    }

    // const dockerImageName = useSelector((state: RootState) => state.clusterBuilder.selectedDockerImageName);
    // const onDockerImageNameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    //     const newDockerImageName = event.target.value;
    //     dispatch(setSelectedDockerImageName(newDockerImageName));
    // };

    return (
        <div>
            <DockerImageCmdArgs />
            <Box mt={2}>
                <TextField
                    fullWidth
                    id="dockerImage"
                    label="Docker Image Name"
                    variant="outlined"
                    // value={dockerImageName}
                    // onChange={onDockerImageNameChange}
                    sx={{ width: '100%' }}
                />
            </Box>
        </div>
    );
}

export function DockerImageCmdArgs() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    if (cluster.componentBases === undefined) {
        return <div></div>
    }
    let show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }

    const skeletonBaseContainerNames = skeletonBaseKeys[selectedSkeletonBaseName];
    show = skeletonBaseContainerNames !== undefined && Object.keys(skeletonBaseContainerNames.containers).length > 0;
    if (!show) {
        return <div></div>
    }

    // const onUpdateDockerCmd = (cmd: string) => {
    //     const input = {
    //         componentBaseKey: selectedComponentBaseName,
    //         skeletonBaseKey: selectedSkeletonBaseName,
    //         containerName: setSelectedContainerName,
    //         cmd: cmd
    //     };
    //     dispatch(setDockerImageCmd(input));
    // };

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
            <AddPortsInputFields />
        </div>
    );
}
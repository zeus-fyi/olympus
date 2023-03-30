import {
    Box,
    Card,
    CardContent,
    Container,
    FormControl,
    FormControlLabel,
    InputLabel,
    MenuItem,
    Select,
    Switch
} from "@mui/material";
import Typography from "@mui/material/Typography";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import {useEffect, useState} from "react";
import TextField from "@mui/material/TextField";
import {AddPortsInputFields} from "./DefinePorts";
import {
    setDockerImage,
    setDockerImageCmd,
    setDockerImageCmdArgs,
    setSelectedContainerName,
    setSelectedDockerImage
} from "../../../../redux/clusters/clusters.builder.reducer";
import {AddVolumeMountsInputFields} from "./DefineVolumeMounts";

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
    const [isInitContainer, setIsInitContainer] = useState(false);
    const label = { inputProps: { label: 'IsInitContainer', 'aria-label': 'IsInitContainer' } };

    const handleClick = () => {
        setIsInitContainer(!isInitContainer);
    };
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
        const containerRef = {
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
        };
        dispatch(setSelectedDockerImage(containerRef));
    };
    return (
        <div>
            {show &&
                <Box mt={2}>
                    <FormControlLabel control={<Switch {...label} defaultChecked={isInitContainer} onClick={handleClick}/>} label={isInitContainer ? 'Init Container [True]' : 'Init Container [False]'} />
                </Box>
            }
            {show &&
                <Box mt={2}>
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
                </Box>

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
    useEffect(() => {
        const containerRef = {
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
        };
        const container = cluster.componentBases[selectedComponentBaseName]?.[selectedSkeletonBaseName]?.containers[selectedContainerName];
        if (!container) {
            return;
        }
        dispatch(setSelectedDockerImage(containerRef));
    }, [dispatch, selectedComponentBaseName, selectedSkeletonBaseName, selectedContainerName, cluster]);

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

    const dockerImageName = skeletonBaseContainerNames?.containers?.[selectedContainerName]?.dockerImage?.imageName;
    const onDockerImageNameChange = (newDockerImageName: string) => {
        const containerRef = {
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            dockerImageKey: newDockerImageName
        };
        dispatch(setDockerImage(containerRef));
        dispatch(setSelectedDockerImage(containerRef));
    };

    return (
        <div>
            <Box mt={2}>
                <TextField
                    fullWidth
                    id="dockerImage"
                    label="Docker Image Name"
                    variant="outlined"
                    value={dockerImageName}
                    onChange={(event) => onDockerImageNameChange(event.target.value as string)}
                    sx={{ width: '100%' }}
                />
            </Box>
            <DockerImageCmdArgs />
            <AddPortsInputFields />
            <AddVolumeMountsInputFields />
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
    const container = skeletonBaseContainerNames.containers[selectedContainerName];
    const cmd = container?.dockerImage?.cmd || [];
    const args = container?.dockerImage?.args || [];
    const onUpdateDockerCmd = (cmd: string) => {
        const input = {
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            cmd: cmd
        };
        dispatch(setDockerImageCmd(input));
    };
    const onUpdateDockerArgs = (args: string) => {
        const input = {
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            args: args
        };
        dispatch(setDockerImageCmdArgs(input));
    };
    return (
        <div>
            <Box mt={2}>
                <TextField
                    fullWidth
                    id="dockerImageCmd"
                    label="Docker Cmd"
                    variant="outlined"
                    onChange={(event) => onUpdateDockerCmd(event.target.value as string)}
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
                    onChange={(event) => onUpdateDockerArgs(event.target.value as string)}
                    sx={{ width: '100%' }}
                />
            </Box>
        </div>
    );
}
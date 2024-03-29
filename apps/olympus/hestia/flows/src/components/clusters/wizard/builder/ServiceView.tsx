import * as React from "react";
import {useMemo} from "react";
import TextField from "@mui/material/TextField";
import {
    Box,
    Card,
    CardContent,
    Container,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    SelectChangeEvent
} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {
    addDockerImagePort,
    removeDockerImagePort,
    setDockerImagePort
} from "../../../../redux/clusters/clusters.builder.reducer";
import {Port} from "../../../../redux/clusters/clusters.types";
import Typography from "@mui/material/Typography";

export function ServiceView(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    const selectedComponentBase = cluster.componentBases?.[selectedComponentBaseName]?.[selectedSkeletonBaseName] ?? '';
    const addDeployment = selectedComponentBase?.addDeployment
    const addStatefulSet = selectedComponentBase?.addStatefulSet
    let selectedDockerImage = cluster.componentBases[selectedComponentBaseName]?.[selectedSkeletonBaseName]?.containers[selectedContainerName]?.dockerImage
    const ports = useMemo(() => {
        const containers = cluster.componentBases[selectedComponentBaseName]?.[selectedSkeletonBaseName]?.containers || {};
        return Object.values(containers).reduce<Port[]>((acc, container) => {
            const dockerImage = container.dockerImage || {};
            const dockerPorts = dockerImage.ports || [{name: "", number: 0, protocol: "TCP", ingressEnabled: false}];
            const filteredPorts = dockerPorts.filter((port) => {
                return port.name !== "" && port.number !== 0;
            });
            return acc.concat(filteredPorts);
        }, []);
    }, [cluster, selectedComponentBaseName, selectedSkeletonBaseName]);

    if (addDeployment === undefined && addStatefulSet === undefined) {
        return <div></div>
    }

    if (addDeployment === false && addStatefulSet === false) {
        return (
            <div>
                <Card sx={{ minWidth: 400, maxWidth: 600}}>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            View Service Ports
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            Add a Deployment or StatefulSet to view service ports. They will be generated from
                            your Docker Images within the Deployment or StatefulSet you created.
                        </Typography>
                    </CardContent>
                </Card>
        </div>
        )
    }
    let workloadType = ''
    if (addDeployment) {
        workloadType = 'Deployment'
    }
    if (addStatefulSet) {
        workloadType = 'StatefulSet'
    }
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

    const handleChange = (index: number, event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const values = [...(selectedDockerImage.ports)];
        values[index] = {...values[index], [event.target.name]: event.target.value};
        dispatch(setDockerImagePort({
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            dockerImageKey: skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
            portIndex: index,
            port: values[index],
        }));
    };

    const handleRemoveField = (index: number) => {
        const values = [...(selectedDockerImage.ports)];
        values.splice(index, 1);
        dispatch(removeDockerImagePort({
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            dockerImageKey: skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
            portIndex: index,
        }));
    };
    const handleChangeSelect = (index: number, event: SelectChangeEvent<string>) => {
        const values = [...(selectedDockerImage.ports)];
        values[index] = { ...values[index], [event.target.name]: event.target.value };
        dispatch(
            setDockerImagePort({
                componentBaseKey: selectedComponentBaseName,
                skeletonBaseKey: selectedSkeletonBaseName,
                containerName: selectedContainerName,
                dockerImageKey:
                skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
                portIndex: index,
                port: values[index],
            })
        );
    };

    const handleAddField = () => {
        const newPort = { name: '', number: 0, protocol: 'TCP', ingressEnabledPort: false}
        dispatch(addDockerImagePort({
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            dockerImageKey: skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
            port: newPort,
        }));
    };
    return (
        <div>
                <Card>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            View Service Ports
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                            These are generated from your docker ports defined in your {workloadType}
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        {ports && ports.map((inputField, index) => (
                            <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <TextField
                                    key={`portName-${index}`}
                                    name="name"
                                    fullWidth
                                    id={`portName-${index}`}
                                    label={`Port Name ${index + 1}`}
                                    variant="outlined"
                                    value={inputField.name}
                                    onChange={(event) => handleChange(index, event)}
                                    sx={{ mr: 1 }}
                                />
                                <TextField
                                    key={`portNumber-${index}`}
                                    name="number"
                                    fullWidth
                                    id={`portNumber-${index}`}
                                    label={`Port Number ${index + 1}`}
                                    variant="outlined"
                                    type="number"
                                    value={inputField.number}
                                    onChange={(event) => handleChange(index, event)}
                                    sx={{ mr: 1 }}
                                    inputProps={{ min: 0 }}
                                />
                                <FormControl fullWidth variant="outlined">
                                    <InputLabel key={`portProtocolLabel-${index}`} id={`portProtocolLabel-${index}`}>Protocol</InputLabel>
                                    <Select
                                        labelId={`portProtocolLabel-${index}`}
                                        id={`portProtocol-${index}`}
                                        name="protocol"
                                        value={inputField.protocol ? inputField.protocol : "TCP"}
                                        onChange={(event) => handleChangeSelect(index, event)}
                                        label="Protocol"
                                    >
                                        <MenuItem value="TCP">TCP</MenuItem>
                                        <MenuItem value="UDP">UDP</MenuItem>
                                    </Select>
                                </FormControl>
                                <Box sx={{ ml: 2 }}>
                                    {/*<Button*/}
                                    {/*    variant="contained"*/}
                                    {/*    onClick={() => handleRemoveField(index)}*/}
                                    {/*>*/}
                                    {/*    Remove*/}
                                    {/*</Button>*/}
                                </Box>
                            </Box>
                        ))}
                        {/*<Button variant="contained" onClick={handleAddField}>*/}
                        {/*    Add Port*/}
                        {/*</Button>*/}
                    </Box>
                    </Container>
                </Card>
        </div>
    );
}
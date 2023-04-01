import * as React from "react";
import {useMemo} from "react";
import TextField from "@mui/material/TextField";
import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select,} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {
    Container as ClustersContainer,
    Port,
    SkeletonBase,
    SkeletonBases
} from "../../../../redux/clusters/clusters.types";
import Typography from "@mui/material/Typography";
import {
    setDockerImagePort,
    setIngressAuthServerURL,
    setIngressHost,
    setIngressPath,
    setIngressPathType
} from "../../../../redux/clusters/clusters.builder.reducer";

export function IngressView(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);

    const ports = useMemo(() => {
        const componentBases = Object.entries(cluster.componentBases);
        const allPorts: {
            [componentBaseName: string]: {
                ports: Port[];
                skeletonBaseName: string;
                componentBaseName: string;
                portIndexToContainer: { [portIndex: number]: { containerName: string, portNumber: number } };
            }
        } = {};
        componentBases.forEach(([componentBaseName, componentBase]: [string, SkeletonBases]) => {
            const skeletonBases = Object.entries(componentBase ?? {});
            skeletonBases.forEach(([skeletonBaseName, skeletonBase]: [string, SkeletonBase]) => {
                if (skeletonBase?.addService){
                    allPorts[componentBaseName] = { ports: [], skeletonBaseName: skeletonBaseName, componentBaseName: componentBaseName, portIndexToContainer: {}};
                    const containers = Object.entries(skeletonBase.containers ?? {})
                    containers.forEach(([containerName, container]: [string,ClustersContainer], containerIndex: number) => {
                        if (container?.dockerImage.imageName !== "") {
                            const dockerPorts = container?.dockerImage?.ports ?? [{ name: "", number: 0, protocol: "TCP", ingressEnabledPort: false }];
                            const filteredPorts = dockerPorts.filter((port) => {
                                return port?.name !== "" && port?.number !== 0;
                            });
                            filteredPorts.forEach((port, portIndex) => {
                                allPorts[componentBaseName].portIndexToContainer[portIndex] = { containerName: containerName, portNumber: port.number };
                            });
                            allPorts[componentBaseName].ports.push(...filteredPorts);
                        }
                    })
                }
            })
        })
        return allPorts;
    }, [cluster]);

    function handleChangeSelect(componentBasePorts: any, selectedPortName: string, selectedPortIndex: number) {
        const containerName= componentBasePorts?.portIndexToContainer[selectedPortIndex].containerName;
        const dockerImage = cluster?.componentBases?.[componentBasePorts?.componentBaseName]?.[componentBasePorts?.skeletonBaseName]?.containers?.[containerName]?.dockerImage;
        if (!dockerImage) {
            console.log('no docker image found');
            return;
        }
        let port = dockerImage.ports[selectedPortIndex];
        const newPort = {name: port.name, number: port.number, protocol: port.protocol, ingressEnabledPort: true};
        dispatch(setDockerImagePort({
            componentBaseKey: componentBasePorts.componentBaseName,
            skeletonBaseKey: componentBasePorts.skeletonBaseName,
            containerName: containerName,
            port: newPort,
            portIndex: selectedPortIndex,
            dockerImageKey: dockerImage.imageName
        }))
    }
    function handleChangeSelectPathType(componentBaseName: string, selectedPathType: string) {
        dispatch(setIngressPathType({componentBaseName: componentBaseName, pathType: selectedPathType}))
    }

    const handleChangePath = (componentBaseName: string, event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const path = event.target.value;
        dispatch(setIngressPath({componentBaseName: componentBaseName, path: path}))
    };

    const handleChangeHost = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const host = event.target.value;
        dispatch(setIngressHost({host: host}))
    };

    const handleChangeAuthURL = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const authServerURL = event.target.value;
        dispatch(setIngressAuthServerURL({authServerURL: authServerURL}))
    };

    return (
        <div>
            <Card>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        View Ingress Paths & Ports
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        Each component base can expose up to one port from a defined service to the Ingress controller.
                        The Ingress controller will then route traffic to the service at port 80.
                    </Typography>
                </CardContent>
                <Container maxWidth="xl" sx={{ mt: 4 }}>
                    <div>
                        <Box mt={2}>
                            <TextField
                                key={`host`}
                                name="host"
                                fullWidth
                                id={`host`}
                                label="Host"
                                variant="outlined"
                                value={cluster.ingressSettings.host}
                                onChange={(event) => handleChangeHost(event)}
                                sx={{ mb: 1 }}
                            />
                            </Box>
                        <Box mt={2}>
                            <TextField
                                key={`authURL`}
                                name="authURL"
                                fullWidth
                                id={`authURL`}
                                label="AuthURL"
                                variant="outlined"
                                value={cluster.ingressSettings.authServerURL}
                                onChange={(event) => handleChangeAuthURL(event)}
                                sx={{ mb: 1 }}
                            />
                        </Box>
                    </div>
                </Container>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box>
                        {ports &&
                            Object.entries(ports).map(([componentBaseName, componentBasePorts]: [string, {}], index) => (
                                <div key={index}>
                                    <Box mt={2} mb={2}>
                                        <Typography variant="body1" color="text.secondary">
                                            Service Component: {componentBaseName}
                                        </Typography>
                                    </Box>
                                        <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                            <FormControl sx={{ mr: 1 }} fullWidth variant="outlined">
                                                <InputLabel key={`portNameLabel-${index}`} id={`portNameLabel-${index}`}>
                                                    Port Name
                                                </InputLabel>
                                                <Select
                                                    labelId={`portNameLabel-${index}`}
                                                    id={`portName-${index}`}
                                                    name="name"
                                                    value={Object.values(ports[componentBaseName].ports).find(port => port.ingressEnabledPort)?.name || ""}
                                                    onChange={(event) => {
                                                        const selectedPortName = event.target.value;
                                                        const selectedPortIndex = Object.values(ports[componentBaseName].ports).findIndex(port => port.name === selectedPortName);
                                                        handleChangeSelect(componentBasePorts, selectedPortName, selectedPortIndex)
                                                    }}
                                                    label="Port Name"
                                                >
                                                    {Object.values(ports[componentBaseName].ports).map((port, portIndex) => (
                                                        <MenuItem selected={port.ingressEnabledPort} key={`menuItem-${port.name}-${portIndex}`} value={port.name}>
                                                            {port.name}
                                                        </MenuItem>
                                                    ))}
                                                </Select>
                                            </FormControl>
                                            <FormControl sx={{ mr: 1 }} fullWidth variant="outlined">
                                                <InputLabel key={`pathTypeLabel-${index}`} id={`pathTypeLabel-${index}`}>
                                                    Path Type
                                                </InputLabel>
                                                <Select
                                                    labelId={`pathTypeLabel-${index}`}
                                                    id={`pathType-${index}`}
                                                    name="pathType"
                                                    value={cluster.ingressPaths[componentBaseName].pathType}
                                                    onChange={(event) => handleChangeSelectPathType(componentBaseName, event.target.value)}
                                                    label="Path Type"
                                                >
                                                    <MenuItem value="ImplementationSpecific">ImplementationSpecific</MenuItem>
                                                    <MenuItem value="Exact">Exact</MenuItem>
                                                </Select>
                                            </FormControl>
                                            <TextField
                                                key={`path-${index}`}
                                                name="path"
                                                fullWidth
                                                id={`path-${index}`}
                                                label="Path"
                                                variant="outlined"
                                                value={cluster.ingressPaths[componentBaseName].path}
                                                onChange={(event) => handleChangePath(componentBaseName, event)}
                                                sx={{ mr: 1 }}
                                            />
                                            <Box sx={{ ml: 2 }}>
                                                {/* <Button */}
                                                {/*   variant="contained" */}
                                                {/*   onClick={() => handleRemoveField(index)} */}
                                                {/* > */}
                                                {/*   Remove */}
                                                {/* </Button> */}
                                            </Box>
                                        </Box>
                                </div>
                            ))}
                        {/* <Button variant="contained" onClick={handleAddField}> */}
                        {/*   Add Port */}
                        {/* </Button> */}
                    </Box>
                </Container>
            </Card>
        </div>
    );
}
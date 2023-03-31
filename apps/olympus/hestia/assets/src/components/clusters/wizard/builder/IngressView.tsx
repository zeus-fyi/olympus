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
import {Container as ClustersContainer, Port, SkeletonBases} from "../../../../redux/clusters/clusters.types";
import Typography from "@mui/material/Typography";

export function IngressView(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const ports = useMemo(() => {
        const componentBases = Object.entries(cluster.componentBases);
        const allPorts: { [componentBaseName: string]: Port[] } = {};
        componentBases.forEach(([componentBaseName, componentBase]: [string, SkeletonBases]) => {
            allPorts[componentBaseName] = [];
            const skeletonBases = Object.values(componentBase ?? {});
            skeletonBases.forEach((skeletonBase) => {
                if (skeletonBase?.addService){
                    const containers = Object.values(skeletonBase.containers ?? {})
                    containers.forEach((container: ClustersContainer) => {
                        const dockerPorts = container?.dockerImage?.ports ?? [{ name: "", number: 0, protocol: "TCP" }];
                        const filteredPorts = dockerPorts.filter((port) => {
                            return port?.name !== "" && port?.number !== 0;
                        });
                        allPorts[componentBaseName].push(...filteredPorts);
                    })
                }
            })
        })
        return allPorts;
    }, [cluster]);
    const handleChangeSelect = (index: number, event: SelectChangeEvent<string>) => {
        // const values = [...(selectedDockerImage.ports)];
        // values[index] = { ...values[index], [event.target.name]: event.target.value };
        //
        console.log(event.target.value);
    };
    const handleChange = (index: number, event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        console.log(event.target.value);
    };

    const handleChangeHost = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        console.log(event.target.value);
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
                    <Box mt={2}>
                        <TextField
                            key={`host`}
                            name="host"
                            fullWidth
                            id={`host`}
                            label="Host"
                            variant="outlined"
                            value="host.zeus.fyi"
                            onChange={(event) => handleChangeHost(event)}
                            sx={{ mb: 1 }}
                        />
                        </Box>
                </Container>
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box>
                        {ports &&
                            Object.entries(ports).map(([componentBaseName, componentBasePorts]: [string, Port[]], index) => (
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
                                                    value={""}
                                                    onChange={(event) => handleChangeSelect(index, event)}
                                                    label="Port Name"
                                                >
                                                    {componentBasePorts.map((port) => (
                                                        <MenuItem key={`menuItem-${port.name}`} value={port.name}>{port.name}</MenuItem>
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
                                                    value={'ImplementationSpecific'}
                                                    onChange={(event) => handleChangeSelect(index, event)}
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
                                                value="/"
                                                onChange={(event) => handleChange(index, event)}
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
import * as React from "react";
import {useMemo} from "react";
import TextField from "@mui/material/TextField";
import {Box, Card, CardContent, Container, FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {Container as ClustersContainer, Port} from "../../../../redux/clusters/clusters.types";
import Typography from "@mui/material/Typography";

export function IngressView(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);

    const ports = useMemo(() => {
        const allPorts: Port[] = [];
        const componentBases = Object.values(cluster.componentBases);
        componentBases.forEach((componentBase) => {
            const skeletonBases = Object.values(componentBase);
            skeletonBases.forEach((skeletonBase) => {
                if (skeletonBase.addService){
                    const containers = Object.values(skeletonBase.containers)
                    containers.forEach((container: ClustersContainer) => {
                        const dockerPorts = container.dockerImage.ports || [{ name: "", number: 0, protocol: "TCP" }];
                        const filteredPorts = dockerPorts.filter((port) => {
                            return port.name !== "" && port.number !== 0;
                        });
                        allPorts.push(...filteredPorts);
                    })
                }
            })
        })
        return allPorts;
    }, [cluster]);

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
                <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                    <Box mt={2}>
                        {ports && ports.map((inputField, index) => (
                            <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                <TextField
                                    key={`path-${index}`}
                                    name="path"
                                    fullWidth
                                    id={`path-${index}`}
                                    label={`Path ${index + 1}`}
                                    variant="outlined"
                                    value={"/"}
                                    sx={{ mr: 1 }}
                                />
                                <FormControl fullWidth variant="outlined">
                                    <InputLabel key={`pathTypeLabel-${index}`} id={`pathTypeLabel-${index}`}>Path Type</InputLabel>
                                    <Select
                                        labelId={`pathTypeLabel-${index}`}
                                        id={`pathType-${index}`}
                                        name="pathType"
                                        value={"ImplementationSpecific"}
                                        // value={inputField.pathType ? inputField.pathType : "ImplementationSpecific"}
                                        //onChange={(event) => handleChangeSelect(index, event)}
                                        label="Path Type"
                                    >
                                        <MenuItem value="ImplementationSpecific">ImplementationSpecific</MenuItem>
                                        <MenuItem value="Exact">Exact</MenuItem>
                                    </Select>
                                </FormControl>
                                <TextField
                                    key={`portName-${index}`}
                                    name="name"
                                    fullWidth
                                    id={`portName-${index}`}
                                    label={`Port Name ${index + 1}`}
                                    variant="outlined"
                                    value={inputField.name}
                                    sx={{ ml: 1 }}
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
                                    sx={{ mr: 1, ml: 1 }}
                                    inputProps={{ min: 0 }}
                                />
                                <FormControl fullWidth variant="outlined">
                                <InputLabel key={`portProtocolLabel-${index}`} id={`portProtocolLabel-${index}`}>Protocol</InputLabel>
                                    <Select
                                        labelId={`portProtocolLabel-${index}`}
                                        id={`portProtocol-${index}`}
                                        name="protocol"
                                        value={inputField.protocol ? inputField.protocol : "TCP"}
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
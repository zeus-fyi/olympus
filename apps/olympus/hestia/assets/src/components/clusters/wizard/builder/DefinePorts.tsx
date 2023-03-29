import * as React from "react";
import {useEffect} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {
    addDockerImagePort,
    removeDockerImagePort,
    setDockerImagePort,
    setSelectedDockerImage
} from "../../../../redux/clusters/clusters.builder.reducer";

export function AddPortsInputFields() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    let selectedDockerImage = useSelector((state: RootState) => state.clusterBuilder.selectedDockerImage);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    const ports = selectedDockerImage.ports || [{name: "", port: "", protocol: "TCP"}];
    useEffect(() => {
        const containerRef = {
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
        };
        dispatch(setSelectedDockerImage(containerRef));

    }, [dispatch, selectedComponentBaseName, selectedSkeletonBaseName, selectedContainerName, cluster, selectedDockerImage, ports]);

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
        const newPort = { name: '', number: 0, protocol: 'TCP' };
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
                            defaultValue=""
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
                            defaultValue={0}
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
                                defaultValue="TCP"
                            >
                                <MenuItem value="TCP">TCP</MenuItem>
                                <MenuItem value="UDP">UDP</MenuItem>
                            </Select>
                        </FormControl>
                        <Box sx={{ ml: 2 }}>
                            <Button
                                variant="contained"
                                onClick={() => handleRemoveField(index)}
                            >
                                Remove
                            </Button>
                        </Box>
                    </Box>
                ))}
                <Button variant="contained" onClick={handleAddField}>
                    Add Item
                </Button>

            </Box>
        </div>
    );
}
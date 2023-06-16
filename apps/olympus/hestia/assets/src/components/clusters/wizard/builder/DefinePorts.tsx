import * as React from "react";
import {useMemo, useState} from "react";
import TextField from "@mui/material/TextField";
import debounce from 'lodash/debounce';
import Button from "@mui/material/Button";
import {Box, FormControl, FormHelperText, InputLabel, MenuItem, Select, SelectChangeEvent, Stack} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {
    addDockerImagePort,
    removeDockerImagePort,
    setDockerImagePort
} from "../../../../redux/clusters/clusters.builder.reducer";

export function AddPortsInputFields() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    let selectedDockerImage = cluster.componentBases[selectedComponentBaseName]?.[selectedSkeletonBaseName]?.containers[selectedContainerName]?.dockerImage
    const ports = useMemo(() => selectedDockerImage?.ports || [{name: "", number: "", protocol: "TCP"}], [selectedDockerImage?.ports]);
    const [error, setError] = useState('');

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
    const debouncedValidation = debounce((value: string) => {
        if (!isValidPortName(value)) {
            setError('Invalid port name. http is a reserved port name for nginx'); // set the error message
            return;
        } else {
            setError(''); // clear the error
        }
    }, 100); // 500ms delay

    const handleChange = (index: number, event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        debouncedValidation(event.target.value);
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
        const newPort = { name: '', number: 0, protocol: 'TCP', ingressEnabledPort: false };
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
                            <Button
                                variant="contained"
                                onClick={() => handleRemoveField(index)}
                            >
                                Remove
                            </Button>
                        </Box>
                    </Box>
                ))}
                <Stack spacing={2} direction="column" sx={{ flex: 1, mr: 2 }}>
                    <div>
                        {error && <FormHelperText error>{error}</FormHelperText>}
                    </div>
                </Stack>
                <Button variant="contained" onClick={handleAddField}>
                    Add Port
                </Button>
            </Box>
        </div>
    );
}

function isValidPortName(name: string): boolean {
    return name !== 'http';
}
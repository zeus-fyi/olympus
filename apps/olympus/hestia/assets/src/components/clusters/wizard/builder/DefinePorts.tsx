import * as React from "react";
import {useState} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {Box, FormControl, InputLabel, MenuItem, Select, SelectChangeEvent} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {setDockerImagePort} from "../../../../redux/clusters/clusters.builder.reducer";

export function AddPortsInputFields() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    const [inputFields, setInputFields] = useState<{ name: string; number: number; protocol: string }[]>([{ name: '', number: 0, protocol: 'TCP' }]);

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
        const values = [...inputFields];
        values[index] = {...values[index], [event.target.name]: event.target.value};
        //console.log("values: " + values[index].name + " " + values[index].number + " " + values[index].protocol)
        setInputFields(values);
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
        const values = [...inputFields];
        values.splice(index, 1);
        setInputFields(values);
    };
    const handleChangeSelect = (index: number, event: SelectChangeEvent<string>) => {
        const values = [...inputFields];
        values[index] = { ...values[index], [event.target.name]: event.target.value };
        setInputFields(values);
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
        setInputFields([...inputFields, { name: '', number: 0, protocol: 'TCP'}]);
    };

    return (
        <div>
            <Box mt={2}>
                {inputFields.map((inputField, index) => (
                    <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                        <TextField
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
                            name="number"
                            fullWidth
                            id={`portNumber-${index}`}
                            label={`Port Number ${index + 1}`}
                            variant="outlined"
                            type="number"
                            value={inputField.number}
                            onChange={(event) => handleChange(index, event)}
                            sx={{ mr: 1 }}
                        />
                        <FormControl fullWidth variant="outlined">
                            <InputLabel id={`portProtocolLabel-${index}`}>Protocol</InputLabel>
                            <Select
                                labelId={`portProtocolLabel-${index}`}
                                id={`portProtocol-${index}`}
                                name="protocol"
                                value={inputField.protocol}
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
                <Button variant="contained" onClick={handleAddField}>
                    Add Item
                </Button>

            </Box>
        </div>
    );
}
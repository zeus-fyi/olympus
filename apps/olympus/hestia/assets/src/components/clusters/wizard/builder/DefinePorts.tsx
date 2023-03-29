import * as React from "react";
import {useState} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {Box} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";

export function AddPortsInputFields() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    const [inputFields, setInputFields] = useState<{ name: string; number: number; protocol: string }[]>([{ name: '', number: 0, protocol: '' }]);

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
        // dispatch(setDockerImagePort({
        //     componentBaseKey: selectedComponentBaseName,
        //     skeletonBaseKey: selectedSkeletonBaseName,
        //     containerName: selectedContainerName,
        //     dockerImageKey: skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
        //     portIndex: index,
        //     port: values[index],
        // }));
    };


    const handleAddField = () => {
        setInputFields([...inputFields, { name: '', number: 0, protocol: ''}]);
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
                        <TextField
                            name="protocol"
                            fullWidth
                            id={`portProtocol-${index}`}
                            label={`Protocol`}
                            variant="outlined"
                            value={inputField.protocol}
                            onChange={(event) => handleChange(index, event)}
                        />
                    </Box>
                ))}
                <Button variant="contained" onClick={handleAddField}>
                    Add Item
                </Button>
            </Box>
        </div>
    );
}
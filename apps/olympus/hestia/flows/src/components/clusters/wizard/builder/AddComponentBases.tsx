import {useState} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {
    addComponentBase,
    removeComponentBase,
    setSelectedComponentBaseName
} from "../../../../redux/clusters/clusters.builder.reducer";
import Box from "@mui/material/Box";
import {FormHelperText, Stack} from "@mui/material";

export function AddComponentBases() {
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBases = useSelector((state: RootState) => state.clusterBuilder.cluster.componentBases);
    const componentBaseKeys = Object.keys(componentBases);
    const dispatch = useDispatch();
    const [inputField, setInputField] = useState('');
    const [error, setError] = useState('');

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputField(event.target.value);
        setError('');
    };

    const handleAddField = () => {
        try {
            if (inputField) {
                if (!isValidLabel(inputField)) {
                    setError('Invalid label. A valid label must not be an empty string have spaces between letters or consist of alphanumeric characters, "-", "_" or ".", and must start and end with an alphanumeric character'); // set the error message
                    return;
                }
                const cb = { componentBaseName: inputField, skeletonBases: {} }
                dispatch(addComponentBase(cb));
                dispatch(setSelectedComponentBaseName(inputField));
                setInputField('');
                setError(''); // clear the error
            }
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError(error.message); // set the error message
            } else {
            }
        }
    };
    const handleRemoveField = (key: string) => {
        const componentBasesCopy = { ...componentBases };
        delete componentBasesCopy[key];
        dispatch(removeComponentBase(key));
        if (cluster.componentBases[key] !== undefined && Object.keys(cluster.componentBases).length > 0) {
            if (Object.keys(cluster.componentBases)[0] === key) {
                dispatch(setSelectedComponentBaseName(Object.keys(cluster.componentBases)[1]));
            } else {
                dispatch(setSelectedComponentBaseName(Object.keys(cluster.componentBases)[0]));
            }
        }
    };
    return (
        <div>
            {componentBaseKeys.map((key, index) => (
                <Box key={index} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                    <TextField
                        fullWidth
                        id={`inputField-${index}`}
                        label={`Cluster Base Name`}
                        variant="outlined"
                        value={key}
                        InputProps={{ readOnly: true }}
                        sx={{ flex: 1, mr: 2, mb: 1 }}
                    />
                    <Button variant="contained" sx={{ width: '100px' }} onClick={() => handleRemoveField(key)}>
                        Remove
                    </Button>
                </Box>))
            }
            <Box key={componentBaseKeys.length} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
            <TextField
                fullWidth
                id="inputField-new"
                label="New Cluster Base Name"
                variant="outlined"
                value={inputField}
                onChange={handleChange}
                sx={{ flex: 1, mr: 2 }}
            />
                <Stack spacing={2} direction="column" sx={{ flex: 1, mr: 2 }}>
                    <div>
                        {error && <FormHelperText error>{error}</FormHelperText>}
                    </div>
                    <Button variant="contained" sx={{ width: '100px' }} onClick={handleAddField}>
                        Add
                    </Button>
                </Stack>
            </Box>
        </div>
    )
}

// https://github.com/kubernetes/apimachinery/blob/master/pkg/util/validation/validation.go
export function isValidLabel(label: string): boolean {
    // Check if the label is empty
    if (label.length === 0) {
        return false;
    }
    // Check for length
    if (label.length > 63) {
        return false;
    }
    // Check if label starts or ends with alphanumeric characters
    if (!label.match(/^[a-z0-9].*[a-z0-9]$/i)) {
        return false;
    }
    // Check if label contains only the allowed characters
    if (!label.match(/^[a-z0-9\-\_\.]*$/i)) {
        return false;
    }
    return true;
}
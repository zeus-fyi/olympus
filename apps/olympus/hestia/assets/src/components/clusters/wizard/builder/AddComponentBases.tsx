import {useState} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {addComponentBase, removeComponentBase} from "../../../../redux/clusters/clusters.builder.reducer";
import Box from "@mui/material/Box";

export function AddComponentBases() {
    const componentBases = useSelector((state: RootState) => state.clusterBuilder.cluster.componentBases);
    const componentBaseKeys = Object.keys(componentBases);
    const dispatch = useDispatch();
    const [inputField, setInputField] = useState('');

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputField(event.target.value);
    };

    const handleAddField = () => {
        if (inputField) {
            dispatch(addComponentBase({ componentBaseName: inputField, skeletonBases: {} }));
            setInputField('');
        }
    };

    const handleRemoveField = (key: string) => {
        dispatch(removeComponentBase(key));
    };

    return (
        <div>
            {componentBaseKeys.map((key, index) => (
                <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                    <TextField
                        fullWidth
                        id={`inputField-${index}`}
                        label={`Component Base Name`}
                        variant="outlined"
                        value={key}
                        InputProps={{ readOnly: true }}
                        sx={{ flex: 1, mr: 2 }}
                    />
                    <Button variant="contained" sx={{ width: '100px' }} onClick={() => handleRemoveField(key)}>
                        Remove
                    </Button>
                </Box>))
            }
            <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
            <TextField
                fullWidth
                id="inputField-new"
                label="New Component Base Name"
                variant="outlined"
                value={inputField}
                onChange={handleChange}
                sx={{ flex: 1, mr: 2 }}
            />
            <Button variant="contained" sx={{ width: '100px' }} onClick={handleAddField}>
                Add
            </Button>
            </Box>
        </div>
    )
}

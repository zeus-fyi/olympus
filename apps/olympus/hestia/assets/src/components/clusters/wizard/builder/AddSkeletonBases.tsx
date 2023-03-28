import {useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {addSkeletonBase, removeSkeletonBase} from "../../../../redux/clusters/clusters.builder.reducer";
import {Box} from "@mui/material";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {RootState} from "../../../../redux/store";

export function AddSkeletonBases(props: any) {
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const dispatch = useDispatch();
    const componentBase = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBase);
    const selectedComponentBaseSkeletonBasesKeys = Object.keys(componentBase);

    const [inputField, setInputField] = useState('');
    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputField(event.target.value);
    };
    const handleAddField = () => {
        if (inputField) {
            let sb = {  dockerImages: {},  };
            dispatch(addSkeletonBase({ componentBaseName: componentBaseName, skeletonBaseName: inputField, skeletonBase: sb }));
            setInputField('');
        }
    };
    const handleRemoveField = (skeletonBaseName: string) => {
        dispatch(removeSkeletonBase({componentBaseName: componentBaseName, skeletonBaseName: skeletonBaseName}));
    };
    return (
        <div>
            {selectedComponentBaseSkeletonBasesKeys.map((key, index) => (
                <Box key={index} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                    <TextField
                        fullWidth
                        id={`inputField-${index}`}
                        label={`Skeleton Base Name`}
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
            <Box key={selectedComponentBaseSkeletonBasesKeys.length} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                <TextField
                    fullWidth
                    id="inputField-new"
                    label="New Skeleton Base Name"
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

import {useState} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {Box} from "@mui/material";

export function AddPortsInputFields() {
    const [inputFields, setInputFields] = useState([{ value: '' }]);

    const handleChange = (index: any, event: any) => {
        const values = [...inputFields];
        values[index].value = event.target.value;
        setInputFields(values);
    };

    const handleAddField = () => {
        setInputFields([...inputFields, { value: '' }]);
    };

    return (
        <div>
            <Box mt={2}>
                {inputFields.map((inputField, index) => (
                    <TextField
                        key={index}
                        fullWidth
                        id={`inputField-${index}`}
                        label={`Port Number ${index + 1}`}
                        variant="outlined"
                        value={inputField.value}
                        onChange={(event) => handleChange(index, event)}
                        sx={{ mb: 1 }}
                    />
                ))}
                <Button variant="contained" onClick={handleAddField}>
                    Add Item
                </Button>
            </Box>
        </div>
    );
}
import {useState} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";

export function AddComponentBases() {
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
            {inputFields.map((inputField, index) => (
                <TextField
                    key={index}
                    fullWidth
                    id={`inputField-${index}`}
                    label={`Component Base Name`}
                    variant="outlined"
                    value={inputField.value}
                    onChange={(event) => handleChange(index, event)}
                    sx={{ mb: 1 }}
                />
            ))}
            <Button variant="contained" onClick={handleAddField}>
                Add Item
            </Button>
        </div>
    );
}
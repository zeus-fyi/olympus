import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import * as React from "react";
import {useState} from "react";
import {removeConfigMapKey, setConfigMapKey,} from "../../../../redux/clusters/clusters.builder.reducer";
import {Box, Card, CardContent, Container, InputLabel, TextareaAutosize} from "@mui/material";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";


export function ConfigMapView(props: any) {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    let selectedConfigMap = cluster.componentBases[selectedComponentBaseName]?.[selectedSkeletonBaseName].configMap
    const [newConfigMapEntry, setNewConfigMapEntry] = useState<{ key: string; value: string }>({ key: "", value: "" });

    const handleNewEntryChange = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>, field: "key" | "value") => {
        setNewConfigMapEntry({ ...newConfigMapEntry, [field]: event.target.value });
    };

    const handleAddField = () => {
        if (newConfigMapEntry.key.trim() === "") {
            alert("Please enter a valid key.");
            return;
        }
        const input = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            key: newConfigMapEntry.key,
            value: newConfigMapEntry.value
        }
        dispatch(setConfigMapKey(input));
        setNewConfigMapEntry({ key: "", value: "" }); // Reset new entry input fields
    };
    if (cluster.componentBases === undefined) {
        return <div></div>
    }
    let show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }

    const handleChange = (key: string, event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>, isKeyField: boolean) => {
        const target = event.target;
        const input = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            key: key,
            value: target.value,
            targetKey: isKeyField ? target.value : undefined
        };
        dispatch(setConfigMapKey(input));
    };

    const handleRemoveField = (key: string) => {
        const input = {
            componentBaseName: selectedComponentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            key: key,
        }
        dispatch(removeConfigMapKey(input));
    };
    return (
        <div>
            <Card sx={{ minWidth: 600 }}>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Configure ConfigMap
                    </Typography>
                </CardContent>
                <Container>
                <Box mt={2}>
                    {Object.entries(selectedConfigMap).map(([key, value], index) => (
                        <Box key={index} sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                            <TextField
                                key={`configMapKey-${index}`}
                                name="configMapKey"
                                fullWidth
                                id={`configMapKey-${index}`}
                                label={`Config Map Key ${index + 1}`}
                                variant="outlined"
                                value={key}
                                onChange={(event) => handleChange(key, event, true)}
                                sx={{ mt: 2, mr: 1 }}
                            />
                            <Box sx={{ mt: 2, mr: 1, flexGrow: 1 }}>
                                <InputLabel htmlFor={`configMapValue-${index}`}>{`Config Map Value ${index + 1}`}</InputLabel>
                                <TextareaAutosize
                                    key={`configMapValue-${index}`}
                                    id={`configMapValue-${index}`}
                                    minRows={3}
                                    value={value}
                                    onChange={(event) => handleChange(key, event, false)}
                                    style={{ resize: "both", width: "100%" }}
                                />
                            </Box>
                            <Box sx={{ mt: 2 }}>
                                <Button
                                    variant="contained"
                                    onClick={() => handleRemoveField(key)}
                                >
                                    Remove
                                </Button>
                            </Box>
                        </Box>
                    ))}
                    <TextField
                        fullWidth
                        id="newConfigMapKey"
                        label="New Config Map Key"
                        variant="outlined"
                        value={newConfigMapEntry.key}
                        onChange={(event) => handleNewEntryChange(event, "key")}
                        sx={{ mt: 2 }}
                    />
                    <Box sx={{ mt: 2, flexGrow: 1 }}>
                        <InputLabel htmlFor="newConfigMapValue">New Config Map Value</InputLabel>
                        <TextareaAutosize
                            id="newConfigMapValue"
                            minRows={3}
                            value={newConfigMapEntry.value}
                            onChange={(event) => handleNewEntryChange(event, "value")}
                            style={{ resize: "both", width: "100%" }}
                        />
                    </Box>
                    <Box sx={{ mt: 2, mb: 4 }}>
                        <Button
                            variant="contained"
                            onClick={handleAddField}
                        >
                            Add ConfigMap Key
                        </Button>
                    </Box>
                </Box>
                </ Container>

            </Card>
        </div>
    );
}
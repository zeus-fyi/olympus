import {useDispatch, useSelector} from "react-redux";
import * as React from "react";
import {useEffect} from "react";
import {
    Box,
    Card,
    CardContent,
    Container,
    FormControl,
    InputLabel,
    MenuItem,
    Select,
    SelectChangeEvent,
    Stack
} from "@mui/material";
import {RootState} from "../../../../redux/store";
import Typography from "@mui/material/Typography";
import {DefineDockerParams} from "./DefineDockerImage";
import {
    addStatefulSetPVC,
    removeStatefulSetPVC,
    setDeploymentReplicaCount,
    setSelectedComponentBaseName,
    setSelectedContainerName,
    setSelectedSkeletonBaseName,
    setStatefulSetPVC,
    setStatefulSetReplicaCount
} from "../../../../redux/clusters/clusters.builder.reducer";
import {AddContainers} from "./AddContainers";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";

export function AddSkeletonBaseDockerConfigs(props: any) {
    const {viewField} = props;

    let viewFieldName = '';
    if (viewField === 'statefulSet') {
        viewFieldName = 'StatefulSet';
    }
    if (viewField === 'deployment') {
        viewFieldName = 'Deployment';
    }
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);

    if (cluster.componentBases === undefined) {
        return <div></div>
    }
    const skeletonBaseKeys = cluster.componentBases[componentBaseName];
    const show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }
    const onChangeStatefulSetReplicaCount = (replicaCount: number) => {
        const stsObjRef = {componentBaseName: componentBaseName, skeletonBaseName: selectedSkeletonBaseName, replicaCount: replicaCount};
        dispatch(setStatefulSetReplicaCount(stsObjRef));
    };
    const onChangeDeploymentReplicaCount = (replicaCount: number) => {
        const stsObjRef = {componentBaseName: componentBaseName, skeletonBaseName: selectedSkeletonBaseName, replicaCount: replicaCount};
        dispatch(setDeploymentReplicaCount(stsObjRef));
    };

    const handleChange = (index: number, event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const values = [...(cluster.componentBases[componentBaseName][selectedSkeletonBaseName].statefulSet.pvcTemplates)];
        values[index] = {...values[index], [event.target.name]: event.target.value};

        let newValues = values[index];
        if (newValues.accessMode == '') {
            newValues.accessMode = 'ReadWriteOnce';
        }
        dispatch(setStatefulSetPVC({
            componentBaseName: componentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            pvcIndex: index,
            pvc: newValues,
        }));
    };

    const handleRemoveField = (index: number) => {
        const values = [...(cluster.componentBases[componentBaseName][selectedSkeletonBaseName].statefulSet.pvcTemplates)];
        const pvc = values[index]
        values.splice(index, 1);
        dispatch(removeStatefulSetPVC({
            componentBaseName: componentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            pvc: pvc,
            pvcIndex: index,
        }));
    };
    const handleChangeSelect = (index: number, event: SelectChangeEvent<string>) => {
        const values = [...(cluster.componentBases[componentBaseName][selectedSkeletonBaseName].statefulSet.pvcTemplates)];
        values[index] = { ...values[index], [event.target.name]: event.target.value };
        dispatch(setStatefulSetPVC({
            componentBaseName: componentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            pvcIndex: index,
            pvc: values[index],
        }));
    };

    const handleAddField = () => {
        const newPVC = { name: '', accessMode: 'ReadWriteOnce', storageSizeRequest: '' };
        dispatch(addStatefulSetPVC({
            componentBaseName: componentBaseName,
            skeletonBaseName: selectedSkeletonBaseName,
            pvc: newPVC,
        }));
    };
    const showPVCs = Object.keys(cluster.componentBases[componentBaseName][selectedSkeletonBaseName].statefulSet?.pvcTemplates).length > 0;
    return (
        <div>
            <Stack direction="row" spacing={2}>
                <Card sx={{ maxWidth: 500 }}>
                    <CardContent>
                        <Typography gutterBottom variant="h5" component="div">
                            Configure {viewFieldName} Workload
                        </Typography>
                    </CardContent>
                    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                        <div>
                            {viewField === 'deployment' &&
                            <Box mt={2}>
                                <TextField
                                    fullWidth
                                    id={"replicaCount"+viewFieldName}
                                    label={"Replica Count"}
                                    variant="outlined"
                                    type={"number"}
                                    InputProps={{ inputProps: { min: 0 } }}
                                    value={cluster.componentBases[componentBaseName][selectedSkeletonBaseName].deployment.replicaCount}
                                    onChange={(event) => onChangeDeploymentReplicaCount(parseInt(event.target.value))}
                                    sx={{ width: '100%' }}
                                />
                            </Box>
                            }
                            {viewField === 'statefulSet' &&
                                <div>
                                    <Box mt={2}>
                                        <TextField
                                            fullWidth
                                            id={"replicaCount"+viewFieldName}
                                            label={"Replica Count"}
                                            variant="outlined"
                                            type={"number"}
                                            InputProps={{ inputProps: { min: 0 } }}
                                            value={cluster.componentBases[componentBaseName][selectedSkeletonBaseName].statefulSet.replicaCount}
                                            onChange={(event) => onChangeStatefulSetReplicaCount(parseInt(event.target.value))}
                                            sx={{ width: '100%' }}
                                        />
                                    </Box>
                                    <Box mt={2}>
                                        {showPVCs && cluster.componentBases[componentBaseName][selectedSkeletonBaseName].statefulSet.pvcTemplates.map((inputField, index) => (
                                            <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                                                <TextField
                                                    key={`name-${index}`}
                                                    name="name"
                                                    fullWidth
                                                    id={`name-${index}`}
                                                    label={`PVC Name ${index + 1}`}
                                                    variant="outlined"
                                                    value={inputField.name}
                                                    onChange={(event) => handleChange(index, event)}
                                                    sx={{ mr: 1 }}
                                                />
                                                <TextField
                                                    key={`storageSizeRequest-${index}`}
                                                    name="storageSizeRequest"
                                                    fullWidth
                                                    id={`storageSizeRequest-${index}`}
                                                    label={`Storage Size Request Number ${index + 1}`}
                                                    variant="outlined"
                                                    value={inputField.storageSizeRequest}
                                                    onChange={(event) => handleChange(index, event)}
                                                    sx={{ mr: 1 }}
                                                    inputProps={{ min: 0 }}
                                                />
                                                <FormControl fullWidth variant="outlined">
                                                    <InputLabel key={`accessModeLabel-${index}`} id={`accessModeLabel-${index}`}>Access Mode</InputLabel>
                                                    <Select
                                                        labelId={`accessModeLabel-${index}`}
                                                        id={`accessMode-${index}`}
                                                        name="accessMode"
                                                        value={inputField.accessMode ? inputField.accessMode : "ReadWriteOnce"}
                                                        onChange={(event) => handleChangeSelect(index, event)}
                                                        label="Access Mode"
                                                    >
                                                        <MenuItem value="ReadWriteOnce">ReadWriteOnce</MenuItem>
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
                                            Add PVC
                                        </Button>
                                    </Box>
                                </div>
                            }
                        {show && cluster.componentBases[componentBaseName] && Object.keys(skeletonBaseKeys).length > 0 &&
                            <Box mt={2}>
                                <AddContainers />
                            </Box>
                        }
                        </div>

                    </Container>
                </Card>
                {show && cluster.componentBases[componentBaseName] && Object.keys(skeletonBaseKeys).length > 0 &&
                    <div>
                        <DefineDockerParams />
                    </div>
                }
            </Stack>
        </div>
    )
}

export function SelectedSkeletonBaseName(props: any) {
    const { onChangeComponentOrSkeletonBase}  = props;
    const dispatch = useDispatch();
    const skeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const componentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);

    useEffect(() => {
        dispatch(setSelectedComponentBaseName(componentBaseName));
        dispatch(setSelectedSkeletonBaseName(skeletonBaseName));
    }, [dispatch,skeletonBaseName, componentBaseName]);

    const onAccessSkeletonBase = (selectedSkeletonBaseName: string) => {
        dispatch(setSelectedSkeletonBaseName(selectedSkeletonBaseName));
        const containerKeys = Object.keys(cluster.componentBases[componentBaseName][selectedSkeletonBaseName].containers)
        if (containerKeys.length > 0) {
            dispatch(setSelectedContainerName(containerKeys[0]));
        }
        onChangeComponentOrSkeletonBase();
    };

    if (cluster.componentBases === undefined) {
        return <div></div>
    }

    const skeletonBaseKeys = cluster.componentBases[componentBaseName];
    const show = skeletonBaseKeys !== undefined && Object.keys(skeletonBaseKeys).length > 0;
    if (!show) {
        return <div></div>
    }
    return (
        <FormControl variant="outlined" style={{ minWidth: '100%' }}>
            <InputLabel id="network-label">Workload Bases</InputLabel>
            <Select
                labelId="skeletonBase-label"
                id="skeletonBase"
                value={skeletonBaseName}
                label="Skeleton Base"
                onChange={(event) => onAccessSkeletonBase(event.target.value as string)}
                sx={{ width: '100%' }}
            >
                {show && Object.keys(skeletonBaseKeys).map((key: any, i: number) => (
                    <MenuItem key={i} value={key}>
                        {key}
                    </MenuItem>))
                }
            </Select>
        </FormControl>
    );
}
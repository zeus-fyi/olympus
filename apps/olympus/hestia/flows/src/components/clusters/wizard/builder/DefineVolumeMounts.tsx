import * as React from "react";
import {useEffect, useMemo} from "react";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {Box} from "@mui/material";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {
    addDockerImageVolumeMount,
    removeDockerImageVolumeMount,
    setDockerImageVolumeMount,
    setSelectedDockerImage
} from "../../../../redux/clusters/clusters.builder.reducer";

export function AddVolumeMountsInputFields() {
    const dispatch = useDispatch();
    const cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    const selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const selectedContainerName = useSelector((state: RootState) => state.clusterBuilder.selectedContainerName);
    let selectedDockerImage = useSelector((state: RootState) => state.clusterBuilder.selectedDockerImage);
    const skeletonBaseKeys = cluster.componentBases[selectedComponentBaseName];
    const volumeMounts = useMemo(() => selectedDockerImage.volumeMounts || [{name: "", mountPath: ""}], [selectedDockerImage.volumeMounts]);

    useEffect(() => {
        const containerRef = {
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
        };
        const container = cluster.componentBases[selectedComponentBaseName]?.[selectedSkeletonBaseName]?.containers[selectedContainerName];
        if (!container) {
            return;
        }
        dispatch(setSelectedDockerImage(containerRef));
    }, [dispatch, selectedComponentBaseName, selectedSkeletonBaseName, selectedContainerName, cluster, selectedDockerImage]);

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
        const values = [...(selectedDockerImage.volumeMounts)];
        values[index] = {...values[index], [event.target.name]: event.target.value};
        dispatch(setDockerImageVolumeMount({
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            dockerImageKey: skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
            volumeMountIndex: index,
            volumeMount: values[index],
        }));
    };

    const handleRemoveField = (index: number) => {
        const values = [...(selectedDockerImage.volumeMounts)];
        values.splice(index, 1);
        dispatch(removeDockerImageVolumeMount({
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            dockerImageKey: skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
            volumeMountIndex: index,
        }));
    };

    const handleAddField = () => {
        const newVolumeMount = { name: '', mountPath: '' };
        dispatch(addDockerImageVolumeMount({
            componentBaseKey: selectedComponentBaseName,
            skeletonBaseKey: selectedSkeletonBaseName,
            containerName: selectedContainerName,
            dockerImageKey: skeletonBaseContainerNames.containers[selectedContainerName].dockerImage.imageName,
            volumeMount: newVolumeMount,
        }));
    };

    return (
        <div>
            <Box mt={2}>
                {volumeMounts && volumeMounts.map((inputField, index) => (
                    <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <TextField
                    key={`volumeMountName-${index}`}
                    name="name"
                    fullWidth
                    id={`volumeMountName-${index}`}
                    label={`Volume Mount Name ${index + 1}`}
                    variant="outlined"
                    value={inputField.name}
                    onChange={(event) => handleChange(index, event)}
                    sx={{ mr: 1 }}
                />
                <TextField
                    key={`volumeMountPath-${index}`}
                    name="mountPath"
                    fullWidth
                    id={`volumeMountPath-${index}`}
                    label={`Volume Mount Path ${index + 1}`}
                    variant="outlined"
                    value={inputField.mountPath}
                    onChange={(event) => handleChange(index, event)}
                    sx={{ mr: 1 }}
                />
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
                Add Volume Mount
            </Button>
        </Box>
        </div>
    );
}
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../../../redux/store";
import {useState} from "react";
import {
    addContainer,
    removeContainer,
    setSelectedContainerName,
} from "../../../../redux/clusters/clusters.builder.reducer";
import {Box} from "@mui/material";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {Container, Port, VolumeMount} from "../../../../redux/clusters/clusters.types";

export function AddContainers(props: any) {
    const dispatch = useDispatch();
    let cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    let selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    let componentBase = cluster.componentBases[selectedComponentBaseName];
    const selectedSkeletonBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedSkeletonBaseName);
    const [inputField, setInputField] = useState('');

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputField(event.target.value);
    };
    const handleAddField = () => {
        if (inputField) {
            const cont = {
                dockerImage:
                    {
                        imageName: '',
                        args: '',
                        cmd: '',
                        resourceRequirements: {cpu: '', memory: ''},
                        volumeMounts: [{name: '', mountPath: ''}] as VolumeMount[],
                        ports: [{name: '', number: 0, protocol: 'TCP', ingressEnabledPort: false} as Port] as Port[]}
                    }
            let contObj = {
                componentBaseKey: selectedComponentBaseName,
                skeletonBaseKey: selectedSkeletonBaseName,
                containerName: inputField,
                container: cont as Container,
            }
            dispatch(addContainer(contObj));
            dispatch(setSelectedContainerName(inputField))
            setInputField('');
        }
    };
    const handleRemoveField = (containerName: string) => {
        dispatch(removeContainer({componentBaseName: selectedComponentBaseName, skeletonBaseName: selectedSkeletonBaseName, containerName: containerName}));
        if (cluster.componentBases[selectedComponentBaseName] !== undefined && Object.keys(cluster.componentBases[selectedComponentBaseName]).length > 0) {
            if (Object.keys(cluster.componentBases[selectedComponentBaseName][selectedSkeletonBaseName].containers)[0] === containerName) {
                dispatch(setSelectedContainerName(Object.keys(cluster.componentBases[selectedComponentBaseName][selectedSkeletonBaseName].containers)[1]));
            } else {
                dispatch(setSelectedContainerName(Object.keys(cluster.componentBases[selectedComponentBaseName][selectedSkeletonBaseName].containers)[0]));
            }
        }
    };
    let showAdd = componentBase !== undefined;
    let show = showAdd && Object.keys(componentBase).length > 0;
    const skeletonBaseContainerNames = cluster.componentBases[selectedComponentBaseName][selectedSkeletonBaseName];
    return (
        <div>
            { show && Object.keys(skeletonBaseContainerNames.containers).map((key, index) => (
                <Box key={index} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                    <TextField
                        fullWidth
                        id={`inputField-${index}`}
                        label={`Container Name`}
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
            { showAdd &&
                <Box key={Object.keys(skeletonBaseContainerNames.containers).length} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                    <TextField
                        fullWidth
                        id="inputField-new"
                        label="New Container Name"
                        variant="outlined"
                        value={inputField}
                        onChange={handleChange}
                        sx={{ flex: 1, mr: 2 }}
                    />
                    <Button variant="contained" sx={{ width: '100px' }} onClick={handleAddField}>
                        Add
                    </Button>
                </Box>
            }
        </div>
    )
}
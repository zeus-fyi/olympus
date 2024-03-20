import {useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {
    addSkeletonBase,
    removeSkeletonBase,
    setSelectedComponentBaseName,
    setSelectedSkeletonBaseName
} from "../../../../redux/clusters/clusters.builder.reducer";
import {Box} from "@mui/material";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import {RootState} from "../../../../redux/store";

export function AddSkeletonBases(props: any) {
    const dispatch = useDispatch();
    let cluster = useSelector((state: RootState) => state.clusterBuilder.cluster);
    let selectedComponentBaseName = useSelector((state: RootState) => state.clusterBuilder.selectedComponentBaseName);
    let componentBase = cluster.componentBases[selectedComponentBaseName];
    const [inputField, setInputField] = useState('');

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setInputField(event.target.value);
    };
    const handleAddField = () => {
        if (inputField) {
            let sb = {
                addStatefulSet: false,
                addDeployment: false,
                addService: false,
                addIngress: false,
                addServiceMonitor: false,
                addConfigMap: false,
                configMap: {},
                statefulSet: {replicaCount: 0, pvcTemplates: [{name: '', storageSizeRequest: '', accessMode: ''}]},
                deployment: {replicaCount: 0},
                containers: {},
                ingress: {host: '',  authServerURL: '', paths: [{path: '', serviceName: '', pathType: ''}]},
            };
            let cbObj = {
                componentBaseName: selectedComponentBaseName,
                skeletonBaseName: inputField,
                skeletonBase: sb,
            }
            dispatch(setSelectedComponentBaseName(selectedComponentBaseName));
            dispatch(setSelectedSkeletonBaseName(inputField))
            dispatch(addSkeletonBase(cbObj));
            setInputField('');
        }
    };
    const handleRemoveField = (skeletonBaseName: string) => {
        dispatch(removeSkeletonBase({componentBaseName: selectedComponentBaseName, skeletonBaseName: skeletonBaseName}));
        if (cluster.componentBases[selectedComponentBaseName] !== undefined && Object.keys(cluster.componentBases[selectedComponentBaseName]).length > 0) {
            if (Object.keys(cluster.componentBases[selectedComponentBaseName])[0] === skeletonBaseName) {
                dispatch(setSelectedSkeletonBaseName(Object.keys(cluster.componentBases[selectedComponentBaseName])[1]));
            } else {
                dispatch(setSelectedSkeletonBaseName(Object.keys(cluster.componentBases[selectedComponentBaseName])[0]));
            }
        }
    };
    let showAdd = componentBase !== undefined;
    let show = showAdd && Object.keys(componentBase).length > 0;
    return (
        <div>
            { show && Object.keys(componentBase).map((key, index) => (
                <Box key={index} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                    <TextField
                        fullWidth
                        id={`inputField-${index}`}
                        label={`Workload Base Name`}
                        variant="outlined"
                        value={key}
                        InputProps={{ readOnly: true }}
                        sx={{ flex: 1, mr: 2, mb: 1}}
                    />
                    <Button variant="contained" sx={{ width: '100px' }} onClick={() => handleRemoveField(key)}>
                        Remove
                    </Button>
                </Box>))
            }
            { showAdd &&
            <Box key={Object.keys(componentBase).length} display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 1 }}>
                <TextField
                    fullWidth
                    id="inputField-new"
                    label="New Workload Base Name"
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

import * as React from "react";
import {useEffect, useState} from "react";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {useNavigate} from "react-router-dom";
import {appsApiGateway} from "../../gateway/apps";
import {setPublicMatrixFamilyApps} from "../../redux/apps/apps.reducer";
import {FormControl, FormControlLabel, FormLabel, Radio, RadioGroup,} from "@mui/material";
import Paper from "@mui/material/Paper";
import {TopologySystemComponents} from "../../redux/apps/apps.types";

export function AppConfigsTable(props: any) {
    const cluster = useSelector((state: RootState) => state.apps.cluster);
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(10);
    const matrixApps = useSelector((state: RootState) => state.apps.publicMatrixFamilyApps);
    const [loading, setLoading] = useState(true);
    const [selectedValue, setSelectedValue] = React.useState<string>('');
    const dispatch = useDispatch();
    let navigate = useNavigate();
    useEffect(() => {
        async function fetchData() {
            try {
                const response = await appsApiGateway.getMatrixPublicAppFamily(cluster.clusterName);
                dispatch(setPublicMatrixFamilyApps(response));
            } catch (e) {
            } finally {
                setLoading(false);
            }
        }
        fetchData();
    }, []);
    if (loading) {
        return null;
    }
    const handleChangeRowsPerPage = (
        event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
    ) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };
    const handleChangePage = (
        event: React.MouseEvent<HTMLButtonElement> | null,
        newPage: number,
    ) => {
        setPage(newPage);
    };

    const handleClick = async (event: any, app: any) => {
        event.preventDefault();
        navigate('/app/'+app.topologySystemComponentName);
    }

    if (matrixApps == null) {
        return (<div></div>)
    }

    const handleRadioChange = (event:any, row:  TopologySystemComponents) => {
        setSelectedValue(event.target.value);
        handleClick(event, row);
    };
    return (
        <Paper elevation={3} style={{ padding: '1rem', marginTop: '12.5px', marginLeft: '8px', marginRight: '8px' }}>
            <FormControl component="fieldset">
                <FormLabel component="legend" style={{ fontWeight: 'normal', color: '#333' }}>Name</FormLabel>
                <RadioGroup aria-label="matrix apps" name="matrixAppsRadio" value={selectedValue}>
                    {matrixApps.map((row, i) => (
                        <FormControlLabel
                            key={i}
                            value={row.topologySystemComponentName}
                            control={<Radio />}
                            label={row.topologySystemComponentName}
                            onChange={(event) => handleRadioChange(event, row)}
                        />
                    ))}
                </RadioGroup>
            </FormControl>
        </Paper>
    );
}

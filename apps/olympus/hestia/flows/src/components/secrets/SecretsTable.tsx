import * as React from "react";
import {useEffect, useState} from "react";
import Container from "@mui/material/Container";
import {Stack, TableContainer, TableRow} from "@mui/material";
import Paper from "@mui/material/Paper";
import Table from "@mui/material/Table";
import TableHead from "@mui/material/TableHead";
import TableCell from "@mui/material/TableCell";
import TableBody from "@mui/material/TableBody";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import {secretsApiGateway, SecretsRequest} from "../../gateway/secrets";

export function SecretKeys() {
    const [isLoading, setIsLoading] = useState(true)
    const [rows, setRows] = useState<SecretsRequest[]>([]);
    const [inputData, setInputData] = useState<SecretsRequest>({name: '', key: '', value: ''});
    const [buttonStatus, setButtonStatus] = useState<any>({});
    const [refresh, setRefresh] = useState<boolean>(false);

    const handleInputChange = (event: any) => {
        setInputData({...inputData, [event.target.name]: event.target.value });
    }

    const handleRequestCreateOrUpdateSecretKeyValue = async (row: SecretsRequest) => {
        try {
            await secretsApiGateway.upsertSecret(row);
            setRefresh(!refresh);
        } catch (exc) {
            console.error('Failed to upsert secret', row);
        } finally {
            setIsLoading(false);
        }
    };

    const handleRequestDeleteSecretKeyValue = async (row: SecretsRequest) => {
        try {
            const response = await secretsApiGateway.deleteSecret(row);
            if (response.status >= 400) {
                return;
            }
            let updatedRows = rows.filter(r => r !== row)
            setRows(updatedRows);
        } catch (exc) {
            console.error('Failed to upsert secret', row);
        } finally {
            setIsLoading(false);
        }
    };


    const getSecretReferenceValue = async (ref: string, rowIndex: number ) => {
        try {
            setIsLoading(true);
            let res = await secretsApiGateway.getSecret(ref);
            const row = res.data as SecretsRequest;
            const newRows = [...rows];
            newRows[rowIndex].value = row.value;
            setRows(newRows);
            setButtonStatus({...buttonStatus, [ref]: true});
        } catch (exc) {
            console.error('Failed to get secret', ref);
        } finally {
            setIsLoading(false);
        }
    }
    const hideSecretValue = (ref: string, rowIndex: number) => {
        const newRows = [...rows];
        newRows[rowIndex].value = '';
        setRows(newRows);

        // update button status
        setButtonStatus({...buttonStatus, [ref]: false});
    }

    useEffect(() => {
        async function getSecretReferences() {
            try {
                setIsLoading(true);
                const res = await secretsApiGateway.getSecrets();
                const rows = res.data as SecretsRequest[];
                setRows(rows)
            } catch (exc) {
            } finally {
                setIsLoading(false);
            }
        }
        getSecretReferences();
    }, [refresh]);

    //     }, [rows]);
    if (isLoading) {
        return <div>Loading...</div>;
    }
    return (
        <div>
            <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
                <TableContainer component={Paper}>
                    <Table sx={{ minWidth: 300 }} aria-label="simple table">
                        <TableHead>
                            <TableRow style={{ backgroundColor: '#333' }}>
                                <TableCell style={{ color: 'white' }} align="left">Secret Name</TableCell>
                                <TableCell style={{ color: 'white' }} align="left">Key</TableCell>
                                <TableCell style={{ color: 'white' }} align="left">Value</TableCell>
                                <TableCell style={{ color: 'white' }} align="left"></TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            <TableRow>
                                <TableCell align="left" component="th" scope="row">
                                    <TextField fullWidth name="name" value={inputData.name} onChange={handleInputChange} placeholder="Enter Secret Name"/>
                                </TableCell>
                                <TableCell align="left">
                                    <TextField fullWidth name="key" value={inputData.key} onChange={handleInputChange} placeholder="Enter Key Name"/>
                                </TableCell>
                                <TableCell align="left" >
                                    <TextField fullWidth name="value" value={inputData.value} onChange={handleInputChange} placeholder="Enter Value"/>
                                </TableCell>
                                <TableCell align="left" sx={{ paddingLeft: 1 }}>
                                    <Button
                                        fullWidth
                                        color="primary"
                                        variant="contained"
                                        onClick={() => handleRequestCreateOrUpdateSecretKeyValue(inputData)}
                                    >
                                        Save Secret
                                    </Button>
                                </TableCell>
                            </TableRow>
                            {rows && rows.map((row,index) => (
                                <TableRow key={index}>
                                    <TableCell>{row.name}</TableCell>
                                    <TableCell>{row.key}</TableCell>
                                    <TableCell>{row.value ? row.value : '***********************************'}</TableCell>
                                    <TableCell>
                                        <Stack direction={'row'} spacing={1}>
                                            <Button
                                                fullWidth
                                                color="primary"
                                                variant="contained"
                                                onClick={() => buttonStatus[row.key]
                                                    ? hideSecretValue(row.key, index)
                                                    : getSecretReferenceValue(row.name, index)
                                                }
                                            >
                                                {buttonStatus[row.key] ? 'Hide Value' : 'Get Value'}
                                            </Button>
                                            <Button
                                                fullWidth
                                                color="secondary"
                                                variant="contained"
                                                onClick={() => handleRequestDeleteSecretKeyValue(row)}
                                            >
                                                Delete Secret
                                            </Button>
                                        </Stack>
                                </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            </Container>
        </div>
    );
}

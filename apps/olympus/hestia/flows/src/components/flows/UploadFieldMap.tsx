import TextField from "@mui/material/TextField";
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import {Card, Stack} from "@mui/material";
import Box from "@mui/material/Box";

export function ContactsTextFieldRows(props: any) {
    const headers = useSelector((state: RootState) => state.flows.csvHeaders);

    return (
        <div>
            <Card sx={{ mb: 2, mt: 4 }}>
                {headers.map((header, index) => (
                    <Stack direction={"row"} key={index} sx={{ mb: 2, mt: 2 }}>
                        <Box flexGrow={1} sx={{ mt: 0, ml: 2 }}>
                            <TextField
                                fullWidth
                                label={header}
                                value={header}
                                inputProps={{ readOnly: true }}
                                variant="outlined"
                            />
                        </Box>
                        <Box flexGrow={1} sx={{ mt:0, ml: 2, mr: 2 }}>
                            <TextField
                                fullWidth
                                label={header}
                                variant="outlined"
                            />
                        </Box>
                    </Stack>
                ))}
            </Card>
        </div>
    );
}

export function PromptsTextFieldRows(props: any) {
    const headers = useSelector((state: RootState) => state.flows.promptHeaders);

    return (
        <div>
            <Card sx={{ mb: 2, mt: 4 }}>
                {headers.map((header, index) => (
                    <Stack direction={"row"} key={index} sx={{ mb: 2, mt: 2 }}>
                        <Box flexGrow={1} sx={{ mt: 0, ml: 2 }}>
                            <TextField
                                fullWidth
                                label={header}
                                value={header}
                                inputProps={{ readOnly: true }}
                                variant="outlined"
                            />
                        </Box>
                        <Box flexGrow={1} sx={{ mt:0, ml: 2, mr: 2 }}>
                            <TextField
                                fullWidth
                                label={header}
                                variant="outlined"
                            />
                        </Box>
                    </Stack>
                ))}
            </Card>
        </div>
    );
}
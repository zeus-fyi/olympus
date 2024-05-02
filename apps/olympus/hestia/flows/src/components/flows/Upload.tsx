import {Card, CardActionArea, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import {useDispatch, useSelector} from "react-redux";
import {setContactsCsvFilename, setCsvHeaders, setUploadContacts} from "../../redux/flows/flows.reducer";
import Container from "@mui/material/Container";
import {ContactsTable} from "./ContactsTable";
import {RootState} from "../../redux/store";
import {ContactsTextFieldRows} from "./UploadFieldMap";
import {UploadButton} from "./Analyze";
import Papa from "papaparse";

export function CsvUploadActionAreaCard(props: any) {
    const dispatch = useDispatch();
    const contacts = useSelector((state: RootState) => state.flows.uploadContentContacts);

    const onUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
        const files = event.target.files;
        const file = files?.item(0);
        if (!file) return;

        if (file.type === "application/json") {
            const reader = new FileReader();
            reader.onload = (e) => {
                const result = e.target?.result;
                if (!result) {
                    console.error("File read resulted in null.");
                    return;
                }

                try {
                    const data = JSON.parse(result as string);
                    // Assuming you want to set headers for JSON as well,
                    // you might need to derive them from the data structure
                    // For example, if data is an array of objects:
                    // headers = Object.keys(data[0]);
                    dispatch(setUploadContacts(data));
                    // You should dispatch the headers here if necessary
                } catch (error) {
                    console.error("Error parsing JSON file:", error);
                    return;
                }
            };
            reader.readAsText(file);
        } else if (file.type === "text/csv") {
            Papa.parse(file, {
                complete: (result) => {
                    try {
                        const data = result.data;
                        const headers = result.meta.fields || [];
                        if (Array.isArray(headers)) {
                            dispatch(setContactsCsvFilename(file.name));
                            dispatch(setCsvHeaders(headers));
                            dispatch(setUploadContacts(data as []));
                        }
                    } catch (error) {
                        console.error("Error parsing CSV file:", error);
                    }
                },
                header: true
            });
        } else {
            console.error("Unsupported file type:", file.type);
        }
    };
    return (
        <div>
            <Card sx={{ maxWidth: 400 }}>
                <CardActionArea>
                    <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
                        <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                            Upload User Contacts CSV
                        </Typography>
                        <UploadButton onUpload={onUpload}/>
                    </CardContent>
                </CardActionArea>
            </Card>
            <Container maxWidth="xl" sx={{ ml: -5, mt: 4}}>
                <ContactsTextFieldRows/>
            </Container>
            <Container maxWidth="xl" sx={{ ml: -5, mt: 4}}>
                <ContactsTable contacts={contacts}/>
            </Container>

        </div>
    );
}

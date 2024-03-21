import * as React from "react";
import {Card, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import {useSelector} from "react-redux";
import {WorkflowAnalysisTable} from "./WorkflowAnalysisTable";

export function Results(props: any) {
    const results = useSelector((state: any) => state.flows.results);
    return (
        <div>
            <Card sx={{ maxWidth: 320 }}>
                <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large', fontWeight: 'thin', marginRight: '15px', color: '#151C2F' }}>
                        Results
                    </Typography>
                </CardContent>
            </Card>
            <Box sx={{ mt: 4 }}>
                <WorkflowAnalysisTable />
            </Box>
        </div>
    );
}
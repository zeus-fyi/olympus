import {Card, CardActionArea, CardContent} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import {useDispatch} from "react-redux";
import {WorkflowAnalysisTable} from "../../flows/WorkflowAnalysisTable";

export function AdminWorkflowViews(props: any) {
    const [loading, setIsLoading] = React.useState(false);
    const dispatch = useDispatch();

    return (
        <div>
            <Card sx={{ maxWidth: 400 }}>
                <CardActionArea>
                    <CardContent style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}>
                        <Typography gutterBottom variant="h5" component="div" style={{ fontSize: 'large',fontWeight: 'thin', marginRight: '15x', color: '#151C2F'}}>
                            Admin Panel
                        </Typography>
                    </CardContent>
                </CardActionArea>
            </Card>
            <WorkflowAnalysisTable csvExport={true} isAdminPanel={true} />
        </div>
    );
}

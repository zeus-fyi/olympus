import {Card, CardActionArea, CardContent, Tab, Tabs} from "@mui/material";
import Typography from "@mui/material/Typography";
import * as React from "react";
import {useDispatch, useSelector} from "react-redux";
import {WorkflowAnalysisTable} from "../../flows/WorkflowAnalysisTable";
import {setAdminFlowsMainTab} from "../../../redux/flows/flows.reducer";
import Box from "@mui/material/Box";
import {UserStatsTable} from "./user_stats/UserStatsTable";

export function AdminWorkflowViews(props: any) {
    const mainTab = useSelector((state: any) => state.flows.adminFlowsMainTab);
    const [loading, setIsLoading] = React.useState(false);
    const dispatch = useDispatch();
    const handleMainTabChange = (event: React.SyntheticEvent, newValue: number) => {
        dispatch(setAdminFlowsMainTab(newValue));
    }
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
            <Box sx={{ mb: 2, mt: 2, ml: 0, mr:0  }}>
                <Tabs value={mainTab} onChange={handleMainTabChange} aria-label="basic tabs">
                    <Tab label="Users" />
                    <Tab label="Workflows" />
                </Tabs>
            </Box>
            {
                mainTab === 0 && (
                    <div>
                        <UserStatsTable />
                    </div>
                )
            }
            {
                mainTab === 1 && (
                    <WorkflowAnalysisTable csvExport={true} isAdminPanel={true} />
                )
            }
        </div>
    );
}

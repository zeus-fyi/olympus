import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import DashboardIcon from '@mui/icons-material/Dashboard';
import {Link} from "react-router-dom";
import AutoFixHighIcon from '@mui/icons-material/AutoFixHigh';
import CloudIcon from '@mui/icons-material/Cloud';
import ViewListIcon from '@mui/icons-material/ViewList';
import {ExpandLess, ExpandMore} from "@mui/icons-material";
import {Collapse, List, ListSubheader} from "@mui/material";
import MiscellaneousServicesIcon from '@mui/icons-material/MiscellaneousServices';
import ConstructionIcon from '@mui/icons-material/Construction';

export default function MainListItems() {
    const [open, setOpen] = React.useState(true);
    const [openClusters, setOpenClusters] = React.useState(true);

    const handleClick = () => {
        setOpen(!open);
    };

    const handleClickClusters = () => {
        setOpenClusters(!openClusters);
    };
    return (
        <List
            sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}
            component="nav"
            aria-labelledby="nested-list-subheader"
            subheader={
                <ListSubheader component="div" id="nested-list-subheader">
                    Zeus Cloud
                </ListSubheader>
            }
        >
            <ListItemButton component={Link} to="/dashboard">
                <ListItemIcon>
                    <DashboardIcon />
                </ListItemIcon>
                <ListItemText primary="Dashboard" />
            </ListItemButton>
            <ListItemButton onClick={handleClickClusters}  component={Link} to="/clusters">
                <ListItemIcon>
                    <CloudIcon />
                </ListItemIcon>
                <ListItemText primary="Clusters"/>
                {openClusters ? <ExpandLess /> : <ExpandMore />}
            </ListItemButton>
            <Collapse in={openClusters} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/clusters/builder">
                        <ListItemIcon>
                            <ConstructionIcon />
                        </ListItemIcon>
                        <ListItemText primary="Builder" />
                    </ListItemButton>
                </List>
            </Collapse>
            <ListItemButton onClick={handleClick}>
                <ListItemIcon>
                    <ViewListIcon />
                </ListItemIcon>
                <ListItemText primary="Services" />
                {open ? <ExpandLess /> : <ExpandMore />}
            </ListItemButton>
            <Collapse in={open} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/ethereum/validators">
                        <ListItemIcon>
                            <MiscellaneousServicesIcon />
                        </ListItemIcon>
                        <ListItemText primary="Validators" />
                    </ListItemButton>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/ethereum/aws">
                        <ListItemIcon>
                            <AutoFixHighIcon />
                        </ListItemIcon>
                        <ListItemText primary="AWS Wizard" />
                    </ListItemButton>
                </List>
            </Collapse>
            {/*<ListItemButton component={Link} to="/access">*/}
            {/*    <ListItemIcon>*/}
            {/*        <SecurityIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="Access"/>*/}
            {/*</ListItemButton>*/}
            {/*<ListItemButton component={Link} to="/billing">*/}
            {/*    <ListItemIcon>*/}
            {/*        <CreditCardIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="Billing" />*/}
            {/*</ListItemButton>*/}
        </List>
    );
}

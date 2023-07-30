import * as React from 'react';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import {Link} from "react-router-dom";
import AutoFixHighIcon from '@mui/icons-material/AutoFixHigh';
import CloudIcon from '@mui/icons-material/Cloud';
import ViewListIcon from '@mui/icons-material/ViewList';
import {ExpandLess, ExpandMore} from "@mui/icons-material";
import {Collapse, List, ListSubheader} from "@mui/material";
import MiscellaneousServicesIcon from '@mui/icons-material/MiscellaneousServices';
import ConstructionIcon from '@mui/icons-material/Construction';
import AppsIcon from '@mui/icons-material/Apps';
import CreditCardIcon from '@mui/icons-material/CreditCard';
import SecurityIcon from '@mui/icons-material/Security';
import ChatIcon from '@mui/icons-material/Chat';
import DnsIcon from "@mui/icons-material/Dns";
import SwapCallsIcon from '@mui/icons-material/SwapCalls';

export default function MainListItems() {
    const [open, setOpen] = React.useState(true);
    const [openClusters, setOpenClusters] = React.useState(true);
    const [openApps, setOpenApps] = React.useState(true);

    const handleClick = () => {
        setOpen(!open);
    };

    const handleClickApps = () => {
        setOpenApps(!openApps);
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
            <ListItemButton onClick={handleClickApps}  component={Link} to="/apps">
                <ListItemIcon>
                    <AppsIcon />
                </ListItemIcon>
                <ListItemText primary="Apps" />
                {openApps ? <ExpandLess /> : <ExpandMore />}
            </ListItemButton>
            <Collapse in={openApps} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/apps/builder">
                        <ListItemIcon>
                            <ConstructionIcon />
                        </ListItemIcon>
                        <ListItemText primary="Builder" />
                    </ListItemButton>
                </List>
            </Collapse>
            <ListItemButton component={Link} to="/compute">
                <ListItemIcon>
                    <DnsIcon />
                </ListItemIcon>
                <ListItemText primary="Compute" />
            </ListItemButton>
            <ListItemButton  component={Link} to="/clusters">
                <ListItemIcon>
                    <CloudIcon />
                </ListItemIcon>
                <ListItemText primary="Clusters"/>
            </ListItemButton>
            <ListItemButton component={Link} to="/loadbalancing/dashboard">
                <ListItemIcon>
                    <SwapCallsIcon />
                </ListItemIcon>
                <ListItemText primary="Load Balancing" />
            </ListItemButton>
            <ListItemButton onClick={handleClick}>
                <ListItemIcon>
                    <ViewListIcon />
                </ListItemIcon>
                <ListItemText primary="Services" />
                {open ? <ExpandLess /> : <ExpandMore />}
            </ListItemButton>
            <Collapse in={open} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/chatgpt">
                        <ListItemIcon>
                            <ChatIcon />
                        </ListItemIcon>
                        <ListItemText primary="ChatGPT" />
                    </ListItemButton>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/ethereum/aws">
                        <ListItemIcon>
                            <AutoFixHighIcon />
                        </ListItemIcon>
                        <ListItemText primary="AWS Wizard" />
                    </ListItemButton>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/ethereum/validators">
                        <ListItemIcon>
                            <MiscellaneousServicesIcon />
                        </ListItemIcon>
                        <ListItemText primary="Validators" />
                    </ListItemButton>
                </List>
            </Collapse>
            <ListItemButton component={Link} to="/billing">
                <ListItemIcon>
                    <CreditCardIcon />
                </ListItemIcon>
                <ListItemText primary="Billing" />
            </ListItemButton>
            <ListItemButton component={Link} to="/access">
                <ListItemIcon>
                    <SecurityIcon />
                </ListItemIcon>
                <ListItemText primary="Access" />
            </ListItemButton>
        </List>
    );
}

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
import {useSelector} from "react-redux";
import {RootState} from "../../redux/store";
import LeaderboardIcon from '@mui/icons-material/Leaderboard';

export default function MainListItems() {
    const [openServices, setOpenServices] = React.useState(true);
    const [openClusters, setOpenClusters] = React.useState(true);
    const [openCompute, setOpenCompute] = React.useState(true);
    const [openApps, setOpenApps] = React.useState(true);
    const isInternal = useSelector((state: RootState) => state.sessionState.isInternal);
    
    const handleClickServices = () => {
        setOpenServices(!openServices);
    };

    const handleClickApps = () => {
        setOpenApps(!openApps);
    };

    const handleClickCompute = () => {
        setOpenCompute(!openCompute);
    };

    return (
        <List
            sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}
            component="nav"
            aria-labelledby="nested-list-subheader"
            subheader={
                <ListSubheader component="div" id="nested-list-subheader">
                    Zeusfyi Universal Cloud
                </ListSubheader>
            }
        >
            {/*<ListItemButton component={Link} to="/ai">*/}
            {/*    <ListItemIcon>*/}
            {/*        <GraphicEqIcon />*/}
            {/*    </ListItemIcon>*/}
            {/*    <ListItemText primary="AI" />*/}
            {/*</ListItemButton>*/}
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
            <ListItemButton onClick={handleClickCompute} component={Link} to="/compute/search">
                <ListItemIcon>
                    <DnsIcon />
                </ListItemIcon>
                <ListItemText primary="Compute" />
                {openCompute ? <ExpandLess /> : <ExpandMore />}
            </ListItemButton>
            <Collapse in={openCompute} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    <ListItemButton sx={{ pl: 4 }} component={Link} to="/compute/summary">
                        <ListItemIcon>
                            <ConstructionIcon />
                        </ListItemIcon>
                        <ListItemText primary="Provisioned" />
                    </ListItemButton>
                </List>
            </Collapse>
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
            <ListItemButton onClick={handleClickServices}>
                <ListItemIcon>
                    <ViewListIcon />
                </ListItemIcon>
                <ListItemText primary="Services" />
                {openServices ? <ExpandLess /> : <ExpandMore />}
            </ListItemButton>
            <Collapse in={openServices} timeout="auto" unmountOnExit>
                <List component="div" disablePadding>
                    {isInternal && (
                            <div>
                                <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/mev">
                                    <ListItemIcon>
                                        <LeaderboardIcon />
                                    </ListItemIcon>
                                    <ListItemText primary="MEV" />
                                </ListItemButton>
                                <ListItemButton sx={{ pl: 4 }} component={Link} to="/services/chatgpt">
                                    <ListItemIcon>
                                        <ChatIcon />
                                    </ListItemIcon>
                                    <ListItemText primary="ChatGPT" />
                                </ListItemButton>
                            </div>
                        )
                    }
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

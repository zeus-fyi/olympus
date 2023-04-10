import * as React from 'react';
import {useEffect} from 'react';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import {AppBar, Drawer} from '../dashboard/Dashboard';
import Toolbar from '@mui/material/Toolbar';
import List from '@mui/material/List';
import Typography from '@mui/material/Typography';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import Container from '@mui/material/Container';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import Button from "@mui/material/Button";
import {useNavigate} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
import MainListItems from "../dashboard/listItems";
import {CustomerOptions, loadStripe, StripeElementsOptionsMode} from "@stripe/stripe-js";
import {Elements, PaymentElement,} from "@stripe/react-stripe-js";
import {configService} from "../../config/config";
import {Card} from "@mui/material";
import {stripeApiGateway} from "../../gateway/stripe";
import {setStripeCustomerID} from "../../redux/billing/billing.reducer";
import {RootState} from "../../redux/store";

const mdTheme = createTheme();

function BillingContent() {
    const [open, setOpen] = React.useState(true);
    const toggleDrawer = () => {
        setOpen(!open);
    };
    let navigate = useNavigate();
    const dispatch = useDispatch();

    const handleLogout = (event: any) => {
        event.preventDefault();
        authProvider.logout()
        dispatch({type: 'LOGOUT_SUCCESS'})
        navigate('/login');
    }

    useEffect(() => {
        async function fetchCustomerID() {
            try {
                const response = await stripeApiGateway.getCustomerID()
                dispatch(setStripeCustomerID(response.data.clientSecret));
            } catch (e) {
            }
        }
        fetchCustomerID().then(r => {});
    }, [dispatch]);


    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{ display: 'flex' }}>
                <CssBaseline />
                <AppBar position="absolute" open={open} style={{ backgroundColor: '#8991B0'}}>
                    <Toolbar
                        sx={{
                            pr: '24px', // keep right padding when drawer closed
                        }}
                    >
                        <IconButton
                            edge="start"
                            color="inherit"
                            aria-label="open drawer"
                            onClick={toggleDrawer}
                            sx={{
                                marginRight: '36px',
                                ...(open && { display: 'none' }),
                            }}
                        >
                            <MenuIcon />
                        </IconButton>
                        <Typography
                            component="h1"
                            variant="h6"
                            color="inherit"
                            noWrap
                            sx={{ flexGrow: 1 }}
                        >
                            Billing
                        </Typography>
                        <Button
                            color="inherit"
                            onClick={handleLogout}
                        >Logout
                        </Button>
                    </Toolbar>
                </AppBar>
                <Drawer variant="permanent" open={open}>
                    <Toolbar
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'flex-end',
                            px: [1],
                        }}
                    >
                        <IconButton onClick={toggleDrawer}>
                            <ChevronLeftIcon />
                        </IconButton>
                    </Toolbar>
                    <Divider />
                    <List component="nav">
                        <MainListItems />
                        <Divider sx={{ my: 1 }} />
                    </List>
                </Drawer>
                <Box
                    component="main"
                    sx={{
                        backgroundColor: (theme) =>
                            theme.palette.mode === 'light'
                                ? theme.palette.grey[100]
                                : theme.palette.grey[900],
                        flexGrow: 1,
                        height: '100vh',
                        overflow: 'auto',
                    }}
                >
                    <Toolbar />
                    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
                        <CheckoutPage />
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function Billing() {
    return <BillingContent />;
}

const stripe = loadStripe(configService.getStripePubKey());


function CheckoutPage() {
    const customerID = useSelector((state: RootState) => state.billing.stripeCustomerID);
    const options: StripeElementsOptionsMode = {
        mode: 'setup',
        currency: 'usd',
        customerOptions: {
            customer: customerID,
            ephemeralKey: ''
        } as CustomerOptions
    };
    return (
        <div>
            <Card>
                <Container maxWidth="lg">
                    <Box sx={{mt: 4, mb: 4}}>
                        <Elements stripe={stripe} options={options}>
                            <CheckoutForm/>
                        </Elements>
                    </Box>
                </Container>
            </Card>
        </div>
    );
}

export function CheckoutForm() {
    return (
        <form>
            <PaymentElement />
            <button>Submit</button>
        </form>
    );
}
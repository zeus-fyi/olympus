import * as React from 'react';
import {FormEvent, useState} from 'react';
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
import {useDispatch} from "react-redux";
import authProvider from "../../redux/auth/auth.actions";
import MainListItems from "../dashboard/listItems";
import {Appearance, loadStripe, StripeElementsOptionsMode, StripeError} from "@stripe/stripe-js";
import {AddressElement, Elements, PaymentElement, useElements, useStripe,} from "@stripe/react-stripe-js";
import {configService} from "../../config/config";
import {Card, CardContent} from "@mui/material";
import {stripeApiGateway} from "../../gateway/stripe";
import ReactGA from "react-ga4";

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

    return (
        <ThemeProvider theme={mdTheme}>
            <Box sx={{display: 'flex'}}>
                <CssBaseline/>
                <AppBar position="absolute" open={open} style={{backgroundColor: '#333'}}>
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
                                ...(open && {display: 'none'}),
                            }}
                        >
                            <MenuIcon/>
                        </IconButton>
                        <Typography
                            component="h1"
                            variant="h6"
                            color="inherit"
                            noWrap
                            sx={{flexGrow: 1}}
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
                            <ChevronLeftIcon/>
                        </IconButton>
                    </Toolbar>
                    <Divider/>
                    <List component="nav">
                        <MainListItems/>
                        <Divider sx={{my: 1}}/>
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
                    <Toolbar/>
                    <Container maxWidth="xl" sx={{mt: 4, mb: 4}}>
                        <CheckoutPage/>
                    </Container>
                </Box>
            </Box>
        </ThemeProvider>
    );
}

export default function Billing() {
    return <BillingContent/>;
}

const stripe = loadStripe(configService.getStripePubKey());

function CheckoutPage() {
    const appearance = {
        theme: 'stripe'
    } as Appearance;

    const options: StripeElementsOptionsMode = {
        paymentMethodTypes: ['card'],
        currency: 'usd',
        mode: 'setup',
        appearance: appearance
    };

    return (
        <div>
            <Card>
                <CardContent>
                    <Typography gutterBottom variant="h5" component="div">
                        Billing Info
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                        To use cluster deployments & other premium features, you must enter your billing information.
                    </Typography>
                </CardContent>
                <Container maxWidth="xl">
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
    const stripe = useStripe();
    const elements = useElements();

    const [errorMessage, setErrorMessage] = useState<string | undefined>();
    const [loading, setLoading] = useState<boolean>(false);

    const handleError = (error: StripeError) => {
        setLoading(false);
        setErrorMessage(error.message);
    };

    const handleSubmit = async (event: FormEvent) => {
        event.preventDefault();

        if (!stripe || !elements) {
            return;
        }

        setLoading(true);
        const { error: submitError } = await elements.submit();
        if (submitError) {
            handleError(submitError);
            return;
        }

        async function fetchCustomerID() {
            try {
                const response = await stripeApiGateway.getClientSecret()
                return response.data.clientSecret
            } catch (e) {
                handleError(error);
            }
        }
        const clientSecret = await fetchCustomerID()
        const { error } = await stripe.confirmSetup({
            elements,
            clientSecret,
            confirmParams: {
                return_url: 'https://cloud.zeus.fyi/dashboard',
            },
        });

        if (error) {
            handleError(error);
        } else {
            ReactGA.gtag('event','add_billing_info', { 'method': 'Stripe' });
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <AddressElement options={{mode: 'billing'}} />
            <PaymentElement />
            <button type="submit" disabled={!stripe || loading}>
                Submit
            </button>
            {errorMessage && <div>{errorMessage}</div>}
        </form>
    );
}
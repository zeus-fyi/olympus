import * as React from 'react';
import {useState} from 'react';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import {ZeusCopyright} from "../copyright/ZeusCopyright";
import {CircularProgress} from "@mui/material";
import {signUpApiGateway} from "../../gateway/signup";
import ReactGA from "react-ga4";

const theme = createTheme();

export default function SignUp() {
    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        setRequestStatus('pending')
        event.preventDefault();
        const data = new FormData(event.currentTarget);
        console.log({
            firstName: data.get('firstName'),
            lastName: data.get('lastName'),
            email: data.get('email'),
            password: data.get('password'),
        });
        const firstName = data.get('firstName') as string
        const lastName = data.get('lastName') as string
        const email = data.get('email') as string
        const password = data.get('password') as string
        try {
            const res = await signUpApiGateway.sendSignUpRequest(firstName, lastName, email, password)
            if (res.status === 400) {
                setRequestStatus('errorFormValidation')
                return
            }
            if (res.status === 200) {
                setRequestStatus('success')
                ReactGA.gtag('event','sign_up', { 'method': 'Website' });
            } else {
                setRequestStatus('error')
            }
        } catch (e) {
            setRequestStatus('error')
            console.error(e)
        }
    };
    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const [requestStatus, setRequestStatus] = useState('');

    switch (requestStatus) {
        case 'pending':
            buttonLabel = <CircularProgress size={20} />;
            buttonDisabled = true;
            break;
        case 'success':
            buttonLabel = 'Check your email to verify your account!';
            buttonDisabled = true;
            statusMessage = 'Check your email to verify your account!';
            break;
        case 'error':
            buttonLabel = 'Resubmit';
            buttonDisabled = false;
            statusMessage = 'A user with that email already exists or an error has occurred. Check your email for a verification code';
            break;
        case 'errorDuplicateUser':
            buttonLabel = 'Resubmit';
            buttonDisabled = false;
            statusMessage = 'A user with that email already exists or an error has occurred. Check your email for a verification code';
            break;
        case 'errorFormValidation':
            buttonLabel = 'Resubmit';
            buttonDisabled = false;
            statusMessage = 'You must provide a valid email and password, and you cannot leave any fields blank';
            break;
        default:
            buttonLabel = 'Sign Up';
            buttonDisabled = false;
            break;
    }
    return (
        <ThemeProvider theme={theme}>
            <Container component="main" maxWidth="xs">
                <CssBaseline />
                <Box
                    sx={{
                        marginTop: 8,
                        display: 'flex',
                        flexDirection: 'column',
                        alignItems: 'center',
                    }}
                >
                    <Typography component="h1" variant="h5">
                        Sign up
                    </Typography>
                    <Box component="form" noValidate onSubmit={handleSubmit} sx={{ mt: 3 }}>
                        <Grid container spacing={2}>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    autoComplete="given-name"
                                    name="firstName"
                                    required
                                    fullWidth
                                    id="firstName"
                                    label="First Name"
                                    autoFocus
                                />
                            </Grid>
                            <Grid item xs={12} sm={6}>
                                <TextField
                                    required
                                    fullWidth
                                    id="lastName"
                                    label="Last Name"
                                    name="lastName"
                                    autoComplete="family-name"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    required
                                    fullWidth
                                    id="email"
                                    label="Email Address"
                                    name="email"
                                    autoComplete="email"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    required
                                    fullWidth
                                    name="password"
                                    label="Password"
                                    type="password"
                                    id="password"
                                    autoComplete="new-password"
                                />
                            </Grid>
                            {/*<Grid item xs={12}>*/}
                            {/*    <FormControlLabel*/}
                            {/*        control={<Checkbox value="allowExtraEmails" color="primary" />}*/}
                            {/*        label="I want to receive inspiration, marketing promotions and updates via email."*/}
                            {/*    />*/}
                            {/*</Grid>*/}
                        </Grid>
                        <Button
                            type="submit"
                            fullWidth
                            color="primary"
                            variant="contained"
                            sx={{ mt: 3, mb: 2, backgroundColor: '#333'}}
                            disabled={buttonDisabled}
                        >
                            {buttonLabel}
                        </Button>
                        {statusMessage && (
                            <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                                {statusMessage}
                            </Typography>
                        )}
                        <Grid container justifyContent="flex-end">
                            <Grid item>
                                <Link href="/login" variant="body2" color="text.primary">
                                    Already have an account? Sign in
                                </Link>
                            </Grid>
                        </Grid>
                    </Box>
                </Box>
                <ZeusCopyright sx={{ mt: 5 }} />
            </Container>
        </ThemeProvider>
    );
}
import * as React from 'react';
import {useState} from 'react';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Paper from '@mui/material/Paper';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import {createTheme, ThemeProvider} from '@mui/material/styles';
import {useNavigate} from "react-router-dom";
import authProvider from "../../redux/auth/auth.actions";
import {useDispatch} from "react-redux";
import {LOGIN_FAIL, LOGIN_SUCCESS,} from "../../redux/auth/auth.types";
import {ZeusCopyright} from "../copyright/ZeusCopyright";
import Link from "@mui/material/Link";
import {CircularProgress} from "@mui/material";
import {setSessionAuth} from "../../redux/auth/session.reducer";

const theme = createTheme();

const Login = () => {
    let navigate = useNavigate();
    const dispatch = useDispatch();
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
            buttonLabel = 'Logged in successfully';
            buttonDisabled = true;
            statusMessage = 'Logged in successfully!';
            break;
        case 'error':
            buttonLabel = 'Retry';
            buttonDisabled = false;
            statusMessage = 'An error occurred while logging in, please try again. If you continue having issues please email support@zeus.fyi';
            break;
        default:
            buttonLabel = 'Login';
            buttonDisabled = false;
            break;
    }
    const handleLogin = async (event: React.FormEvent<HTMLFormElement>) =>  {
        event.preventDefault();
        const data = new FormData(event.currentTarget);
        let email = data.get('email') as string
        let password = data.get('password') as string
        try {
            setRequestStatus('pending');
            let res: any = await authProvider.login(email, password)
            const statusCode = res.status;
            if (statusCode === 200 || statusCode === 204) {
                setRequestStatus('success');
                dispatch(setSessionAuth(true))
                dispatch({type: 'LOGIN_SUCCESS', payload: res.data})
                navigate('/dashboard');
            } else {
                dispatch(setSessionAuth(false))
                dispatch({type: 'LOGIN_FAIL', payload: res.data})
                setRequestStatus('error');
            }
        } catch (e) {
            dispatch(setSessionAuth(false))
            setRequestStatus('error');
        }
    }
    return (
        <ThemeProvider theme={theme}>
            <Grid container component="main" sx={{height: '100vh'}}>
                <CssBaseline/>
                <Grid
                    item
                    xs={false}
                    sm={4}
                    md={7}
                    sx={{
                        backgroundColor: (t) =>
                            t.palette.mode === 'dark' ? t.palette.grey[50] : t.palette.grey[900],
                        backgroundSize: 'auto 100%',
                        backgroundPosition: 'center',
                    }}
                />
                <Grid item xs={12} sm={8} md={5} component={Paper} elevation={6} square>
                    <Box
                        sx={{
                            my: 8,
                            mx: 4,
                            display: 'flex',
                            flexDirection: 'column',
                            alignItems: 'center',
                        }}
                    >
                        <Typography component="h1" variant="h5">
                            Sign in
                        </Typography>
                        <Box component="form" noValidate onSubmit={handleLogin} sx={{mt: 1}}>
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                id="email"
                                label="Email Address"
                                name="email"
                                autoComplete="email"
                                autoFocus
                            />
                            <TextField
                                margin="normal"
                                required
                                fullWidth
                                name="password"
                                label="Password"
                                type="password"
                                id="password"
                                autoComplete="current-password"
                            />
                            {/*<FormControlLabel*/}
                            {/*    control={<Checkbox value="remember" color="primary"/>}*/}
                            {/*    label="Remember me"*/}
                            {/*/>*/}
                            <Button
                                type="submit"
                                fullWidth
                                variant="contained"
                                color="primary"
                                sx={{mt: 3, mb: 2, backgroundColor: '#333'}}
                                disabled={buttonDisabled}>{buttonLabel}</Button>
                            {statusMessage && (
                                <Typography variant="body2" color={requestStatus === 'error' ? 'error' : 'success'}>
                                    {statusMessage}
                                </Typography>
                            )}
                            <Grid container>
                                {/*<Grid item xs>*/}
                                {/*    <Link href="#" variant="body2">*/}
                                {/*        Forgot password?*/}
                                {/*    </Link>*/}
                                {/*</Grid>*/}
                                <Grid item>
                                    <Link href="/signup" variant="body2" color="text.primary">
                                        {"Don't have an account? Sign Up"}
                                    </Link>
                                </Grid>
                            </Grid>
                            <ZeusCopyright sx={{mt: 5}}/>
                        </Box>
                    </Box>
                </Grid>
            </Grid>
        </ThemeProvider>
    );
}

export default Login;

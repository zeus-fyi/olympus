import {useNavigate} from "react-router-dom";
import {useDispatch} from "react-redux";
import * as React from "react";
import {useEffect, useState} from "react";
import {CircularProgress} from "@mui/material";

export const TwitterLogin = () => {
    let navigate = useNavigate();
    const dispatch = useDispatch();
    let buttonLabel;
    let buttonDisabled;
    let statusMessage;
    const [requestStatus, setRequestStatus] = useState('');
    const [loading, setLoading] = useState(false);
    switch (requestStatus) {
        case 'pending':
            buttonLabel = <CircularProgress size={20}/>;
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
    useEffect(() => {
        if (requestStatus === 'success') {
            navigate('/ai');
        }
    }, [requestStatus]);

    return <div>Redirecting</div>
}
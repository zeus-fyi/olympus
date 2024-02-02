import {useLocation, useNavigate} from "react-router-dom";
import * as React from "react";
import {useEffect, useState} from "react";
import {accessApiGateway} from "../../gateway/access";

export const TwitterLogin = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [requestStatus, setRequestStatus] = useState('');
    const location = useLocation();
    const navigate = useNavigate();

    useEffect(() => {
        const fetchData = async () => {
            try {
                setIsLoading(true); // Set loading to true
                const queryParams = new URLSearchParams(location.search);
                const code = queryParams.get('code');
                const state = queryParams.get('state');
                if (code && state) {
                    // Assuming accessApiGateway.callbackPlatformAuthFlow has been updated to accept code and state
                    const response = await accessApiGateway.callbackPlatformAuthFlow('twitter', code, state);
                    console.log("response", response);
                    console.log("response.data", response.data);

                    setRequestStatus('success'); // Update request status to success
                }
            } catch (error) {
                console.log("error", error);
                setRequestStatus('error'); // Update request status to error
            } finally {
                setIsLoading(false); // Set loading to false regardless of success or failure.
            }
        };

        fetchData();
    }, [location]);

    useEffect(() => {
        // This separate useEffect watches for changes in requestStatus and navigates accordingly
        if (requestStatus === 'success') {
            navigate('/ai');
        }
    }, [requestStatus, navigate, isLoading]);

    if (isLoading) {
        return <div>Loading...</div>
    }
    return <div>Redirecting</div>
}
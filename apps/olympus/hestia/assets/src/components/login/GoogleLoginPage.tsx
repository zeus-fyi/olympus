import React from 'react';
import {GoogleLogin, GoogleOAuthProvider} from '@react-oauth/google';
import {authApiGateway} from "../../gateway/login";
import {setSessionAuth} from "../../redux/auth/session.reducer";
import {setUserPlanDetails} from "../../redux/loadbalancing/loadbalancing.reducer";
import {useDispatch} from "react-redux";
import {useNavigate} from "react-router-dom";

const GoogleLoginPage = (props: any) => {
    const dispatch = useDispatch();
    const navigate = useNavigate();
    const {loading, setLoading, setRequestStatus,} = props;
    const handleGoogleLogin = async (credentialResponse: any) =>  {
        try {
            setLoading(true);
            setRequestStatus('pending');
            console.log('credentialResponse', credentialResponse)
            const res = await authApiGateway.sendGoogleLoginRequest(credentialResponse);
            const statusCode = res.status;
            if (statusCode < 300) {
                setRequestStatus('success');
                dispatch(setSessionAuth(true))
                if (res.data.planUsageDetails != null && res.data.planUsageDetails.plan != undefined){
                    dispatch(setUserPlanDetails(res.data.planUsageDetails))
                }
                dispatch({type: 'LOGIN_SUCCESS', payload: res.data})
                navigate('/apps');
            } else {
                dispatch(setSessionAuth(false))
                dispatch({type: 'LOGIN_FAIL', payload: res.data})
            }
        } catch (e) {
            dispatch(setSessionAuth(false))
            setRequestStatus('error');
        } finally {
            setLoading(false);
            setRequestStatus('done')
        }
    }
    if (loading) return (<div></div>)
    return (
        <div className="">
            <div className="">
                <GoogleOAuthProvider
                    clientId={`${process.env.REACT_APP_GOOGLE_API_TOKEN}`}
                >
                    <GoogleLogin
                        onSuccess={async credentialResponse => {
                            await handleGoogleLogin(credentialResponse)
                        }}
                        onError={() => {
                            console.log('Login Failed');
                        }}
                    />
                </GoogleOAuthProvider>
            </div>
        </div>
    )
}

export default GoogleLoginPage
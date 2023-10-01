import React from 'react';
import {GoogleLogin, GoogleOAuthProvider} from '@react-oauth/google';
import {authApiGateway} from "../../gateway/login";

const GoogleLoginPage = () => {
    return (
        <div className="">
            <div className="">
                <GoogleOAuthProvider
                    clientId={`${process.env.REACT_APP_GOOGLE_API_TOKEN}`}
                >
                    <GoogleLogin
                        onSuccess={async credentialResponse => {
                            const verificationResponse = await authApiGateway.sendGoogleLoginRequest(credentialResponse);
                            // TODO add redirect to dashboard
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
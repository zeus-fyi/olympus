import React from 'react';
// @ts-ignore
import {GoogleLogin, GoogleOAuthProvider} from '@react-oauth/google';

import jwt_decode from "jwt-decode";
import {useNavigate} from "react-router-dom";


const GoogleLoginPage = () => {
    const navigate = useNavigate();
    // @ts-ignore
    const responseGoogle = (response) => {
        console.log(response);
        var userObject = jwt_decode(response.credential);
        console.log(userObject);
        localStorage.setItem('user', JSON.stringify(userObject));
        // @ts-ignore
        const { name, sub, picture } = userObject;
        const doc = {
            _id: sub,
            _type: 'user',
            userName: name,
            image: picture,
        };
        console.log(doc);
        // @ts-ignore
        client.createIfNotExists(doc).then(() => {
            navigate('/', { replace: true });
        });

    }
    // @ts-ignore
    return (
        <div className="">
            <div className="">
                <GoogleOAuthProvider
                    clientId={`${process.env.REACT_APP_GOOGLE_API_TOKEN}`}
                >
                    <GoogleLogin
                        // @ts-ignore
                        render={(renderProps) => (
                            <button
                                type="button"
                                className=""
                                onClick={renderProps.onClick}
                                disabled={renderProps.disabled}
                            >Sign in with google
                            </button>
                        )}
                        onSuccess={responseGoogle}
                        onFailure={responseGoogle}
                        cookiePolicy="single_host_origin"
                    />
                </GoogleOAuthProvider>
            </div>
        </div>
    )
}

export default GoogleLoginPage
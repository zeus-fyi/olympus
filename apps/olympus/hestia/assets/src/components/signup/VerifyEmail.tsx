import {useNavigate, useParams} from "react-router-dom";
import React, {useEffect} from "react";
import {signUpApiGateway} from "../../gateway/signup";

export function VerifyEmail() {
    const params = useParams();
    let navigate = useNavigate();

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await signUpApiGateway.verifyEmail(params.id);
                const statusCode = response.status;
                if (statusCode === 200 || statusCode === 204) {
                    navigate('/login');
                } else {
                    navigate('/signup');
                }
            } catch (error) {
                navigate('/signup');
                console.log("error", error);
            }}
        fetchData(params);
    }, []);
    return (<div></div>)
}
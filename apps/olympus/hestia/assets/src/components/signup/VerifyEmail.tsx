import {Navigate, useParams} from "react-router-dom";
import React, {useEffect} from "react";
import {signUpApiGateway} from "../../gateway/signup";

export function VerifyEmail() {
    const params = useParams();
    useEffect(() => {
        const fetchData = async (params: any) => {
            try {
                const response = await signUpApiGateway.verifyEmail(params.id);
                if (response.status === 200) {
                    return <Navigate to="/login" />;
                } else {
                    return <Navigate to="/signup" />;
                }
            } catch (error) {
                console.log("error", error);
            }}
        fetchData(params);
    }, []);
    return (<div></div>)
}
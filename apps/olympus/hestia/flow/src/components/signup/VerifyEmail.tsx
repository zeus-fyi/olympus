import {useNavigate, useParams} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {signUpApiGateway} from "../../gateway/signup";

export function VerifyEmail() {
    const params = useParams();
    let navigate = useNavigate();
    const [loading, setLoading] = useState(true);

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
            } finally {
                setLoading(false); // Set loading to false regardless of success or failure.
            }
        }
        fetchData(params);
    }, []);
    if (loading) {
        return <div>Loading...</div> // Display loading message while data is fetching
    }
    return (<div></div>)
}
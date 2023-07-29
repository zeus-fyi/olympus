import {useLocation, useNavigate} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {signUpApiGateway} from "../../gateway/signup";

export function VerifyQuickNodeLoginJWT() {
    let navigate = useNavigate();
    let location = useLocation(); // Missing in your code
    const [loading, setLoading] = useState(true);
    const parseJwtFromSearch = () => {
        const searchParams = new URLSearchParams(location.search);
        return searchParams.get('jwt');
    }
    const jwtToken = parseJwtFromSearch();
    if (!jwtToken) {
        throw new Error('JWT token is missing in the query parameters');
    }

    useEffect(() => {
        const fetchData = async (params: any) => {
            try {

                const response = await signUpApiGateway.verifyJWT(jwtToken);
                const statusCode = response.status;
                if (statusCode === 200 || statusCode === 204) {
                    navigate('/services/quicknode/dashboard');
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
        fetchData(jwtToken);
    }, [jwtToken]);

    if (loading) {
        return <div>Loading...</div> // Display loading message while data is fetching
    }
    return (<div></div>)
}
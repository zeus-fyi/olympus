import {useLocation, useNavigate} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {signUpApiGateway} from "../../gateway/signup";
import {setSessionAuth} from "../../redux/auth/session.reducer";
import {useDispatch} from "react-redux";
import inMemoryJWT from "../../auth/InMemoryJWT";
import {pipe, prop} from "ramda";
import {getAxiosResponse} from "../../helpers/get-axios-response";
import {setUserPlanDetails} from "../../redux/loadbalancing/loadbalancing.reducer";

const sessionIDParse = pipe(getAxiosResponse,prop('sessionID'));
const ttlSeconds = pipe(getAxiosResponse, prop('ttl'));
const userIDParse = pipe(getAxiosResponse, prop('userID'));
const planUsageDetailsParse = pipe(getAxiosResponse, prop('planUsageDetails'));

export function VerifyQuickNodeLoginJWT() {
    let navigate = useNavigate();
    let location = useLocation(); // Missing in your code
    const [loading, setLoading] = useState(true);
    const dispatch = useDispatch();

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
                setLoading(true)
                const response = await signUpApiGateway.verifyJWT(jwtToken);
                const statusCode = response.status;
                if (statusCode === 200 || statusCode === 204) {
                    const sessionID = sessionIDParse(response);
                    const tokenExpiry = ttlSeconds(response);
                    const userID = userIDParse(response);
                    const planDetails = planUsageDetailsParse(response);
                    dispatch(setUserPlanDetails(planDetails));
                    inMemoryJWT.setToken(sessionID, tokenExpiry);
                    localStorage.setItem("userID", userID);
                    dispatch(setSessionAuth(true))
                    dispatch({type: 'LOGIN_SUCCESS', payload: response.data})
                    navigate('/loadbalancing/dashboard');
                } else {
                    dispatch(setSessionAuth(false))
                    dispatch({type: 'LOGIN_FAIL', payload: response.data})
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
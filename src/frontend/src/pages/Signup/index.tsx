import React from "react";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";

import SignUpPage from "./SignupPage";

const SignUp = () => {
    let navigate = useNavigate();
    return ( 
        <SignUpPage navigate={navigate}/>
    )
}

export default SignUp;
import React, {useEffect} from "react";
import {Outlet, Navigate, useLocation, useNavigate} from "react-router-dom";


const PrivateRoute = () => {
    const location = useLocation()
    const navigate = useNavigate()

    // const token = localStorage.getItem("token")
    //
    // if (token) {
    //     return (
    //         <div className={"app"}>
    //             <Outlet/>
    //         </div>
    //     )
    // } else {
    //     return <Navigate to="/login"/>;
    // }

    return (
        <div className={"app"}>
            <nav>
                <ul>
                    <li className={!location.pathname.includes("list") ? "active" : ""}
                        onClick={() => navigate("/")}>Поиск</li>
                    <li className={location.pathname.includes("list") ? "active" : ""}
                        onClick={() => navigate("/list")}>Список адресов с файлами</li>
                </ul>
            </nav>
            <Outlet/>
        </div>
    )
};

export default PrivateRoute
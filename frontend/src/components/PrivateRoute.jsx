import React, {useEffect, useState} from "react";
import {Outlet, Navigate, useNavigate, useLocation} from "react-router-dom";


const PrivateRoute = () => {
    const location = useLocation()
    const navigate = useNavigate()
    const [token, setToken] = useState(localStorage.getItem("token"))

    useEffect(() => {
        setToken(localStorage.getItem("token"))
    }, [location.pathname]);

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

    if (token) {
        return (
            <div className={"app"}>
                <nav>
                    <ul>
                        <li className={!location.pathname.includes("list") ? "active" : ""}
                            onClick={() => navigate("/")}>Поиск</li>
                        <li className={location.pathname.includes("list") ? "active" : ""}
                            onClick={() => navigate("/list")}>Список адресов с файлами</li>
                        <li>
                            <span>123</span>
                            <button>Выход</button>
                        </li>
                    </ul>
                </nav>
                <Outlet/>
            </div>
        )
    } else {
        return <Navigate to="/login" />;
    }
};

export default PrivateRoute
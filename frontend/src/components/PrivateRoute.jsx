import React, {useEffect, useState} from "react";
import {Outlet, Navigate, useNavigate, useLocation} from "react-router-dom";
import API_DOMAIN from "../config";
import FetchRequest from "../fetchRequest";


const PrivateRoute = ({requiredAdmin}) => {
    const location = useLocation()
    const navigate = useNavigate()
    const [token, setToken] = useState(localStorage.getItem("token"))
    const [user, setUser] = useState(null)
    const [isLoaded, setIsLoaded] = useState(false)

    // useEffect(() => {
    //     setToken(localStorage.getItem("token"))
    // }, [location.pathname]);

    useEffect(() => {
        FetchRequest("GET", "/get_auth", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setToken(localStorage.getItem("token"))
                    setUser(response.data)
                }

                setIsLoaded(true)
            })
    }, []);


    if (token) {
        return (
            <div className={"app"}>
                {isLoaded && <><nav>
                    <ul>
                        <li className={!location.pathname.includes("list") && !location.pathname.includes("users") ? "active" : ""}
                            onClick={() => navigate("/")}>Поиск</li>
                        <li className={location.pathname.includes("list") ? "active" : ""}
                            onClick={() => navigate("/list")}>Список адресов с файлами</li>
                        {user.Role.Value === "admin" && <li className={location.pathname.includes("users") ? "active" : ""}
                                                 onClick={() => navigate("/users")}>Пользователи</li>}
                    </ul>
                    <div>
                        <span style={{color: "#fff"}}>{user.Login}</span>
                        <button>Выход</button>
                    </div>
                </nav>
                {requiredAdmin ? user.Role.Value === "admin" ? <Outlet/> : <Navigate to="/" /> : <Outlet/>}
                </>}
            </div>
        )
    } else {
        return <Navigate to="/login" />;
    }
};

export default PrivateRoute
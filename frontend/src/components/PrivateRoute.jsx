import React, {useEffect, useState} from "react";
import {Outlet, Navigate, useNavigate, useLocation, Link} from "react-router-dom";
import API_DOMAIN from "../config";
import FetchRequest from "../fetchRequest";


const PrivateRoute = ({requiredAdmin}) => {
    const location = useLocation()
    const navigate = useNavigate()
    const [token, setToken] = useState(localStorage.getItem("token"))
    const [user, setUser] = useState(null)
    const [isLoaded, setIsLoaded] = useState(false)
    const [activeTab, setActiveTab] = useState(1)

    useEffect(() => {
        if (location.pathname === "/") {
            setActiveTab(1)
        } else if (location.pathname === "/list" || location.pathname.includes("house")) {
            setActiveTab(2)
        } else if (location.pathname.includes("nodes")) {
            setActiveTab(3)
        } else if (location.pathname.includes("/references")) {
            setActiveTab(6)
        } else if (location.pathname.includes("/hardware")) {
            setActiveTab(4)
        } else if (location.pathname.includes("/users")) {
            setActiveTab(5)
        } else if (location.pathname.includes("/events")) {
            setActiveTab(7)
        }
    }, [location.pathname]);

    useEffect(() => {
        FetchRequest("GET", "/auth/me", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setToken(localStorage.getItem("token"))
                    setUser(response.data)
                    setIsLoaded(true)
                }
            })
    }, []);

    const handlerExit = () => {
        FetchRequest("GET", "/auth/logout", null)
            .then(response => {
                if (response.success) {
                    setToken("")
                    localStorage.removeItem("token")
                    navigate("/login")
                }
            })
    }

    if (token) {
        return (
            <div className={"app"}>
                {isLoaded && <><nav>
                    <ul>
                        <Link to={"/"}>
                            <li className={activeTab === 1 ? "active" : ""} onClick={() => {setActiveTab(1)}}>Поиск</li>
                        </Link>
                        <Link to={"/list"}>
                            <li className={activeTab === 2 ? "active" : ""} onClick={() => {setActiveTab(2)}}>Дома</li>
                        </Link>
                        <Link to={"/nodes"}>
                            <li className={activeTab === 3 ? "active" : ""} onClick={() => {setActiveTab(3)}}>Узлы</li>
                        </Link>
                        <Link to={"/hardware"}>
                            <li className={activeTab === 4 ? "active" : ""} onClick={() => {setActiveTab(4)}}>Оборудование</li>
                        </Link>
                        {user.Role.Value === "admin" &&
                            <Link to={"/users"}>
                                <li className={activeTab === 5 ? "active" : ""} onClick={() => {setActiveTab(5)}}>Пользователи</li>
                            </Link>}
                        {user.Role.Value === "admin" &&
                            <Link to={"/references"}>
                                <li className={activeTab === 6 ? "active" : ""} onClick={() => {setActiveTab(6)}}>Справочники</li>
                            </Link>}
                        <Link to={"/events"}>
                            <li className={activeTab === 7 ? "active" : ""} onClick={() => {setActiveTab(7)}}>События</li>
                        </Link>
                    </ul>
                    <div>
                        <span style={{color: "#fff"}}>{user.Login}</span>
                        <button onClick={handlerExit}>Выход</button>
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
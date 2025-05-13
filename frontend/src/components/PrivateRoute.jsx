import React, {useContext, useEffect, useState} from "react";
import {Outlet, Navigate, useNavigate, useLocation, Link} from "react-router-dom";
import API_DOMAIN from "../config";
import FetchRequest from "../fetchRequest";
import AuthContext from "../context/AuthContext";

const PrivateRoute = ({requiredAdmin}) => {
    const { token, user, isLoaded, logout } = useContext(AuthContext);
    const location = useLocation()
    const navigate = useNavigate()
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

    const handlerExit = () => {
        logout()
        navigate("/login")
    }

    if (!token) return <Navigate to="/login" />;

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
                    <Link to={"/references"}>
                        <li className={activeTab === 6 ? "active" : ""} onClick={() => {setActiveTab(6)}}>Справочники</li>
                    </Link>
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
};

export default PrivateRoute
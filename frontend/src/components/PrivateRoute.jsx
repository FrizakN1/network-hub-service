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
    const [activeTab, setActiveTab] = useState(1)

    // useEffect(() => {
    //     setToken(localStorage.getItem("token"))
    // }, [location.pathname]);

    useEffect(() => {
        if (location.pathname.includes("list")) {
            setActiveTab(2)
        } else if (location.pathname.includes("nodes")) {
            setActiveTab(3)
        } else if (location.pathname.includes("hardware")) {
            setActiveTab(4)
        } else if (location.pathname.includes("users")) {
            setActiveTab(5)
        } else if (location.pathname.includes("references")) {
            setActiveTab(6)
        }

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
                        <li className={activeTab === 1 ? "active" : ""}
                            onClick={() => {
                                setActiveTab(1)
                                navigate("/")
                            }}>Поиск</li>
                        <li className={activeTab === 2 ? "active" : ""}
                            onClick={() => {
                                setActiveTab(2)
                                navigate("/list")
                            }}>Дома</li>
                        <li className={activeTab === 3 ? "active" : ""}
                            onClick={() => {
                                setActiveTab(3)
                                navigate("/nodes")
                            }}>Узлы</li>
                        <li className={activeTab === 4 ? "active" : ""}
                            onClick={() => {
                                setActiveTab(4)
                                navigate("/hardware")
                            }}>Оборудование</li>
                        {user.Role.Value === "admin" && <li className={activeTab === 5 ? "active" : ""}
                                                 onClick={() => {
                                                     setActiveTab(5)
                                                     navigate("/users")
                                                 }}>Пользователи</li>}
                        {user.Role.Value === "admin" && <li className={activeTab === 6 ? "active" : ""}
                                                            onClick={() => {
                                                                setActiveTab(6)
                                                                navigate("/references")
                                                            }}>Справочники</li>}
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
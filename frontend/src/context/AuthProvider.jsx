import {useEffect, useState} from "react";
import AuthContext from "./AuthContext";
import FetchRequest from "../fetchRequest";

const AuthProvider = ({ children }) => {
    const [token, setToken] = useState(localStorage.getItem("token"));
    const [user, setUser] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);

    useEffect(() => {
        if (token) {
            FetchRequest("GET", "/auth/me", null)
                .then(response => {
                    if (response.success && response.data) {
                        setUser(response.data);
                        setIsLoaded(true);
                    } else {
                        setToken(null);
                        setUser(null);
                        setIsLoaded(true);
                    }
                });
        } else {
            setIsLoaded(true);
        }
    }, [token]);

    const logout = () => {
        FetchRequest("GET", "/auth/logout", null)
            .then(response => {
                if (response.success) {
                    setToken(null);
                    setUser(null);
                    localStorage.removeItem("token");
                }
            });
    };

    return (
        <AuthContext.Provider value={{ token, user, isLoaded, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthProvider;
import React, {useEffect, useState} from "react";
import {faRightToBracket} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import FetchRequest from "../fetchRequest";

const Authorization = () => {
    const [loginData, setLoginData] = useState({
        Login: "",
        Password: "",
    })
    const [failureMessage, setFailureMessage] = useState("")

    const handlerPressEnter = (e) => {
        if (e.key === 'Enter') {
            handlerSendData();
        }
    };

    useEffect(() => {
        document.addEventListener('keydown', handlerPressEnter);
        return () => {
            document.removeEventListener('keydown', handlerPressEnter);
        };
    }, [loginData]);

    const handlerChange = (e) => {
        let { name, value } = e.target

        setLoginData(prevState => ({...prevState, [name]: value}))
    }

    const handlerSendData = () => {
        FetchRequest("POST", "/auth/login", loginData)
            .then(response => {
                if (response.success) {
                    if ("failure" in response.data) {
                        setFailureMessage(response.data.failure)
                    } else {
                        const authToken = response.data.token
                        if (authToken) {
                            localStorage.setItem("token", authToken);
                            window.location.href = "/"
                        }
                    }
                } else {
                    setFailureMessage("Произошла ошибка при отправке запроса")
                }
            })
    }

    return(
        <section className="login">
            <div className="box">
                <div className="form">
                    <FontAwesomeIcon icon={faRightToBracket} />
                        <div>
                            <label>
                                <input type="text" name="Login" required="required" value={loginData.Login} onInput={handlerChange}/>
                                    <span>Логин</span>
                                    <i></i>
                                    <p>{failureMessage}</p>
                            </label>
                            <label>
                                <input type="password" name="Password" required="required" value={loginData.Password} onInput={handlerChange}/>
                                    <span>Пароль</span>
                                    <i></i>
                            </label>
                        </div>
                        <button id="btn" onClick={handlerSendData}>Войти</button>
                </div>
            </div>
        </section>
    )
}

export default Authorization
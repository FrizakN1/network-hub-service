import React, {useEffect, useRef, useState} from "react";
import CustomSelect from "./CustomSelect";
import FetchRequest from "../fetchRequest";
import InputErrorDescription from "./InputErrorDescription";

const UserModalCreate = ({action, setState, returnUser, editUser}) => {
    const validateDebounceTimer = useRef(0)
    const [roles, setRoles] = useState([])
    const [fields, setFields] = useState({
        Login: "",
        Name: "",
        Password: "",
        PasswordConfirm: "",
        Role: {
            ID: 0,
            Value: "",
            TranslateValue: "",
        },
    })
    const [validation, setValidation] = useState({
        Login: true,
        Name: true,
        Password: true,
        PasswordConfirm: true,
        Role: true,
    })

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        if (action === "edit") {
            setFields(prevState => ({
                ...prevState,
                Login: editUser.Login,
                Role: editUser.Role,
                Name: editUser.Name,
            }))
        }
    }, [action, editUser]);

    useEffect(() => {
        FetchRequest("GET", "/get_roles", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setRoles(response.data)
                }
            })
    }, []);


    const validateField = (name, value) => {
        let isValid

        switch (name) {
            case "Name":
            case "Login":
                isValid = value.trim().length > 0
                break
            case "Password":
                isValid = action === "create" ? value.trim().length > 5 : value.trim().length === 0 || value.trim().length > 5
                break
            case "PasswordConfirm":
                isValid = fields.Password === value
                break
            case "Role":
                isValid = value.ID > 0
                break
            default: isValid = true
        }

        setValidation((prevValidation) => ({ ...prevValidation, [name]: isValid }));

        return isValid
    }

    const handlerChange = (e) => {
        let { name, value } = e.target

        setFields(prevState => ({...prevState, [name]: value}))

        clearTimeout(validateDebounceTimer.current)

        validateDebounceTimer.current = setTimeout(() =>  validateField(name, value), 500)
    }

    const handlerSelectRole = (role) => {
        setFields(prevState => ({...prevState, Role: role}))
    }

    const checkChange = (field) => {
        switch (field) {
            case "Login":
            case "Name":
            case "Password":
                return fields[field] !== editUser[field]
            case "Role":
                return fields.Role.ID !== editUser.Role.ID
            default: return false
        }
    }

    const handlerSendData = () => {
        let isFormValid = true;
        let hasChanges = action === "create";

        Object.keys(fields).forEach((field) => {
            if (!validateField(field, fields[field])) {
                isFormValid = false
            }

            if (action === "edit") {
                if (checkChange(field)) {
                    hasChanges = true
                }
            }
        });

        if (!hasChanges) {
            setState(false)
        }

        if (!isFormValid || !hasChanges) {
            return
        }

        let body = {
            Role: fields.Role,
            Login: fields.Login,
            Name: fields.Name,
            Password: fields.Password,
        }

        if (action === "edit") {body = {...editUser, ...fields}}

        FetchRequest("POST", `/${action}_user`, body)
            .then(response => {
                if (response.success && response.data != null) {
                    returnUser(response.data)
                    setState(false)
                }
            })
    }

    return (
        <div className={"modal-window"} onMouseDown={handlerModalCreateClose}>
            <div className="form">
                <h2>{action === "create" ? "Создание пользователя" : "Изменение пользователя"}</h2>
                <div className="fields">
                    <label>
                        <span>ФИО</span>
                        <input type="text" name="Name" value={fields.Name} onChange={handlerChange}/>
                        {!validation.Name && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>
                    <label>
                        <span>Логин</span>
                        <input type="text" name="Login" value={fields.Login} onChange={handlerChange}/>
                        {!validation.Login && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>
                    <label>
                        <span>Пароль</span>
                        <input type="password" name="Password" value={fields.Password} onChange={handlerChange}/>
                        {!validation.Password && <InputErrorDescription text={"Пароль должен состоять минимум из 6 символов"}/>}
                    </label>
                    <label>
                        <span>Подтверждение пароля</span>
                        <input type="password" name="PasswordConfirm" value={fields.PasswordConfirm} onChange={handlerChange}/>
                        {!validation.PasswordConfirm && <InputErrorDescription text={"Пароль не совпадает"}/>}
                    </label>
                    <label>
                        <span>Роль</span>
                        <CustomSelect placeholder="Выбрать" value={fields.Role.TranslateValue} values={roles} setValue={handlerSelectRole}/>
                        {!validation.Role && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>

                    <div className="buttons">
                        <button className={"bg-blue"} onClick={handlerSendData}>{action === "create" ? "Создать" : "Сохранить"}</button>
                        <button className={"bg-red"} onClick={() => setState(false)}>Отмена</button>
                    </div>

                </div>
            </div>
        </div>
    )
}

export default UserModalCreate
import React, {useEffect, useRef, useState} from "react";
import CustomSelect from "./CustomSelect";
import FetchRequest from "../fetchRequest";
import InputErrorDescription from "./InputErrorDescription";

const UserModalCreate = ({action, setState, returnUser, editUser}) => {
    const validateDebounceTimer = useRef(0)
    const [roles, setRoles] = useState([])
    const [fields, setFields] = useState({
        login: "",
        name: "",
        password: "",
        passwordConfirm: "",
        role: {
            id: 0,
            key: "",
            value: "",
        },
    })
    const [validation, setValidation] = useState({
        login: true,
        name: true,
        password: true,
        passwordConfirm: true,
        role: true,
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
                login: editUser.login,
                role: editUser.role,
                name: editUser.name,
            }))
        }
    }, [action, editUser]);

    useEffect(() => {
        FetchRequest("GET", "/users/roles", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setRoles(response.data)
                }
            })
    }, []);


    const validateField = (name, value) => {
        let isValid

        switch (name) {
            case "name":
            case "login":
                isValid = value.trim().length > 0
                break
            case "password":
                isValid = action === "create" ? value.trim().length > 5 : value.trim().length === 0 || value.trim().length > 5
                break
            case "passwordConfirm":
                isValid = fields.password === value
                break
            case "role":
                isValid = value.id > 0
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
        setFields(prevState => ({...prevState, role: role}))
    }

    const checkChange = (field) => {
        switch (field) {
            case "login":
            case "name":
            case "password":
                return fields[field] !== editUser[field]
            case "role":
                return fields.role.id !== editUser.role.id
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
            role_id: fields.role.id,
            login: fields.login,
            name: fields.name,
            password: fields.password,
        }

        if (action === "edit") {body = {...editUser, ...body}}

        FetchRequest(action === "create" ? "POST" : "PUT", `/users`, body)
            .then(response => {
                if (response.success && response.data != null) {
                    let isActive = action === "edit" ? editUser.is_active : true
                    returnUser({...editUser, ...fields, ...response.data, is_active: isActive})
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
                        <input type="text" name="name" value={fields.name} onChange={handlerChange}/>
                        {!validation.name && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>
                    <label>
                        <span>Логин</span>
                        <input type="text" name="login" value={fields.login} onChange={handlerChange}/>
                        {!validation.login && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>
                    <label>
                        <span>Пароль</span>
                        <input type="password" name="password" value={fields.password} onChange={handlerChange}/>
                        {!validation.password && <InputErrorDescription text={"Пароль должен состоять минимум из 6 символов"}/>}
                    </label>
                    <label>
                        <span>Подтверждение пароля</span>
                        <input type="password" name="passwordConfirm" value={fields.passwordConfirm} onChange={handlerChange}/>
                        {!validation.passwordConfirm && <InputErrorDescription text={"Пароль не совпадает"}/>}
                    </label>
                    <label>
                        <span>Роль</span>
                        <CustomSelect placeholder="Выбрать" value={fields.role.value} values={roles} setValue={handlerSelectRole}/>
                        {!validation.role && <InputErrorDescription text={"Поле не может быть пустым"}/>}
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
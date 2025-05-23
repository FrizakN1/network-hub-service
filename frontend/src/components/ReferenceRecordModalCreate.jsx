import React, {useEffect, useRef, useState} from "react";
import InputErrorDescription from "./InputErrorDescription";
import FetchRequest from "../fetchRequest";

const ReferenceRecordModalCreate = ({setState, action, reference, returnRecord, editRecord, withKey}) => {
    const validateDebounceTimer = useRef(0)
    const [fields, setFields] = useState({
        Key: "",
        Value: ""
    })
    const [validation, setValidation] = useState({
        Key: true,
        Value: true
    })

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        if (action === "edit") {
            setFields({
                Key: editRecord.Key,
                Value: editRecord.Value,
            })
        }
    }, [action]);

    const validateField = (name, value) => {
        let isValid = value.trim() !== ""

        if (name === "Key") isValid = isValid || !withKey

        setValidation((prevValidation) => ({ ...prevValidation, [name]: isValid }));

        return isValid
    }

    const handlerChange = (e) => {
        const { name, value } = e.target

        setFields(prevState => ({...prevState, [name]: value}))

        clearTimeout(validateDebounceTimer.current)

        validateDebounceTimer.current = setTimeout(() => validateField(name, value), 500)
    }

    const handlerSendData = () => {
        let isFormValid = true;
        let hasChanges = action === "create";

        Object.keys(fields).forEach((field) => {
            if (!validateField(field, fields[field])) {
                isFormValid = false
            }

            if (action === "edit") {
                if (fields[field] !== editRecord[field]) {
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
            ...fields
        }

        if (action === "edit") {body = {...editRecord, ...body}}

        FetchRequest(action === "create" ? "POST" : "PUT", `/references/${reference}`, body)
            .then(response => {
                if (response.success && response.data != null) {
                    returnRecord(response.data)
                    setState(false)
                }
            })
    }

    return (
        <div className={"modal-window"} onMouseDown={handlerModalCreateClose}>
            <div className="form">
                <h2>
                    {reference === "hardware_types" && <>{action === "create" ? "Создание типа оборудования" : "Изменение типа оборудования"}</>}
                    {reference === "operation_modes" && <>{action === "create" ? "Создание режима работы" : "Изменение режима работы"}</>}
                    {reference === "node_types" && <>{action === "create" ? "Создание типа узла" : "Изменение типа узла"}</>}
                    {reference === "owners" && <>{action === "create" ? "Создание владельца" : "Изменение владельца"}</>}
                    {reference === "roof_types" && <>{action === "create" ? "Создание типа крыши" : "Изменение типа крыши"}</>}
                    {reference === "wiring_types" && <>{action === "create" ? "Создание типа разводки" : "Изменение типа разводки"}</>}
                </h2>
                <div className="fields">
                    {withKey &&
                        <label>
                            <span>Ключ</span>
                            <input type="text" name="Key" value={fields.Key} onChange={handlerChange}/>
                            {!validation.Key && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                        </label>
                    }

                    <label>
                        <span>Значение</span>
                        <input type="text" name="Value" value={fields.Value} onChange={handlerChange}/>
                        {!validation.Value && <InputErrorDescription text={"Поле не может быть пустым"}/>}
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

export default ReferenceRecordModalCreate
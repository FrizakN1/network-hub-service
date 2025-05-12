import React, {useEffect, useRef, useState} from "react";
import InputErrorDescription from "./InputErrorDescription";
import FetchRequest from "../fetchRequest";

const HardwareReferenceRecordModalCreate = ({setState, action, reference, returnRecord, editRecord}) => {
    const validateDebounceTimer = useRef(0)
    const [fields, setFields] = useState({
        Value: "",
        TranslateValue: ""
    })
    const [validation, setValidation] = useState({
        Value: true,
        TranslateValue: true
    })

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        if (action === "edit") {
            setFields({
                Value: editRecord.Value,
                TranslateValue: editRecord.TranslateValue,
            })
        }
    }, [action]);

    const validateField = (name, value) => {
        let isValid = value.trim() !== ""

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
                <h2>{action === "create" ? reference === "hardware_type" ? "Создание типа оборудования" : "Создание режима работы" : reference === "hardware_type" ? "Изменение типа оборудования" : "Изменение режима работы"}</h2>
                <div className="fields">
                    <label>
                        <span>Ключ</span>
                        <input type="text" name="Value" value={fields.Value} onChange={handlerChange}/>
                        {!validation.Value && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>

                    <label>
                        <span>Значение</span>
                        <input type="text" name="TranslateValue" value={fields.TranslateValue} onChange={handlerChange}/>
                        {!validation.TranslateValue && <InputErrorDescription text={"Поле не может быть пустым"}/>}
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

export default HardwareReferenceRecordModalCreate
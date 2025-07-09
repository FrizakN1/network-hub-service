import React, {useEffect, useState} from "react";
import InputErrorDescription from "./InputErrorDescription";
import FetchRequest from "../fetchRequest";

const ReportDataModal = ({setState, editRecord, returnRecord}) => {
    const [fields, setFields] = useState({
        Key: "",
        Value: "",
        Description: "",
    })
    const [validation, setValidation] = useState({
        Value: true
    })

    const handlerModalClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        setFields({
            Key: editRecord.Key,
            Value: editRecord.Value,
            Description: editRecord.Description.String
        })
    }, [editRecord]);

    const handlerChange = (e) => {
        let { name, value } = e.target

        setFields(prevState => ({...prevState, [name]: value}))

        if (name === "Value") {
            setValidation({Value: value.trim().length > 0})
        }
    }

    const checkChange = (field) => {
        switch (field) {
            case "Value":
                return fields[field] !== editRecord[field]
            case "Description":
                return fields[field] !== editRecord[field].String
            default: return false
        }
    }

    const handlerSendData = () => {
        let hasChanges = false;

        Object.keys(fields).forEach((field) => {
            if (checkChange(field)) {
                hasChanges = true
            }
        });

        if (!hasChanges) {
            setState(false)
        }

        if (!validation.Value || !hasChanges) {
            return
        }

        let body = {
            ID: editRecord.ID,
            Key: editRecord.Key,
            Value: fields.Value,
            Description: {
                String: fields.Description,
                Valid: fields.Description.trim().length > 0
            }
        }

        FetchRequest("PUT", "/report", body)
            .then(response => {
                if (response.success) {
                    returnRecord(response.data)
                    setState(false)
                }
            })
    }

    return (
        <div className={"modal-window"} onMouseDown={handlerModalClose}>
            <div className="form">
                <h2>Изменение значения ключа {editRecord.Key}</h2>
                <div className="fields">
                    <label>
                        <span>Ключ</span>
                        <input type="text" value={fields.Key} disabled={true}/>
                    </label>

                    <label>
                        <span>Значение</span>
                        <input type="text" name="Value" value={fields.Value} onChange={handlerChange}/>
                        {!validation.Value && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>

                    <label>
                        <span>Описание</span>
                        <input type="text" name="Description" value={fields.Description} onChange={handlerChange}/>
                    </label>

                    <div className="buttons">
                        <button className={"bg-blue"} onClick={handlerSendData}>Сохранить</button>
                        <button className={"bg-red"} onClick={() => setState(false)}>Отмена</button>
                    </div>

                </div>
            </div>
        </div>
    )
}

export default ReportDataModal
import React, {useEffect, useRef, useState} from "react";
import FetchRequest from "../fetchRequest";
import InputErrorDescription from "./InputErrorDescription";

const NodeReferenceRecordModalCreate = ({action, setState, returnRecord, editRecord, reference}) => {
    const [name, setName] = useState({
        Value: "",
        Valid: true,
    })

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        if (action === "edit") {
            setName({
                Value: editRecord.Name,
                Valid: true,
            })
        }
    }, [action, editRecord]);


    const handlerChange = (e) => {
        setName({
            Value: e.target.value,
            Valid: e.target.value.trim().length > 0
        })
    }

    const handlerSendData = () => {
        if (name.Value.trim().length === 0) {
            setName(prevState => ({Value: prevState.Value, Valid: false}))
            return
        }

        let body = {
            Name: name.Value
        }

        if (action === "edit") {body = {...editRecord, Name: name.Value}}

        FetchRequest("POST", `/${action}_${reference.slice(0, reference.length-1)}`, body)
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
                <h2>{action === "create" ? reference === "node_types" ? "Создание тип узла" : "Создание владельца" : reference === "node_types" ? "Изменение тип узла" : "Изменение владельца"}</h2>
                <div className="fields">
                    <label>
                        <span>Наименование</span>
                        <input type="text" name="Name" value={name.Value} onChange={handlerChange}/>
                        {!name.Valid && <InputErrorDescription text={"Поле не может быть пустым"}/>}
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

export default NodeReferenceRecordModalCreate
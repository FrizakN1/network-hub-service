import React, {useEffect, useState} from "react";
import ModalSelectTable from "./ModalSelectTable";
import InputErrorDescription from "./InputErrorDescription";
import CustomSelect from "./CustomSelect";
import FetchRequest from "../fetchRequest";
import {useParams} from "react-router-dom";

const HouseParamsModalEdit = ({setState, editData, returnData}) => {
    const [roofTypes, setRoofTypes] = useState([])
    const [wiringTypes, setWiringTypes] = useState([])
    const [fields, setFields] = useState({
        RoofType: {ID: 0, Value: ""},
        WiringType: {ID: 0, Value: ""},
    })
    const { id } = useParams()

    const handlerModalClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        setFields({
            RoofType: editData.RoofType,
            WiringType: editData.WiringType
        })
    }, [editData])

    useEffect(() => {
        FetchRequest("GET", "/references/roof_types")
            .then(response => {
                if (response.success) {
                    setRoofTypes(response.data != null ? response.data : [])
                }
            })

        FetchRequest("GET", "/references/wiring_types")
            .then(response => {
                if (response.success) {
                    setWiringTypes(response.data != null ? response.data : [])
                }
            })
    }, []);

    const handlerSelectRoofType = (roofType) => {
        setFields(prevState => ({...prevState, RoofType: roofType}))
    }

    const handlerSelectWiringType = (wiringType) => {
        setFields(prevState => ({...prevState, WiringType: wiringType}))
    }

    const handlerSendData = () => {
        if (editData.WiringType.ID === fields.WiringType.ID && editData.RoofType.ID === fields.RoofType.ID) {
            setState(false)
            return
        }

        FetchRequest("POST", `/houses/${id}/params`, fields)
            .then(response => {
                if (response.success) {
                    returnData(response.data)
                }
            })
    }

    return (
        <div className={"modal-window"} onMouseDown={handlerModalClose}>
            <div className="form">
                <h2>Изменение параметров дома</h2>
                <div className="fields">
                    <label>
                        <span>Тип крыши</span>
                        <CustomSelect placeholder="Выбрать" value={fields.RoofType.Value} values={roofTypes} setValue={handlerSelectRoofType}/>
                    </label>
                    <label>
                        <span>Тип разводки</span>
                        <CustomSelect placeholder="Выбрать" value={fields.WiringType.Value} values={wiringTypes} setValue={handlerSelectWiringType}/>
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

export default HouseParamsModalEdit
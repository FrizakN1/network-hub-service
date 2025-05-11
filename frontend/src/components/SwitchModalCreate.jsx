import React, {useEffect, useRef, useState} from "react";
import InputErrorDescription from "./InputErrorDescription";
import CustomSelect from "./CustomSelect";
import FetchRequest from "../fetchRequest";

const SwitchModalCreate = ({action, setState, returnSwitch, editSwitch}) => {
    const validateDebounceTimer = useRef(0)
    const [fields, setFields] = useState({
        Name: "",
        OperationMode: {ID: 0, Value: "", TranslateValue: ""},
        CommunityRead: "",
        CommunityWrite: "",
        PortAmount: 0,
        FirmwareOID: "",
        SystemNameOID: "",
        SerialNumberOID: "",
        SaveConfigOID: "",
        PortDescOID: "",
        VlanOID: "",
        PortUntaggedOID: "",
        SpeedOID: "",
        BatteryStatusOID: "",
        BatteryChargeOID: "",
        PortModeOID: "",
        UptimeOID: "",
    })
    const [validation, setValidation] = useState({
        Name: true,
        PortAmount: true,
    })
    const [operationModes, setOperationModes] = useState([])

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        FetchRequest("GET", "/references/operation_modes", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setOperationModes(response.data)
                }
            })
    }, [])

    useEffect(() => {
        if (action === "edit") {
            setFields({
                Name: editSwitch.Name,
                OperationMode: editSwitch.OperationMode,
                CommunityRead: editSwitch.CommunityRead.String,
                CommunityWrite: editSwitch.CommunityWrite.String,
                PortAmount: editSwitch.PortAmount,
                FirmwareOID: editSwitch.FirmwareOID.String,
                SystemNameOID: editSwitch.SystemNameOID.String,
                SerialNumberOID: editSwitch.SerialNumberOID.String,
                SaveConfigOID: editSwitch.SaveConfigOID.String,
                PortDescOID: editSwitch.PortDescOID.String,
                VlanOID: editSwitch.VlanOID.String,
                PortUntaggedOID: editSwitch.PortUntaggedOID.String,
                SpeedOID: editSwitch.SpeedOID.String,
                BatteryStatusOID: editSwitch.BatteryStatusOID.String,
                BatteryChargeOID: editSwitch.BatteryChargeOID.String,
                PortModeOID: editSwitch.PortDescOID.String,
                UptimeOID: editSwitch.UptimeOID.String,
            })
        }
    }, [action, editSwitch]);

    const validateField = (name, value) => {
        let isValid

        switch (name) {
            case "Name":
                isValid = value.trim().length > 0
                break
            case "PortAmount":
                isValid = value > 0
                break
            default: isValid = true
        }

        setValidation(prevState => ({...prevState, [name]: value}))

        return isValid
    }

    const handlerChange = (e) => {
        const { name, value } = e.target

        setFields(prevState => ({...prevState, [name]: value}))

        clearTimeout(validateDebounceTimer.current)

        validateDebounceTimer.current = setTimeout(() => validateField(name, value), 500)
    }

    const handlerSelectOperationMode = (operationMode) => {
        setFields(prevState => ({...prevState, OperationMode: operationMode}))
    }

    const checkChange = (field) => {
        if (field === "OperationMode") {
            return fields[field].ID !== editSwitch[field].ID
        }

        return fields[field] !== editSwitch[field]
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
            Name: fields.Name,
            OperationMode: fields.OperationMode,
            CommunityRead: {String: fields.CommunityRead, Valid: fields.CommunityRead !== ""},
            CommunityWrite: {String: fields.CommunityWrite, Valid: fields.CommunityWrite !== ""},
            PortAmount: Number(fields.PortAmount),
            FirmwareOID: {String: fields.FirmwareOID, Valid: fields.FirmwareOID !== ""},
            SystemNameOID: {String: fields.SystemNameOID, Valid: fields.SystemNameOID !== ""},
            SerialNumberOID: {String: fields.SerialNumberOID, Valid: fields.SerialNumberOID !== ""},
            SaveConfigOID: {String: fields.SaveConfigOID, Valid: fields.SaveConfigOID !== ""},
            PortDescOID: {String: fields.PortDescOID, Valid: fields.PortDescOID !== ""},
            VlanOID: {String: fields.VlanOID, Valid: fields.VlanOID !== ""},
            PortUntaggedOID: {String: fields.PortUntaggedOID, Valid: fields.PortUntaggedOID !== ""},
            SpeedOID: {String: fields.SpeedOID, Valid: fields.SpeedOID !== ""},
            BatteryStatusOID: {String: fields.BatteryStatusOID, Valid: fields.BatteryStatusOID !== ""},
            BatteryChargeOID: {String: fields.BatteryChargeOID, Valid: fields.BatteryChargeOID !== ""},
            PortModeOID: {String: fields.PortModeOID, Valid: fields.PortModeOID !== ""},
            UptimeOID: {String: fields.UptimeOID, Valid: fields.UptimeOID !== ""},
        }

        if (action === "edit") {body = {...editSwitch, ...body}}

        FetchRequest(action === "create" ? "POST" : "PUT", `/switches`, body)
            .then(response => {
                if (response.success && response.data != null) {
                    returnSwitch(response.data)
                    setState(false)
                }
            })
    }

    return (
        <div className={"modal-window"} onMouseDown={handlerModalCreateClose}>
            <div className="form switch">
                <h2>{action === "create" ? "Создание коммутатора" : "Изменение коммутатора"}</h2>
                <div className="fields">
                    <div className="row">
                        <label>
                            <span>Название</span>
                            <input type="text" name="Name" value={fields.Name} onChange={handlerChange}/>
                            {!validation.Name && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                        </label>

                        <label>
                            <span>Количество портов</span>
                            <input type="number" name="PortAmount" value={fields.PortAmount} onChange={handlerChange}/>
                            {!validation.PortAmount && <InputErrorDescription text={"Значение не может быть меньше единицы"}/>}
                        </label>
                    </div>

                    <div className="row">
                        <div className="column">
                            <label>
                                <span>Режим работы</span>
                                <CustomSelect placeholder="Выбрать" value={fields.OperationMode.TranslateValue} values={operationModes} setValue={handlerSelectOperationMode}/>
                            </label>
                            <label>
                                <span>SNMP-комьюнити (только чтение)</span>
                                <input type="text" name="CommunityRead" value={fields.CommunityRead} onChange={handlerChange}/>
                            </label>
                            <label>
                                <span>SNMP-комьюнити (чтение и запись)</span>
                                <input type="text" name="CommunityWrite" value={fields.CommunityWrite} onChange={handlerChange}/>
                            </label>
                            <label>
                                <span>Firmware OID</span>
                                <input type="text" name="FirmwareOID" value={fields.FirmwareOID} onChange={handlerChange}/>
                            </label>
                            <label>
                                <span>System Name OID</span>
                                <input type="text" name="SystemNameOID" value={fields.SystemNameOID} onChange={handlerChange}/>
                            </label>
                        </div>

                        <div className="column">
                            <label>
                                <span>Serial Number OID</span>
                                <input type="text" name="SerialNumberOID" value={fields.SerialNumberOID} onChange={handlerChange}/>
                            </label>
                            <label>
                                <span>Сохранение конфигурации OID</span>
                                <input type="text" name="SaveConfigOID" value={fields.SaveConfigOID} onChange={handlerChange}/>
                            </label>
                            <label>
                                <span>Описание портов OID</span>
                                <input type="text" name="PortDescOID" value={fields.PortDescOID} onChange={handlerChange}/>
                            </label>
                            <label>
                                <span>Скорость портов OID</span>
                                <input type="text" name="SpeedOID" value={fields.SpeedOID} onChange={handlerChange}/>
                            </label>
                            <label>
                                <span>Время работы OID</span>
                                <input type="text" name="UptimeOID" value={fields.UptimeOID} onChange={handlerChange}/>
                            </label>
                        </div>

                        <div className="column">
                            <label>
                                <span>VLAN OID</span>
                                <input type="text" name="VlanOID" value={fields.VlanOID} onChange={handlerChange}/>
                            </label>
                            {fields.OperationMode.Value === "dlink" &&
                                <label>
                                    <span>Untagged порты OID</span>
                                    <input type="text" name="PortUntaggedOID" value={fields.PortUntaggedOID} onChange={handlerChange}/>
                                </label>
                            }
                            {fields.OperationMode.Value === "eltex" && <>
                                <label>
                                    <span>Статус батареи OID</span>
                                    <input type="text" name="BatteryStatusOID" value={fields.BatteryStatusOID} onChange={handlerChange}/>
                                </label>
                                <label>
                                    <span>Заряд батареи OID</span>
                                    <input type="text" name="BatteryChargeOID" value={fields.BatteryChargeOID} onChange={handlerChange}/>
                                </label>
                                <label>
                                    <span>Режим портов OID</span>
                                    <input type="text" name="PortModeOID" value={fields.PortModeOID} onChange={handlerChange}/>
                                </label>
                            </>}
                        </div>
                    </div>


                    <div className="buttons">
                        <button className={"bg-blue"} onClick={handlerSendData}>{action === "create" ? "Создать" : "Сохранить"}</button>
                        <button className={"bg-red"} onClick={() => setState(false)}>Отмена</button>
                    </div>

                </div>
            </div>
        </div>
    )
}

export default SwitchModalCreate
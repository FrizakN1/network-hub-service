import React, {useEffect, useRef, useState} from "react";
import InputErrorDescription from "./InputErrorDescription";
import CustomSelect from "./CustomSelect";
import FetchRequest from "../fetchRequest";
import ModalSelectTable from "./ModalSelectTable";

const HardwareModalCreate = ({action, setState, returnHardware, editHardwareID}) => {
    const validateDebounceTimer = useRef(0)
    const [fields, setFields] = useState({
        Node: {ID: 0, Name: ""},
        Type: {ID: 0, Key: "", Value: ""},
        Switch: {ID: 0, Name: ""},
        IpAddress: "",
        MgmtVlan: "",
        Description: "",
    })
    const [validation, setValidation] = useState({
        Node: true,
        Type: true,
        Switch: true,
        IpAddress: true,
    })
    const [hardwareTypes, setHardwareTypes] = useState([])
    const [switches, setSwitches] = useState([])
    const [modalSelectTable, setModalSelectTable] = useState(false)
    const [editHardware, setEditHardware] = useState({})

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        if (action === "edit") {
            FetchRequest("GET", `/hardware/${editHardwareID}`)
                .then(response => {
                    if (response.success) {
                        setEditHardware(response.data)

                        setFields({
                            Node: response.data.Node,
                            Type: response.data.Type,
                            Switch: response.data.Switch,
                            IpAddress: response.data.IpAddress.String,
                            MgmtVlan: response.data.MgmtVlan.String,
                            Description: response.data.Description.String,
                        })
                    }
                })
        }
    }, [action, editHardwareID]);

    useEffect(() => {
        FetchRequest("GET", "/switches", null)
            .then(response => {
                if (response.success) {
                    setSwitches(response.data != null ? response.data : [])
                }
            })

        FetchRequest("GET", "/references/hardware_types", null)
            .then(response => {
                if (response.success) {
                    setHardwareTypes(response.data != null ? response.data : [])
                }
            })
    }, []);

    const handlerSelectHardware = (hardwareType) => {
        setFields(prevState => ({...prevState, Type: hardwareType}))
    }

    const handlerSelectSwitch = (_switch) => {
        setFields(prevState => ({...prevState, Switch: _switch}))
    }

    const validateField = (name, value) => {
        let isValid

        switch (name) {
            case "Type":
            case "Node":
                isValid = value.ID !== 0
                break
            case "Switch":
                isValid = fields.Type.Key !== "switch" || value.ID !== 0
                break
            case "IpAddress":
                isValid = fields.Type.Key !== "switch" || value.trim() !== ""
                break
            default: isValid = true
        }

        setValidation(prevState => ({...prevState, [name]: isValid}))

        return isValid
    }

    const handlerChange = (e) => {
        const { name, value } = e.target

        setFields(prevState => ({...prevState, [name]: value}))

        clearTimeout(validateDebounceTimer.current)

        validateDebounceTimer.current = setTimeout(() => validateField(name, value), 500)
    }

    const checkChange = (field) => {
        switch (field) {
            case "Type":
            case "Switch":
            case "Node":
                return fields[field].ID !== editHardware[field].ID
            case "IpAddress":
            case "MgmtVlan":
            case "Description":
                return fields[field] !== editHardware[field].String
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
            Node: fields.Node,
            Type: fields.Type,
            Switch: fields.Switch,
            IpAddress: {String: fields.IpAddress, Valid: fields.IpAddress.trim() !== ""},
            MgmtVlan: {String: fields.MgmtVlan, Valid: fields.MgmtVlan.trim() !== ""},
            Description: {String: fields.Description, Valid: fields.Description.trim() !== ""},
        }

        if (action === "edit") {body = {...editHardware, ...body}}

        FetchRequest(action === "create" ? "POST" : "PUT", `/hardware`, body)
            .then(response => {
                if (response.success && response.data != null) {
                    returnHardware(response.data)
                    setState(false)
                }
            })
    }

    const handlerSelectNode = (node) => {
        setFields(prevState => ({...prevState, Node: node}))
    }

    return (
        <div className={"modal-window"} onMouseDown={handlerModalCreateClose}>
            {modalSelectTable && <ModalSelectTable setState={setModalSelectTable} alreadySelect={fields.Node} selectRecord={handlerSelectNode} />}
            <div className="form">
                <h2>{action === "create" ? "Создание оборудования" : "Изменение оборудования"}</h2>
                <div className="fields">
                    <label>
                        <span>Узел</span>
                        <div className="select-field" onClick={() => setModalSelectTable(true)}>{fields.Node.Name === "" ? "Выбрать..." : fields.Node.Name}</div>
                        {!validation.Node && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>
                    <label>
                        <span>Тип оборудования</span>
                        <CustomSelect placeholder="Выбрать" value={fields.Type.Value} values={hardwareTypes} setValue={handlerSelectHardware}/>
                        {!validation.Type && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                    </label>
                    {fields.Type.Key === "switch" && <>
                        <label>
                            <span>Модель</span>
                            <CustomSelect placeholder="Выбрать" value={fields.Switch.Name} values={switches} setValue={handlerSelectSwitch}/>
                            {!validation.Switch && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                        </label>
                        <label>
                            <span>IP адрес</span>
                            <input type="text" name="IpAddress" value={fields.IpAddress} onChange={handlerChange}/>
                            {!validation.IpAddress && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                        </label>
                        <label>
                            <span>Управляющий VLAN</span>
                            <input type="text" name="MgmtVlan" value={fields.MgmtVlan} onChange={handlerChange}/>
                        </label>
                    </>}
                    <label>
                        <span>Описание</span>
                        <textarea name="Description" cols="30" rows="7" value={fields.Description} onChange={handlerChange}></textarea>
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

export default HardwareModalCreate
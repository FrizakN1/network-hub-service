import React, {useEffect, useRef, useState} from "react";
import InputErrorDescription from "./InputErrorDescription";
import CustomSelect from "./CustomSelect";
import ModalSelectTable from "./ModalSelectTable";
import FetchRequest from "../fetchRequest";
import {useParams} from "react-router-dom";
import SearchInput from "./SearchInput";

const NodeModalCreate = ({action, setState, editNode, returnNode, defaultAddress = null}) => {
    const validateDebounceTimer = useRef(0)
    const [modalSelectTable, setModalSelectTable] = useState({
        State: false,
        Uri: "",
        Type: "",
        SelectRecord: null
    })
    const [fields, setFields] = useState({
        Parent: {ID: 0, Name: ""},
        Type: {ID: 0, Name: ""},
        Owner: {ID: 0, Name: ""},
        HouseID: 0,
        Address: {
            Street: {
                Name: "",
                Type: {ShortName: ""}
            },
            House: {
                Name: "",
                Type: {ShortName: ""}
            }
        },
        Name: "",
        Zone: "",
        Placement: "",
        Supply: "",
        Access: "",
        Description: "",
    })
    const [validation, setValidation] = useState({
        Parent: true,
        Type: true,
        Owner: true,
        Name: true,
        Address: true
    })
    const [generalNode, setGeneralNode] = useState(false)
    const [nodeTypes, setNodeTypes] = useState([])
    const [owners, setOwners] = useState([])

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        FetchRequest("GET", "/get_node_types", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setNodeTypes(response.data)
                }
            })

        FetchRequest("GET", "/get_owners", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setOwners(response.data)
                }
            })
    }, []);

    useEffect(() => {
        if (action === "edit") {
            setGeneralNode(editNode.Parent == null)

            setFields({
                Parent: editNode.Parent != null ? editNode.Parent : {ID: 0, Name: ""},
                Address: editNode.Address,
                Type: editNode.Type,
                Owner: editNode.Owner,
                Name: editNode.Name,
                Zone: editNode.Zone.String,
                Placement: editNode.Placement.String,
                Supply: editNode.Supply.String,
                Access: editNode.Access.String,
                Description: editNode.Description.String,
            })
        } else if (defaultAddress != null) {
            setFields(prevState => ({...prevState, Address: defaultAddress}))
        }
    }, [action, editNode, defaultAddress]);

    const validateField = (name, value) => {
        let isValid = true

        switch (name) {
            case "Name":
                isValid = value.trim().length > 0
                break
            case "Type":
            case "Owner":
                isValid = value.ID > 0
                break
            case "Parent":
                isValid = generalNode || value.ID > 0
                break
            case "Address":
                isValid = value.House.ID > 0
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

    const checkChange = (field) => {
        switch (field) {
            case "Name":
            case "Zone":
            case "Placement":
            case "Supply":
            case "Access":
            case "Description":
                return fields[field] !== editNode[field]
            case "Type":
            case "Owner":
                return fields[field].ID !== editNode[field].ID
            case "Parent":
                return editNode.Parent != null ? fields.Parent.ID !== editNode[field].ID : fields.Parent !== editNode.Parent
            case "Address":
                return editNode.Address.House.ID !== fields.Address.House.ID
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
            Parent: generalNode ? null : fields.Parent,
            Address: fields.Address,
            Type: fields.Type,
            Owner: fields.Owner,
            Name: fields.Name,
            Zone: {String: fields.Zone, Valid: fields.Zone !== ""},
            Placement: {String: fields.Placement, Valid: fields.Placement !== ""},
            Supply: {String: fields.Supply, Valid: fields.Supply !== ""},
            Access: {String: fields.Access, Valid: fields.Access !== ""},
            Description: {String: fields.Description, Valid: fields.Description !== ""},
        }

        if (action === "edit") {body = {...editNode, ...body}}

        FetchRequest("POST", `/${action}_node`, body)
            .then(response => {
                if (response.success && response.data != null) {
                    returnNode(response.data)
                    setState(false)
                }
            })
    }

    const handlerSelectParent = (node) => {
        setFields(prevState => ({...prevState, Parent: node}))
    }

    const handlerSelectNodeType = (nodeType) => {
        setFields(prevState => ({...prevState, Type: nodeType}))
    }

    const handlerSelectOwner = (owner) => {
        setFields(prevState => ({...prevState, Owner: owner}))
    }

    const handlerSelectAddress = (address) => {
        setFields(prevState => ({...prevState, Address: address}))
    }

    return (
        <div className={"modal-window"} onMouseDown={handlerModalCreateClose}>
            {modalSelectTable.State && <ModalSelectTable uri={modalSelectTable.Uri} alreadySelect={fields.Parent} type={modalSelectTable.Type} selectRecord={modalSelectTable.SelectRecord} setState={(state) => setModalSelectTable(prevState => ({...prevState, State: state}))}/>}
            <div className="form node">
                <h2>{action === "create" ? "Создание узла" : "Изменение узла"}</h2>
                <div className="fields">
                    <label>
                        <span>Адрес</span>
                        <SearchInput defaultValue={fields.Address.House.ID > 0 ? `${fields.Address.Street.Type.ShortName} ${fields.Address.Street.Name}, ${fields.Address.House.Type.ShortName} ${fields.Address.House.Name}` : ""} action="select" returnAddress={handlerSelectAddress}/>
                        {!validation.Address && <InputErrorDescription text={"Некорректный адрес"}/>}
                    </label>
                   <div className="row">
                       <div className="column">
                           <label>
                               <span>Родительский узел</span>
                               <div className="select-field" onClick={() => setModalSelectTable({State: true, Uri: "", Type: "node", SelectRecord: handlerSelectParent})}>{fields.Parent.Name === "" ? "Выбрать..." : fields.Parent.Name}</div>
                               {!validation.Parent && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                           </label>

                           <label className="checkbox">
                               <input type="checkbox" checked={generalNode} onChange={() => setGeneralNode(prevState => !prevState)}/>
                               <span>Все свое мужское ставлю, что у этого узла нет родителя, а не потому что мне лень</span>
                           </label>

                           <label>
                               <span>Тип узла</span>
                               <CustomSelect placeholder="Выбрать" value={fields.Type.Name} values={nodeTypes} setValue={handlerSelectNodeType}/>
                               {!validation.Type && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                           </label>
                           {/*<label>*/}
                           {/*    <span>Тип узла</span>*/}
                           {/*    <div className="select-field" onClick={() => setModalSelectTable({State: true, Uri: "/get_node_types", Type: "node_type", SelectRecord: handlerSelectNodeType})}>{fields.Type.Name === "" ? "Выбрать..." : fields.Type.Name}</div>*/}
                           {/*    {!validation.Type && <InputErrorDescription text={"Поле не может быть пустым"}/>}*/}
                           {/*</label>*/}

                           <label>
                               <span>Владелец узла</span>
                               <CustomSelect placeholder="Выбрать" value={fields.Owner.Name} values={owners} setValue={handlerSelectOwner}/>
                               {!validation.Owner && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                           </label>
                           {/*<label>*/}
                           {/*    <span>Владелец узла</span>*/}
                           {/*    <div className="select-field" onClick={() => setModalSelectTable({State: true, Uri: "/get_owners", Type: "owner", SelectRecord: handlerSelectOwner})}>{fields.Owner.Name === "" ? "Выбрать..." : fields.Owner.Name}</div>*/}
                           {/*    {!validation.Owner && <InputErrorDescription text={"Поле не может быть пустым"}/>}*/}
                           {/*</label>*/}

                           <label>
                               <span>Название</span>
                               <input type="text" name="Name" value={fields.Name} onChange={handlerChange}/>
                               {!validation.Name && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                           </label>

                           <label>
                               <span>Район</span>
                               <input type="text" name="Zone" value={fields.Zone} onChange={handlerChange}/>
                           </label>
                       </div>

                       <div className="column">
                           <label>
                               <span>Расположение узла</span>
                               <textarea name="Placement" cols="30" rows="7" value={fields.Placement} onChange={handlerChange}></textarea>
                           </label>

                           <label>
                               <span>Питание узла</span>
                               <textarea name="Supply" cols="30" rows="7" value={fields.Supply} onChange={handlerChange}></textarea>
                           </label>
                       </div>

                       <div className="column">
                           <label>
                               <span>Доступ к узлу</span>
                               <textarea name="Access" cols="30" rows="7" value={fields.Access} onChange={handlerChange}></textarea>
                           </label>

                           <label>
                               <span>Описание</span>
                               <textarea name="Description" cols="30" rows="7" value={fields.Description} onChange={handlerChange}></textarea>
                           </label>
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

export default NodeModalCreate
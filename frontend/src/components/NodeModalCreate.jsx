import React, {useEffect, useRef, useState} from "react";
import InputErrorDescription from "./InputErrorDescription";
import CustomSelect from "./CustomSelect";
import ModalSelectTable from "./ModalSelectTable";
import FetchRequest from "../fetchRequest";
import SearchInput from "./SearchInput";

const NodeModalCreate = ({action, setState, editNodeID, returnNode, defaultAddress = null}) => {
    const validateDebounceTimer = useRef(0)
    const [modalSelectTable, setModalSelectTable] = useState({
        State: false,
        Uri: "",
        Type: "",
        SelectRecord: null
    })
    const [fields, setFields] = useState({
        Parent: {ID: 0, Name: ""},
        Type: {ID: 0, Value: ""},
        Owner: {ID: 0, Value: ""},
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
        IsPassive: false
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
    const [editNode, setEditNode] = useState({})

    const handlerModalCreateClose = (e) => {
        if (e.target.className === "modal-window") {
            setState(false)
        }
    }

    useEffect(() => {
        FetchRequest("GET", "/references/node_types", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setNodeTypes(response.data)
                }
            })

        FetchRequest("GET", "/references/owners", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setOwners(response.data)
                }
            })
    }, []);

    useEffect(() => {
        if (action === "edit") {
            FetchRequest("GET", `/nodes/${editNodeID}`)
                .then(response => {
                    if (response.success) {
                        setEditNode(response.data)

                        setGeneralNode(response.data.Parent == null)
                        
                        setFields({
                            Parent: response.data.Parent != null ? response.data.Parent : {ID: 0, Name: ""},
                            Address: response.data.Address,
                            Type: response.data.Type != null ? response.data.Type : {ID: 0, Value: ""},
                            Owner: response.data.Owner,
                            Name: response.data.Name,
                            Zone: response.data.Zone.String,
                            Placement: response.data.Placement.String,
                            Supply: response.data.Supply.String,
                            Access: response.data.Access.String,
                            Description: response.data.Description.String,
                            IsPassive: response.data.IsPassive
                        })
                    }
                })
        } else if (defaultAddress != null) {
            setFields(prevState => ({...prevState, Address: defaultAddress}))
        }
    }, [action, editNodeID, defaultAddress]);

    const validateField = (name, value) => {
        let isValid = true

        switch (name) {
            case "Name":
                isValid = value.trim().length > 0
                break
            case "Type":
                isValid = value.ID > 0 || fields.IsPassive
                break
            case "Owner":
                isValid = value.ID > 0
                break
            case "Parent":
                isValid = generalNode || fields.IsPassive || value.ID > 0
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
            case "IsPassive":
                return fields[field] !== editNode[field]
            case "Type":
                return fields.IsPassive || (editNode.Type != null ? fields.Type.ID !== editNode[field].ID : fields.Type !== editNode.Type)
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
            Parent: generalNode || fields.IsPassive ? null : fields.Parent,
            Address: fields.Address,
            Type: fields.IsPassive ? null : fields.Type,
            Owner: fields.Owner,
            Name: fields.Name,
            Zone: {String: fields.Zone, Valid: fields.Zone !== ""},
            Placement: {String: fields.Placement, Valid: fields.Placement !== ""},
            Supply: {String: fields.Supply, Valid: fields.Supply !== ""},
            Access: {String: fields.Access, Valid: fields.Access !== ""},
            Description: {String: fields.Description, Valid: fields.Description !== ""},
            IsPassive: fields.IsPassive,
        }

        if (action === "edit") {body = {...editNode, ...body}}

        FetchRequest(action === "create" ? "POST" : "PUT", `/nodes`, body)
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
                           <label className="checkbox">
                               <input type="checkbox" name="IsPassive" checked={fields.IsPassive} onChange={() => setFields(prevState => ({...prevState, IsPassive: !prevState.IsPassive}))}/>
                               <span>Узел пассивный</span>
                           </label>

                           {!fields.IsPassive && <>
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
                                   <CustomSelect placeholder="Выбрать" value={fields.Type.Value} values={nodeTypes} setValue={handlerSelectNodeType}/>
                                   {!validation.Type && <InputErrorDescription text={"Поле не может быть пустым"}/>}
                               </label>
                           </>}
                           {/*<label>*/}
                           {/*    <span>Тип узла</span>*/}
                           {/*    <div className="select-field" onClick={() => setModalSelectTable({State: true, Uri: "/get_node_types", Type: "node_type", SelectRecord: handlerSelectNodeType})}>{fields.Type.Name === "" ? "Выбрать..." : fields.Type.Name}</div>*/}
                           {/*    {!validation.Type && <InputErrorDescription text={"Поле не может быть пустым"}/>}*/}
                           {/*</label>*/}

                           <label>
                               <span>Владелец узла</span>
                               <CustomSelect placeholder="Выбрать" value={fields.Owner.Value} values={owners} setValue={handlerSelectOwner}/>
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
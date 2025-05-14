import React, {useContext, useEffect, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faPen, faPlus} from "@fortawesome/free-solid-svg-icons";
import FetchRequest from "../fetchRequest";
import NodeReferenceRecordModalCreate from "./NodeReferenceRecordModalCreate";
import HardwareReferenceRecordModalCreate from "./HardwareReferenceRecordModalCreate";
import SwitchModalCreate from "./SwitchModalCreate";
import AuthContext from "../context/AuthContext";

const ReferencePage = ({reference}) => {
    const [records, setRecords] = useState([])
    const [nodeReferenceModalCreate, setNodeReferenceModalCreate] = useState(false)
    const [nodeReferenceModalEdit, setNodeReferenceModalEdit]= useState({
        State: false,
        EditRecord: null
    })
    const [hardwareReferenceModalCreate, setHardwareReferenceModalCreate] = useState(false)
    const [hardwareReferenceModalEdit, setHardwareReferenceModalEdit] = useState({
        State: false,
        EditRecord: null
    })
    const [switchModalCreate, setSwitchModalCreate] = useState(false)
    const [switchModalEdit, setSwitchModalEdit] = useState({
        State: false,
        EditRecord: null
    })
    const [isLoaded, setIsLoaded] = useState(true)
    const { user } = useContext(AuthContext)

    useEffect(() => {
        FetchRequest("GET", `/references/${reference}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    setRecords(response.data)
                }

                setIsLoaded(true)
            })
    }, []);

    const handlerOpenModalCreate = () => {
        switch (reference) {
            case "node_types":
            case "owners":
                setNodeReferenceModalCreate(true)
                break
            case "hardware_types":
            case "operation_modes":
                setHardwareReferenceModalCreate(true)
                break
            case "switches":
                setSwitchModalCreate(true)
                break
        }
    }

    const handlerOpenModalEdit = (record) => {
        switch (reference) {
            case "node_types":
            case "owners":
                setNodeReferenceModalEdit({State: true, EditRecord: record})
                break
            case "hardware_types":
            case "operation_modes":
                setHardwareReferenceModalEdit({State: true, EditRecord: record})
                break
            case "switches":
                setSwitchModalEdit({State: true, EditRecord: record})
                break
        }
    }

    const handlerAddRecord = (record) => {
        setRecords(prevState => [...prevState, record])
    }

    const handlerEditRecord = (record) => {
        setRecords(prevState => prevState.map(_record => record.ID === _record.ID ? record : _record))
    }

    return (
        <section className="references">
            {user.Role.Value !== "user" && <>
                {switchModalCreate && <SwitchModalCreate action="create" setState={setSwitchModalCreate} returnSwitch={handlerAddRecord}/>}
                {nodeReferenceModalCreate && <NodeReferenceRecordModalCreate action="create" setState={setNodeReferenceModalCreate} returnRecord={handlerAddRecord} reference={reference}/>}
                {hardwareReferenceModalCreate && <HardwareReferenceRecordModalCreate action="create" setState={setHardwareReferenceModalCreate} returnRecord={handlerAddRecord} reference={reference}/>}
                {nodeReferenceModalEdit.State && <NodeReferenceRecordModalCreate action="edit"
                                                                                 setState={(state) => setNodeReferenceModalEdit(prevState => ({...prevState, State: state}))}
                                                                                 returnRecord={handlerEditRecord}
                                                                                 editRecord={nodeReferenceModalEdit.EditRecord}
                                                                                 reference={reference}
                />}
                {hardwareReferenceModalEdit.State && <HardwareReferenceRecordModalCreate action="edit"
                                                                                         setState={(state) => setHardwareReferenceModalEdit(prevState => ({...prevState, State: state}))}
                                                                                         returnRecord={handlerEditRecord}
                                                                                         editRecord={hardwareReferenceModalEdit.EditRecord}
                                                                                         reference={reference}
                />}
                {switchModalEdit.State && <SwitchModalCreate action="edit"
                                                             setState={(state) => setSwitchModalEdit(prevState => ({...prevState, State: state}))}
                                                             returnRecord={handlerEditRecord}
                                                             editSwitch={switchModalEdit.EditRecord}
                />}

                <div className="buttons">
                    <button onClick={handlerOpenModalCreate}><FontAwesomeIcon icon={faPlus}/>
                        {reference === "node_types" && "Создать тип узла"}
                        {reference === "owners" && "Создать владельца"}
                        {reference === "hardware_types" && "Создать тип оборудования"}
                        {reference === "operation_modes" && "Создать режим работы"}
                        {reference === "switches" && "Создать модель коммутатора"}
                    </button>
                </div>
            </>}
            {isLoaded && <>{records.length > 0 ? (
                    <table>
                        <thead>
                        <tr className={"row-type-2"}>
                            <th>ID</th>
                            <th>Наименование</th>
                            <th>Дата создания</th>
                            <th></th>
                        </tr>
                        </thead>
                        <tbody>
                        {records.map((record, index) => (
                            <tr key={"record"+index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                <td>{record.ID}</td>
                                <td>{record.Name || record.TranslateValue}</td>
                                <td>{new Date(record.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>
                                <td>
                                    {user.Role.Value !== "user" && <FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => handlerOpenModalEdit(record)}/>}
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                )
                :
                <div className="empty">Таблица пуста</div>
            }</>}
        </section>
    )
}

export default ReferencePage
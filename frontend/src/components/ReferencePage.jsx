import React, {useContext, useEffect, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faPen, faPlus} from "@fortawesome/free-solid-svg-icons";
import FetchRequest from "../fetchRequest";
import NodeReferenceRecordModalCreate from "./NodeReferenceRecordModalCreate";
import ReferenceRecordModalCreate from "./ReferenceRecordModalCreate";
import SwitchModalCreate from "./SwitchModalCreate";
import AuthContext from "../context/AuthContext";

const ReferencePage = ({reference}) => {
    const [records, setRecords] = useState([])
    const [referenceModalCreate, setReferenceModalCreate] = useState(false)
    const [referenceModalEdit, setReferenceModalEdit]= useState({
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

    const handlerAddRecord = (record) => {
        setRecords(prevState => [...prevState, record])
    }

    const handlerEditRecord = (record) => {
        setRecords(prevState => prevState.map(_record => record.ID === _record.ID ? record : _record))
    }

    return (
        <section className="references">
            {user.role.key !== "user" && <>
                {referenceModalCreate && <ReferenceRecordModalCreate action="create" setState={setReferenceModalCreate} returnRecord={handlerAddRecord} reference={reference} withKey={reference === "hardware_types" || reference === "operation_modes"}/>}
                {referenceModalEdit.State && <ReferenceRecordModalCreate action="edit"
                                                                             setState={(state) => setReferenceModalEdit(prevState => ({...prevState, State: state}))}
                                                                             returnRecord={handlerEditRecord}
                                                                             editRecord={referenceModalEdit.EditRecord}
                                                                             reference={reference}
                                                                             withKey={reference === "hardware_types" || reference === "operation_modes"}
                />}

                <div className="buttons">
                    <button onClick={() => setReferenceModalCreate(true)}><FontAwesomeIcon icon={faPlus}/>
                        {reference === "node_types" && "Создать тип узла"}
                        {reference === "owners" && "Создать владельца"}
                        {reference === "hardware_types" && "Создать тип оборудования"}
                        {reference === "operation_modes" && "Создать режим работы"}
                        {reference === "roof_types" && "Создать тип крыши"}
                        {reference === "wiring_types" && "Создать тип разводки"}
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
                                <td>{record.Value}</td>
                                <td>{new Date(record.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>
                                <td>
                                    {user.role.key !== "user" && <FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => setReferenceModalEdit({State: true, EditRecord: record})}/>}
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
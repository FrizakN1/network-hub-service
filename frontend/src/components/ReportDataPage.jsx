import React, {useContext, useEffect, useState} from "react";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faPen} from "@fortawesome/free-solid-svg-icons";
import AuthContext from "../context/AuthContext";
import FetchRequest from "../fetchRequest";
import ReportDataModal from "./ReportDataModal";

const ReportDataPage = () => {
    const [isLoaded, setIsLoaded] = useState(false)
    const [records, setRecords] = useState([])
    const [modalEdit, setModalEdit] = useState({
        State: false,
        Record: null,
    })
    const { user } = useContext(AuthContext)

    useEffect(() => {
        FetchRequest("GET", "/report")
            .then(response => {
                if (response.success) {
                    setRecords(response.data)
                    setIsLoaded(true)
                }
            })
    }, []);

    const returnRecord = (record) => {
        setRecords(prevState => prevState.map(oldRecord => oldRecord.Key === record.Key ? record : oldRecord))
    }

    return (
        <section className="report">
            {modalEdit.State && <ReportDataModal
                setState={(state) => setModalEdit(prevState => ({...prevState, State: state}))}
                editRecord={modalEdit.Record}
                returnRecord={returnRecord}
            />}
            {isLoaded && <>{records.length > 0 ? (
                    <table>
                        <thead>
                        <tr className={"row-type-2"}>
                            <th>ID</th>
                            <th>Ключ</th>
                            <th>Значение</th>
                            <th>Описание</th>
                        </tr>
                        </thead>
                        <tbody>
                        {records.map((record, index) => (
                            <tr key={"record"+index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                <td>{record.ID}</td>
                                <td>{record.Key}</td>
                                <td>{record.Value}</td>
                                <td>{record.Description.String}</td>
                                <td>
                                    {user.role.key !== "user" && <FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => setModalEdit({State: true, Record: record})}/>}
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

export default ReportDataPage
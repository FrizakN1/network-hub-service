import React, {useContext, useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import SearchInput from "./SearchInput";
import FilesTable from "./FilesTable";
import FetchRequest from "../fetchRequest";
import NodesTable from "./NodesTable";
import HardwareTable from "./HardwareTable";
import EventsTable from "./EventsTable";
import {faFileExcel, faPen} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import HouseParamsModalEdit from "./HouseParamsModalEdit";
import AuthContext from "../context/AuthContext";

const HousePage = () => {
    const { id } = useParams()
    const [address, setAddress] = useState({})
    const [addressParams, setAddressParams] = useState({})
    const [isLoaded, setIsLoaded] = useState(false)
    const [activeTab, setActiveTab] = useState(1)
    const [modalEdit, setModalEdit] = useState({
        State: false,
        EditData: {},
    })
    const { user } = useContext(AuthContext)

    useEffect(() => {
        FetchRequest("GET", `/houses/${id}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    setAddress(response.data?.Address || {})
                    setAddressParams(response.data?.Params || {})
                    setIsLoaded(true)
                }
            })
    }, []);

    const handlerSetModalEdit = () => {
        setModalEdit({
            State: true,
            EditData: addressParams
        })
    }

    const handlerSetParams = (params) => {
        setAddressParams(params)
        setModalEdit({State: false, EditData: {}})
    }

    const handlerGetExcel = () => {
        FetchRequest("GET", `/houses/${id}/excel`)
            .then(response => {
                if (response.success) {
                    const binaryString = atob(response.data);
                    const byteArray = new Uint8Array(binaryString.length);

                    for (let i = 0; i < binaryString.length; i++) {
                        byteArray[i] = binaryString.charCodeAt(i);
                    }

                    // Создаем Blob и скачиваем файл
                    const blob = new Blob([byteArray], {
                        type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
                    });
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');

                    a.href = url;
                    a.download = 'Узлы.xlsx';
                    a.click();

                    URL.revokeObjectURL(url);
                }
            })
    }

    return (
        <section className="house">
            {modalEdit.State && <HouseParamsModalEdit editData={modalEdit.EditData} setState={(s) => setModalEdit(
                prevState => ({...prevState, State: s}))}
                      returnData={handlerSetParams}/>}
            {isLoaded &&
                <div>
                    <SearchInput defaultValue={`${address.street.type.short_name} ${address.street.name}, ${address.house.type.short_name} ${address.house.name}`}/>
                    <div className="contain">
                        <div className="info column" style={{alignItems: "center"}}>
                            {user.role.key !== "user" && <div className="buttons">
                                <button onClick={handlerGetExcel}>
                                    <FontAwesomeIcon icon={faFileExcel} /> Получить Excel
                                </button>

                                <button onClick={handlerSetModalEdit}>
                                    <FontAwesomeIcon icon={faPen}/> Редактировать
                                </button>
                            </div>}
                            <div className="row">
                                <div className="column">
                                    <div className="block">
                                        <span>Тип крыши</span>
                                        <p>{addressParams.RoofType.Value !== "" ? addressParams.RoofType.Value : "-"}</p>
                                    </div>
                                </div>
                                <div className="column">
                                    <div className="block">
                                        <span>Тип разводки</span>
                                        <p>{addressParams.WiringType.Value ? addressParams.WiringType.Value : "-"}</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="tabs-contain">
                        <div className="tabs">
                            <div className={activeTab === 1 ? "tab active" : "tab"} onClick={() => setActiveTab(1)}>Файлы</div>
                            <div className={activeTab === 2 ? "tab active" : "tab"} onClick={() => setActiveTab(2)}>Узлы</div>
                            <div className={activeTab === 3 ? "tab active" : "tab"} onClick={() => setActiveTab(3)}>Оборудование</div>
                            <div className={activeTab === 4 ? "tab active" : "tab"} onClick={() => setActiveTab(4)}>События</div>
                        </div>
                    </div>
                    {activeTab === 1 && <FilesTable type="houses"/>}
                    {activeTab === 2 &&
                        <div>
                            <NodesTable houseID={Number(id)} canCreate={true} defaultAddress={address}/>
                        </div>
                    }
                    {activeTab === 3 &&
                        <div>
                            <HardwareTable type={"houses"} id={Number(id)} canCreate={true}/>
                        </div>
                    }
                    {activeTab === 4 &&
                        <div>
                            <EventsTable from={"houses"}/>
                        </div>
                    }
                </div>
            }
        </section>
    )
}

export default HousePage
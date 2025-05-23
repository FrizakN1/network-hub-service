import React, {useContext, useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import SearchInput from "./SearchInput";
import FilesTable from "./FilesTable";
import FetchRequest from "../fetchRequest";
import NodesTable from "./NodesTable";
import HardwareTable from "./HardwareTable";
import EventsTable from "./EventsTable";
import {faPen} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import HouseParamsModalEdit from "./HouseParamsModalEdit";
import AuthContext from "../context/AuthContext";

const HousePage = () => {
    const { id } = useParams()
    const [address, setAddress] = useState({})
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
                    setAddress(response.data)
                    setIsLoaded(true)
                }
            })
    }, []);

    const handlerSetModalEdit = () => {
        setModalEdit({
            State: true,
            EditData: {
                RoofType: address.RoofType,
                WiringType: address.WiringType
            }
        })
    }

    const handlerSetParams = (params) => {
        setAddress(prevState => ({...prevState, RoofType: params.RoofType, WiringType: params.WiringType}))
        setModalEdit({State: false, EditData: {}})
    }

    return (
        <section className="house">
            {modalEdit.State && <HouseParamsModalEdit editData={modalEdit.EditData} setState={(s) => setModalEdit(
                prevState => ({...prevState, State: s}))}
                      returnData={handlerSetParams}/>}
            {isLoaded &&
                <div>
                    <SearchInput defaultValue={`${address.Street.Type.ShortName} ${address.Street.Name}, ${address.House.Type.ShortName} ${address.House.Name}`}/>
                    <div className="contain">
                        <div className="info column" style={{alignItems: "center"}}>
                            {user.role.key !== "user" && <div className="buttons">
                                <button onClick={handlerSetModalEdit}>
                                    <FontAwesomeIcon icon={faPen}/> Редактировать
                                </button>
                            </div>}
                            <div className="row">
                                <div className="column">
                                    <div className="block">
                                        <span>Тип крыши</span>
                                        <p>{address.RoofType.Value !== "" ? address.RoofType.Value : "-"}</p>
                                    </div>
                                </div>
                                <div className="column">
                                    <div className="block">
                                        <span>Тип разводки</span>
                                        <p>{address.WiringType.Value ? address.WiringType.Value : "-"}</p>
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
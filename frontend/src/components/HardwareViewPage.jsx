import React, {useEffect, useState} from "react";
import {Link, useParams} from "react-router-dom";
import FetchRequest from "../fetchRequest";
import FilesTable from "./FilesTable";
import EventsTable from "./EventsTable";

const HardwareViewPage = () => {
    const { id } = useParams()
    const [hardware, setHardware] = useState(null)
    const [activeTab, setActiveTab] = useState(1)
    const [isLoaded, setIsLoaded] = useState(false)

    useEffect(() => {
        FetchRequest("GET", `/hardware/${id}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    setHardware(response.data)
                    setIsLoaded(true)
                }
            })
    }, [id]);

    return (
        <section className="hardware-view">
            <h2>Оборудование</h2>
            <div className="contain">
                {isLoaded &&
                    <div className="info">
                        <div className="column">
                            <div className="block">
                                <span>Узел</span>
                                <Link to={`/nodes/view/${hardware.Node.ID}`}><p>{hardware.Node.Name}</p></Link>
                            </div>
                            <div className="block">
                                <span>Адрес</span>
                                <p>{`${hardware.Node.Address.street.type.short_name} ${hardware.Node.Address.street.name}, ${hardware.Node.Address.house.type.short_name} ${hardware.Node.Address.house.name}`}</p>
                            </div>
                            <div className="block">
                                <span>Тип оборудования</span>
                                <p>{hardware.Type.Value}</p>
                            </div>
                            <div className="block textarea">
                                <span>Описание</span>
                                <p>{hardware.Description.String}</p>
                            </div>
                        </div>
                        <div className="column">
                            <div className="block">
                                <span>Модель</span>
                                <p>{hardware.Type.Key === "switch" && hardware.Switch.ID !== 0 ? hardware.Switch.Name : "-"}</p>
                            </div>
                            <div className="block">
                                <span>IP адрес</span>
                                <p>{hardware.Type.Key === "switch" && hardware.IpAddress.Valid ? hardware.IpAddress.String : "-"}</p>
                            </div>
                            <div className="block">
                                <span>Управляющий VLAN</span>
                                <p>{hardware.Type.Key === "switch" && hardware.MgmtVlan.Valid ? hardware.MgmtVlan.String : "-"}</p>
                            </div>
                            <div className="block">
                                <span>Дата создания</span>
                                <p>{new Date(hardware.CreatedAt * 1000).toLocaleString().slice(0, 17)}</p>
                            </div>
                            <div className="block">
                                <span>Дата изменения</span>
                                <p>{hardware.UpdatedAt.Valid ? new Date(hardware.UpdatedAt.Int64 * 1000).toLocaleString().slice(0, 17) : "-"}</p>
                            </div>
                        </div>
                    </div>
                }
            </div>
            <div className="tabs-contain">
                <div className="tabs">
                    <div className={activeTab === 1 ? "tab active" : "tab"} onClick={() => setActiveTab(1)}>Файлы</div>
                    <div className={activeTab === 2 ? "tab active" : "tab"} onClick={() => setActiveTab(2)}>События</div>
                </div>
            </div>
            {activeTab === 1 && <FilesTable type="hardware"/>}
            {activeTab === 2 && <EventsTable type="hardware"/>}
        </section>
    )
}

export default HardwareViewPage
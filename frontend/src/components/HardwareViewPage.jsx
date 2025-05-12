import React, {useEffect, useState} from "react";
import {Link, useParams} from "react-router-dom";
import FetchRequest from "../fetchRequest";
import FilesTable from "./FilesTable";

const HardwareViewPage = () => {
    const { id } = useParams()
    const [hardware, setHardware] = useState(null)
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
                                <p>{`${hardware.Node.Address.Street.Type.ShortName} ${hardware.Node.Address.Street.Name}, ${hardware.Node.Address.House.Type.ShortName} ${hardware.Node.Address.House.Name}`}</p>
                            </div>
                            <div className="block">
                                <span>Тип оборудования</span>
                                <p>{hardware.Type.TranslateValue}</p>
                            </div>
                            <div className="block textarea">
                                <span>Описание</span>
                                <p>{hardware.Description.String}</p>
                            </div>
                        </div>
                        <div className="column">
                            <div className="block">
                                <span>Модель</span>
                                <p>{hardware.Type.Value === "switch" && hardware.Switch.ID !== 0 ? hardware.Switch.Name : "-"}</p>
                            </div>
                            <div className="block">
                                <span>IP адрес</span>
                                <p>{hardware.Type.Value === "switch" && hardware.IpAddress.Valid ? hardware.IpAddress.String : "-"}</p>
                            </div>
                            <div className="block">
                                <span>Управляющий VLAN</span>
                                <p>{hardware.Type.Value === "switch" && hardware.MgmtVlan.Valid ? hardware.MgmtVlan.String : "-"}</p>
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
            <FilesTable type="hardware"/>
        </section>
    )
}

export default HardwareViewPage
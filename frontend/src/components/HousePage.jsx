import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import API_DOMAIN from "../config";
import SearchInput from "./SearchInput";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faDownload, faEye, faFileWord, faPen, faPlus, faTrash} from "@fortawesome/free-solid-svg-icons";
import FilesTable from "./FilesTable";
import UploadFile from "./UploadFile";
import fetchRequest from "../fetchRequest";
import FetchRequest from "../fetchRequest";
import NodeModalCreate from "./NodeModalCreate";
import NodesTable from "./NodesTable";

const HousePage = () => {
    const { id } = useParams()
    const [address, setAddress] = useState({})
    const [isLoaded, setIsLoaded] = useState(false)
    const [activeTab, setActiveTab] = useState(1)

    useEffect(() => {
        FetchRequest("GET", `/get_house/${id}`, null)
            .then(response => {
                if (response.success && response.data != null) {
                    setAddress(response.data)
                    setIsLoaded(true)
                }
            })

        // fetch(`${API_DOMAIN}/get_house/${houseID}`, {method: "GET"})
        //     .then(response => response.json())
        //     .then(data => {
        //         if (data != null) {
        //             setAddress(data)
        //             setIsLoaded(true)
        //         }
        //     })
        //     .catch(error => console.error(error))
    }, []);

    return (
        <section className="house">
            {isLoaded &&
                <div>
                    <SearchInput defaultValue={`${address.Street.Type.ShortName} ${address.Street.Name}, ${address.House.Type.ShortName} ${address.House.Name}`}/>
                    <div className="tabs-contain">
                        <div className="tabs">
                            <div className={activeTab === 1 ? "tab active" : "tab"} onClick={() => setActiveTab(1)}>Файлы</div>
                            <div className={activeTab === 2 ? "tab active" : "tab"} onClick={() => setActiveTab(2)}>Узлы</div>
                            <div className={activeTab === 3 ? "tab active" : "tab"} onClick={() => setActiveTab(3)}>Оборудование</div>
                        </div>
                    </div>
                    {activeTab === 1 && <FilesTable type="house"/>}
                    {activeTab === 2 &&
                        <div>
                            <NodesTable id={id} canCreate={true}/>
                        </div>
                    }
                </div>
            }
        </section>
    )
}

export default HousePage
import React, {useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import SearchInput from "./SearchInput";
import FilesTable from "./FilesTable";
import FetchRequest from "../fetchRequest";
import NodesTable from "./NodesTable";
import HardwareTable from "./HardwareTable";

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
                            <NodesTable id={Number(id)} canCreate={true}/>
                        </div>
                    }
                    {activeTab === 3 &&
                        <div>
                            <HardwareTable type={"house"} id={Number(id)} canCreate={true}/>
                        </div>
                    }
                </div>
            }
        </section>
    )
}

export default HousePage
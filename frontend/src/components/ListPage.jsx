import React, {useEffect, useState} from "react";
import API_DOMAIN from "../config";
import SearchInput from "./SearchInput";
import AddressesTable from "./AddressesTable";
import fetchRequest from "../fetchRequest";
import FetchRequest from "../fetchRequest";

const ListPage = () => {
    const [addresses, setAddresses] = useState([])
    const [count, setCount] = useState(0)
    const [offset, setOffset] = useState(1)
    const [isLoaded, setIsLoaded] = useState(false)

    useEffect(() => {
        // let options = {
        //     method: "POST",
        //     body: JSON.stringify((offset-1)*20)
        // }

        FetchRequest("POST", "/get_list", (offset-1)*20)
            .then(response => {
                if (response.success) {
                    if (response.data != null) {
                        setAddresses(response.data?.Addresses || [])
                        setCount(response.data?.Count || 0)
                    }
                    setIsLoaded(true)
                }
            })

        // fetch(`${API_DOMAIN}/get_list`, options)
        //     .then(response => response.json())
        //     .then(data => {
        //         if (data != null) {
        //             setAddresses(data?.Addresses || [])
        //             setCount(data?.Count || 0)
        //         }
        //         setIsLoaded(true)
        //     })
        //     .catch(error => console.error(error))
    }, [offset]);

    return (
        <section className="result">
            {isLoaded && <AddressesTable addresses={addresses} count={count} setOffset={setOffset} h={"Адреса, содержащие файлы: "}/>}
        </section>
    )
}

export default ListPage
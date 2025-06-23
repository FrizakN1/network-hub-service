import React, {useEffect, useState} from "react";
import AddressesTable from "./AddressesTable";
import FetchRequest from "../fetchRequest";

const ListPage = () => {
    const [addresses, setAddresses] = useState([])
    const [count, setCount] = useState(0)
    const [offset, setOffset] = useState(1)
    const [isLoaded, setIsLoaded] = useState(false)
    const [addressAmounts, setAddressAmounts] = useState({})

    useEffect(() => {
        let params = new URLSearchParams({
            offset: String((offset-1)*20)
        })

        FetchRequest("GET", `/houses?${params.toString()}`, null)
            .then(response => {
                if (response.success) {
                    if (response.data != null) {
                        setAddresses(response.data?.Addresses || [])
                        setCount(response.data?.Count || 0)
                        setAddressAmounts(response.data?.AddressAmounts || {})
                    }
                    setIsLoaded(true)
                }
            })
    }, [offset]);

    return (
        <section className="result">
            {isLoaded && <AddressesTable addresses={addresses} count={count} addressAmounts={addressAmounts} setOffset={setOffset} h={"Адреса, содержащие файлы, узлы или оборудование: "}/>}
        </section>
    )
}

export default ListPage
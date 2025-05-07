import React from "react";
import NodesTable from "./NodesTable";

const NodesPage = () => {
    return (
        <section className="nodes">
            <NodesTable canCreate={true}/>
        </section>
    )
}

export default NodesPage
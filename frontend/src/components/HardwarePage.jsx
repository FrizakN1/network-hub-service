import React from "react";
import HardwareTable from "./HardwareTable";

const HardwarePage = () => {
    return (
        <section className="hardware">
            <HardwareTable canCreate={true} />
        </section>
    )
}

export default HardwarePage
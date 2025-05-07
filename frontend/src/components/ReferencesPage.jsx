import React from "react";
import {Link} from "react-router-dom";

const ReferencesPage = () => {
    return (
        <section className="references">
            <div className="contain">
                <div className="list">
                    <Link to="/references/node_types">
                        Типы узлов
                    </Link>
                    <Link to="/references/owners">
                        Владельцы узлов
                    </Link>
                    <Link to="/references/hardware_types">
                        Типы оборудования
                    </Link>
                    <Link to="/references/operation_modes">
                        Режимы работы коммутаторов
                    </Link>
                    <Link to="/references/switches">
                        Модели коммутаторов
                    </Link>
                </div>
            </div>
        </section>
    )
}

export default ReferencesPage
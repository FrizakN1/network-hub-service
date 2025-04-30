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
                </div>
            </div>
        </section>
    )
}

export default ReferencesPage
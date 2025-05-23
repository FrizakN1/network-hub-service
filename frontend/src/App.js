import './App.css';
import PrivateRoute from "./components/PrivateRoute";
import {Route, BrowserRouter as Router, Routes} from "react-router-dom";
import MainPage from "./components/MainPage";
import HousePage from "./components/HousePage";
import ResultPage from "./components/ResultPage";
import ListPage from "./components/ListPage";
import Authorization from "./components/Authorization";
import UsersPage from "./components/UsersPage";
import ReferencesPage from "./components/ReferencesPage";
import ReferencePage from "./components/ReferencePage";
import NodeViewPage from "./components/NodeViewPage";
import NodesPage from "./components/NodesPage";
import HardwarePage from "./components/HardwarePage";
import HardwareViewPage from "./components/HardwareViewPage";
import SwitchesPage from "./components/SwitchesPage";
import EventsPage from "./components/EventsPage";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/login" element={<Authorization />} />
                <Route path="/" element={<PrivateRoute />}>
                    <Route path="" element={<MainPage />} />
                    <Route path="/house/:id" element={<HousePage />} />
                    <Route path="/result" element={<ResultPage />} />
                    <Route path="/list" element={<ListPage />} />
                </Route>
                <Route path="/" element={<PrivateRoute requiredAdmin={true} />}>
                    <Route path="/users" element={<UsersPage />} />
                </Route>
                <Route path="/references/" element={<PrivateRoute />}>
                    <Route path="" element={<ReferencesPage />} />
                    <Route path="owners/" element={<ReferencePage reference={"owners"}/>} />
                    <Route path="node_types/" element={<ReferencePage reference={"node_types"} />} />
                    <Route path="hardware_types/" element={<ReferencePage reference={"hardware_types"} />} />
                    <Route path="operation_modes/" element={<ReferencePage reference={"operation_modes"} />} />
                    <Route path="roof_types/" element={<ReferencePage reference={"roof_types"} />} />
                    <Route path="wiring_types/" element={<ReferencePage reference={"wiring_types"} />} />
                </Route>
                <Route path="/switches/" element={<PrivateRoute />}>
                    <Route path="" element={<SwitchesPage />} />
                </Route>
                <Route path="/nodes/" element={<PrivateRoute />}>
                    <Route path="" element={<NodesPage />} />
                    <Route path="view/:id" element={<NodeViewPage />} />
                </Route>
                <Route path="/hardware/" element={<PrivateRoute />}>
                    <Route path="" element={<HardwarePage />} />
                    <Route path="view/:id" element={<HardwareViewPage />} />
                </Route>
                <Route path="/events/" element={<PrivateRoute />}>
                    <Route path="" element={<EventsPage />} />
                </Route>
            </Routes>
        </Router>
    )
}

export default App;

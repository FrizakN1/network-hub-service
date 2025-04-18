import './App.css';
import PrivateRoute from "./components/PrivateRoute";
import {Route, BrowserRouter as Router, Routes} from "react-router-dom";
import MainPage from "./components/MainPage";
import HousePage from "./components/HousePage";
import ResultPage from "./components/ResultPage";
import ListPage from "./components/ListPage";
import Authorization from "./components/Authorization";
import UsersPage from "./components/UsersPage";

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/login" element={<Authorization />} />
                <Route path="/" element={<PrivateRoute />}>
                    <Route path="" element={<MainPage />} />
                    <Route path="/house/:houseID" element={<HousePage />} />
                    <Route path="/result" element={<ResultPage />} />
                    <Route path="/list" element={<ListPage />} />
                </Route>
                <Route path="/" element={<PrivateRoute requiredAdmin={true} />}>
                    <Route path="/users" element={<UsersPage />} />
                </Route>
            </Routes>
        </Router>
    )
}

export default App;

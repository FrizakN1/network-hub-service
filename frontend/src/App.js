import './App.css';
import PrivateRoute from "./components/PrivateRoute";
import {Route, BrowserRouter as Router, Routes} from "react-router-dom";
import MainPage from "./components/MainPage";
import HousePage from "./components/HousePage";
import ResultPage from "./components/ResultPage";
import ListPage from "./components/ListPage";

function App() {
    return (
        <Router>
            <Routes>
                {/*<Route path="/login" element={<Authorization />} />*/}
                <Route path="/" element={<PrivateRoute />}>
                    <Route path="" element={<MainPage />} />
                    <Route path="/house/:houseID" element={<HousePage />} />
                    <Route path="/result" element={<ResultPage />} />
                    <Route path="/list" element={<ListPage />} />
                </Route>
            </Routes>
        </Router>
    )
}

export default App;

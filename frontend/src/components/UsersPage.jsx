import React, {useEffect, useState} from "react";
import FetchRequest from "../fetchRequest";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faBan, faDownload, faEye, faFolderMinus, faFolderPlus, faPen, faTrash} from "@fortawesome/free-solid-svg-icons";

const UsersPage = () => {
    const [users, setUsers] = useState([])

    useEffect(() => {
        FetchRequest("GET", "/get_users", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setUsers(response.data)
                    console.log(response.data)
                }
            })
    }, []);

    return (
        <section className="users">
            {users.length > 0 ? (
                    <table>
                        <thead>
                        <tr className={"row-type-2"}>
                            <th>ID</th>
                            <th>Логин</th>
                            <th>ФИО</th>
                            <th>Роль</th>
                            <th>Статус</th>
                            <th>Дата создания</th>
                            <th></th>
                        </tr>
                        </thead>
                        <tbody>
                        {users.map((user, index) => (
                            <tr key={index} className={index % 2 === 0 ? 'row-type-1' : 'row-type-2'}>
                                <td>{user.ID}</td>
                                <td>{user.Login}</td>
                                <td>{user.Name}</td>
                                <td>{user.Role.TranslateValue}</td>
                                <td>{user.Baned ? <span className={"bg-red"}>Заблокирован</span> : <span className={"bg-green"}>Активен</span>}</td>
                                <td>{new Date(user.CreatedAt * 1000).toLocaleString().slice(0, 17)}</td>
                                <td>
                                    <FontAwesomeIcon icon={faEye} title="Просмотр" />
                                    <FontAwesomeIcon icon={faPen} title="Редактировать" />
                                    <FontAwesomeIcon icon={faBan} title="Заблокировать" />
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                )
                :
                <div className="empty">Таблица пуста</div>
            }
        </section>
    )
}

export default UsersPage
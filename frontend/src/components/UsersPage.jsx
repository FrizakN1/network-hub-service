import React, {useEffect, useState} from "react";
import FetchRequest from "../fetchRequest";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {
    faBan,
    faCircleCheck,
    faPen, faPlus,
} from "@fortawesome/free-solid-svg-icons";
import UserModalCreate from "./UserModalCreate";

const UsersPage = () => {
    const [users, setUsers] = useState([])
    const [isLoaded, setIsLoaded] = useState(false)
    const [modalCreate, setModalCreate] = useState(false)
    const [modalEdit, setModalEdit] = useState({
        State: false,
        EditUser: {}
    })

    useEffect(() => {
        FetchRequest("GET", "/users", null)
            .then(response => {
                if (response.success && response.data != null) {
                    setUsers(response.data)
                }

                setIsLoaded(true)
            })
    }, []);

    const changeUserStatus = (userID) => {
        FetchRequest("PATCH", "/users/status", {ID: Number(userID)})
            .then(response => {
                if (response.success && response.data != null) {
                    setUsers(prevState => prevState.map(user =>
                        user.ID === response.data.ID ? response.data : user
                    ))
                }
            })
    }

    const handlerAddUser = (user) => {
        setUsers(prevState => [...prevState, user])
    }

    const handlerEditUser = (user) => {
        setUsers(prevState => prevState.map(_user => user.ID === _user.ID ? user : _user))
    }

    return (
        <section className="users">
            {modalCreate && <UserModalCreate action="create" setState={setModalCreate} returnUser={handlerAddUser}/>}
            {modalEdit.State && <UserModalCreate action="edit"
                                                 setState={(state) => setModalEdit(prevState => ({...prevState, State: state}))}
                                                 returnUser={handlerEditUser}
                                                 editUser={modalEdit.EditUser}
            />}
            <div className="buttons">
                <button onClick={() => setModalCreate(true)}><FontAwesomeIcon icon={faPlus}/> Создать пользователя</button>
            </div>
            {isLoaded && <>{users.length > 0 ? (
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
                                    {/*<FontAwesomeIcon icon={faEye} title="Просмотр" />*/}
                                    <FontAwesomeIcon icon={faPen} title="Редактировать" onClick={() => setModalEdit({State: true, EditUser: user})}/>
                                    {user.Baned ?
                                        <FontAwesomeIcon icon={faCircleCheck} className="eye" title="Разаблокировать" onClick={() => changeUserStatus(user.ID)}/>
                                        :
                                        <FontAwesomeIcon icon={faBan} className="delete" title="Заблокировать" onClick={() => changeUserStatus(user.ID)}/>
                                    }
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                )
                :
                <div className="empty">Таблица пуста</div>
            }</>}
        </section>
    )
}

export default UsersPage
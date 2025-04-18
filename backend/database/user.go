package database

import (
	"backend/settings"
	"backend/utils"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Role struct {
	ID             int
	Value          string
	TranslateValue string
}
type User struct {
	ID        int
	Role      Role
	Login     string
	Name      string
	Password  string
	Baned     bool
	CreatedAt int64
	UpdatedAt sql.NullInt64
}

type Session struct {
	Hash      string
	User      User
	CreatedAt int64
}

var sessionMap map[string]Session
var roleMap map[int]Role

func prepareUsers() []string {
	var e error
	errorsList := make([]string, 0)
	sessionMap = make(map[string]Session)
	roleMap = make(map[int]Role)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["GET_ROLES"], e = Link.Prepare(`
		SELECT * FROM "Role"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_USERS"], e = Link.Prepare(`
		SELECT id, role_id, login, name, baned, created_at, updated_at FROM "User"
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_USER"], e = Link.Prepare(`
		SELECT * FROM "User" WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_USER"], e = Link.Prepare(`
		INSERT INTO "User"(role_id, login, name, password, baned, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["EDIT_USER"], e = Link.Prepare(`
		UPDATE "User" SET role_id = $2, login = $3, name = $4, password = $5, updated_at = $6
		WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_AUTHORIZED_USER"], e = Link.Prepare(`
		SELECT id, role_id, name, baned, created_at, updated_at FROM "User" WHERE login = $1 AND password = $2
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SESSIONS"], e = Link.Prepare(`
		SELECT s.*, u.role_id, u.login, u.name, u.baned, u.created_at, u.updated_at 
		FROM "Session" AS s
		JOIN "User" AS u ON u.id = s.user_id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CREATE_SESSION"], e = Link.Prepare(`
		INSERT INTO "Session" (hash, user_id, created_at) VALUES ($1, $2, $3)
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["DELETE_SESSION"], e = Link.Prepare(`
		DELETE FROM "Session" WHERE hash = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_USER_SESSIONS"], e = Link.Prepare(`
		SELECT hash FROM "Session" WHERE user_id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["CHANGE_SUPER_ADMIN_PASSWORD"], e = Link.Prepare(`
		UPDATE "User" SET password = $2 WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SUPER_ADMIN"], e = Link.Prepare(`
		SELECT id, password FROM "User" WHERE login = 'SuperAdmin'
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func GetUsers() ([]User, error) {
	stmt, ok := query["GET_USERS"]
	if !ok {
		err := errors.New("запрос GET_USERS не подготовлен")
		utils.Logger.Println(err)
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		utils.Logger.Println(err)
		return nil, err
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		if err = rows.Scan(
			&user.ID,
			&user.Role.ID,
			&user.Login,
			&user.Name,
			&user.Baned,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		user.Role = roleMap[user.Role.ID]

		users = append(users, user)
	}

	return users, nil
}

func (user *User) EditUser() error {
	stmt, ok := query["EDIT_USER"]
	if !ok {
		err := errors.New("запрос EDIT_USER не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	if err := stmt.QueryRow(
		user.Role.ID,
		user.Login,
		user.Name,
		user.Password,
		user.Baned,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID); err != nil {
		utils.Logger.Println(err)
		return err
	}

	user.Password = ""

	return nil
}

func (user *User) CreateUser() error {
	stmt, ok := query["CREATE_USER"]
	if !ok {
		err := errors.New("запрос CREATE_USER не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	if err := stmt.QueryRow(
		user.Role.ID,
		user.Login,
		user.Name,
		user.Password,
		false,
		user.CreatedAt,
		nil,
	).Scan(&user.ID); err != nil {
		utils.Logger.Println(err)
		return err
	}

	user.Role = roleMap[user.Role.ID]

	user.Password = ""

	return nil
}

func (user *User) GetAuthorize() error {
	stmt, ok := query["GET_AUTHORIZED_USER"]
	if !ok {
		err := errors.New("запрос GET_AUTHORIZED_USER не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	if err := stmt.QueryRow(user.Login, user.Password).Scan(
		&user.ID,
		&user.Role.ID,
		&user.Name,
		&user.Baned,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user.ID = 0
			return nil
		}

		utils.Logger.Println(err)
		return err
	}

	user.Role = roleMap[user.Role.ID]

	return nil
}

func DeleteSession(s *Session) error {
	stmt, ok := query["DELETE_SESSION"]
	if !ok {
		return errors.New("запрос DELETE_SESSION не подготовлен")
	}

	_, e := stmt.Exec(s.Hash)
	if e != nil {
		return e
	}

	delete(sessionMap, s.Hash)

	return nil
}

func GetSession(hash string) *Session {
	session, ok := sessionMap[hash]
	if ok {
		return &session
	}

	return nil
}

func CreateSession(user User) (string, error) {
	stmt, ok := query["CREATE_SESSION"]
	if !ok {
		err := errors.New("запрос CREATE_SESSION не подготовлен")
		utils.Logger.Println(err)
		return "", err
	}

	hash, err := utils.GenerateHash(fmt.Sprintf("%s-%d", user.Login, time.Now().Unix()))
	if err != nil {
		utils.Logger.Println(err)
		return "", err
	}

	if _, err = stmt.Exec(hash, user.ID, time.Now().Unix()); err != nil {
		utils.Logger.Println(err)
		return "", err
	}

	if sessionMap != nil {
		sessionMap[hash] = Session{
			Hash:      hash,
			User:      user,
			CreatedAt: time.Now().Unix(),
		}
	}

	return hash, nil
}

func LoadSession(m map[string]Session) {
	stmt, ok := query["GET_SESSIONS"]
	if !ok {
		return
	}

	rows, e := stmt.Query()
	if e != nil {
		fmt.Println(e)
		utils.Logger.Println(e)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var session Session

		e = rows.Scan(
			&session.Hash,
			&session.User.ID,
			&session.CreatedAt,
			&session.User.Role.ID,
			&session.User.Login,
			&session.User.Name,
			&session.User.Baned,
			&session.User.CreatedAt,
			&session.User.UpdatedAt,
		)
		if e != nil {
			fmt.Println(e)
			utils.Logger.Println(e)
			return
		}

		session.User.Role = roleMap[session.User.Role.ID]

		m[session.Hash] = session
	}
}

func LoadRole(m map[int]Role) {
	stmt, ok := query["GET_ROLES"]
	if !ok {
		return
	}

	rows, e := stmt.Query()
	if e != nil {
		fmt.Println(e)
		utils.Logger.Println(e)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var role Role

		e = rows.Scan(
			&role.ID,
			&role.Value,
			&role.TranslateValue,
		)
		if e != nil {
			fmt.Println(e)
			utils.Logger.Println(e)
			return
		}

		m[role.ID] = role
	}
}

func CheckAdmin(config *settings.Setting) error {
	stmt, ok := query["GET_SUPER_ADMIN"]
	if !ok {
		err := errors.New("запрос GET_SUPER_ADMIN не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var admin User

	e := stmt.QueryRow().Scan(&admin.ID, &admin.Password)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			if e = CreateAdmin(config); e != nil {
				utils.Logger.Println(e)
				return e
			}
		} else {
			utils.Logger.Println(e)
			return e
		}
	}

	encryptPass, e := utils.Encrypt(config.SuperAdminPassword)
	if e != nil {
		utils.Logger.Println(e)
		return e
	}

	if encryptPass != admin.Password {
		if e = admin.ChangeSuperAdminPassword(encryptPass); e != nil {
			utils.Logger.Println(e)
			return e
		}

		if e = DeleteUserSessions(admin.ID); e != nil {
			utils.Logger.Println(e)
			return e
		}
	}

	return nil
}

func (user *User) ChangeSuperAdminPassword(newPassword string) error {
	stmt, ok := query["CHANGE_SUPER_ADMIN_PASSWORD"]
	if !ok {
		return errors.New("запрос CHANGE_SUPER_ADMIN_PASSWORD не подготовлен")
	}

	_, e := stmt.Exec(user.ID, newPassword)
	if e != nil {
		return e
	}

	return nil
}

func CreateAdmin(config *settings.Setting) error {
	var admin User

	encryptPass, e := utils.Encrypt(config.SuperAdminPassword)
	if e != nil {
		utils.Logger.Println(e)
		return e
	}

	for _, role := range roleMap {
		if role.Value == "admin" {
			admin = User{
				Login:     "SuperAdmin",
				Name:      "SuperAdmin",
				Role:      role,
				Password:  encryptPass,
				CreatedAt: time.Now().Unix(),
			}

			break
		}
	}

	if e = admin.CreateUser(); e != nil {
		utils.Logger.Println(e)
		return e
	}
	return nil
}

func DeleteUserSessions(id int) error {
	stmt, ok := query["GET_USER_SESSIONS"]
	if !ok {
		return errors.New("запрос GET_USER_SESSIONS не подготовлен")
	}

	rows, e := stmt.Query(id)
	if e != nil {
		return e
	}

	defer rows.Close()

	for rows.Next() {
		var session Session
		e = rows.Scan(&session.Hash)
		if e != nil {
			return e
		}

		if e = DeleteSession(&session); e != nil {
			utils.Logger.Println(e)
			return e
		}
	}

	return nil
}

package database

import (
	"backend/settings"
	"backend/utils"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID        int
	Role      Reference
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

type UserService interface {
	ChangeStatus(user *User) error
	GetUser(user *User) error
	GetUsers() ([]User, error)
	EditUser(user *User) error
	CreateUser(user *User) error
	GetAuthorize(user *User) error
	DeleteSession(s *Session) error
	CreateSession(user User) (string, error)
	ChangeUserPassword(user *User) error
	CreateAdmin(config *settings.Setting) error
	DeleteUserSessions(userID int) error
	ValidateUser(user User, action string) bool
	CheckAdmin(config *settings.Setting) error
}

type DefaultUserService struct {
	Encrypt utils.Encrypt
}

var sessionMap map[string]Session

func prepareUsers() []string {
	var e error
	errorsList := make([]string, 0)
	sessionMap = make(map[string]Session)

	if query == nil {
		query = make(map[string]*sql.Stmt)
	}

	query["GET_USERS"], e = Link.Prepare(`
		SELECT u.id, u.role_id, u.login, u.name, u.baned, u.created_at, u.updated_at, r.value, r.translate_value
		FROM "User" AS u
		JOIN "Role" AS r ON r.id = u.role_id
		ORDER BY u.id
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_USER"], e = Link.Prepare(`
		SELECT u.role_id, u.login, u.name, u.baned, u.created_at, u.updated_at, r.value, r.translate_value
		FROM "User" AS u
		JOIN "Role" AS r ON r.id = u.role_id
		WHERE u.id = $1
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
		UPDATE "User" SET role_id = $2, login = $3, name = $4, updated_at = $5
		WHERE id = $1
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_AUTHORIZED_USER"], e = Link.Prepare(`
		SELECT u.id, u.role_id, u.name, u.baned, u.created_at, u.updated_at, r.value, r.translate_value
		FROM "User" AS u
		JOIN "Role" AS r ON r.id = u.role_id
		WHERE login = $1 AND password = $2
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	query["GET_SESSIONS"], e = Link.Prepare(`
		SELECT s.*, u.role_id, u.login, u.name, u.baned, u.created_at, u.updated_at, r.value, r.translate_value
		FROM "Session" AS s
		JOIN "User" AS u ON u.id = s.user_id
		JOIN "Role" AS r ON r.id = u.role_id
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

	query["CHANGE_USER_PASSWORD"], e = Link.Prepare(`
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

	query["CHANGE_USER_STATUS"], e = Link.Prepare(`
		UPDATE "User" SET baned = NOT baned WHERE id = $1
		RETURNING baned
    `)
	if e != nil {
		errorsList = append(errorsList, e.Error())
	}

	return errorsList
}

func (us *DefaultUserService) ChangeStatus(user *User) error {
	stmt, ok := query["CHANGE_USER_STATUS"]
	if !ok {
		err := errors.New("запрос CHANGE_USER_STATUS не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	if err := stmt.QueryRow(user.ID).Scan(&user.Baned); err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (us *DefaultUserService) GetUser(user *User) error {
	stmt, ok := query["GET_USER"]
	if !ok {
		err := errors.New("запрос GET_USER не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	if err := stmt.QueryRow(user.ID).Scan(
		&user.Role.ID,
		&user.Login,
		&user.Name,
		&user.Baned,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.Value,
		&user.Role.TranslateValue,
	); err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (us *DefaultUserService) GetUsers() ([]User, error) {
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
			&user.Role.Value,
			&user.Role.TranslateValue,
		); err != nil {
			utils.Logger.Println(err)
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (us *DefaultUserService) EditUser(user *User) error {
	stmt, ok := query["EDIT_USER"]
	if !ok {
		err := errors.New("запрос EDIT_USER не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	_, err := stmt.Exec(
		user.ID,
		user.Role.ID,
		user.Login,
		user.Name,
		user.UpdatedAt,
	)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (us *DefaultUserService) CreateUser(user *User) error {
	stmt, ok := query["CREATE_USER"]
	if !ok {
		err := errors.New("запрос CREATE_USER не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var err error

	user.Password, err = us.Encrypt.Encrypt(user.Password)
	if err != nil {
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

	return nil
}

func (us *DefaultUserService) GetAuthorize(user *User) error {
	stmt, ok := query["GET_AUTHORIZED_USER"]
	if !ok {
		err := errors.New("запрос GET_AUTHORIZED_USER не подготовлен")
		utils.Logger.Println(err)
		return err
	}

	var err error

	user.Password, err = us.Encrypt.Encrypt(user.Password)
	if err != nil {
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
		&user.Role.Value,
		&user.Role.TranslateValue,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user.ID = 0
			return nil
		}

		utils.Logger.Println(err)
		return err
	}

	return nil
}

func (us *DefaultUserService) DeleteSession(s *Session) error {
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

func (us *DefaultUserService) CreateSession(user User) (string, error) {
	stmt, ok := query["CREATE_SESSION"]
	if !ok {
		err := errors.New("запрос CREATE_SESSION не подготовлен")
		utils.Logger.Println(err)
		return "", err
	}

	hash, err := us.Encrypt.GenerateHash(fmt.Sprintf("%s-%d", user.Login, time.Now().Unix()))
	if err != nil {
		utils.Logger.Println(err)
		return "", err
	}

	if _, err = stmt.Exec(hash, user.ID, time.Now().Unix()); err != nil {
		utils.Logger.Println(err)
		return "", err
	}

	if sessionMap == nil {
		sessionMap = make(map[string]Session)
	}

	sessionMap[hash] = Session{
		Hash:      hash,
		User:      user,
		CreatedAt: time.Now().Unix(),
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
			&session.User.Role.Value,
			&session.User.Role.TranslateValue,
		)
		if e != nil {
			fmt.Println(e)
			utils.Logger.Println(e)
			return
		}

		m[session.Hash] = session
	}
}

func (us *DefaultUserService) CheckAdmin(config *settings.Setting) error {
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
			if e = us.CreateAdmin(config); e != nil {
				utils.Logger.Println(e)
				return e
			}
		} else {
			utils.Logger.Println(e)
			return e
		}
	}

	encryptPass, e := us.Encrypt.Encrypt(config.SuperAdminPassword)
	if e != nil {
		utils.Logger.Println(e)
		return e
	}

	if encryptPass != admin.Password {
		admin.Password = encryptPass

		if e = us.ChangeUserPassword(&admin); e != nil {
			utils.Logger.Println(e)
			return e
		}

		if e = us.DeleteUserSessions(admin.ID); e != nil {
			utils.Logger.Println(e)
			return e
		}
	}

	return nil
}

func (us *DefaultUserService) ChangeUserPassword(user *User) error {
	stmt, ok := query["CHANGE_USER_PASSWORD"]
	if !ok {
		return errors.New("запрос CHANGE_USER_PASSWORD не подготовлен")
	}

	var err error

	user.Password, err = us.Encrypt.Encrypt(user.Password)
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	_, e := stmt.Exec(user.ID, user.Password)
	if e != nil {
		return e
	}

	return nil
}

func (us *DefaultUserService) CreateAdmin(config *settings.Setting) error {
	var admin User

	encryptPass, e := us.Encrypt.Encrypt(config.SuperAdminPassword)
	if e != nil {
		utils.Logger.Println(e)
		return e
	}

	roles, err := GetReferenceRecords("ROLES")
	if err != nil {
		utils.Logger.Println(err)
		return err
	}

	for _, role := range roles {
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

	if e = us.CreateUser(&admin); e != nil {
		utils.Logger.Println(e)
		return e
	}
	return nil
}

func (us *DefaultUserService) DeleteUserSessions(userID int) error {
	stmt, ok := query["GET_USER_SESSIONS"]
	if !ok {
		return errors.New("запрос GET_USER_SESSIONS не подготовлен")
	}

	rows, e := stmt.Query(userID)
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

		if e = us.DeleteSession(&session); e != nil {
			utils.Logger.Println(e)
			return e
		}
	}

	return nil
}

func (us *DefaultUserService) ValidateUser(user User, action string) bool {
	if len(user.Name) == 0 || len(user.Login) == 0 {
		return false
	}

	roles, err := GetReferenceRecords("ROLES")
	if err != nil {
		utils.Logger.Println(err)
		return false
	}

	validRole := false
	for _, role := range roles {
		if role.ID == user.Role.ID {
			validRole = true
			break
		}
	}

	if !validRole {
		return false
	}

	if action == "create" {
		if len(user.Password) < 6 {
			return false
		}
	} else if len(user.Password) != 0 {
		if len(user.Password) < 6 {
			return false
		}
	}

	return true
}

func NewUserService() UserService {
	return &DefaultUserService{
		Encrypt: &utils.DefaultEncrypt{},
	}
}

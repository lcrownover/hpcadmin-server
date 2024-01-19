package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/exp/slices"
)

type Pirg struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	OwnerId    int       `json:"owner_id"`
	AdminIds   []int     `json:"admin_ids"`
	UserIds    []int     `json:"user_ids"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type PirgRequest struct {
	Name     string `json:"name"`
	OwnerId  int    `json:"owner_id"`
	AdminIds []int  `json:"admin_ids"`
	UserIds  []int  `json:"user_ids"`
}

func GetAllPirgs(db *sql.DB) ([]*Pirg, error) {
	var pirgs []*Pirg
	rows, err := db.Query("SELECT id FROM pirgs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		pirg, err := GetPirgById(db, id)
		if err != nil {
			return nil, err
		}
		pirgs = append(pirgs, pirg)
	}
	return pirgs, nil
}

func GetPirgById(db *sql.DB, id int) (*Pirg, error) {
	slog.Debug("querying database for pirg", "id", id, "package", "data", "method", "GetPirgById")
	var pirg Pirg
	err := db.QueryRow("SELECT id, name, owner_id, created_at, modified_at FROM pirgs WHERE id = $1", id).Scan(&pirg.Id, &pirg.Name, &pirg.OwnerId, &pirg.CreatedAt, &pirg.ModifiedAt)
	if err != nil {
		slog.Error("failed to look up pirg from database", "package", "data", "method", "GetPirgById", "error", err)
		return nil, err
	}
	adminIds, err := getPirgAdminIds(db, id)
	if err != nil {
		return nil, err
	}
	pirg.AdminIds = adminIds
	userIds, err := getPirgUserIds(db, id)
	if err != nil {
		return nil, err
	}
	pirg.UserIds = userIds
	return &pirg, err
}

func getPirgAdminIds(db *sql.DB, id int) ([]int, error) {
	slog.Debug("getting pirg admin ids from database", "package", "data", "method", "getPirgAdminIds")
	var adminIds []int
	rows, err := db.Query("SELECT user_id FROM pirgs_admins WHERE pirg_id = $1", id)
	if err != nil {
		slog.Error("failed to look up pirg admins from database", "package", "data", "method", "getPirgAdminIds", "error", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var adminId int
		err := rows.Scan(&adminId)
		if err != nil {
			return nil, err
		}
		adminIds = append(adminIds, adminId)
	}
	return adminIds, err
}

func getPirgUserIds(db *sql.DB, id int) ([]int, error) {
	slog.Debug("getting pirg user ids from database", "package", "data", "method", "getPirgUserIds")
	var userIds []int
	rows, err := db.Query("SELECT user_id FROM pirgs_users WHERE pirg_id = $1", id)
	if err != nil {
		slog.Error("failed to look up pirg users from database", "package", "data", "method", "getPirgUserIds", "error", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userId int
		err := rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, userId)
	}
	return userIds, err
}

func CreatePirg(db *sql.DB, pirg *PirgRequest) (*Pirg, error) {
	slog.Debug("creating new pirg in database", "package", "data", "method", "CreatePirg")
	var newId int

	// verify that owner_id is a valid user
	err := validateUserId(db, pirg.OwnerId)
	if err != nil {
		return nil, fmt.Errorf("validating owner_id failed: %v", err)
	}
	// verify that all the admin_ids are users
	for _, adminId := range pirg.AdminIds {
		err = validateUserId(db, adminId)
		if err != nil {
			return nil, fmt.Errorf("validating admin_id failed: %v", err)
		}
	}
	// verify that all the user_ids are users
	for _, userId := range pirg.UserIds {
		err = validateUserId(db, userId)
		if err != nil {
			return nil, fmt.Errorf("validating user_id failed: %v", err)
		}
	}
	err = db.QueryRow("INSERT INTO pirgs (name, owner_id) VALUES ($1, $2) RETURNING id", pirg.Name, pirg.OwnerId).Scan(&newId)
	if err != nil {
		return nil, err
	}
	for adminId := range pirg.AdminIds {
		db.QueryRow("INSERT INTO pirgs_admins (pirg_id, user_id) VALUES ($1, $2)", newId, adminId)
	}
	for userId := range pirg.UserIds {
		db.QueryRow("INSERT INTO pirgs_users (pirg_id, user_id) VALUES ($1, $2)", newId, userId)
	}
	newPirg, err := GetPirgById(db, newId)
	if err != nil {
		return nil, err
	}
	return newPirg, err
}

func UpdatePirg(db *sql.DB, id int, pr *PirgRequest) (*Pirg, error) {
	slog.Debug("updating pirg in database", "package", "data", "method", "UpdatePirg")
	existingPirg, err := GetPirgById(db, id)
	if err != nil {
		return nil, err
	}
	// Updates name and owner_id if changed
	if pr.Name != existingPirg.Name || pr.OwnerId != existingPirg.OwnerId {
		slog.Debug("updating pirg name and owner_id", "name", pr.Name, "owner_id", pr.OwnerId, "package", "data", "method", "UpdatePirg")
		res, err := db.Exec("UPDATE pirgs SET name = $1, owner_id = $2 WHERE id = $3", pr.Name, pr.OwnerId, id)
		if err = checkAffectedRows(res, err); err != nil {
			return nil, err
		}
	}
	existingAdminIds, err := getPirgAdminIds(db, id)
	if err != nil {
		return nil, err
	}
	// Adds new admin ids
	for _, adminId := range pr.AdminIds {
		if !slices.Contains(existingAdminIds, adminId) {
			if err = addPirgAdmin(db, id, adminId); err != nil {
				return nil, err
			}
		}
	}
	// Removes admin ids not present in request
	for _, existingAdminId := range existingAdminIds {
		if !slices.Contains(pr.AdminIds, existingAdminId) {
			if err = deletePirgAdmin(db, id, existingAdminId); err != nil {
				return nil, err
			}
		}
	}
	existingUserIds, err := getPirgUserIds(db, id)
	if err != nil {
		return nil, err
	}
	// Adds new User ids
	for _, UserId := range pr.UserIds {
		if !slices.Contains(existingUserIds, UserId) {
			if err = addPirgUser(db, id, UserId); err != nil {
				return nil, err
			}
		}
	}
	// Removes User ids not present in request
	for _, existingUserId := range existingUserIds {
		if !slices.Contains(pr.UserIds, existingUserId) {
			if err = deletePirgUser(db, id, existingUserId); err != nil {
				return nil, err
			}
		}
	}
	newPirg, err := GetPirgById(db, id)
	if err != nil {
		return nil, err
	}
	return newPirg, err
}

func DeletePirg(db *sql.DB, id int) error {
	slog.Debug("deleting pirg from database", "package", "data", "method", "DeletePirg")
	res, err := db.Exec("DELETE FROM pirgs WHERE id = $1", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return err
	}
	return nil
}

func checkAffectedRows(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("expected to update 1 row, updated %d rows", count)
	}
	return nil
}

func addPirgAdmin(db *sql.DB, pirgId int, userId int) error {
	slog.Debug("adding pirg admin to database", "package", "data", "method", "addPirgAdmin")
	_, err := db.Exec("INSERT INTO pirgs_admins (pirg_id, user_id) VALUES ($1, $2)", pirgId, userId)
	return err
}

func deletePirgAdmin(db *sql.DB, pirgId int, userId int) error {
	slog.Debug("deleting pirg admin from database", "package", "data", "method", "deletePirgAdmin")
	_, err := db.Exec("DELETE FROM pirgs_admins WHERE pirg_id = $1 AND user_id = $2", pirgId, userId)
	return err
}

func addPirgUser(db *sql.DB, pirgId int, userId int) error {
	slog.Debug("adding pirg user to database", "package", "data", "method", "addPirgUser")
	_, err := db.Exec("INSERT INTO pirgs_users (pirg_id, user_id) VALUES ($1, $2)", pirgId, userId)
	return err
}

func deletePirgUser(db *sql.DB, pirgId int, userId int) error {
	slog.Debug("deleting pirg user from database", "package", "data", "method", "deletePirgUser")
	_, err := db.Exec("DELETE FROM pirgs_users WHERE pirg_id = $1 AND user_id = $2", pirgId, userId)
	return err
}

func updateModifiedDatestamp(db *sql.DB, table string, id int) error {
	slog.Debug("updating modified_at in database", "package", "data", "method", "updateModifiedDatestamp")
	_, err := db.Exec(fmt.Sprintf("UPDATE %s SET modified_at = NOW() WHERE id = $1", table), id)
	return err
}

func validateUserId(db *sql.DB, userId int) error {
	slog.Debug("validating user id", "id", userId, "package", "data", "method", "validateUserIds")
	_, err := GetUserById(db, userId)
	if err != nil {
		return fmt.Errorf("user does not exist with id: %d", userId)
	}
	return nil
}

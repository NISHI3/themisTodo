package module

import (
	"database/sql"
	"log"
	"../utils"
	"../models"
)

type AccountModule struct {
	db *sql.DB
}

func NewAccountModule(db *sql.DB) *AccountModule {
	return &AccountModule{db}
}

func (self *AccountModule) Add(name, password string) bool {
	stmt, err := self.db.Prepare("INSERT INTO `users` (`displayName`, `name`, `password`) VALUES(?, ?, ?);")

	if err != nil {
		log.Printf("AccountModule.Add Error: %+v", err)
		return true
	}

	defer stmt.Close()

	_, err = stmt.Exec(name, name, utils.SHA512(password))
	if err != nil {
		log.Printf("AccountModule.Add Error: %+v", err)
		return true
	}

	return false
}

func (self *AccountModule) Get(name string) int {
	rows, err := self.db.Query("SELECT `uuid` FROM `users` WHERE `name` = ?;", name)

	if err != nil {
		log.Printf("AccountModule.Get Error: %+v", err)
		return 0
	}

	defer rows.Close()

	var userId int

	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			log.Printf("AccountModule.Get Error: %+v", err)
			return 0
		}

		return userId
	}

	return 0
}

func (self *AccountModule) Search(searchObject *models.AccountSearchModel) (isError bool, users []models.AccountSearchResultModel) {
	queryText := "SELECT `uuid`, `name`, `displayName` FROM `users` WHERE "
	execArgs := []interface{}{}

	if searchObject.ProjectId > 0 {
		isInCmmand := "IN"
		if !searchObject.IsInProject {
			isInCmmand = "NOT IN"
		}
		queryText += "`users`.`uuid` " + isInCmmand + " (SELECT `user_id` FROM `users_in_projects` WHERE `project_id` = ? ORDER BY `user_id`) "
		execArgs = append(execArgs, searchObject.ProjectId)
	}

	if searchObject.ProjectId > 0 && (searchObject.Name != "" || searchObject.DisplayName != "" || searchObject.Uuid > 0) {
		queryText += " AND ( "
	}

	if searchObject.Name != "" {
		queryText += "`name` LIKE ? "
		execArgs = append(execArgs, "%"+searchObject.Name+"%")
	}

	if searchObject.Name != "" && searchObject.DisplayName != "" {
		queryText += " OR "
	}

	if searchObject.DisplayName != "" {
		queryText += "`displayName` LIKE ? "
		execArgs = append(execArgs, "%"+searchObject.DisplayName+"%")
	}

	if (searchObject.Name != "" || searchObject.DisplayName != "") && searchObject.Uuid > 0 {
		queryText += " OR "
	}

	if searchObject.Uuid > 0 {
		queryText += "`uuid` = ? "
		execArgs = append(execArgs, searchObject.Uuid)
	}

	if searchObject.ProjectId > 0 && (searchObject.Name != "" || searchObject.DisplayName != "" || searchObject.Uuid > 0) {
		queryText += " ) "
	}

	if len(execArgs) < 1 {
		queryText += "1"
	}

	queryText += " LIMIT 0,? ;"
	execArgs = append(execArgs, searchObject.Max)

	resultModel := make([]models.AccountSearchResultModel, 0)
	rows, err := self.db.Query(queryText, execArgs...)

	if err != nil {
		log.Printf("SQL Text is  [%s]", queryText)
		log.Printf("AccountModule.Search Error: %+v\n", err)
		return true, nil
	}

	for rows.Next() {
		resultModelOne := models.AccountSearchResultModel{}
		if err := rows.Scan(&resultModelOne.Uuid, &resultModelOne.Name, &resultModelOne.DisplayName); err != nil {
			log.Printf("AccountModule.Search Error: %+v\n", err)
			return true, nil
		}
		resultModel = append(resultModel, resultModelOne)
	}

	return false, resultModel
}

package repository

import (
	"backend/database"
	"backend/model"
	"fmt"
)

func CreateTask(task model.Task) error {
	_, err := database.DB.Exec(
		`INSERT INTO tasks(user_id, title, done)
		VALUES (?, ?, ?)`,
		task.UserID,
		task.Title,
		false,
	)

	return err
}

func GetTasks(UserID int) ([]model.Task, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, title, done, created_at
		FROM tasks
		WHERE user_id = ?
	`, UserID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Done,
			&task.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func UpdateTask(id int, task model.Task) error {
	// Capturamos o resultado em res
	res, err := database.DB.Exec(`
		UPDATE tasks
		SET title = ?, done = ?
		WHERE id = ? AND user_id = ?
	`, task.Title, task.Done, id, task.UserID)

	if err != nil {
		return err
	}

	// Se rows == 0, disparamos aviso de erro com a biblioteca fmt
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("tarefa não encontrada ou sem permissão")
	}

	return nil
}

func DeleteTask(id int, userID int) error {
	res, err := database.DB.Exec(
		`DELETE FROM tasks
		WHERE id = ? AND user_id = ?`,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	// Verificação de segurança:
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("tarefa não encontrada ou sem permissão")
	}

	return nil
}
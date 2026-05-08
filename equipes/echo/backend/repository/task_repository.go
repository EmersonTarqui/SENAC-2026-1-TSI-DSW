package repository

import (
	"backend/database"
	"backend/model"
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

func GetTasks() ([]model.Task, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, title, done, created_at
		FROM tasks
	`)
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
	_, err := database.DB.Exec(`
		UPDATE tasks
		SET title = ?, done = ?
		WHERE id = ? AND user_id = ?
	`,
		task.Title,
		task.Done,
		id,
		task.UserID,
	)

	return err
}

func DeleteTask(id int, userID int) error {
	_, err := database.DB.Exec(
		`DELETE FROM tasks
		WHERE id = ? AND user_id = ?`,
		id,
		userID,
	)

	return err
}
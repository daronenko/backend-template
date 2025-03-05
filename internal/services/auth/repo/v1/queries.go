package repo

const (
	createUserQuery = `
		INSERT INTO "user" (
			username, email, password, role, avatar, created_at, updated_at, login_date
		) VALUES (
		 	$1, $2, $3, COALESCE(NULLIF($4, ''), 'user'), $5, now(), now(), now()
		) RETURNING *
	`

	updateUserQuery = `
		UPDATE "user"
		SET
			username = COALESCE(NULLIF($1, ''), username),
			email = COALESCE(NULLIF($2, ''), email),
			role = COALESCE(NULLIF($3, ''), role),
			avatar = COALESCE(NULLIF($4, ''), avatar),
			updated_at = now()
		WHERE id = $5
		RETURNING *
	`

	deleteUserQuery = `
		DELETE FROM "user" WHERE id = $1
	`

	getUserQuery = `
		SELECT id, email, role, avatar, created_at, updated_at, login_date
		FROM "user"
		WHERE id = $1
	`

	getUserByEmailQuery = `
		SELECT id, username, email, role, avatar, created_at, updated_at, login_date, password
		FROM "user"
		WHERE email = $1
	`

	getTotalUsersQuery = `
		SELECT COUNT(id) FROM "user"
	`

	findUsersByNameQuery = `
		SELECT id, username, email, role, avatar, created_at, updated_at, login_date 
		FROM "user"
		WHERE username ILIKE '%' || $1 || '%'
		ORDER BY username
		OFFSET $2 LIMIT $3
	`

	getUsers = `
		SELECT id, username, email, role, avatar, created_at, updated_at, login_date
		FROM "user"
		ORDER BY COALESCE(NULLIF($1, ''), username)
		OFFSET $2 LIMIT $3
	`
)

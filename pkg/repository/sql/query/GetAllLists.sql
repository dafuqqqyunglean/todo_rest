SELECT tl.id, tl.title, tl.description FROM todo_lists tl INNER JOIN users_lists ul on tl.id = ul.list_id WHERE ul.user_id = $1

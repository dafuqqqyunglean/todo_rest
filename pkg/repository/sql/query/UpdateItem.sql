UPDATE todo_items ti SET %s FROM lists_items li, users_lists ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d
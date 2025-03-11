package models

import (
	"database/sql"
	"fmt"

	"reat-time-forum/utils"
)

func GetCategory(category string) int {
	var categoryId int
	getCategoryIdQuery := "SELECT id FROM categories WHERE name = ?"
	if err := db.QueryRow(getCategoryIdQuery, category).Scan(&categoryId); err != nil {
		if err == sql.ErrNoRows {
			return 0
		} else {
			fmt.Println(err)
			return -1
		}
	}
	return categoryId
}

func FilterPostsByCategory(user_id, categoryId, offset int) []utils.Posts {
	getPostsQuery := `
		SELECT 
			p.id, 
			p.title, 
			p.content,
			u.nickname,
			p.createdAt,
			COALESCE(COUNT(CASE WHEN r.reaction_type = 'like' THEN 1 END), 0) AS likes,
			COALESCE(COUNT(CASE WHEN r.reaction_type = 'dislike' THEN 1 END), 0) AS dislikes,
			COALESCE((
				SELECT reaction_type
				FROM reactions
				WHERE user_id = ? AND post_id = p.id
			), '') AS user_reaction
		FROM posts p
		JOIN users u ON u.id = p.user_id
		JOIN posts_categories pc ON p.id = pc.post_id
		LEFT JOIN reactions r ON r.post_id = p.id
		WHERE pc.category_id = ?
		GROUP BY p.id, p.title, p.content, u.nickname, p.createdAt
		ORDER BY p.id DESC
		LIMIT 20 OFFSET ?;
	`

	rows, err := db.Query(getPostsQuery, user_id, categoryId, offset)
	if err != nil {
		fmt.Println("erro while getting posts:", err)
	}

	defer rows.Close()

	var posts []utils.Posts

	for rows.Next() {
		var post utils.Posts
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Poster, &post.CreatedAt, &post.Reactions.Likes, &post.Reactions.Dislikes, &post.Reactions.Action); err != nil {
			fmt.Println("Row scan error:", err)
		}
		posts = append(posts, post)
	}

	defer rows.Close()
	GetPostsCategories(posts)
	return posts
}

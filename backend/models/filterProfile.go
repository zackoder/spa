package models

import (
	"fmt"
	"reat-time-forum/utils"
)

func FilterProfile(user_id int, nickname string, offset int) []utils.Posts {
	query := `
			SELECT
			    p.id,
			    p.title,
			    p.content,
			    u.nickname,
			    p.createdAt,
			    COALESCE(SUM(CASE WHEN r.reaction_type = 'like' THEN 1 ELSE 0 END), 0) AS likes,
			    COALESCE(SUM(CASE WHEN r.reaction_type = 'dislike' THEN 1 ELSE 0 END), 0) AS dislikes,
			    COALESCE((
			        SELECT reaction_type
			        FROM reactions
			        WHERE user_id = ? AND post_id = p.id
			    ), '') AS user_reaction
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN reactions r ON r.post_id = p.id
			WHERE u.nickname = ?
			GROUP BY p.id
			ORDER BY p.id DESC
			LIMIT 20 OFFSET ?;
	`
	rows, err := db.Query(query, user_id, nickname, offset)
	if err != nil {
		fmt.Printf("error while geting user %s profile\n", nickname)
	}

	var posts []utils.Posts

	for rows.Next() {
		var post utils.Posts
		if err := rows.Scan(&post.Id, &post.Title, &post.Content, &post.Poster, &post.CreatedAt, &post.Reactions.Likes, &post.Reactions.Dislikes, &post.Reactions.Action); err != nil {
			fmt.Println(err)
		}
		posts = append(posts, post)
	}

	GetPostsCategories(posts)

	defer rows.Close()
	return posts
}

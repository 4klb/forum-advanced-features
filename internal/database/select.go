package database

import (
	"database/sql"
	forum "forum/internal"
	"log"
	"net/http"
	"strconv"
)

//GetPostsByCategories ..
func SelectPostsByCategories(db *sql.DB, tag string) ([]forum.Post, error) {
	var posts []forum.Post
	rows, err := db.Query(`SELECT Post.id, Post.Title , Post.Post , Post.AuthorId
	FROM Categories 
	JOIN Post ON PostId = Post.id AND Categories.Name =? 
	ORDER BY Post.id DESC`, tag)
	if err != nil {
		log.Println(err)
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		var postid, authorid int
		var title, text, authorname string
		rows.Scan(&postid, &title, &text, &authorid)
		tags := SelectOnePostCategories(db, postid)
		if err := db.QueryRow(`SELECT Name FROM SignInUser 
		WHERE id = ?`, authorid).Scan(&authorname); err != nil {
			log.Println(err)
			return posts, err
		}
		post := forum.Post{
			ID:     postid,
			Title:  title,
			Text:   text,
			Author: authorname,
			Tags:   tags,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

//GetPostByUser ..
func SelectPostByUser(db *sql.DB, user forum.User) ([]forum.Post, error) {
	var postArr []forum.Post
	rows, err := db.Query(`SELECT id, Title FROM POST 
	WHERE AuthorId = ? 
	ORDER BY id DESC`, user.ID)
	if err != nil {
		log.Println(err)
		return postArr, err
	}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			log.Println(err)
			return postArr, err
		}
		curPost := forum.Post{
			ID:    id,
			Title: title,
		}
		postArr = append(postArr, curPost)
	}
	return postArr, nil
}

//GetCommentsFromPost ..
func SelectCommentsFromPost(db *sql.DB, postID int) ([]forum.Comment, error) {
	var comments []forum.Comment
	rows, err := db.Query(`SELECT Comments.id, Comment, SignInUser.Name FROM Comments
	JOIN SignInUser ON CommentatorId = SignInUser.Id 
	WHERE PostId = ? 
	ORDER BY Comments.id DESC`, postID)
	if err != nil {
		log.Println(err)
		return comments, err
	}

	defer rows.Close()

	for rows.Next() {
		var comment, userName string
		var id, like, dislike int
		if err := rows.Scan(&id, &comment, &userName); err != nil {
			return comments, err
		}
		db.QueryRow(`SELECT SUM(Like),SUM(DisLike) FROM CommentRaiting 
		WHERE CommentsId =? `, id).Scan(&like, &dislike) //err?

		Comment := forum.Comment{
			ID:              id,
			Text:            comment,
			Author:          userName,
			CountOfLikes:    like,
			CountOfDisLikes: dislike,
		}
		comments = append(comments, Comment)
	}
	return comments, err
}

//GetLikedPostsByUser ..
func SelectLikedPostsByUser(db *sql.DB, user forum.User) ([]forum.Post, error) {
	var postArr []forum.Post
	rows, err := db.Query(`SELECT Post.id, Post.Title FROM PostRaiting 
	JOIN Post ON PostId=Post.id WHERE LikerId=? AND Like=true 
	ORDER BY Post.id DESC`, user.ID)
	if err != nil {
		log.Println(err)
		return postArr, err
	}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			log.Println(err)
			return postArr, err
		}
		curPost := forum.Post{
			ID:    id,
			Title: title,
		}
		postArr = append(postArr, curPost)
	}
	return postArr, nil
}

func SelectDislikedPostsByUser(db *sql.DB, user forum.User) ([]forum.Post, error) {
	var postArr []forum.Post
	rows, err := db.Query(`SELECT Post.id, Post.Title FROM PostRaiting 
	JOIN Post ON PostId=Post.id WHERE LikerId=? AND Like=false 
	ORDER BY Post.id DESC`, user.ID)
	if err != nil {
		log.Println(err)
		return postArr, err
	}
	for rows.Next() {
		var id int
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			log.Println(err)
			return postArr, err
		}
		curPost := forum.Post{
			ID:    id,
			Title: title,
		}
		postArr = append(postArr, curPost)
	}
	return postArr, nil
}

//GetOnePostCategories ..
func SelectOnePostCategories(db *sql.DB, postid int) []string {
	var tags []string
	row, err := db.Query(`
	SELECT Name FROM Categories
	WHERE PostId=?;`, postid)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer row.Close()
	for row.Next() {
		var catname string
		if err := row.Scan(&catname); err != nil {
			log.Println(err)
			return nil
		}
		tags = append(tags, catname)
	}
	if row.Err() != nil {
		return nil
	}
	return tags
}

//GetPosts ..
func SelectPosts(db *sql.DB) ([]forum.Post, error) {
	var posts []forum.Post
	rows, err := db.Query(`
	SELECT  Post.id, Post.Title, Post.Post, SignInUser.Name FROM Post  
	JOIN SignInUser ON AuthorId = SignInUser.id 
	ORDER BY Post.id DESC `)
	if err != nil {
		log.Println(err)
		return posts, err
	}
	defer rows.Close()
	for rows.Next() {
		var postid int
		var title, text, author string
		if err := rows.Scan(&postid, &title, &text, &author); err != nil {
			log.Println(err)
			return posts, err
		}
		tags := SelectOnePostCategories(db, postid)

		post := forum.Post{
			ID:     postid,
			Title:  title,
			Text:   text,
			Author: author,
			Tags:   tags,
		}
		posts = append(posts, post)

	}

	if rows.Err() != nil {
		return posts, err
	}
	return posts, nil
}

func SelectUserCommentsByComment(db *sql.DB, commentid int, cookie *http.Cookie) ([]forum.Comment, error) { //
	var comments []forum.Comment
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return comments, err
	}
	rows, err := db.Query(`SELECT id, Comment FROM Comments
	WHERE CommentatorId = ? AND id = ?`, user.ID, commentid)
	if err != nil {
		log.Println(err)
		return comments, err
	}
	for rows.Next() {
		var id int
		var commentText string
		if err := rows.Scan(&id, &commentText); err != nil {
			return comments, err
		}
		comment := forum.Comment{
			ID:   id,
			Text: commentText,
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func SelectOneUserCommentsByComment(db *sql.DB, commentid int, cookie *http.Cookie) (forum.Comment, error) { //
	var comment forum.Comment
	var id int
	var commentText string
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return comment, err
	}
	if err := db.QueryRow(`SELECT id, Comment FROM Comments
	WHERE CommentatorId = ? AND id = ?`, user.ID, commentid).Scan(&id, &commentText); err != nil {
		log.Println(err)
		return comment, err
	}
	comment = forum.Comment{
		ID:   id,
		Text: commentText,
	}
	return comment, nil
}

//GetOnePost ..
func SelectOnePost(db *sql.DB, postid int) (forum.Post, error) {
	var post forum.Post
	row := db.QueryRow(`SELECT Post.Title, Post.Post, SignInUser.Name 
	FROM Post
	JOIN SignInUser ON AuthorId = SignInUser.id 
	WHERE Post.id =?`, postid)
	var title, text, author string
	if err := row.Scan(&title, &text, &author); err != nil {
		return post, err
	}
	rows, err := db.Query(`
		SELECT Name FROM Categories
		WHERE PostId=?;`, postid)
	if err != nil {
		return post, err
	}
	var tags []string
	for rows.Next() {
		var categname string
		if err := rows.Scan(&categname); err != nil {
			return post, err
		}
		tags = append(tags, categname)
	}
	result, err := SelectCommentsFromPost(db, postid)
	if err != nil {
		log.Println(err)
		return post, err
	}
	countlikes, countdislikes, err := SelectPostLikes(db, postid)
	if err != nil {
		log.Println(err)
		return post, err
	}
	post = forum.Post{
		ID:              postid,
		Title:           title,
		Text:            text,
		Author:          author,
		Tags:            tags,
		Comments:        result,
		CountOfLikes:    countlikes,
		CountOfDisLikes: countdislikes,
	}
	return post, nil
}

//IsPostInDB ..
func IsPostInDB(db *sql.DB, id string) (bool, int) {
	postid, err := strconv.Atoi(id)
	if err != nil {
		return false, postid
	}
	var val int
	err = db.QueryRow("SELECT id FROM Post WHERE id = ?", postid).Scan(&val)
	if err == nil && err == sql.ErrNoRows {
		return false, postid
	}
	if val == 0 {
		return false, postid
	}
	return true, postid
}

//GetPostLikes ..
func SelectPostLikes(db *sql.DB, postID int) (int, int, error) {
	rows, err := db.Query("SELECT Like, DisLike FROM PostRaiting WHERE PostId = ?", postID)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	var like, dislike int
	for rows.Next() {
		var pos, neg int
		if err := rows.Scan(&pos, &neg); err != nil {
			return 0, 0, err
		}
		like += pos
		dislike += neg
	}
	return like, dislike, nil
}

func SelectPostIdByUser(db *sql.DB, cookie *http.Cookie) ([]int, error) {
	var pArr []int
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return pArr, err
	}
	rows, err := db.Query(`SELECT Post.Id FROM PostRaiting 
	JOIN Post ON PostId = Post.id
	WHERE Post.AuthorId = ?`, user.ID)
	if err != nil {
		log.Println(err)
		return pArr, err
	}
	for rows.Next() {
		var postId int
		if err = rows.Scan(&postId); err != nil {
			log.Println(err)
			return pArr, err
		}
		pArr = append(pArr, postId)
	}
	return pArr, nil
}

func IsLikeInUserPost(db *sql.DB, user forum.User) ([]bool, error) {
	var values []bool

	rows, err := db.Query(`SELECT Like FROM PostRaiting 
	JOIN POST ON Postid = Post.id 
	WHERE AuthorId = ?`, user.ID)

	if err != nil {
		log.Println(err)
		return values, err
	}

	for rows.Next() {
		var like bool
		if err := rows.Scan(&like); err != nil {
			return values, err
		}
		values = append(values, like)
	}
	return values, nil
}

func IsPostRaitingTableEmpty(db *sql.DB, user forum.User) bool {
	var id sql.NullString
	row := db.QueryRow(`SELECT PostRaiting.id FROM PostRaiting 
	JOIN POST ON Postid = Post.id 
	WHERE AuthorId = ?`, user.ID)
	if err := row.Scan(&id); err != nil {
		log.Println(err, "There are no any values in PostRaiting table for the user")
		return true
	}
	if id.String == "" {
		return true
	}
	return false
}

func IsCommentsTableEmpty(db *sql.DB, user forum.User) bool {
	var id sql.NullString
	row := db.QueryRow(`SELECT Comments.Id FROM Comments 
	JOIN POST ON Postid = Post.id 
	WHERE AuthorId = ?`, user.ID)
	if err := row.Scan(&id); err != nil {
		log.Println(err, "There are no any values in Comments table for the user")
		return true
	}
	if id.String == "" {
		return true
	}
	return false
}

func SelectIncomingLikedPostNotifications(db *sql.DB, user forum.User, visit bool) ([]forum.PostNotification, []int, error) {
	var notifications []forum.PostNotification
	var postRaitingIdArr []int
	rows, err := db.Query(`SELECT Post.id, Post.Title, SignInUser.Id, SignInUser.Name, PostRaiting.id
	FROM VisitedNotificationsPost 
	JOIN PostRaiting ON PostRaitingId = PostRaiting.id
	JOIN Post ON PostId = Post.id 
	JOIN SignInUser ON LikerId = SignInUser.id 
	WHERE Like = true AND Visited = ? AND AuthorId = ? AND LikerId != ?
	ORDER BY PostRaiting.id DESC`, visit, user.ID, user.ID)
	if err != nil {
		log.Println(err)
		return notifications, postRaitingIdArr, err
	}
	for rows.Next() {
		var postId, likerId, postRaitingId int
		var postTitle, likerName string
		if err := rows.Scan(&postId, &postTitle, &likerId, &likerName, &postRaitingId); err != nil {
			return notifications, postRaitingIdArr, err
		}
		notification := forum.PostNotification{
			Post: forum.Post{
				ID:    postId,
				Title: postTitle,
			},
			Liker: forum.SingInUsers{
				ID:   likerId,
				Name: likerName,
			},
		}
		postRaitingIdArr = append(postRaitingIdArr, postRaitingId)
		notifications = append(notifications, notification)
	}
	return notifications, postRaitingIdArr, nil
}

func SelectIncomingDislikedPostNotifications(db *sql.DB, user forum.User, visit bool) ([]forum.PostNotification, []int, error) {
	var notifications []forum.PostNotification
	var postRaitingIdArr []int
	rows, err := db.Query(`SELECT Post.id, Post.Title, SignInUser.Id, SignInUser.Name, PostRaiting.id
	FROM VisitedNotificationsPost 
	JOIN PostRaiting ON PostRaitingId = PostRaiting.id
	JOIN Post ON PostId = Post.id 
	JOIN SignInUser ON LikerId = SignInUser.id 
	WHERE Like = false AND Visited = ? AND AuthorId = ? AND LikerId != ?
	ORDER BY PostRaiting.id DESC`, visit, user.ID, user.ID)
	if err != nil {
		log.Println(err)
		return notifications, postRaitingIdArr, err
	}
	for rows.Next() {
		var postId, likerId, postRaitingId int
		var postTitle, likerName string
		if err := rows.Scan(&postId, &postTitle, &likerId, &likerName, &postRaitingId); err != nil {
			return notifications, postRaitingIdArr, err
		}
		notification := forum.PostNotification{
			Post: forum.Post{
				ID:    postId,
				Title: postTitle,
			},
			Liker: forum.SingInUsers{
				ID:   likerId,
				Name: likerName,
			},
		}
		postRaitingIdArr = append(postRaitingIdArr, postRaitingId)
		notifications = append(notifications, notification)
	}
	return notifications, postRaitingIdArr, nil
}

func SelectCommentedPostsByUser(db *sql.DB, user forum.User) ([]forum.CommentedPosts, error) {
	var commentedPosts []forum.CommentedPosts
	rows, err := db.Query(`SELECT Post.Id, Title, Comments.id, Comment FROM Comments
	JOIN SignInUser ON CommentatorId = SignInUser.id
	JOIN Post ON PostId = Post.id
	WHERE CommentatorId = ? 
	ORDER BY Comments.id DESC`, user.ID)
	if err != nil {
		log.Println(err)
		return commentedPosts, err
	}
	for rows.Next() {
		var postId, commentId int
		var postTitle, commentText string
		if err := rows.Scan(&postId, &postTitle, &commentId, &commentText); err != nil {
			return commentedPosts, err
		}
		commentedPost := forum.CommentedPosts{
			Post: forum.Post{
				ID:    postId,
				Title: postTitle,
			},
			Comment: forum.Comment{
				ID:   commentId,
				Text: commentText,
			},
		}
		commentedPosts = append(commentedPosts, commentedPost)
	}
	return commentedPosts, nil
}

func SelectIncomingCommentedPostNotifications(db *sql.DB, user forum.User, visit bool) ([]forum.CommentNotification, []int, error) {
	var notifications []forum.CommentNotification
	var commentsIdArr []int
	rows, err := db.Query(`SELECT Post.id, Post.Title, SignInUser.Id, SignInUser.Name, Comments.id
	FROM VisitedNotificationsComment
	JOIN Comments ON CommentsId = Comments.id
	JOIN Post ON PostId = Post.id 
	JOIN SignInUser ON CommentatorId = SignInUser.id 
	WHERE Visited = ? AND AuthorId = ? AND CommentatorId != ?
	ORDER BY Comments.id DESC`, visit, user.ID, user.ID)
	if err != nil {
		log.Println(err)
		return notifications, commentsIdArr, err
	}
	for rows.Next() {
		var postId, commentatorId, commentId int
		var postTitle, commentatorName string
		if err := rows.Scan(&postId, &postTitle, &commentatorId, &commentatorName, &commentId); err != nil {
			return notifications, commentsIdArr, err
		}
		notification := forum.CommentNotification{
			Post: forum.Post{
				ID:    postId,
				Title: postTitle,
			},
			Commentator: forum.SingInUsers{
				ID:   commentatorId,
				Name: commentatorName,
			},
		}
		commentsIdArr = append(commentsIdArr, commentId)
		notifications = append(notifications, notification)
	}
	return notifications, commentsIdArr, nil
}

func SelectPostByComment(db *sql.DB, commentId int) (int, error) {
	var postId int
	if err := db.QueryRow(`SELECT PostId FROM Comments 
	WHERE id = ?`, commentId).Scan(&postId); err != nil {
		log.Println(err)
		return postId, err
	}
	return postId, nil
}

func SelectOneComment(db *sql.DB, commentId int, cookie *http.Cookie) (forum.Comment, error) {
	var comment forum.Comment
	var id int
	var commentText string
	user, err := GetUserFromDB(db, cookie)
	if err != nil {
		log.Println(err)
		return comment, err
	}
	if err := db.QueryRow(`SELECT id, Comment FROM Comments 
	WHERE id = ? AND CommentatorId = ?`, commentId, user.ID).Scan(&id, &commentText); err != nil {
		log.Println(err)
		return comment, err
	}
	comment = forum.Comment{
		ID:   id,
		Text: commentText,
	}
	return comment, nil
}

//omitempty ????

func SelectNotificationsHistory(db *sql.DB, user forum.User) ([]forum.Notification, error) {
	var notification []forum.Notification
	rows, err := db.Query(selectNotifications, user.ID, user.ID)
	if err != nil {
		log.Println(err)
		return notification, err
	}
	for rows.Next() {
		var n forum.Notification
		if err := rows.Scan(&n.ID, &n.Title, &n.Comment, &n.LikerName, &n.CommentatorName, &n.CommentedPost); err != nil {
			if err := rows.Scan(&n); err != nil {
				log.Println(err)
				return notification, err
			}
		}
		notification = append(notification, n)
	}
	// for _, v := range notification {
	// 	if v.Title != nil {
	// 		fmt.Println("v.Title", *v.Title)
	// 	} else {
	// 		fmt.Println("v.Comment => ", *v.Comment)
	// 	}
	// }

	return notification, nil
}

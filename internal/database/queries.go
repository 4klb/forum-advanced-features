package database

const selectNotifications = `SELECT UserNotificationsHistory.id, p.Title, c.Comment, su.Name, s.Name, pp.Title
FROM UserNotificationsHistory
LEFT JOIN VisitedNotificationsPost ON VisitedNotificationsPostId = VisitedNotificationsPost.id
LEFT JOIN VisitedNotificationsComment ON VisitedNotificationsCommentId = VisitedNotificationsComment.id
LEFT JOIN PostRaiting pr ON PostRaitingId = pr.id
LEFT JOIN Post p ON pr.PostId = p.id
LEFT JOIN Comments c ON CommentsId = c.id
LEFT JOIN SignInUser su ON LikerId = su.id
LEFT JOIN SignInUser s ON CommentatorId = s.id
LEFT JOIN Post pp ON c.PostId = pp.id
WHERE p.AuthorId = ? OR pp.AuthorId = ?
ORDER BY UserNotificationsHistory.id DESC`

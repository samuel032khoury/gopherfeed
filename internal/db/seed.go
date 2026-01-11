package db

import (
	"context"
	"log"
	"strconv"

	"github.com/samuel032khoury/gopherfeed/internal/store"
	"github.com/samuel032khoury/gopherfeed/internal/utils"
)

func Seed(store *store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, nil, user); err != nil {
			log.Println("failed to create user:", err)
			return
		}
	}

	posts := generatePosts(users, 200)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("failed to create post:", err)
			return
		}
	}

	comments := generateComments(users, posts, 500)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("failed to create comment:", err)
			return
		}
	}
	log.Println("Database seeding completed successfully.")
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	hashedPassword, _ := utils.EncryptPassword("password")
	for i := range n {
		users[i] = &store.User{
			Username: "user" + strconv.Itoa(i+1),
			Email:    "user" + strconv.Itoa(i+1) + "@example.com",
			Password: hashedPassword,
			RoleID:   1, // Default role ID
		}
	}
	return users
}

func generatePosts(users []*store.User, n int) []*store.Post {
	posts := make([]*store.Post, n)
	for i := range n {
		user := users[i%len(users)]
		posts[i] = &store.Post{
			Title:   "Post Title " + strconv.Itoa(i),
			Content: "This is the content of post number " + strconv.Itoa(i),
			UserID:  user.ID,
			Tags:    []string{"tag1", "tag2"},
		}
	}
	return posts
}

func generateComments(users []*store.User, posts []*store.Post, n int) []*store.Comment {
	comments := make([]*store.Comment, n)
	for i := range n {
		user := users[i%len(users)]
		post := posts[i%len(posts)]
		comments[i] = &store.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: "This is comment number " + strconv.Itoa(i),
		}
	}
	return comments
}

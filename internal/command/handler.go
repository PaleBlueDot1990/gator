package command

import (
	"context"
	"fmt"
	"time"

	"github.com/PaleBlueDot1990/gator/internal/config"
	"github.com/PaleBlueDot1990/gator/internal/database"
	"github.com/PaleBlueDot1990/gator/internal/rssfeed"
	"github.com/google/uuid"
)

func HandlerLogin(state *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login command expects a single argument: the username")
	}

	username := cmd.Args[0]
	user, err := state.DbQueries.GetUser(context.Background(), username)
	if err != nil {
		return err 
	}

	err = state.Cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s have been logged in!\n", user.Name)
	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Printf("User Name: %s\n", user.Name)
	fmt.Printf("User Created Timestamp: %s\n", user.CreatedAt)
	fmt.Printf("User Updated Timestamp: %s\n", user.UpdatedAt)
	return nil
}

func HandlerRegister(state *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the register command expetcs a single argument: the username")
	}

	dbQueryArgs := database.CreateUserParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Args[0],
	}

	user, err := state.DbQueries.CreateUser(context.Background(), dbQueryArgs)
	if err != nil {
		return err
	}

	err = state.Cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s has been registered!\n", user.Name)
	fmt.Printf("User ID: %s\n", user.ID)
	fmt.Printf("User Name: %s\n", user.Name)
	fmt.Printf("User Created Timestamp: %s\n", user.CreatedAt)
	fmt.Printf("User Updated Timestamp: %s\n", user.UpdatedAt)
	return nil
}

func HandleReset(state *config.State, cmd Command) error {
	err := state.DbQueries.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("All the users have been deleted!\n")
	return nil
}

func HandleUsers(state *config.State, cmd Command) error {
	users, err := state.DbQueries.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name != state.Cfg.CURRENT_USER_NAME {
			fmt.Printf("* %s\n", user.Name)
		} else {
			fmt.Printf("* %s (current)\n", user.Name)
		}
	}
	return nil 
}

func HandleAddFeed(state *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("the addfeed command expects two arguments: name of the feed, url of the feed")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	dbQueryArgs := database.CreateFeedParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: feedName,
		Url: feedUrl,
		UserID: user.ID,
	}

	feed, err := state.DbQueries.CreateFeed(context.Background(), dbQueryArgs)
	if err != nil {
		return err 
	}

	dbQueryArgs2 := database.CreateFeedFollowsParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	_, err = state.DbQueries.CreateFeedFollows(context.Background(), dbQueryArgs2)
	if err != nil {
		return err
	}

	fmt.Printf("ID of Feed: %v\n", feed.ID)
	fmt.Printf("Name of Feed: %s\n", feed.Name)
	fmt.Printf("Feed Created Time Stamp: %v\n", feed.CreatedAt)
	fmt.Printf("Feed Updated Time Stamp: %v\n", feed.UpdatedAt)
	fmt.Printf("URL of Feed: %s\n", feed.Url)
	fmt.Printf("User ID of Feed Owner: %v\n", feed.UserID)
	return nil 
}

func HandleFeeds(state *config.State, cmd Command) error {
	feeds, err := state.DbQueries.GetFeeds(context.Background())
	if err != nil {
		return err 
	}

	for _, feed := range feeds {
		user, err := state.DbQueries.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err 
		}

		fmt.Printf("Name of Feed: %s\n", feed.Name)
		fmt.Printf("URL of Feed: %s\n", feed.Url)
		fmt.Printf("Author of Feed: %s\n\n", user.Name)
	}

	return nil
}

func HandleFollow(state *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the follow command expects one argument: url of the feed to follow")
	}

	feedUrl := cmd.Args[0]
	feed, err := state.DbQueries.GetFeedsByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	dbQueryArgs := database.CreateFeedFollowsParams {
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	feedFollows, err := state.DbQueries.CreateFeedFollows(context.Background(), dbQueryArgs)
	if err != nil {
		return err 
	}

	fmt.Printf("Name of Current User: %s\n", feedFollows.UserName)
	fmt.Printf("Name of Feed Followed by User: %s\n", feedFollows.FeedName)
	return nil;
}

func HandleFollowing(state *config.State, cmd Command, user database.User) error {
	feedFollowsForUser, err := state.DbQueries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return nil 
	}

	if len(feedFollowsForUser) == 0 {
		fmt.Printf("You are not following any feed!\n")
		return nil
	}

	fmt.Printf("Feeds followed by User: %s\n", user.Name)
	for _, feedFollow := range feedFollowsForUser {
		fmt.Printf(" - %s\n", feedFollow.FeedName)
	}
	return nil
}

func HandleUnfollow(state *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the unfollow command expetcs a single argument: the url of the feed to unfollow")
	}

	feedUrl := cmd.Args[0]
	feed, err := state.DbQueries.GetFeedsByUrl(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	dbQueryArgs := database.DeleteFeedFollowsParams {
		FeedID: feed.ID,
		UserID: user.ID,
	}

	err = state.DbQueries.DeleteFeedFollows(context.Background(), dbQueryArgs)
	if err != nil {
		return err
	}

	fmt.Printf("You have unfollowed the feed!\n")
	return nil 
}

func HandleAgg(state *config.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("the agg command expetcs a single argument: the duration between each feed scraping job")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	
	for ; ; <-ticker.C {
		err = rssfeed.ScrapeFeeds(context.Background(), state)
		if err != nil {
			return err
		}
	}
}
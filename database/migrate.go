package database

import (
	"context"
	"database/sql"
	"log"
)

func Migrate(conn *sql.Conn) {
	log.Println("Start Migrate")
	createUserStm := `CREATE TABLE IF NOT EXISTS VkUserModel 
	(
		Id INTEGER,
		Avatar Text,
		Name Text,
		PRIMARY KEY(Id)
	)
	`
	log.Println("Create VkUserModel")
	log.Println(createUserStm)
	_, createRes := conn.ExecContext(context.Background(), createUserStm)
	if createRes != nil {
		log.Fatal(createRes)
	}

	log.Println("Created VkUserModel")

	createUserStm = `
	CREATE TABLE IF NOT EXISTS messages(
		Id Integer,
		FromId Integer,
		Date DateTime,
		Images TExt,
		LikesCount integer,
		Owner Text,
		OwnerId Integer,
		RepostedFrom integer,
		RepostsCount Integer,
		UserReposted Boolean,
		Text text,
		Primary Key(Id, OwnerId) 
		)`
	_, crecreateRes := conn.ExecContext(context.Background(), createUserStm)

	if crecreateRes != nil {
		log.Fatal(createRes)
	}

	log.Println("Create Fulltext search table")

	createUserStm = `
	CREATE VIRTUAL TABLE IF NOT EXISTS messages_search USING fts5(Id, OwnerId, Text)
	`

	_, crecreateRes = conn.ExecContext(context.Background(), createUserStm)

	if crecreateRes != nil {
		log.Fatal("Error occured during creating virtual table: ", createUserStm)
	}

	log.Println("Created messages")

	log.Println("Create on create trigger for full text search")
	createUserStm = `
	CREATE TRIGGER IF NOT EXISTS TR_messages_AI AFTER INSERT on messages
	BEGIN
		INSERT INTO messages_search (Id, OwnerId, Text) VALUES (new.Id, new.OwnerId, new.Text);
	END;
	`
	_, crecreateRes = conn.ExecContext(context.Background(), createUserStm)

	if crecreateRes != nil {
		log.Fatalln("Error creating TRIGGER on message: ", crecreateRes)
	}

	log.Println("Trigger created")

	log.Println("Stop migrate")
}

package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	config  *Config
	db      *mongo.Database
	userCol *UserCol
	docCol  *DocCol
}

func New(config *Config) *Database {
	return &Database{config: config}
}

//connect to db and ping it
func (c *Database) Open() error {
	clientOptions := options.Client().ApplyURI(c.config.DatabaseURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}

	c.db = client.Database(c.config.DatabaseName)
	return nil
}

//close db connection
func (c *Database) Close() error {
	return c.db.Client().Disconnect(context.TODO())
}

//access to userCol
func (c *Database) User() *UserCol {
	if c.userCol == nil {
		c.userCol = c.NewUserCol()
	}

	return c.userCol
}

//access to docCol
func (c *Database) Document() *DocCol {
	if c.docCol == nil {
		c.docCol = c.NewDocCol()
	}

	return c.docCol
}

# Book app - a GO project
This project is used to create backend APIs for users to follow other users's feeds, and the user can upload Post on to those feeds
## General API Usage
- For user: User can create an account with email and password, user can get their own information and can update or delete their account if needed
- For feeds: User can create a feed, get all the feed or a specific feed from the database, and user can update or delete their own feed
- Actions:
  - User can follow feeds, 1 user can follow many feeds but you can only follow that feed once
  - User can unfollowed a feed that they have followed.

## How to set up the project
- Clone the project:
```bash
git clone https://github.com/Darkred69/GO-Book-Project.git
```
- Install dependencies
```bash
go mod download
```
- Create env file:
```bash
PORT= {web's port}
DB_URL=postgres://postgres:{username}@{database_IP}:{database_port}/{databasename}?sslmode=disable
SECRET_KEY={create your own secret key}
EXPIRATION_MINUTES={add expiration minutes}
```
- Run goose command with terminal in sql/schema
```bash
goose postgres://postgres:{username}@{database_IP}:{database_port}/{databasename}?sslmode=disable
```
- Run the application
```bash
go build && GO-Book-Project.exe
```
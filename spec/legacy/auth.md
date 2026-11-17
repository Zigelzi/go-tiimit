# Authentication
For the web app to be useful, it needs to be deployed to a server that is accessible from public web. This allows the coach to access the app in the beginning of the practices, so that they can hand out the vests to the players and assign the players who haven’t marked their attendance if needed.

Some pages contain PII (name) that shouldn’t be accessible in public web. To prevent unauthorized access to PII these pages need authentication.

## Objective

Protect routes with authentication and deploy the app to server accessible from public web.

### Registering a new user

Credentials created via CLI to hash the passwords correctly.
- [x] Add command to create user to DB via CLI.

Register page requires that the users, practices and players belong to a club.

### Logging in

Login page, username / password.
- [x] Add login page
- [x] Return session after successful login
- [x] Limit access to specific routes to logged in users.

**Expected behaviour**
1. You can log in with username and password that exist in the database.
2. You need to log in again after 14 days.
3. You get feedback if you try to log in with incorrect credentials.
4. You get feedback if there's an error in the server while processing the login request.
5. You get redirected to index page after successfully logging in.

### Logging out
You can log out from the product, so that you can prevent others using your account.

**Expected behaviour**
1. You can logout.
2. You get redirected to index page after logging out from the device the session belongs to.
3. You can access any pages requiring authentication after logging out from the same device.

### Viewing protected routes
Only logged in users can view pages that require authorization, so that visitors can't access sensitive data.
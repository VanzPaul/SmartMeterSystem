The issue you're encountering is due to the fact that the `Logout` endpoint requires both a valid session token **and** a CSRF token for authentication. In your `curl` request, you provided the `session_token` cookie but did not include the `X-CSRF-Token` header. As a result, the server detected a mismatch between the expected CSRF token and the one provided (or lack thereof).

### Problem Analysis
1. **CSRF Token Missing**:
   - The `authorize` function in the `services` package checks for the presence of the `X-CSRF-Token` header and validates it against the stored CSRF token for the user.
   - Since the `X-CSRF-Token` header was not included in your `curl` request, the server logged a CSRF token mismatch error.

2. **Session Token Validity**:
   - The `session_token` cookie you provided (`ucGFf4MyeGjNdN5dWl7hGSA9OsiH6cNud8wKEoQ87TU=`) appears to be valid, as the server did not log a session token mismatch error.

3. **Error Response**:
   - The server returned `"Unauthorized"` with a `401 Unauthorized` status code because the CSRF token validation failed.

---

### Solution
To resolve this issue, you need to include the `X-CSRF-Token` header in your `curl` request. The CSRF token should match the one stored in the server's database for the user associated with the session token.

---

### Steps to Fix

#### 1. Retrieve the CSRF Token
After logging in, the server sets the `csrf_token` cookie. You can extract this value from the login response or inspect it using browser developer tools or the `-v` flag with `curl`.

For example:
```bash
curl -v -X POST http://localhost:8080/login \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "username=testuser&password=securepassword123"
```

Look for the `Set-Cookie` header in the response:
```
Set-Cookie: csrf_token=yzpI3I7G4LPXMWCytv7yFEpbcCAhVS7bhp8OxD1aIwk=; Expires=...; HttpOnly=false
```

In this case, the CSRF token is `yzpI3I7G4LPXMWCytv7yFEpbcCAhVS7bhp8OxD1aIwk=`.

---

#### 2. Include the CSRF Token in the Logout Request
Update your `curl` command to include the `X-CSRF-Token` header with the correct CSRF token value.

```bash
curl -X GET http://localhost:8080/logout \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -H "Cookie: session_token=ucGFf4MyeGjNdN5dWl7hGSA9OsiH6cNud8wKEoQ87TU=" \
     -H "X-CSRF-Token: yzpI3I7G4LPXMWCytv7yFEpbcCAhVS7bhp8OxD1aIwk=" \
     -d "username=testuser"
```

---

### Expected Response
If both the session token and CSRF token are valid, the server will respond with:
```
Logged Out Successfully!
```

Additionally, the server will clear the `session_token` and `csrf_token` cookies by setting their expiration dates to a past time.

---

### Debugging Tips
1. **Inspect Cookies**:
   - Use the `-v` flag with `curl` to view the `Set-Cookie` headers in the login response:
     ```bash
     curl -v -X POST http://localhost:8080/login \
          -H "Content-Type: application/x-www-form-urlencoded" \
          -d "username=testuser&password=securepassword123"
     ```

2. **Verify CSRF Token**:
   - Ensure that the `X-CSRF-Token` header matches the `csrf_token` cookie value.

3. **Server Logs**:
   - Check the server logs for detailed error messages. For example, if the CSRF token is invalid, the server will log:
     ```
     CSRF token mismatch: <provided_token> != <expected_token>
     ```

---

### Example Workflow

#### Step 1: Register a User
```bash
curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "username=testuser&password=securepassword123&accountNo=123456789"
```

#### Step 2: Login and Extract Tokens
```bash
curl -v -X POST http://localhost:8080/login \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "username=testuser&password=securepassword123"
```

Extract the `session_token` and `csrf_token` from the response headers.

#### Step 3: Access a Protected Route
```bash
curl -X POST http://localhost:8080/client/dashboard \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -H "Cookie: session_token=<SESSION_TOKEN>" \
     -H "X-CSRF-Token: <CSRF_TOKEN>" \
     -d "username=testuser"
```

#### Step 4: Logout
```bash
curl -X GET http://localhost:8080/logout \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -H "Cookie: session_token=<SESSION_TOKEN>" \
     -H "X-CSRF-Token: <CSRF_TOKEN>" \
     -d "username=testuser"
```

---

### Final Notes
- Always ensure that both the `session_token` cookie and `X-CSRF-Token` header are included in requests to protected routes and logout endpoints.
- If you encounter further issues, check the server logs for detailed error messages and verify that the tokens are being correctly extracted and used in your requests.
From the logs and your `curl` commands, it seems that the issue lies in how the `/client/dashboard` route is being accessed. Specifically:

1. **Incorrect HTTP Method**:
   - The `Protected` handler is registered as a `GET` route (`GET /dashboard`) in your `protectedMux`. However, you are sending a `POST` request to `/client/dashboard`.

2. **Route Mismatch**:
   - Since the `ServeMux` does not find a matching `POST` route for `/client/dashboard`, it defaults to returning an `Unauthorized` error (likely from middleware or default behavior).

---

### Solution:

You need to send a `GET` request to `/client/dashboard` instead of a `POST` request. Here's the corrected workflow:

---

### Corrected Workflow:

#### 1. **Register a User**:
```bash
curl -X POST "http://localhost:8080/register" \
     -d "username=testuser" \
     -d "password=securepassword123" \
     -d "accountNo=123456789"
```

**Expected Output:**
```
User Registered Successfully
```

---

#### 2. **Login**:
```bash
curl -X POST "http://localhost:8080/login" \
     -d "username=testuser" \
     -d "password=securepassword123" \
     -v
```

**Sample Response Headers (Success):**
```
Set-Cookie: session_token=<session_token_value>; HttpOnly; Expires=<expiry_time>
Set-Cookie: csrf_token=<csrf_token_value>; Expires=<expiry_time>
```

Save the `session_token` and `csrf_token` values from the response headers.

---

#### 3. **Access Protected Route**:
Use a `GET` request to access `/client/dashboard` with the correct query parameter, cookies, and headers:

```bash
curl -X GET "http://localhost:8080/client/dashboard?username=testuser" \
     -H "X-CSRF-Token: <csrf_token_value>" \
     -b "session_token=<session_token_value>"
```

Replace `<csrf_token_value>` and `<session_token_value>` with the values obtained from the login response.

**Example Command:**
```bash
curl -X GET "http://localhost:8080/client/dashboard?username=testuser" \
     -H "X-CSRF-Token: KlZ3dcPkl7bQsLE2U82lYrvdMiN_Oif31SMfMJn67fk=" \
     -b "session_token=Re_xlsB8IcPVMWrtUMfqL-YKAyQOc4X-yZOXq5waub8="
```

**Expected Output (Success):**
```
CSRF validation successful! Welcome, testuser
```

---

### Key Changes:

1. **Switch to `GET` Request**:
   - The `Protected` handler is registered as a `GET` route in your `protectedMux`. Sending a `POST` request will not match this route.

2. **Pass `username` as a Query Parameter**:
   - Use `?username=testuser` in the URL instead of passing it as form data (`-d`). This ensures that the `username` is correctly retrieved using `r.FormValue("username")`.

3. **Verify Middleware**:
   - Ensure that the `AuthMiddleware` is applied only to the `protectedMux` and does not block valid requests.

---

### Debugging Tips:

1. **Log Incoming Requests**:
   Add debug logs to the `authorize` function to confirm the values of `username`, `session_token`, and `csrf_token` during the request:
   ```go
   func authorize(r *http.Request) error {
       username := r.FormValue("username")
       fmt.Println("Username:", username)
       user, ok := users[username]
       if !ok {
           fmt.Println("User not found:", username)
           return ErrAuth
       }
       st, err := r.Cookie("session_token")
       if err != nil || st.Value != user.SessionToken {
           fmt.Println("Session token mismatch:", st.Value, "!=", user.SessionToken)
           return ErrAuth
       }
       csrf := r.Header.Get("X-CSRF-Token")
       if csrf != user.CSRFToken || csrf == "" {
           fmt.Println("CSRF token mismatch:", csrf, "!=", user.CSRFToken)
           return ErrAuth
       }
       return nil
   }
   ```

2. **Check Middleware Configuration**:
   Verify that the `ChainMiddleware` function applies the `AuthMiddleware` correctly to the `protectedMux`.

3. **Test with Verbose Output**:
   Use the `-v` flag with `curl` to inspect the request and response headers:
   ```bash
   curl -X GET "http://localhost:8080/client/dashboard?username=testuser" \
        -H "X-CSRF-Token: KlZ3dcPkl7bQsLE2U82lYrvdMiN_Oif31SMfMJn67fk=" \
        -b "session_token=Re_xlsB8IcPVMWrtUMfqL-YKAyQOc4X-yZOXq5waub8=" \
        -v
   ```

---

### Example Logs After Fix:

If everything is configured correctly, you should see logs like this on the server:

```
2025/02/07 13:23:59 Request received: GET /client/dashboard
Username: testuser
Session Token: Re_xlsB8IcPVMWrtUMfqL-YKAyQOc4X-yZOXq5waub8=
CSRF Token: KlZ3dcPkl7bQsLE2U82lYrvdMiN_Oif31SMfMJn67fk=
CSRF validation successful! Welcome, testuser
```

---

### Conclusion:

The main issue was sending a `POST` request to a `GET` route. By switching to a `GET` request and ensuring the `username` is passed as a query parameter, the `Protected` handler should work as expected. Always double-check the HTTP method and route configuration when debugging such issues.
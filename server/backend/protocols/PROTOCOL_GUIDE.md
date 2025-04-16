# 📡 Communication Protocol Guide

This document describes the communication protocol used to interact with the Flag Checker service.

---

### 🔁 Response Format

The Flag Checker must return a **JSON array** of objects with the following structure:

```json
{
  "msg": "[<flag>] <message>",
  "flag": "<flag>",
  "status": "ACCEPTED" | "DENIED" | "RESUBMIT" | "ERROR"
}
```

- `msg`: A human-readable message that includes the flag and additional context or result.
- `flag`: The original flag that was submitted.
- `status`: The result of the submission. It must be one of:
  - `"ACCEPTED"` – The flag is valid and was accepted.
  - `"DENIED"` – The flag is invalid or already submitted.
  - `"RESUBMIT"` – Temporary failure, the client should retry later.
  - `"ERROR"` – A generic error occurred during processing.

---

### 📤 Submit Function

To submit flags to the Flag Checker, the client must use the following function signature:

```go
func Submit(host string, token string, flags []string) ([]byte, error)
```

- `host`: The address of the Flag Checker (e.g., `http://localhost:8080`).
- `token`: Authentication token for the team.
- `flags`: A slice of strings containing flags to be submitted.

The function returns:
- The raw response as a byte slice (which should be a JSON array as described above).
- An error if the request fails or if the response cannot be processed.

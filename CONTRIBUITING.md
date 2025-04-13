# Contributor

## 👷 Development Workflow

1. Create a new branch from `develop` using the following naming convention:
   ```
   dev-{your_name}-{feature_name}
   ```
   _Example: `dev-john-login_page`_

2. Make your changes in this branch

3. Push your branch to the remote repository:
   ```bash
   git push origin dev-{your_name}-{feature_name}
   ```

4. Create a Pull Request (PR):
   - Go to the repository on GitHub
   - Click "New Pull Request"
   - Set base branch to `develop`
   - Set compare branch to your feature branch
   - Add a descriptive title and description
   - Submit the PR

5. Wait for review and approval

## 🗒️ Important Notes
- Never push directly to `develop` branch!!
- NEVER PUSH DIRECTLY TO `main` BRANCH!!
- Make sure your branch is up to date with `develop` before creating a PR
- Delete your branch after it has been merged

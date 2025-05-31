# Contributor

## 👷 Development Workflow

1. **Fork the repository**

   * Go to the [CookieFarm GitHub page](https://github.com/ByteTheCookies/CookieFarm)
   * Click the **"Fork"** button in the top-right corner
   * Clone the forked repository to your local machine:

     ```bash
     git clone https://github.com/your-username/your-forked-repo.git
     cd your-forked-repo
     ```

2. Create a new branch from `dev` using the following naming convention:

   ```
   dev-{your_name}-{feature_name}
   ```

   *Example: `dev-john-login_page`*

3. Make your changes in this branch

4. Push your branch to the remote repository:

   ```bash
   git push origin dev-{your_name}-{feature_name}
   ```

5. Create a Pull Request (PR):

   * Go to the repository on GitHub
   * Click **"New Pull Request"**
   * Set base branch to `dev`
   * Set compare branch to your feature branch
   * Add a descriptive title and description
   * Submit the PR

6. Wait for review and approval

## 🗒️ Important Notes

* Never push directly to `dev` branch!!
* NEVER PUSH DIRECTLY TO `main` BRANCH!!
* Test your code before pushing (test environment in `/tests`)
* Make sure your branch is up to date with `dev` before creating a PR
* Delete your branch after it has been merged

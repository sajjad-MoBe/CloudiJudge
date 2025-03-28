# CloudiJudge

This project aims to create an online judging system similar to platforms like Quora or Codeforces, where users can view problems, submit solutions (codes), and see the results of their submissions. This implementation is built using **Golang** and offers a simplified version of such systems.

## Features

### User Authentication
- Users can register and log in to the system.
- Passwords are securely hashed using the `bcrypt` package and stored in the database. The original passwords are never exposed.

### User Roles and Permissions
- Two types of users: regular users and admin.
- Admins have special permissions to publish questions and change user roles.
- Access control checks are performed on both the backend and the frontend.

### User Profile Page
- Each user has a profile displaying general information (username, submission status, etc.).
- Displays statistics on total questions attempted, success rate, and solved questions.
- Admins can modify user roles through the profile page.

### Question Listing
- A list of published questions displayed in reverse chronological order.
- Pagination support (e.g., 10 questions per page) using Query Parameters.

### Question Management
- Each question has an owner and starts as a draft until published by an admin.
- Each question includes:
  - Title
  - Statement
  - Time limit (in milliseconds)
  - Memory limit (in megabytes)
  - Test Input
  - Expected Output

### Submissions
- Users can submit Golang code for questions.
- Submission results are initially marked as "uncategorized" and processed by a judging service afterward.
- Possible outcomes include:
  - Correct (OK)
  - Compile Error
  - Wrong Answer
  - Memory Limit Exceeded
  - Time Limit Exceeded
  - Runtime Error

### Question Pages and Submissions
- Users can see the list of published questions and click on any question to view its full content.
- An "Submit Response" page allows uploading or pasting code.
- A "My Submissions" page displays the user's submission history with corresponding statuses.

### Question Creation
- Both admin and regular users can access a page to create new questions.
- Newly created questions are saved as drafts.
- Fields include title, statement, limits, etc.

### Admin Panel
- Admins can view and manage all questions, edit, or publish them.

### Internal APIs
- Internal APIs check submissions via a separate process called the "runner."
- Ensures that submissions are processed safely and efficiently without exposing the internal network.

### Commands
The compiled project contains the following commands:

- **serve**
  - Starts an HTTP Server based on a given configuration file.
  - `listen` address is provided as an argument (e.g., `--listen :8080`).
  - Configuration read via the `viper` package.

- **code-runner**
  - Compiles and executes submitted codes (currently only Golang).
  - Runs in a Docker container with restrictions (e.g. one CPU, no network access).

- **create-admin**
  - Takes a username and password to create an admin user or elevate an existing user to admin status.

## Getting Started

### Prerequisites
- Go (version 1.18 or later) installed.
- Docker (for code execution).

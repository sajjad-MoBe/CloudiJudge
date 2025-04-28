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
  - `listen` address is provided as an argument (e.g., `--listen=80`).
  - Configuration read via the `viper` package.

- **code-runner**
  - Compiles and executes submitted codes (currently only Golang).
  - Runs in a Docker container with restrictions (e.g. one CPU, no network access).

- **create-admin**
  - Takes a username and password to create an admin user or elevate an existing user to admin status.

## Getting Started

### Prerequisites
Before you begin, ensure you have the following installed:

- [Go](https://go.dev/) (version 1.23 or later)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### USAGE

Follow these steps to set up and run the **CloudiJudge** project:

  

1.  **Clone the Repository**

  

	Clone the project to your local machine:
	```bash

	git clone https://github.com/sajjad-MoBe/CloudiJudge

	cd CloudiJudge

	```

  

2.  **Set Up Environment Variables**

	Create a `.env` file by copying the provided `.env.example`:

	```bash

	cp .env.example .env

	```

	Edit the `.env` file to set the following variables:

	-  `PORT`: Application port

	-  `POSTGRES_HOST`: Your PostgreSQL service name(default is `db` if you changed it on docker-comppose.yml you should also change this variable)

	-  `POSTGRES_USER`: Your PostgreSQL username

	-  `POSTGRES_PASSWORD`: Your PostgreSQL password

	-  `POSTGRES_DB`: Your PostgreSQL database name

	-  `POSTGRES_DATA_FOLDER`: The folder for PostgreSQL data (e.g., `./postgres/data`)

	-  `PROBLEM_UPLOAD_FOLDER_SRC`: The folder for save problem related files(input, output, codes) on host (e.g., `./problems`)

	-  `PROBLEM_UPLOAD_FOLDER`: The folder for save problem related files on services (e.g., `/problems`)

	-  `MAX_CONCURRENT_RUNS`: Max concurrent code runs in a single code-runner service


3. **Build go code runner**
	build go code runner using this command:

	```bash

	docker build -t go-code-runner ./go-runner

	```

4.  **Run and Deploy**

	Start the application using Docker Compose (you can change scale of code runners):

	```bash

	docker-compose up --scale code-runner=3

	```

  
	- To run it in the background (detached mode), use:

		``docker-compose up --build --scale code-runner=3 -d``

	- To force a rebuild of the application, add the `--build` flag:

	  ``docker-compose up --scale code-runner=3 -d --build``

4.  **Verify the Setup**
	Once the containers are up, the application should be running and accessible as configured.
	- To create an Admin user you can use this command:
		``docker-compose exec judge create-admin --email=sajjad@beigi.com``
	
	- You also can fill database with test data using (erase=true for delete test datas):
		``docker-compose exec judge load-test-data erase=false``


## Contributors

- Sajjad: Implemented the backend part of the website.
  - Email: [sajjad.mohammadbeigi@gmail.com](mailto:sajjad.mohammadbeigi@gmail.com)
  - GitHub: [@sajjad-MoBe](https://github.com/sajjad-MoBe)

- Mohammad: Implemented the frontend part of the website.
  - Email: [mohammadmohammadbeigi1381@gmail.com](mailto:mohammadmohammadbeigi1381@gmail.com)
  - GitHub: [@mbmohammad](https://github.com/mbmohammad)
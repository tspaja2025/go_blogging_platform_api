# Go Blogging Platform API

A simple blogging platform API built with Go.

## Roadmap.sh beginner project
This project was created as a part of Blogging Platform API beginner project.
Check out the project details [roadmap.sh](https://roadmap.sh/projects/blogging-platform-api)

## Features

* Create a new blog post
* Update an existing blog post
* Delete an existing blog post
* Get a single blog post
* Get all blog posts
* Filter blog posts by a search term

---

## Installation

Clone the repository:

```bash
git clone https://github.com/tspaja2025/go_blogging_platform_api.git
cd go_blogging_platform_api
```

Run the application:

```bash
go run main.go
```

---

## Usage

### Create blog post

```bash
POST /posts
{
  "title": "My First Blog Post",
  "content": "This is the content of my first blog post.",
  "category": "Technology",
  "tags": ["Tech", "Programming"]
}
```

### Update blog post

```bash
PUT /posts/1
{
  "title": "My Updated Blog Post",
  "content": "This is the updated content of my first blog post.",
  "category": "Technology",
  "tags": ["Tech", "Programming"]
}
```

### Delete blog post

```bash
DELETE /posts/1
```

### Get blog post

```bash
GET /posts/1
```

### Get all posts

```bash
GET /posts
```

### Filter

```bash
GET /posts?term=tech
```

---

## Example Output

```text
-
```

---

## Data Storage

-

Example:

```json
-
```

---

## Technologies Used

* Go
* JSON file storage
* Standard library packages:

  * `-`

## Learning Goals

This project was built to practice:

* Creation RESTful APis including best practices and conventions
* Common HTTP methods GET, POST, PUT, PATCH, DELETE
* Status codes and error handling in APIs
* CRUD operations using an API

---

## License

This project is open source and available under the MIT License.

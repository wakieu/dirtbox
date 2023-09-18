# drtbox

This project is a clone of [dontpad](https://dontpad.com/) made in Go for practice purposes.

## Implementation

When run, the application will expose two endpoints:
### Page server @ *localhost:3030*
Responsible for serving the browser pages. This endpoint requires no parameters, instead it uses the request path string to fetch the content from the server and renders it via HTML templates.

### Content server @ *localhost:3131*
Responsible for serving content and interacting with the database. Basically this is a CRUD API, but it uses the client's request path as the *primary key* inside the database.

The database is a `.sqlite` file in the root directory, created automatically, if non-existent.

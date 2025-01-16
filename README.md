# Animorank API

Animorank API is an experimental code execution service designed to facilitate the remote execution of C code. This is build for the Animorank project.

## Usage
### Endpoint

**POST** `/run`
### Request Body
The Request Body should be a JSON object with the following fields.
```
{
  "language": "c", 
  "code": "/* Your C code string goes here */",
  "stdin": ["/* Array of inputs to standard input */"]
}
```
### Response

**Example Success Response**
```
{
  "results": [
    {
      "status": "successfully compiled and run",
      "stdout": "/* Standard output from the executed program */"
    }
  ]
}
```

**Example Error Response**
```
{
  "results": [
    {
      "status": "successfully compiled and run",
      "stdout": "time limit exceeded"
    }
  ]
}
```

## Notes
 - Currently, **C** is the only supported programming language.
 - The stdin field accepts an array of inputs, which are executed separately from each other, simulating a single test case.
 - Additional support for other languages may be added in the future.
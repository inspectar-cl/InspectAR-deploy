# FastAPI Microservice

This is a simple FastAPI microservice project that includes a basic "Hello, World!" endpoint.

## Project Structure

```
fastapi-microservice
├── api
│   ├── controllers
│   │   └── hello_controller.py
│   ├── middleware
│   └── router
│       └── hello_router.py
├── cmd
│   └── main.py
├── config
│   └── settings.py
├── internal
│   ├── dto
│   ├── handlers
│   ├── models
│   │   └── __init__.py
│   ├── repository
│   └── services
├── requirements.txt
└── README.md
```

## Setup Instructions

1. Clone the repository:
   ```
   git clone <repository-url>
   cd fastapi-microservice
   ```

2. Install the required dependencies:
   ```
   pip install -r requirements.txt
   ```

3. Run the application:
   ```
   uvicorn cmd.main:app --reload
   ```

## Usage

Once the application is running, you can access the "Hello, World!" endpoint by navigating to:

```
http://127.0.0.1:8000/hello
```

You should see a response with the message:

```
{"message": "Hello, World!"}
```

## License

This project is licensed under the MIT License.
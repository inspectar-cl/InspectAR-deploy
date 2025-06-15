from fastapi import APIRouter
from api.controllers.hello_controller import get_hello

hello_router = APIRouter()

@hello_router.get("/hello")
def hello():
    return get_hello()
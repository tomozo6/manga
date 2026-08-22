from pydantic import BaseModel


class MangaSummary(BaseModel):
    id: str
    title: str
    author: str

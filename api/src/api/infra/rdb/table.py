from sqlmodel import SQLModel

from api.modules.manga.models import Manga, Volume  # noqa: F401

from .session import engine


async def create_tables() -> None:
    async with engine.begin() as connection:
        await connection.run_sync(SQLModel.metadata.create_all)

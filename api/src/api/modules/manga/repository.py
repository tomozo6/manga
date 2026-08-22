from sqlmodel import select
from sqlmodel.ext.asyncio.session import AsyncSession

from .models import Manga


class MangaRepository:
    def __init__(self, session: AsyncSession) -> None:
        self.session = session

    async def list_all(self) -> list[Manga]:
        # fmt: off
        statement = (
            select(Manga)
            .order_by(Manga.title, Manga.id)
        )
        # fmt: on
        result = await self.session.exec(statement)
        return list(result)

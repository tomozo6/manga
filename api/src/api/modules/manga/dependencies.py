from typing import Annotated

from fastapi import Depends

from api.infra.rdb.session import SessionDep

from .repository import MangaRepository


def get_manga_repository(session: SessionDep):
    return MangaRepository(session)


RepoDep = Annotated[MangaRepository, Depends(get_manga_repository)]

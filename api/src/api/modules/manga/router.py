from fastapi import APIRouter

from api.core.auth.firebase import UserDep

from .dependencies import RepoDep
from .schemas import MangaSummary

router = APIRouter(prefix="/api/manga", tags=["manga"])


@router.get("")
async def list_manga(_user: UserDep, repo: RepoDep) -> list[MangaSummary]:
    mangas = await repo.list_all()

    manga_summaries: list[MangaSummary] = []

    for manga in mangas:
        manga_summaries.append(MangaSummary(**manga.model_dump()))

    return manga_summaries

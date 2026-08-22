from fastapi import APIRouter

from api.core.auth.firebase import UserDep

from .schemas import MangaSummary

router = APIRouter(prefix="/api/manga", tags=["manga"])


@router.get("")
async def list_manga(_user: UserDep) -> list[MangaSummary]:
    return [
        MangaSummary(
            id="historie",
            title="ヒストリエ",
            author="岩明均",
        )
    ]

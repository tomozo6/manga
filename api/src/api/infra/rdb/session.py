from collections.abc import AsyncIterator
from pathlib import Path
from typing import Annotated

from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncEngine, async_sessionmaker, create_async_engine
from sqlmodel.ext.asyncio.session import AsyncSession

# ローカル開発用の SQLite データベース。環境ごとの接続先設定は後でここから分離する。
api_directory = Path(__file__).resolve().parents[4]
database_path = api_directory / "manga.db"
database_url = f"sqlite+aiosqlite:///{database_path}"

# aiosqlite ドライバを使う非同期 SQLAlchemy エンジン。
engine: AsyncEngine = create_async_engine(database_url)

# 各リクエストで利用する非同期セッションのファクトリ。
# commit 後も取得済みモデルの属性を参照できるよう expire_on_commit=False にする。
session_factory = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False,
)


async def get_session() -> AsyncIterator[AsyncSession]:
    """リクエスト処理用の非同期データベースセッションを提供する。"""

    async with session_factory() as session:
        yield session


# エンドポイントの引数に指定する非同期データベースセッションの依存性。
SessionDep = Annotated[AsyncSession, Depends(get_session)]

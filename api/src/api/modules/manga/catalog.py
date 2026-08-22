from pathlib import Path

import yaml
from pydantic import BaseModel, Field
from sqlalchemy import delete

from api.infra.rdb.session import session_factory

from .models import Manga, Volume

api_directory = Path(__file__).resolve().parents[4]
catalog_directory = api_directory / "catalog" / "mangas"


class MangaCatalog(BaseModel):
    id: str
    title: str
    author: str
    cover_object_key: str = ""
    volumes: list[VolumeCatalog]


class VolumeCatalog(BaseModel):
    id: str
    number: int = Field(ge=1)
    title: str
    cover_object_key: str = ""
    page_count: int = Field(ge=1)
    page_extension: str


def load_mangas_from_yaml(catalog_dir: Path = catalog_directory) -> list[MangaCatalog]:
    # 読み込んだ全作品を、このリストに順番に追加する。
    mangas: list[MangaCatalog] = []

    # catalog_dir 内の .yaml ファイルをファイル名順に取得する。
    # sort しておくと、毎回同じ順番でカタログを読み込める。
    for path in sorted(catalog_dir.glob("*.yaml")):
        # YAML ファイルを UTF-8 の文字列として読み込む。
        contents = path.read_text(encoding="utf-8")

        # YAML の文字列を Python の辞書・リストへ変換する。
        data = yaml.safe_load(contents)

        # 辞書が MangaCatalog の形式を満たすか検証し、作品としてリストへ追加する。
        # 必須項目の不足や page_count が 1 未満の場合は、ここでエラーになる。
        mangas.append(MangaCatalog.model_validate(data))

    # 読み込んで検証済みの全作品を呼び出し元へ返す。
    return mangas


async def import_catalog() -> None:
    mangas = load_mangas_from_yaml()

    async with session_factory.begin() as session:
        # 外部キーがあるため子テーブルから消す。
        await session.exec(delete(Volume))
        await session.exec(delete(Manga))

        for manga in mangas:
            session.add(Manga(**manga.model_dump(exclude={"volumes"})))

        # Manga を先に INSERT して、Volume の外部キー参照を成立させる。
        await session.flush()

        for manga in mangas:
            for volume in manga.volumes:
                session.add(
                    Volume(
                        manga_id=manga.id,
                        **volume.model_dump(),
                    )
                )

from sqlmodel import Field, SQLModel


class Manga(SQLModel, table=True):
    """作品カタログの一作品を表すテーブル。"""

    __tablename__ = "manga"

    id: str = Field(primary_key=True)
    title: str
    author_name: str
    cover_object_key: str = ""


class Volume(SQLModel, table=True):
    """作品に属する一巻と、そのページ画像を復元するための情報を表すテーブル。"""

    __tablename__ = "volume"

    manga_id: str = Field(foreign_key="manga.id", primary_key=True)
    id: str = Field(primary_key=True)
    number: int = Field(ge=1)
    title: str
    cover_object_key: str = ""
    page_count: int = Field(ge=1)
    page_extension: str

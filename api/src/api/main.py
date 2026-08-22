from contextlib import asynccontextmanager
from pathlib import Path

from fastapi import FastAPI
from fastapi.responses import FileResponse
from fastapi.staticfiles import StaticFiles
from rich import panel, print
from scalar_fastapi import get_scalar_api_reference

from api.core.auth.firebase import init_firebase
from api.infra.rdb.table import create_tables
from api.modules.manga.catalog import import_catalog

from .modules.manga.router import router as manga_router

static_directory = Path(__file__).parent / "static"


@asynccontextmanager
async def lifespan_handler(_: FastAPI):
    print(panel.Panel("Server started...", border_style="green"))
    init_firebase()
    await create_tables()
    await import_catalog()
    yield
    print(panel.Panel("...Server stopped", border_style="red"))


app = FastAPI(lifespan=lifespan_handler)

app.include_router(manga_router)


# @app.get("/")
# def index() -> FileResponse:
#     return FileResponse(static_directory / "index.html")


@app.get("/library", include_in_schema=False)
def library() -> FileResponse:
    return FileResponse(static_directory / "library.html")


@app.get("/manga/{manga_id}", include_in_schema=False)
def manga(manga_id: str) -> FileResponse:
    return FileResponse(static_directory / "manga.html")


@app.get("/manga/{manga_id}/volumes/{volume_id}", include_in_schema=False)
def reader(manga_id: str, volume_id: str) -> FileResponse:
    return FileResponse(static_directory / "reader.html")


@app.get("/scalar", include_in_schema=False)
def get_scalar_docs():
    return get_scalar_api_reference(
        openapi_url=app.openapi_url,
        title="Scalar API",
        persist_auth=True,
    )


# APIルートを優先させたうえで、静的なフロントエンドを配信する。
app.mount("/", StaticFiles(directory=static_directory, html=True), name="static")

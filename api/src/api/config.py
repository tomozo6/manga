from pathlib import Path

from pydantic import EmailStr
from pydantic_settings import BaseSettings, SettingsConfigDict


class AuthSettings(BaseSettings):
    firebase_project_id: str
    allowed_emails: list[EmailStr]

    model_config = SettingsConfigDict(
        env_file=Path(__file__).parent / ".env",
        extra="ignore",
    )


auth_settings = AuthSettings()  # type: ignore[call-arg]

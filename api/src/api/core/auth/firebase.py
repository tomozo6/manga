from typing import Annotated

import firebase_admin  # type: ignore[import-untyped]
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from firebase_admin import auth as firebase_auth
from pydantic import BaseModel, EmailStr

from api.config import auth_settings


class CurrentUser(BaseModel):
    """Firebase ID token の検証後に API へ渡す利用者情報。"""

    uid: str
    name: str = ""  # Firebase ID token には name が存在する保証がないため
    email: EmailStr


# Authorization ヘッダーから Bearer トークンを取得する FastAPI 依存性。
bearer_scheme = HTTPBearer(auto_error=False)


def init_firebase() -> None:
    """Firebase Admin SDK のデフォルトアプリを一度だけ初期化する。"""

    if not firebase_admin._apps:
        firebase_admin.initialize_app(options={"projectId": auth_settings.firebase_project_id})


def get_current_user(credentials: HTTPAuthorizationCredentials | None = Depends(bearer_scheme)) -> CurrentUser:
    """Firebase ID token を検証し、許可リストにある利用者を返す。"""

    if credentials is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="authentication is required",
        )

    try:
        token = firebase_auth.verify_id_token(credentials.credentials)
        user = CurrentUser(
            uid=token["uid"],
            name=token.get("name", ""),
            email=token["email"],
        )
    # InvalidIdTokenError は期限切れ・失効した ID token も含む。
    except firebase_auth.InvalidIdTokenError, KeyError, ValueError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="invalid Firebase ID token",
        ) from None

    # メールアドレスの大文字・小文字の表記ゆれを無視して許可リストと照合する。
    if user.email.lower() not in {email.lower() for email in auth_settings.allowed_emails}:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="this account is not allowed",
        )

    return user


# 保護された API の引数に指定する、認証済み利用者の依存性。
UserDep = Annotated[CurrentUser, Depends(get_current_user)]

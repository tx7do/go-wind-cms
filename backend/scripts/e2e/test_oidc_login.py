"""OIDC 授权码登录 E2E：内置 mock IdP，全链路验证。

链路：GetOidcAuthorizeUrl 取授权地址（state 落 Redis）→ 以
grant_type=authorization_code 调 Login → core 向 mock IdP 换 token →
取 userinfo → 自动建号（首登）/凭证映射（复登）→ 签发本仓 JWT。

前置：
  - core/admin 已以 OIDC_* 环境变量启动（指向本脚本同款 mock IdP 端口 6900）：
      OIDC_ISSUER=http://127.0.0.1:6900
      OIDC_CLIENT_ID=gowind-e2e
      OIDC_CLIENT_SECRET=e2e-secret
      OIDC_AUTHORIZE_URL=http://127.0.0.1:6900/authorize
      OIDC_TOKEN_URL=http://127.0.0.1:6900/token
      OIDC_USERINFO_URL=http://127.0.0.1:6900/userinfo
  - mock IdP 由本脚本自行启动（127.0.0.1:6900），无需单独进程。

用法：python test_oidc_login.py
"""
import base64
import json
import re
import threading
import time
import urllib.error
import urllib.parse
import urllib.request
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

from common import ADMIN_BASE, admin_login

IDP_PORT = 6900
IDP_CLIENT_ID = "gowind-e2e"
IDP_CLIENT_SECRET = "e2e-secret"
IDP_ACCESS_TOKEN = "idp-access-token-e2e"

# mock IdP 的用户身份（userinfo 返回）；sub 每次运行刷新，
# 避免上一轮已建号的 (issuer, sub) 命中映射路径、跳过自动建号分支。
IDP_SUB = "e2e-sub-42"
IDP_EMAIL = "sso.user@e2e.example.com"
IDP_USERNAME = "sso.e2e.user"

passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


class MockIdpHandler(BaseHTTPRequestHandler):
    """最小 OIDC IdP：POST /token（授权码换令牌）+ GET /userinfo。"""

    def _json(self, code, payload):
        body = json.dumps(payload).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_POST(self):
        if self.path != "/token":
            self._json(404, {"error": "not found"})
            return
        form = urllib.parse.parse_qs(
            self.rfile.read(int(self.headers.get("Content-Length", 0))).decode()
        )
        if form.get("grant_type", [""])[0] != "authorization_code":
            self._json(400, {"error": "unsupported_grant_type"})
            return
        if (
            form.get("client_id", [""])[0] != IDP_CLIENT_ID
            or form.get("client_secret", [""])[0] != IDP_CLIENT_SECRET
        ):
            self._json(401, {"error": "invalid_client"})
            return
        if form.get("code", [""])[0] != "e2e-auth-code":
            self._json(400, {"error": "invalid_grant"})
            return
        self._json(200, {"access_token": IDP_ACCESS_TOKEN, "token_type": "Bearer"})

    def do_GET(self):
        if self.path != "/userinfo":
            self._json(404, {"error": "not found"})
            return
        if self.headers.get("Authorization") != f"Bearer {IDP_ACCESS_TOKEN}":
            self._json(401, {"error": "invalid_token"})
            return
        self._json(
            200,
            {
                "sub": IDP_SUB,
                "email": IDP_EMAIL,
                "preferred_username": IDP_USERNAME,
                "name": "SSO E2E User",
            },
        )

    def log_message(self, *args):
        pass


def start_idp():
    server = ThreadingHTTPServer(("127.0.0.1", IDP_PORT), MockIdpHandler)
    threading.Thread(target=server.serve_forever, daemon=True).start()
    return server


def raw(method, path, payload=None, headers=None):
    """裸请求（需要读错误响应体/非 2xx 场景）。返回 (status, body dict)。"""
    req = urllib.request.Request(ADMIN_BASE + path, method=method)
    req.add_header("Content-Type", "application/json")
    for k, v in (headers or {}).items():
        req.add_header(k, v)
    data = json.dumps(payload).encode() if payload is not None else None
    try:
        with urllib.request.urlopen(req, data=data, timeout=15) as resp:
            return resp.status, json.loads(resp.read().decode() or "{}")
    except urllib.error.HTTPError as e:
        try:
            return e.code, json.loads(e.read().decode() or "{}")
        except Exception:
            return e.code, {}


def jwt_uid(access_token):
    """解 JWT payload 的 uid（不含验签，仅测试断言用）。"""
    try:
        payload = access_token.split(".")[1]
        payload += "=" * (-len(payload) % 4)
        return json.loads(base64.urlsafe_b64decode(payload)).get("uid")
    except Exception:
        return None


def new_state():
    """取一次合法授权地址，返回 (status, state)。"""
    status, body = raw(
        "POST",
        "/admin/v1/oidc/authorize-url",
        {"redirect_uri": "http://localhost:15999/auth/login"},
    )
    url = body.get("authorizeUrl") or ""
    state = (urllib.parse.parse_qs(urllib.parse.urlparse(url).query).get("state") or [""])[0]
    return status, state


def oidc_login(state, code="e2e-auth-code"):
    return raw(
        "POST",
        "/admin/v1/login",
        {
            "grant_type": "authorization_code",
            "code": code,
            "oidc_state": state,
            "redirect_uri": "http://localhost:15999/auth/login",
        },
    )


def main():
    global IDP_SUB, IDP_USERNAME, IDP_EMAIL
    IDP_SUB = f"e2e-sub-{int(time.time())}"
    IDP_USERNAME = f"sso.e2e.user{IDP_SUB.split('-')[-1]}"
    IDP_EMAIL = f"sso.user{IDP_SUB.split('-')[-1]}@e2e.example.com"
    idp = start_idp()
    print("🧪 mock IdP 已启动 :6900")

    # ===== 1. authorize-url：未带 redirect_uri 拒绝 =====
    status, _ = raw("POST", "/admin/v1/oidc/authorize-url", {})
    check("authorize-url 缺 redirect_uri 被拒", status == 400, f"status={status}")

    # ===== 2. authorize-url：返回含 state 的完整地址 =====
    status, body = raw(
        "POST",
        "/admin/v1/oidc/authorize-url",
        {"redirect_uri": "http://localhost:15999/auth/login"},
    )
    check("authorize-url 200", status == 200, f"status={status}")
    url = body.get("authorizeUrl") or ""
    check("authorizeUrl 非空", bool(url))
    qs = urllib.parse.parse_qs(urllib.parse.urlparse(url).query)
    check("含 client_id", qs.get("client_id") == [IDP_CLIENT_ID], f"{qs.get('client_id')}")
    check("含 response_type=code", qs.get("response_type") == ["code"])
    state = (qs.get("state") or [""])[0]
    check("state 为 64 位 hex", bool(re.fullmatch(r"[0-9a-f]{64}", state)), f"{state[:16]}…")

    # ===== 3. 伪造 state 拒绝（未签发）=====
    status, body = oidc_login("f" * 64)
    check("伪造 state 登录被拒 400", status == 400, f"status={status}")

    # ===== 4. 正常链路：state 换真登录态（首登自动建号）=====
    status, body = oidc_login(state)
    check("首登 200", status == 200, f"status={status} body={json.dumps(body, ensure_ascii=False)[:120]}")
    token1 = body.get("access_token") or ""
    check("签发 access_token", bool(token1))
    uid1 = jwt_uid(token1)
    check("JWT 载荷含非零 uid", bool(uid1), f"uid={uid1}")

    # ===== 5. state 一次性：复用同一 state 二次登录被拒 =====
    status, _ = oidc_login(state)
    check("state 复用被拒 400", status == 400, f"status={status}")

    # ===== 6. 复登映射：新 state + 同一 sub → 同一用户（不重复建号）=====
    status, state2 = new_state()
    status, body = oidc_login(state2)
    check("复登 200", status == 200, f"status={status}")
    uid2 = jwt_uid(body.get("access_token") or "")
    check("复登映射同一用户", uid1 is not None and uid1 == uid2, f"uid1={uid1} uid2={uid2}")

    # ===== 7. IdP 拒绝错误 code → 登录 400 =====
    status, state3 = new_state()
    status, _ = oidc_login(state3, code="wrong-code")
    check("IdP 拒绝错误 code → 登录 400", status == 400, f"status={status}")

    # ===== 8. 密码登录回归：OIDC 接入不影响既有链路 =====
    try:
        admin_login("admin", "admin")
        check("密码登录回归", True)
    except Exception as e:  # noqa: BLE001 —— 断言脚本，失败要打印而非吞
        check("密码登录回归", False, str(e)[:60])

    idp.shutdown()
    print(f"\n通过 {passed} / 失败 {failed}")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()

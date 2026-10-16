"""E2E 公共助手：登录（验证码+AES 密码）、admin API 调用、curl 透传。

前置条件：
- admin(:6600)/core/app(:6700) 服务已启动且清代理（HTTP_PROXY='' NO_PROXY='*'）
- redis 容器名 redis-server（密码 *Abcd123456），验证码存 gowind-cms:captcha:{id}
- python 依赖：pycryptodome
- 本机 curl 可用（app 侧请求走 curl --noproxy '*'，规避代理劫持）

用法：
    from e2e_common import admin_login, AdminApi
    hdr = admin_login('admin', 'admin')          # -> Authorization header
    api = AdminApi(hdr)
    r = api('/admin/v1/posts?pageSize=10')
"""
import base64
import json
import subprocess
import urllib.error
import urllib.request

from Crypto.Cipher import AES
from Crypto.Util.Padding import pad

AES_KEY = b"f51d66a73d8a0927"  # DefaultAESKey（前端同款）
ADMIN_BASE = "http://localhost:6600"
APP_BASE = "http://localhost:6700"


def _admin_api_raw(path, method="GET", body=None, headers=None):
    h = {"Content-Type": "application/json"}
    h.update(headers or {})
    r = urllib.request.Request(
        ADMIN_BASE + path,
        data=json.dumps(body).encode() if body is not None else None,
        method=method,
        headers=h,
    )
    try:
        resp = urllib.request.urlopen(r)
        return resp.read().decode(), resp.headers
    except urllib.error.HTTPError as e:
        return e.read().decode(), e.headers


def admin_login(username="admin", password="admin"):
    """验证码登录，返回 {'Authorization': 'Bearer ...'}；失败抛 AssertionError。"""
    raw, _ = _admin_api_raw("/admin/v1/captcha?type=login")
    d = json.loads(raw)
    ans = subprocess.check_output([
        "docker", "exec", "redis-server", "redis-cli",
        "-a", "*Abcd123456", "--no-auth-warning", "GET",
        "gowind-cms:captcha:" + d["captchaId"],
    ]).decode().strip()
    enc = base64.b64encode(
        AES.new(AES_KEY, AES.MODE_CBC, AES_KEY).encrypt(pad(password.encode(), 16))
    ).decode()
    raw, _ = _admin_api_raw(
        "/admin/v1/login", "POST",
        {"grantType": "password", "username": username, "password": enc},
        {"X-Captcha-Id": d["captchaId"], "X-Captcha-Value": ans},
    )
    login = json.loads(raw)
    token = login.get("access_token")
    assert token, f"login failed: {raw[:200]}"
    return {"Authorization": "Bearer " + token}


class AdminApi:
    """admin REST 薄封装：api(path, method, body) -> dict（HTTP 错误也返回 JSON）。"""

    def __init__(self, headers):
        self.headers = headers

    def __call__(self, path, method="GET", body=None):
        raw, _ = _admin_api_raw(path, method, body, self.headers)
        try:
            return json.loads(raw)
        except json.JSONDecodeError:
            return {"_raw": raw}


def app_curl(path, method="GET", body=None, cookie=None, user_agent=None):
    """app(:6700) 请求走 curl --noproxy '*'，返回 (stdout, set_cookie_header)。"""
    cmd = ["curl", "-s", "--noproxy", "*", "--max-time", "15", "--retry", "2", "--retry-connrefused", "-i", "-X", method]
    if cookie:
        cmd += ["-H", f"Cookie: {cookie}"]
    if user_agent:
        cmd += ["-H", f"User-Agent: {user_agent}"]
    if body is not None:
        cmd += ["-H", "Content-Type: application/json", "-d", json.dumps(body)]
    cmd.append(APP_BASE + path)
    r = subprocess.run(cmd, capture_output=True, text=True, timeout=30)
    if r.returncode != 0 or not r.stdout:
        # 透出 curl 层失败（超时/连接拒绝等），避免静默空响应
        print(f"[app_curl] curl rc={r.returncode} stderr={r.stderr[:200]!r}")
    out = r.stdout
    sep = "\r\n\r\n" if "\r\n\r\n" in out else "\n\n"
    header_block, _, payload = out.partition(sep)
    set_cookie = None
    for line in header_block.splitlines():
        if line.lower().startswith("set-cookie:"):
            set_cookie = line.split(":", 1)[1].strip()
    return payload, set_cookie

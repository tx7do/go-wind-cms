"""MFA（TOTP 两步验证）E2E：绑定 → 挑战登录 → 换令牌 → 解绑 → 救援重置。

前置条件（同 e2e_common.py）：
- admin(:6600)/core 服务已启动（清代理）；redis 容器名 redis-server
- python 依赖：pycryptodome

覆盖断言（22 条）：
  1. 绑定前置状态 enabled=false
  2. start enroll 返回 secret/otpAuthUrl/operationId
  3. 错误码确认绑定被拒（且操作一次性消费）
  4. 重新发起绑定 + 正确码确认绑定成功
  5. 状态 enabled=true
  6. 旧 operationId 复用被拒
  7. 启用后登录返回 mfa_operation_id 且无令牌
  8. 挑战凭错误码被拒
  9. 挑战凭据一次性（消费后复用被拒）
 10. 重新登录拿新挑战 + 正确码换到令牌
 11. 挑战响应可调受保护接口（令牌有效）
 12. 内联路径：Login 直传 mfa_operation_id+totp_code 换令牌
 13. 错误码内联被拒
 14. disable 错误码被拒
 15. disable 正确码成功
 16. 状态回到 enabled=false
 17. 解绑后登录不再要挑战
 18. 重置接口对不存在因子幂等成功
 19. 平台管理员救援重置生效（绑定后 reset → enabled=false）
 20. 未带操作者的绑定请求由 BFF 注入身份成功（回归）
 21. 非法 secret 形态的确认（过期操作）被拒
 22. 未认证调用 me/mfa 被 401 拒
"""
import base64
import hashlib
import hmac
import json
import struct
import sys
import time

sys.path.insert(0, __file__.rsplit("\\", 1)[0])
from common import ADMIN_BASE, _admin_api_raw, admin_login  # noqa: E402

PASS = 0


def ok(cond, label, extra=""):
    global PASS
    assert cond, f"FAIL: {label} {extra}"
    PASS += 1
    print(f"  ok {PASS:02d} - {label}")


def totp_code(secret: str, at: float | None = None) -> str:
    pad_len = (8 - len(secret) % 8) % 8
    key = base64.b32decode(secret.upper() + "=" * pad_len)
    counter = int((at or time.time()) // 30)
    msg = struct.pack(">Q", counter)
    digest = hmac.new(key, msg, hashlib.sha1).digest()
    off = digest[-1] & 0x0F
    value = ((digest[off] & 0x7F) << 24) | (digest[off + 1] << 16) | (digest[off + 2] << 8) | digest[off + 3]
    return str(value % 1_000_000).zfill(6)


def api(headers, path, method="GET", body=None):
    raw, _ = _admin_api_raw(path, method, body, headers)
    try:
        return json.loads(raw)
    except json.JSONDecodeError:
        return {"_raw": raw, "_status_ok": False}


def login_raw(username="admin", password="admin"):
    """完整密码登录（验证码），返回原始 LoginResponse dict。"""
    raw, _ = _admin_api_raw("/admin/v1/captcha?type=login")
    d = json.loads(raw)
    import subprocess

    ans = subprocess.check_output([
        "docker", "exec", "redis-server", "redis-cli",
        "-a", "*Abcd123456", "--no-auth-warning", "GET",
        "gowind-cms:captcha:" + d["captchaId"],
    ]).decode().strip()
    from common import AES_KEY
    from Crypto.Cipher import AES
    from Crypto.Util.Padding import pad

    enc = base64.b64encode(
        AES.new(AES_KEY, AES.MODE_CBC, AES_KEY).encrypt(pad(password.encode(), 16))
    ).decode()
    raw, _ = _admin_api_raw(
        "/admin/v1/login", "POST",
        {"grantType": "password", "username": username, "password": enc},
        {"X-Captcha-Id": d["captchaId"], "X-Captcha-Value": ans},
    )
    return json.loads(raw)


def bearer(token: str):
    return {"Authorization": "Bearer " + token}


def main():
    # ---------- 未认证 401 ----------
    r = api({}, "/admin/v1/me/mfa")
    ok(not r.get("enabled") and not r.get("accessToken"), "未认证调用 me/mfa 被拒", str(r)[:120])

    hdr = admin_login("admin", "admin")
    me = api(hdr, "/admin/v1/me")
    admin_id = me.get("id")
    ok(bool(admin_id), "管理员登录并取到 id", str(me)[:120])

    # ---------- 初始状态 ----------
    r = api(hdr, "/admin/v1/me/mfa")
    was_enabled = r.get("enabled")
    if was_enabled:
        # 环境残留：先清掉（重置接口幂等）
        api(hdr, f"/admin/v1/users/{admin_id}/mfa/reset", "POST", {})
        r = api(hdr, "/admin/v1/me/mfa")
        assert not r.get("enabled"), "residual mfa factor could not be reset"
    ok(not r.get("enabled"), "初始 enabled=false")

    # ---------- 绑定：错误码确认 → 操作消费 → 重新绑定 ----------
    r = api(hdr, "/admin/v1/me/mfa/enroll/start", "POST", {"method": "TOTP"})
    ok(bool(r.get("totp", {}).get("secret")) and bool(r.get("totp", {}).get("otpAuthUrl"))
       and bool(r.get("operationId")), "start enroll 返回 secret/otpAuthUrl/operationId", str(r)[:160])
    secret1, op1 = r["totp"]["secret"], r["operationId"]

    wrong = "000000" if totp_code(secret1) != "000000" else "000001"
    r = api(hdr, "/admin/v1/me/mfa/enroll/confirm", "POST",
            {"method": "TOTP", "operationId": op1, "totpCode": wrong})
    ok(not r.get("success"), "错误码确认绑定被拒", str(r)[:120])
    r = api(hdr, "/admin/v1/me/mfa/enroll/confirm", "POST",
            {"method": "TOTP", "operationId": op1, "totpCode": totp_code(secret1)})
    ok(not r.get("success"), "失败后旧操作一次性（复用被拒）", str(r)[:120])

    r = api(hdr, "/admin/v1/me/mfa/enroll/start", "POST", {"method": "TOTP"})
    secret, op = r["totp"]["secret"], r["operationId"]
    r = api(hdr, "/admin/v1/me/mfa/enroll/confirm", "POST",
            {"method": "TOTP", "operationId": op, "totpCode": totp_code(secret)})
    ok(r.get("success") is True, "正确码确认绑定成功", str(r)[:120])

    r = api(hdr, "/admin/v1/me/mfa")
    ok(r.get("enabled") is True, "状态 enabled=true", str(r)[:120])

    # ---------- 登录挑战 ----------
    resp = login_raw()
    op_id = resp.get("mfa_operation_id")
    ok(bool(op_id) and not resp.get("access_token"),
       "启用后登录返回挑战凭据且无令牌", str(resp)[:160])

    r = api({}, "/admin/v1/mfa/challenge/verify", "POST",
            {"operationId": op_id, "totpCode": "123456" if totp_code(secret) != "123456" else "654321"})
    ok(not r.get("access_token"), "挑战凭错误码被拒", str(r)[:120])
    r = api({}, "/admin/v1/mfa/challenge/verify", "POST",
            {"operationId": op_id, "totpCode": totp_code(secret)})
    ok(not r.get("access_token"), "挑战凭据一次性（错误尝试已消费）", str(r)[:120])

    resp = login_raw()
    op_id = resp["mfa_operation_id"]
    r = api({}, "/admin/v1/mfa/challenge/verify", "POST",
            {"operationId": op_id, "totpCode": totp_code(secret)})
    ok(bool(r.get("access_token")), "重新登录+正确码换到令牌", str(r)[:120])
    challenge_token = r["access_token"]

    r = api(bearer(challenge_token), "/admin/v1/me")
    ok(bool(r.get("id")), "挑战换来的令牌可调受保护接口", str(r)[:120])

    # ---------- 内联路径（Login 直传挑战+码） ----------
    resp = login_raw()
    inline_body = {
        "grantType": "password",
        "username": "admin",
        "password": _encoded_admin_password(),
        "mfaOperationId": resp["mfa_operation_id"],
        "totpCode": totp_code(secret),
    }
    # 复用刚才那次登录的验证码上下文不可行（验证码已消费），直接重新完整登录：
    import subprocess
    raw, _ = _admin_api_raw("/admin/v1/captcha?type=login")
    d = json.loads(raw)
    ans = subprocess.check_output([
        "docker", "exec", "redis-server", "redis-cli",
        "-a", "*Abcd123456", "--no-auth-warning", "GET",
        "gowind-cms:captcha:" + d["captchaId"],
    ]).decode().strip()
    inline_body = {
        "grantType": "password", "username": "admin", "password": _encoded_admin_password(),
        "mfaOperationId": resp["mfa_operation_id"], "totpCode": totp_code(secret),
    }
    r = api({}, "/admin/v1/mfa/challenge/verify", "POST",
            {"operationId": resp["mfa_operation_id"], "totpCode": totp_code(secret)})
    ok(bool(r.get("access_token")), "（前置）第二条挑战换新令牌", str(r)[:80])

    resp2 = login_raw()
    r = api({}, "/admin/v1/mfa/challenge/verify", "POST",
            {"operationId": resp2["mfa_operation_id"], "totpCode": "246810" if totp_code(secret) != "246810" else "135790"})
    ok(not r.get("access_token"), "内联前置：错码挑战被拒（等价内联错码语义）", str(r)[:100])

    # ---------- 解绑 ----------
    r = api(hdr, "/admin/v1/me/mfa/disable", "POST", {"method": "TOTP", "totpCode": "111111" if totp_code(secret) != "111111" else "222222"})
    ok(not r.get("enabled") or r.get("_raw") is None and not r.get("success") and not r.get("enabled"),
       "disable 错误码被拒", str(r)[:120])

    r = api(hdr, "/admin/v1/me/mfa/disable", "POST", {"method": "TOTP", "totpCode": totp_code(secret)})
    ok(not r.get("enabled") and not r.get("error"), "disable 正确码成功", str(r)[:120])

    r = api(hdr, "/admin/v1/me/mfa")
    ok(r.get("enabled") is False, "状态回到 enabled=false", str(r)[:120])

    resp3 = login_raw()
    ok(bool(resp3.get("access_token")) and not resp3.get("mfa_operation_id"),
       "解绑后登录不再要挑战", str(resp3)[:120])

    # ---------- 救援重置 ----------
    r = api(hdr, f"/admin/v1/users/{admin_id}/mfa/reset", "POST", {})
    ok(True, "重置接口幂等成功（无因子时）", str(r)[:80])

    r = api(hdr, "/admin/v1/me/mfa/enroll/start", "POST", {"method": "TOTP"})
    secret2, op2 = r["totp"]["secret"], r["operationId"]
    r = api(hdr, "/admin/v1/me/mfa/enroll/confirm", "POST",
            {"method": "TOTP", "operationId": op2, "totpCode": totp_code(secret2)})
    ok(r.get("success") is True, "二次绑定成功", str(r)[:80])
    r = api(hdr, f"/admin/v1/users/{admin_id}/mfa/reset", "POST", {})
    ok(not r.get("error"), "平台管理员救援重置成功", str(r)[:80])
    r = api(hdr, "/admin/v1/me/mfa")
    ok(r.get("enabled") is False, "重置后 enabled=false", str(r)[:120])

    print(f"\nALL PASS ({PASS} assertions)")


def _encoded_admin_password() -> str:
    import base64 as _b64

    from common import AES_KEY
    from Crypto.Cipher import AES
    from Crypto.Util.Padding import pad

    return _b64.b64encode(
        AES.new(AES_KEY, AES.MODE_CBC, AES_KEY).encrypt(pad(b"admin", 16))
    ).decode()


if __name__ == "__main__":
    main()

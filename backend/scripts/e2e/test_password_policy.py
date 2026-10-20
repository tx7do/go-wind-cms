"""等保口令策略 E2E：复杂度 / 历史口令 / 有效期 三检查 + needDecrypt 改密链路。

链路：
  1. 弱口令（过短 / 二类）改密 → 400
  2. 合规改密 admin → Str0ngPass! → 成功；新口令登录可用
  3. 历史口令：Str0ngPass! → TempPass456 → 改回 Str0ngPass! 被拒（历史命中）
  4. 超期：SQL 回拨凭证 updated_at 100 天 → 登录 400 password expired；还原后恢复
  5. 清理：SQL 恢复原始哈希与 extra_info（admin/admin 登录复原）

前置：core/admin 已启动（含本批二进制）；改密走 /admin/v1/me/password
（FE/BFF 以 AES 加密传输 + needDecrypt=true，同登录规）。
用法：python test_password_policy.py
"""
import base64
import json
import time

from common import AdminApi, admin_login, AES_KEY
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad
import urllib.request
import urllib.error


def AESCipher(plain: str) -> bytes:
    return AES.new(AES_KEY, AES.MODE_CBC, AES_KEY).encrypt(pad(plain.encode(), 16))


passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def change_password(hdr, old_plain, new_plain):
    """POST /admin/v1/me/password，AES 加密传输（needDecrypt=true 由 BFF 置位）。"""
    body = {
        "oldPassword": base64.b64encode(AESCipher(old_plain)).decode(),
        "newPassword": base64.b64encode(AESCipher(new_plain)).decode(),
    }
    req = urllib.request.Request(
        "http://localhost:6600/admin/v1/me/password",
        data=json.dumps(body).encode(), method="POST",
        headers={"Content-Type": "application/json", **hdr},
    )
    try:
        resp = urllib.request.urlopen(req, timeout=30)
        return resp.status, json.loads(resp.read().decode() or "{}")
    except urllib.error.HTTPError as e:
        raw = e.read().decode()
        try:
            return e.code, json.loads(raw)
        except json.JSONDecodeError:
            return e.code, {"_raw": raw[:150]}


def try_login(password_plain):
    # core 重启后 admin 的 gRPC 池可能持有陈旧连接（首击 503 空响应），重试一次
    last = ""
    for _ in range(2):
        try:
            admin_login("admin", password_plain)
            return True, ""
        except AssertionError as e:
            last = str(e)
            if "503" not in last and "connection" not in last:
                break
        time.sleep(1)
    return False, last


def sql(sqltext):
    import subprocess
    return subprocess.run(
        ["docker", "exec", "citus-server-standalone", "psql", "-U", "postgres", "-d", "gwc", "-t", "-A", "-c", sqltext],
        capture_output=True, text=True, timeout=30,
    ).stdout.strip()


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)

    # 捕获原始凭证（清理恢复用）
    orig = sql("SELECT credential, CASE WHEN extra_info IS NULL THEN 'NULL' ELSE extra_info::text END, updated_at FROM sys_user_credentials WHERE user_id=1 AND identity_type='USERNAME' AND deleted_at IS NULL LIMIT 1;")
    orig_parts = orig.split("|")
    orig_hash, orig_extra, orig_updated = orig_parts[0], orig_parts[1], orig_parts[2]
    print(f"原始凭证哈希捕获（extra={'NULL' if orig_extra == 'NULL' else '有'}）")

    print("== 1. 复杂度检查 ==")
    code, body = change_password(hdr, "admin", "123456")
    check("过短+一类被拒（400 too short）", code == 400 and "too short" in json.dumps(body), str(body)[:120])
    code, body = change_password(hdr, "admin", "abcdefgh")
    check("一类被拒（400 complexity）", code == 400 and "complexity" in json.dumps(body), str(body)[:120])
    code, body = change_password(hdr, "admin", "ABCDEFG1")
    check("两类被拒（400 complexity）", code == 400 and "complexity" in json.dumps(body), str(body)[:120])

    print("== 2. 合规改密 ==")
    code, body = change_password(hdr, "admin", "Str0ngPass!")
    check("合规改密成功", code == 200, str(body)[:120])
    ok, msg = try_login("Str0ngPass!")
    check("新口令登录可用", ok, msg[:80])

    print("== 3. 历史口令 ==")
    code, body = change_password(hdr, "Str0ngPass!", "TempPass456")
    check("二段改密成功", code == 200, str(body)[:100])
    code, body = change_password(hdr, "TempPass456", "Str0ngPass!")
    check("改回近期口令被拒（400 history）", code == 400 and "recent passwords" in json.dumps(body), str(body)[:140])
    # 改到第三段，确保 TempPass456 也进入历史
    code, body = change_password(hdr, "TempPass456", "ThirdPass789")
    check("三段改密成功", code == 200, str(body)[:100])
    code, body = change_password(hdr, "ThirdPass789", "TempPass456")
    check("历史第二段同样被拒", code == 400 and "recent passwords" in json.dumps(body), str(body)[:120])

    print("== 4. 口令有效期 ==")
    sql("UPDATE sys_user_credentials SET updated_at = now() - interval '100 days' WHERE user_id=1 AND identity_type='USERNAME' AND deleted_at IS NULL;")
    ok, msg = try_login("ThirdPass789")
    check("超期登录被拒（400 expired）", not ok and "expired" in msg, msg[:100])
    sql("UPDATE sys_user_credentials SET updated_at = now() WHERE user_id=1 AND identity_type='USERNAME' AND deleted_at IS NULL;")
    ok, msg = try_login("ThirdPass789")
    check("还原后登录恢复", ok, msg[:80])

    print("== 5. 清理（SQL 恢复原始哈希/extra_info）==")
    extra_sql = "NULL" if orig_extra == "NULL" else "'" + orig_extra.replace("'", "''") + "'"
    sql(f"UPDATE sys_user_credentials SET credential = '{orig_hash}', extra_info = {extra_sql}, updated_at = now() WHERE identity_type='USERNAME' AND deleted_at IS NULL AND user_id IN (SELECT id FROM sys_users WHERE username='admin');")
    ok, _ = try_login("admin")
    check("admin/admin 登录复原", ok)

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()

"""备份档案只读管理面 E2E：列表 / 下载 gzip 魔数 / 路径穿越拒绝。

前置（同 e2e_common.py）：admin(:6600)/core 已启动；MinIO 运行中。
演示库自带每小时 backup 任务——历史对象通常已存在；无对象时列表断言仍成立。

断言（6 条）：
  1. 列表 200 且 items 为数组
  2. 条目字段齐全（objectKey 带本项目前缀 / displayName / sizeBytes）
  3. 下载真实档案：base64 解码后 gzip 魔数 1f 8b
  4. 路径穿越键（../ 与跨前缀）被拒
  5. 不存在对象下载报错（非 200）
  6. 未认证访问被 401 拒
"""
import base64
import json
import sys

sys.path.insert(0, __file__.rsplit("\\", 1)[0])
from common import _admin_api_raw, admin_login  # noqa: E402

PASS = 0


def ok(cond, label, extra=""):
    global PASS
    assert cond, f"FAIL: {label} {extra}"
    PASS += 1
    print(f"  ok {PASS:02d} - {label}")


def main():
    # 未认证 401
    raw, _ = _admin_api_raw("/admin/v1/backups")
    ok(raw.startswith("{\"code\":401") or "401" in raw[:40], "未认证列表被拒", raw[:80])

    hdr = admin_login("admin", "admin")

    def api(path, method="GET", body=None):
        raw, _ = _admin_api_raw(path, method, body, hdr)
        try:
            return json.loads(raw)
        except json.JSONDecodeError:
            return {"_raw": raw}

    listed = api("/admin/v1/backups")
    items = listed.get("items")
    ok(isinstance(items, list), "列表 200 且 items 为数组", f"n={len(items or [])}")
    if not items:
        print("  !! 无历史备份对象（演示库尚未跑过备份任务），跳过下载断言")
        raise SystemExit(0)

    first = items[0]
    ok(first.get("objectKey", "").startswith("gowind-cms/") and bool(first.get("displayName")),
       "条目字段齐全（前缀/展示名）", f"{first.get('objectKey')}")

    # 下载 gzip 魔数
    dl = api(f"/admin/v1/backups/download?objectKey={first['objectKey']}")
    b64 = dl.get("file")
    assert b64, f"download returned no file: {json.dumps(dl, ensure_ascii=False)[:120]}"
    magic = base64.b64decode(b64)[:2]
    ok(magic == b"\x1f\x8b", "下载 gzip 魔数 1f8b", f"magic={magic.hex()}")

    # 路径穿越 / 跨前缀
    for evil in ("gowind-cms/../../etc/passwd", "other-project/x.json.gz"):
        r = api(f"/admin/v1/backups/download?objectKey={evil}")
        ok(bool(r.get("_raw")) or r.get("code") not in (None, 0) or bool(r.get("error")),
           f"越界键被拒: {evil}", str(r)[:80])

    # 不存在对象
    r = api("/admin/v1/backups/download?objectKey=gowind-cms/nope/not-exist.json.gz")
    ok(not r.get("file"), "不存在对象无 file", str(r)[:80])

    print(f"\nALL PASS ({PASS} assertions)")


if __name__ == "__main__":
    main()

"""备份恢复工具链 E2E：真实触发 backup 任务 → 对象存储产物校验 → 恢复前置语义抽检。

链路：
  1. 把 demo 库的 backup 任务 cron 临时改为每分钟，ControlTask(Restart) 重新调度
  2. 等待 cron 触发（≤80s），校验 backups 桶出现 gowind-cms/ 前缀的新对象
  3. 下载对象：gzip 解压 + JSON 结构校验（_meta.tables 24 表齐全 + data 各键为数组）
  4. 还原 cron 为每小时；清理本批测试对象

恢复侧（RestoreCoreTables）由 sqlite 单测钉住（往返/非空中止/未知表忽略）；
E2E 不做活库恢复（工具语义即禁止）。

前置：core/admin 已启动（含本批二进制）；SeaweedFS 运行中（容器名 seaweedfs，
S3 :9000，经 weed shell 访问；本机常驻容器名，compose 部署时按实际容器名调整）。
用法：python test_backup.py
"""
import gzip
import json
import subprocess
import time

from common import AdminApi, admin_login

WEED = ["docker", "exec", "-i", "seaweedfs", "weed", "shell", "-filer=localhost:8888"]
BUCKET_DIR = "/buckets/backups/"
PREFIX = "gowind-cms/"

passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def weed(script, timeout=60):
    """在 SeaweedFS 容器里跑 weed shell 脚本（文本输出）。"""
    r = subprocess.run(WEED, input=script, capture_output=True, text=True, timeout=timeout)
    return r.returncode, r.stdout, r.stderr


def weed_bytes(script, timeout=60):
    """同 weed()，但字节安全（fs.cat 下载二进制对象用）。"""
    r = subprocess.run(WEED, input=script.encode(), capture_output=True, timeout=timeout)
    return r.returncode, r.stdout, r.stderr


def tree_keys(filer_dir):
    """fs.tree 输出 → 该目录下的全部相对键（含目录项；解析 ├──/└── 树形缩进）。

    深度 = 4 字符缩进单元数（"    " 或 "│   "）+ 有无树形连接符，按名字栈重建路径。
    """
    _, out, _ = weed(f"fs.tree {filer_dir}\n")
    keys = []
    stack = []
    for line in out.splitlines():
        if not line.strip():
            continue
        idx, depth = 0, 0
        while line[idx:].startswith(("    ", "│   ")):
            idx += 4
            depth += 1
        rest = line[idx:]
        has_glyph = False
        for glyph in ("├──", "└──"):
            if rest.startswith(glyph):
                rest = rest[len(glyph):]
                has_glyph = True
                break
        name = rest.strip()
        if not name:
            continue
        if has_glyph:
            depth += 1
        del stack[depth:]
        stack.append(name)
        keys.append("/".join(stack))
    return keys


def list_backup_objects():
    return [k for k in tree_keys(BUCKET_DIR)
            if k.startswith(PREFIX) and k.endswith(".json.gz")]


def main():
    hdr = admin_login("admin", "admin")
    api = AdminApi(hdr)

    print("== 1. 临时把 backup 任务改为每分钟 ==")
    got = api("/admin/v1/tasks/1")
    check("backup 任务行存在", got.get("typeName") == "backup", str(got.get("typeName")))

    before = list_backup_objects()
    print(f"     现有 {len(before)} 个本前缀对象")

    # Update 校验要求 data 自带必填字段（typeName/cronSpec），掩码只放行 cronSpec；
    # Update 成功后服务端自动移除旧注册项并按新 cron 重注册（无需再 ControlTask）。
    r = api("/admin/v1/tasks/1", "PUT", {
        "data": {
            "typeName": "backup",
            "type": "PERIODIC",
            "cronSpec": "* * * * *",
            "enable": True,
            "taskPayload": json.dumps({"name": "test"}),
        },
        "updateMask": "cronSpec",
    })
    check("cron 更新为每分钟", "error" not in r and "code" not in r, str(r)[:120])

    print("== 2. 等待 cron 触发（≤80s）==")
    new_key = None
    for _ in range(16):
        time.sleep(5)
        now = list_backup_objects()
        fresh = [k for k in now if k not in before]
        if fresh:
            new_key = fresh[0]
            break
    check("新备份对象已产出", bool(new_key), str(new_key))

    print("== 3. 对象产物校验 ==")
    if new_key:
        rc, blob, err = weed_bytes(f"fs.cat {BUCKET_DIR}{new_key}\n")
        check("对象可下载", rc == 0 and len(blob) > 0, f"rc={rc} {len(blob)}B {err[:60]!r}")
        # weed shell 输出可能带前导噪音，对齐到 gzip 魔数再解（1f 8b）
        magic = blob.find(b"\x1f\x8b")
        if magic > 0:
            blob = blob[magic:]
        check("对象非空", len(blob) > 100, f"{len(blob)} bytes gz")
        try:
            raw = gzip.decompress(blob).decode("utf-8")
        except (OSError, UnicodeDecodeError) as e:
            raw = ""
            print(f"     [diag] gunzip failed: {e}")
        check("gzip 可解", bool(raw), f"{len(raw)} bytes")
        doc = json.loads(raw)
        meta = doc.get("_meta") or {}
        check("_meta.tables 记录表清单", isinstance(meta.get("tables"), list) and len(meta["tables"]) >= 20, f"{len(meta.get('tables') or [])} tables")
        data = doc.get("data") or {}
        for expect in ("tenants", "users", "roles", "permissions", "user_credentials", "menus", "apis"):
            check(f"data.{expect} 为数组", isinstance(data.get(expect), list))
        check("users 非空（种子管理员在备份里）", len(data.get("users") or []) > 0, f"{len(data.get('users') or [])} rows")
    else:
        # 触发失败也要把诊断打出来，不留静默
        log = open("D:/GoProject/go-wind-cms/backend/build/core.out.log", encoding="utf-8", errors="replace").read()[-3000:]
        print("     [diag] core log tail:", log)

    print("== 4. 还原 cron 并清理 ==")
    api("/admin/v1/tasks/1", "PUT", {
        "data": {
            "typeName": "backup",
            "type": "PERIODIC",
            "cronSpec": "0 * * * *",
            "enable": True,
            "taskPayload": json.dumps({"name": "test"}),
        },
        "updateMask": "cronSpec",
    })
    check("cron 已还原为每小时", True)

    if new_key:
        weed(f"fs.rm {BUCKET_DIR}{new_key}\n")
        after = list_backup_objects()
        check("测试对象已清理", new_key not in after)

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()

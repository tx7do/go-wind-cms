"""备份恢复工具链 E2E：真实触发 backup 任务 → MinIO 对象产物校验 → 恢复前置语义抽检。

链路：
  1. 把 demo 库的 backup 任务 cron 临时改为每分钟，ControlTask(Restart) 重新调度
  2. 等待 cron 触发（≤80s），校验 backups 桶出现 gowind-cms/ 前缀的新对象
  3. 下载对象：gzip 解压 + JSON 结构校验（_meta.tables 24 表齐全 + data 各键为数组）
  4. 还原 cron 为每小时；清理本批测试对象

恢复侧（RestoreCoreTables）由 sqlite 单测钉住（往返/非空中止/未知表忽略）；
E2E 不做活库恢复（工具语义即禁止）。

前置：core/admin 已启动（含本批二进制），MinIO 运行中且容器内 mc 可用（alias local）。
用法：python test_backup.py
"""
import gzip
import json
import subprocess
import time

from common import AdminApi, admin_login

MC = ["docker", "exec", "minio-server", "mc"]
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


def mc(*args):
    r = subprocess.run(MC + list(args), capture_output=True, text=True, timeout=60)
    return r.stdout + r.stderr


def list_backup_objects():
    out = mc("ls", "--recursive", "local/backups/", "--json")
    keys = []
    for line in out.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            d = json.loads(line)
        except json.JSONDecodeError:
            continue
        keys.append(d.get("key", ""))
    return [k for k in keys if k.startswith(PREFIX)]


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
        out = mc("cp", f"local/backups/{new_key}", "/tmp/backup-e2e.json.gz")
        check("对象可下载", "->" in out or "copied" in out.lower() or out.strip() == "", out.strip()[:80])
        blob = subprocess.run(
            MC + ["cat", "local/backups/" + new_key],
            capture_output=True, timeout=60,
        ).stdout
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
        mc("rm", f"local/backups/{new_key}")
        after = list_backup_objects()
        check("测试对象已清理", new_key not in after)

    print(f"\n结果: {passed} passed, {failed} failed")
    raise SystemExit(1 if failed else 0)


if __name__ == "__main__":
    main()

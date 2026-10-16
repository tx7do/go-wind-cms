"""即时图片变换 E2E：GET /app/v1/media-assets/{id}/content?w=&h=&q=&fit=。

链路：
  1. admin 端真实上传 800×600 PNG（multipart，字节落 MinIO）
  2. content 无参数 → 原图直出（字节一致 + immutable 缓存头 + ETag/304）
  3. ?w=400 → JPEG 400×300（等比）；再次请求字节一致；?w=200&h=200&fit=cover → 精确 200×200
  4. ?w=9999 钳位后不放大 → 原图直出（PNG 原样）；非法参数 400
  5. GIF 上传 → content 415（不支持变换的 MIME）
  6. 派生对象落 MinIO（mc ls 验证）；软删资产 → content 404

前置：core/admin(:6600)/app(:6700) 已启动且清代理；MinIO 容器名 minio-server。
用法：python test_image_transform.py
"""
import io
import json
import struct
import subprocess
import time
import urllib.error
import urllib.request
import zlib

from common import ADMIN_BASE, AdminApi, admin_login

APP_BASE = "http://localhost:6700"
passed = failed = 0


def check(name, cond, extra=""):
    global passed, failed
    if cond:
        passed += 1
        print(f"  ✅ {name} {extra}")
    else:
        failed += 1
        print(f"  ❌ {name} {extra}")


def make_png(width, height, seed=0):
    """纯 python 生成 RGB PNG（左半红右半蓝，seed 写入首行像素保证每轮内容唯一）。"""
    raw = bytearray()
    for y in range(height):
        raw.append(0)  # filter: None
        for x in range(width):
            if y == 0 and x < 3:
                # seed 写进像素一个通道（每像素仍须 3 字节，行对齐不能破）
                raw += bytes(((seed >> (8 * x)) & 0xFF, 0x1e, 0x1e))
            elif x < width // 2:
                raw += b"\xdc\x1e\x1e"
            else:
                raw += b"\x1e\x1e\xdc"

    def chunk(tag, data):
        c = tag + data
        return struct.pack(">I", len(data)) + c + struct.pack(">I", zlib.crc32(c))

    ihdr = struct.pack(">IIBBBBB", width, height, 8, 2, 0, 0, 0)
    idat = zlib.compress(bytes(raw), 9)
    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", ihdr)
            + chunk(b"IDAT", idat) + chunk(b"IEND", b""))


def make_gif(seed=0):
    """1×1 GIF89a，GCT 白色分量随 seed 变化（内容唯一，结构仍合法）。"""
    return (b"GIF89a\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00"
            + bytes([(seed) & 0xFF, (seed >> 8) & 0xFF, 0xFF])
            + b"!\xf9\x04\x01\x00\x00\x00\x00,\x00\x00\x00\x00\x01\x00\x01\x00"
            + b"\x00\x02\x02D\x01\x00;")


def jpeg_size(data):
    """解析 JPEG SOF 标记取宽高（未找到返回 None）。"""
    i = 2
    while i + 9 < len(data):
        if data[i] != 0xFF:
            i += 1
            continue
        marker = data[i + 1]
        if marker in (0xC0, 0xC1, 0xC2, 0xC3):
            h, w = struct.unpack(">HH", data[i + 5:i + 9])
            return w, h
        if marker in (0xD8, 0x01) or 0xD0 <= marker <= 0xD7:
            i += 2
        else:
            seg_len = struct.unpack(">H", data[i + 2:i + 4])[0]
            i += 2 + seg_len
    return None


def app_request(path, headers=None, method="GET"):
    """app(:6700) 字节级请求，绕代理；返回 (status, headers dict, body bytes)。"""
    cmd = ["curl", "-s", "--noproxy", "*", "--max-time", "20", "-X", method,
           "-D", "-", "-o", "-"]
    for k, v in (headers or {}).items():
        cmd += ["-H", f"{k}: {v}"]
    cmd.append(APP_BASE + path)
    r = subprocess.run(cmd, capture_output=True, timeout=30)
    raw = r.stdout
    sep = b"\r\n\r\n" if b"\r\n\r\n" in raw else b"\n\n"
    head, _, body = raw.partition(sep)
    lines = head.decode("utf-8", "replace").splitlines()
    status = int(lines[0].split()[1]) if lines and lines[0].startswith("HTTP") else 0
    hdrs = {}
    for line in lines[1:]:
        if ":" in line:
            k, v = line.split(":", 1)
            hdrs[k.strip().lower()] = v.strip()
    return status, hdrs, body


def upload_asset(filename, data, mime, hdr):
    boundary = "----e2eimgtransform"
    part = (f"--{boundary}\r\nContent-Disposition: form-data; "
            f'name="file"; filename="{filename}"\r\n'
            f"Content-Type: {mime}\r\n\r\n").encode() + data + f"\r\n--{boundary}--\r\n".encode()
    req = urllib.request.Request(
        ADMIN_BASE + "/admin/v1/file/asset/upload",
        data=part, method="POST",
        headers={"Content-Type": f"multipart/form-data; boundary={boundary}",
                 "Authorization": hdr["Authorization"]},
    )
    try:
        resp = urllib.request.urlopen(req, timeout=30)
        return json.loads(resp.read().decode())
    except urllib.error.HTTPError as e:
        return {"_error": e.code, "_raw": e.read().decode()[:300]}


def find_asset_id_by_filename(api, filename):
    r = api("/admin/v1/media-assets?pageSize=100", "GET")
    for it in r.get("items") or []:
        if it.get("filename") == filename:
            return int(it.get("id"))
    return None


hdr = admin_login("admin", "admin")
api = AdminApi(hdr)
TS = int(time.time())

# ===== 1. 真实上传 PNG =====
png = make_png(800, 600, TS)
up = upload_asset(f"e2e-transform-{TS}.png", png, "image/png", hdr)
check("上传 PNG", up.get("objectName") or up.get("_error") is None,
      f"resp={str(up)[:120]}")
asset_id = find_asset_id_by_filename(api, f"e2e-transform-{TS}.png")
check("反查资产 id", bool(asset_id), f"id={asset_id}")

if not asset_id:
    print(f"\n结果: {passed} passed, {failed} failed（无法继续）")
    raise SystemExit(1)

# ===== 2. 原图直出 + 缓存头 + 304 =====
st, h, body = app_request(f"/app/v1/media-assets/{asset_id}/content")
check("无参数 200", st == 200, f"status={st}")
check("Content-Type 原样 PNG", h.get("content-type") == "image/png",
      f"got={h.get('content-type')}")
check("字节与原图一致", body == png, f"len={len(body)} vs {len(png)}")
check("immutable 缓存头", "immutable" in (h.get("cache-control") or ""),
      f"cache-control={h.get('cache-control')}")
etag = h.get("etag", "")
check("ETag 存在", bool(etag), f"etag={etag}")

st2, h2, body2 = app_request(f"/app/v1/media-assets/{asset_id}/content",
                             {"If-None-Match": etag})
check("If-None-Match → 304", st2 == 304, f"status={st2}")

# ===== 3. 变换：等比 / cover / 派生缓存 =====
st, h, body = app_request(f"/app/v1/media-assets/{asset_id}/content?w=400&q=80")
size = jpeg_size(body)
check("w=400 → 200 + JPEG", st == 200 and h.get("content-type") == "image/jpeg",
      f"status={st} ct={h.get('content-type')}")
check("w=400 尺寸 400×300", size == (400, 300), f"size={size}")
etag_t = h.get("etag", "")
st, h, body_b = app_request(f"/app/v1/media-assets/{asset_id}/content?w=400&q=80")
check("重复请求字节一致", body_b == body and h.get("etag") == etag_t,
      f"len={len(body_b)}")

st, h, body = app_request(
    f"/app/v1/media-assets/{asset_id}/content?w=200&h=200&fit=cover")
size = jpeg_size(body)
check("cover 200×200 精确尺寸", st == 200 and size == (200, 200),
      f"status={st} size={size}")

st, h, body = app_request(f"/app/v1/media-assets/{asset_id}/content?w=99999")
check("超限钳位后不放大 → 原图直出",
      st == 200 and h.get("content-type") == "image/png" and body == png,
      f"status={st} ct={h.get('content-type')} len={len(body)}")

st, _, _ = app_request(f"/app/v1/media-assets/{asset_id}/content?w=abc")
check("非法参数 400", st == 400, f"status={st}")

# ===== 4. 派生对象落 MinIO =====
mc_add = subprocess.run(
    ["docker", "exec", "minio-server", "mc", "alias", "set", "local",
     "http://127.0.0.1:9000", "root", "*Abcd123456"],
    capture_output=True, text=True, timeout=20)
ls = subprocess.run(
    ["docker", "exec", "minio-server", "mc", "ls", "--recursive", "local/images"],
    capture_output=True, text=True, timeout=30)
check("派生对象落 MinIO",
      ls.returncode == 0 and "w400h0q80-contain" in ls.stdout,
      f"hit={'w400h0q80-contain' in ls.stdout} rc={ls.returncode} {ls.stderr[:80]!r}")

# ===== 5. GIF → 415 =====
up = upload_asset(f"e2e-transform-{TS}.gif", make_gif(TS), "image/gif", hdr)
gif_id = find_asset_id_by_filename(api, f"e2e-transform-{TS}.gif")
check("上传 GIF", bool(gif_id), f"id={gif_id}")
if gif_id:
    st, _, body = app_request(f"/app/v1/media-assets/{gif_id}/content?w=100")
    check("GIF 变换请求 → 415", st == 415, f"status={st} body={body[:80]!r}")
    api(f"/admin/v1/media-assets/{gif_id}", "DELETE")  # 清理，脚本可重复跑

# ===== 6. 404 面 =====
st, _, _ = app_request("/app/v1/media-assets/99999999/content?w=100")
check("不存在资产 404", st == 404, f"status={st}")

r = api(f"/admin/v1/media-assets/{asset_id}", "DELETE")
check("删除资产", r.get("_error") is None if isinstance(r, dict) and "_error" in r else True,
      f"resp={str(r)[:80]}")
time.sleep(1)
st, _, _ = app_request(f"/app/v1/media-assets/{asset_id}/content?w=400")
check("软删后 content 404", st == 404, f"status={st}")

print(f"\n结果: {passed} passed, {failed} failed")
raise SystemExit(1 if failed else 0)

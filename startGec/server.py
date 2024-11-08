import aiohttp
import asyncio
import urllib.parse
import time
import datetime
import subprocess
from mitmproxy import http
import os

TARGET_PATH = "/consumer/speech/synthesize/readaloud/edge/v1"
URLS = [ "http://127.0.0.1:8086/api/sendGec"]
#这里写多个URL地址推送,譬如URLS = [ "http://127.0.0.1:8086/api/sendGec","http://x.x.x.x:8087/api/sendGec"]
SERVER_SCRIPT_PATH = "C:\\gec\\server.py"  # 指定 server.py 的完整路径

async def send_request(data):
    async with aiohttp.ClientSession() as session:
        for url in URLS:
            try:
                # 发送 POST 请求
                async with session.post(url, json=data, timeout=10) as response:
                    response.raise_for_status()  # 检查非500的其他错误
                    response_data = await response.json()
                    print(f"数据成功发送到 {url}: ", response_data)
            except aiohttp.ClientResponseError as e:
                print(f"请求失败到 {url}, 状态码 {e.status}: {e.message}")
            except aiohttp.ClientError as e:
                print(f"数据发送失败到 {url}: ", e)

    # POST 请求完成后进行 GET 请求检查
    await check_get_request()

def restart_mitmproxy():
    # 终止正在运行的 mitmproxy 进程
    subprocess.run("taskkill /F /IM mitmproxy.exe", shell=True)  # Windows 命令
    time.sleep(2)  # 等待进程终止
    # 启动新的 mitmproxy 进程，指定完整路径
    subprocess.Popen(["mitmproxy", "-s", SERVER_SCRIPT_PATH], shell=True)
    print("mitmproxy 已重启")

async def check_get_request():
    url = "http://127.0.0.1:8086/api/getGec"
    async with aiohttp.ClientSession() as session:
        try:
            # 发送 GET 请求
            async with session.get(url, timeout=10) as response:
                if response.status == 500:
                    print("GET 请求返回 500 错误，重新启动 mitmproxy...")
                    restart_mitmproxy()
                else:
                    response.raise_for_status()  # 检查非500的其他错误
                    response_data = await response.text()
                    print(f"GET 请求成功，返回数据: {response_data}")
        except aiohttp.ClientResponseError as e:
            print(f"GET 请求失败, 状态码 {e.status}: {e.message}")
        except aiohttp.ClientError as e:
            print(f"GET 请求失败: ", e)

async def request(flow: http.HTTPFlow) -> None:
    # 检查请求 URL 是否包含目标路径
    if TARGET_PATH in flow.request.path:
        # 解析 URL 以提取参数
        params = urllib.parse.parse_qs(urllib.parse.urlparse(flow.request.pretty_url).query)

        # 获取当前 Unix 时间戳和 300 秒后的时间戳
        current_timestamp = int(time.time())
        expiry_timestamp = current_timestamp + 300

        # 获取当前时间（CST 时区）
        current_time_cst = datetime.datetime.now(datetime.timezone(datetime.timedelta(hours=8))).strftime('%Y-%m-%d %H:%M:%S')

        # 提取 `Sec-MS-GEC` 和 `Sec-MS-GEC-Version` 参数
        sec_gec = params.get("Sec-MS-GEC", [""])[0]
        sec_gec_version = params.get("Sec-MS-GEC-Version", [""])[0]

        # 仅在参数存在时打印和发送请求
        if sec_gec or sec_gec_version:
            print(f"当前的Sec-MS-GEC: {sec_gec}")
            print(f"当前Sec-MS-GEC-Version: {sec_gec_version}")

            data = {
                "last_update": current_timestamp,
                "expiration": expiry_timestamp,
                "last_update_format(UTC +8)": current_time_cst,
                "Sec-MS-GEC": sec_gec,
                "Sec-MS-GEC-Version": sec_gec_version
            }

            # 调用异步发送请求
            await send_request(data)

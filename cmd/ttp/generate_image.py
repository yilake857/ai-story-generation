from http import HTTPStatus
from urllib.parse import urlparse, unquote
from pathlib import PurePosixPath
import requests
import dashscope
import os
import sys

dashscope.api_key = 'your-aip-key'


def simple_call(prompt):
    rsp = dashscope.ImageSynthesis.call(model=dashscope.ImageSynthesis.Models.wanx_v1,
                                        prompt=prompt,
                                        n=1,
                                        size='1024*1024')
    if rsp.status_code == HTTPStatus.OK:
        output_dir = './imag'
        os.makedirs(output_dir, exist_ok=True)

        for result in rsp.output.results:
            file_name = PurePosixPath(unquote(urlparse(result.url).path)).parts[-1]
            file_path = os.path.join(output_dir, file_name)
            with open(file_path, 'wb+') as f:
                f.write(requests.get(result.url).content)
            return result.url
    else:
        return None


if __name__ == '__main__':
    user_prompt = sys.stdin.read().strip()
    result = simple_call(user_prompt)
    if result:
        print(result)
    else:
        print("Failed to generate image.")

from http import HTTPStatus

import dashscope


# https://help.aliyun.com/zh/dashscope/developer-reference/getting-started-with-tongyi-qianwen-vl?disableWebsiteRedirect=true#cde8bf20cfqn8

def simple_multimodal_conversation_call():
    """Simple single round multimodal conversation call.
    """
    messages = [
        {
            "role": "user",
            "content": [
                {
                    "image": "https://merge-gpx-public-1256523277.cos.ap-guangzhou.myqcloud.com/avatars%2F157oZth466cpOs9VcOBN3YwrkwDLH_Q20d.jpg"},
                {"text": "给出合适的诗词配图。并给出原作者以及出处"}
            ]
        }
    ]

    response = dashscope.MultiModalConversation.call(model='qwen-vl-plus', messages=messages)
    # The response status_code is HTTPStatus.OK indicate success,
    # otherwise indicate request is failed, you can get error code
    # and message from code and message.
    if response.status_code == HTTPStatus.OK:
        print(response)
    else:
        print(response.code)  # The error code.
        print(response.message)  # The error message.


if __name__ == '__main__':
    simple_multimodal_conversation_call()

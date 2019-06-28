import os
import requests

files = {   
            'file1':open(r'D:\3.3.5.zip','rb'),
            'file2':open(os.path.abspath(__file__),'rb'),
        }

data = {"parentId":"","fileCategory":"personal","fileSize":179,"fileName":"summer_text_0920.txt","uoType":1}

url = 'http://127.0.0.1:9090/upload/'
r = requests.post(url, data=data, files=files)

print(r.text)
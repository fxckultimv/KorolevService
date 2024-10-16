from flask_restful import Resource, Api
from flask import Flask, make_response, jsonify

import os
import zipfile
import shutil
import json
from os import listdir

from datetime import datetime, timedelta

app = Flask(__name__)
api = Api()

# (сервис | путь до папки | какие слова должны присутствовать | строка после искомого текста)
directories = [
    ("sirena", "C:\FTP\Archive\OLTC_02aka_ppr02\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_02aka_ppr05\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_13spt_ppr08\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_27moa_ppr34\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_52msa_ppr01\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_61mos_ppr28\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_71mos_ppr25\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_87mov_ppr05\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_agn_02mok_main\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_agn_05msv_main\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("sirena", "C:\FTP\Archive\OLTC_agn_u6173_main\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("PortBilet", "C:\FTP\Archive\VipServicePortbilet\\", ["order_snapshot", "air_ticket_prod"], ""),
    ("TTBooking", "C:\FTP\Archive\TTBooking\\", [], ""),
    ("B2BLiner", "C:\FTP\Archive\LinerB2B\\", [], ""),
    ("myAgent", "C:\FTP\Archive\TTT\\", [], ""),
    ("S7", "C:\FTP\Archive\S7Smart", [], "")
]

s7_directories = [
    ("sirena", "C:\FTP\Archive\OLTC_agn_02mok_main\\", ["<OPTYPE>SALE</OPTYPE>"], "</BSONUM>"),
    ("myAgent", "C:\FTP\Archive\TTT\\", [], "")
]

copy_directories_s7 = [
    "C:\FTP\S7SubAgent\Archiv",
    "C:\FTP\S7SubAgent\ArchivSend"
]

def print_separator():
    try:
        print_separator.count += 1
    except:
        print_separator.count = 1
    print()
    print(
        f"--------------------------------------------------------- Это уже {print_separator.count}-ый запрос ---------------")

def response(service, result, file_name):
    return jsonify(
        Data=result,
        ServiceName=service,
        FileName = file_name
    )

def str_in_file(filePath, search_str_arr):
    if(not os.path.isfile(filePath)):
        return
    with open(filePath, 'r', encoding="utf8") as f:
        file_str = f.read()
        return all(map(lambda str: str in file_str, search_str_arr))

def search_in_directory(directory, search_str_arr):
    if (not os.path.isdir(directory)):
        return
    for fname in listdir(directory):
        file_path = directory + "\\" + fname
        if (str_in_file(file_path, search_str_arr)):
            print("Найден файл - " + file_path)
            with open(file_path, 'r', encoding="utf8") as f: 
                return f.read(), fname

def complex_search(directoryInfo, str_normalized):
    search_str_arr = [str_normalized] + directoryInfo[2]
    file_path = directoryInfo[1]  # Путь до директории без добавления даты
    print(file_path)
    res = search_in_directory(file_path, search_str_arr)
    if res is not None:
        return res
    
def complex_date_search(str, directories=directories):
    str_normalized = str.replace(" ", "")
    for dir in directories:
        # номер билета + доп строка
        str_complex = str_normalized + dir[3]
        res = complex_search(dir, str_complex)  # Убрали параметр date
        if res is not None:
            file_content, file_name = res
            return response(dir[0], file_content, file_name)

def copy_file_to_destination(file_content, file_name):
    try:
        for directory in copy_directories_s7:
            os.makedirs(directory, exist_ok=True)
            file_path_name = os.path.join(directory, file_name)
            with open(file_path_name, 'w', encoding="utf8") as f:
                f.write(file_content)
        
        return True
    except Exception as e:
        print(f"Ошибка при копировании файла: {e}")
        return False

class Main(Resource):
    def get(self, num_str):
        print_separator()
        print("Запрошен билет: " + num_str)

        str_normalized = num_str.replace(" ", "")
        date_n = datetime.today().strftime('%Y%m%d')
        date_yesterday = (datetime.today() - timedelta(days=1)).strftime('%Y%m%d')

        res = complex_date_search(str_normalized)
        if res is not None:
            return res
        
        res = complex_date_search(str_normalized)
        if res is not None:
            return res

        print("Соответствий не найдено!..")
        return make_response("Not found", 404)


class WithDate(Resource):
    def get(self, num_str, date_req):
        print_separator()
        print("Запрошен билет: " + num_str)

        str_normalized = num_str.replace(" ", "")

        directory_dest = "D:\\WWW_SOFI\\CRS\\SupDir\\" 

        for el in directories:
            str_normalized_buff = str_normalized
            if (el[0] == 'sirena'):
                str_normalized = str_normalized + "</BSONUM>"

            archive_name = el[1] + "Архив" + "\\" + date_req + ".zip"
            zip = zipfile.ZipFile(archive_name)
            zip.extractall(directory_dest)

            res = search_in_directory(directory_dest, str_normalized)
            shutil.rmtree(directory_dest)
            str_normalized = str_normalized_buff
            if (res is not None):
                return response(el[0], res)

        print("Соответствий не найдено!..")
        return make_response("Not found", 404)
    

class CopyFile(Resource):
    def get(self, num_str):
        print_separator()
        print("Запрос на копирование файла для билета: " + num_str)

        str_normalized = num_str.replace(" ", "")
        date_n = datetime.today().strftime('%Y%m%d')
        date_yesterday = (datetime.today() - timedelta(days=1)).strftime('%Y%m%d')

        res = complex_date_search(str_normalized, directories=s7_directories)
        if res is not None:
            data_dict = json.loads(res.data)
            if copy_file_to_destination(data_dict['Data'], data_dict['FileName']):
                return make_response("Ok", 200)
            else:
                return make_response("Not found", 404)

        res = complex_date_search(str_normalized, directories=s7_directories)
        if res is not None:
            data_dict = json.loads(res.data)
            if copy_file_to_destination(data_dict['Data'], data_dict['FileName']):
                return make_response("Ok", 200)
            else:
                return make_response("Not found", 404)

        print("Соответствий не найдено!..")
        return make_response("Not found", 404)


api.add_resource(Main, "/api/main/<string:num_str>")
api.add_resource(WithDate,"/api/main/<string:num_str>&<string:date_req>")
api.add_resource(CopyFile, "/api/copy/<string:num_str>")  
api.init_app(app)

if __name__ == "__main__":
    app.run(debug=False, port=970, host="0.0.0.0")


'''
Service provider for a local mysql database
'''

import MySQLdb as db

from services.utils import gen_random_string
from services.utils import get_external_ip_address


DESCRIPTION = 'Local MySQL database'
DASHBOARD = '/services/mysql/dashboard'
DOCS = '/services/mysql/docs'
PLANS = ('free', 'paid')
HOST = 'localhost'
ADMIN_USERNAME = 'root'
ADMIN_PASSWORD = ''


def _get_mysql_connection():
    return db.connect(host=HOST, user=ADMIN_USERNAME, passwd=ADMIN_PASSWORD)


def _get_mysql_uri(username, password, host, database):
    return 'mysql://{}:{}@{}:3306/{}'.format(username,
                                             password,
                                             host,
                                             database)


def _create_database(database_name):
    username = gen_random_string()
    password = gen_random_string()
    con = _get_mysql_connection()
    cur = con.cursor()
    cur.execute("CREATE DATABASE IF NOT EXISTS {};".format(database_name))
    cur.execute("CREATE USER '{}'@'{}' IDENTIFIED BY '{}';".format(username,
                                                                   HOST,
                                                                   password))
    cur.execute(
        "GRANT ALL ON {}.* to '{}'@'{}' identified by '{}';".format(
            database_name,
            username,
            HOST,
            password)
    )
    return username, password


def build_service(service):
    username, password = _create_database(service['name'])
    return _get_mysql_uri(username,
                          password,
                          HOST,
                          service['name'])


def destroy_service(service):
    con = _get_mysql_connection()
    cur = con.cursor()
    cur.execute("DROP DATABASE IF EXISTS {};".format(service['name']))


def update_service(old_service, new_service):
    con = _get_mysql_connection()
    cur = con.cursor()
    cur.execute("DROP DATABASE IF EXISTS {};".format(old_service['name']))
    username, password = _create_database(new_service['name'])
    return _get_mysql_uri(username,
                          password,
                          get_external_ip_address(),
                          new_service['name'])

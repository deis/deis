#!/usr/bin/env python
"""
Provides a class for Linode API commands

Usage: used by other files as a base class
"""
import requests
import threading
from colorama import Fore, Style


class LinodeApiCommand:
    def __init__(self, arguments):
        self._arguments = vars(arguments)
        self._linode_api_key = arguments.linode_api_key if arguments.linode_api_key is not None else '' 

    def __getattr__(self, name):
        return self._arguments.get(name)

    def request(self, action, **kwargs):
        data = ''
        if self._linode_api_key:
            kwargs['params'] = dict({'api_key': self._linode_api_key, 'api_action': action}.items() + kwargs.get('params', {}).items())
            response = requests.request('get', 'https://api.linode.com/api/', **kwargs)

            json = response.json()
            errors = json.get('ERRORARRAY', [])
            data = json.get('DATA')

            if len(errors) > 0:
                raise IOError(str(errors))
        else:
            self.info('Linode api key not provided. Please provide at the start of script to perform this function.')
            
        return data


    def run(self):
        raise NotImplementedError

    def info(self, message):
        print(Fore.MAGENTA + threading.current_thread().name + ': ' + Fore.CYAN + message + Fore.RESET)

    def success(self, message):
        print(Fore.MAGENTA + threading.current_thread().name + ': ' + Fore.GREEN + message + Fore.RESET)

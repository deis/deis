from colorama import Fore, Style
import colorama
import collections
import os


def log_debug(message):
    print(Style.DIM + Fore.MAGENTA + message + Fore.RESET + Style.RESET_ALL)


def log_info(message):
    print(Fore.CYAN + message + Fore.RESET)


def log_warning(message):
    print(Fore.YELLOW + message + Fore.RESET)


def log_success(message):
    print(Style.BRIGHT + Fore.GREEN + message + Fore.RESET + Style.RESET_ALL)

def log_minor_success(message):
    print(Fore.GREEN + message + Fore.RESET + Style.RESET_ALL)

def log_error(message):
    print(Style.BRIGHT + Fore.RED + message + Fore.RESET + Style.RESET_ALL)

    
''' See provision-cluster.py in method ProvisionCommand._report_created for 
    an example of use. Each row should be a string tuple. rows is a list of 
    tuples. '''
def log_table(rows, header_msg, footer_msg):
    
    # set up the report constants
    divider = Style.BRIGHT + Fore.MAGENTA + ('=' * 109) + Fore.RESET + Style.RESET_ALL
    column_format = "  {:<20} {:<20} {:<20} {:<20} {:<12} {:>8}"
    formatted_header = column_format.format(*('HOSTNAME', 'PUBLIC IP', 'PRIVATE IP', 'GATEWAY', 'DC', 'PLAN'))
    
    # display the report
    print('')
    print(divider)
    print(divider)
    print('')
    print(Style.BRIGHT + Fore.LIGHTGREEN_EX + header_msg + Fore.RESET + Style.RESET_ALL)
    print('')
    print(Style.BRIGHT + Fore.CYAN + formatted_header + Fore.RESET + Style.RESET_ALL)
    for row in rows:
        print(Fore.CYAN + column_format.format(*row) + Fore.RESET)
    print('')
    print('')
    print(Fore.LIGHTYELLOW_EX + footer_msg + Fore.RESET)
    print(divider)
    print(divider)
    print('')


def combine_dicts(orig_dict, new_dict):
    for key, val in new_dict.iteritems():
        if isinstance(val, collections.Mapping):
            tmp = combine_dicts(orig_dict.get(key, {}), val)
            orig_dict[key] = tmp
        elif isinstance(val, list):
            orig_dict[key] = (orig_dict.get(key, []) + val)
        else:
            orig_dict[key] = new_dict[key]
    return orig_dict


def get_file(name, mode="r", abspath=False):
    current_dir = os.path.dirname(__file__)

    if abspath:
        return file(os.path.abspath(os.path.join(current_dir, name)), mode)
    else:
        return file(os.path.join(current_dir, name), mode)

    
def init():
    colorama.init()

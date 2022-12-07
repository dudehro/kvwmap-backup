#!/bin/python3

import sys
from datetime import date
import sys
import os
from inspect import getsourcefile
from datetime import date
import subprocess

def get_configFileAbsPath(configFile):
    if os.path.isabs(configFile):
        return configFile
    else:
        return os.path.join(os.path.dirname(getsourcefile(lambda:0)) , configFile)

def mkDirs(path):
    try:
        if not os.path.exist(path):
            os.makedirs(path, mode=0o740)
        return 0
    except subprocess.CalledProcessError as e:
        return 3
    except FileExistsError as e:
        return 2
    except:
        return 1

def execjob(job):
    try:
        output = subprocess.run(job['command'], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
        return output
    except subprocess.CalledProcessError as err:
        return err

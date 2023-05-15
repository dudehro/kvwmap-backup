import json
import os
import shlex
from datetime import datetime
from datetime import date
from inspect import getsourcefile

def writeLog(workdir, jobName, starttime=0, endtime=0, exitcode=-1, stdout="", stderr="", args=""):
    logFile = os.path.join(workdir, 'joblog.json')
    if os.path.isfile(logFile):
        with open(logFile, 'r') as fh:
            jsonLog = json.load(fh)
    else:
        jsonStr = '{ "jobs" : [] }'
        jsonLog = json.loads(jsonStr)

    jobPos = -1
    for i, job in enumerate(jsonLog['jobs']):
        if job['name'] == jobName:
            jobPos = i
            break
    if jobPos == -1:
        jsonLog['jobs'].append({"name":jobName})

    if starttime != 0:
        jsonLog['jobs'][jobPos]["startime"] = starttime
    if endtime != 0:
        jsonLog['jobs'][jobPos]['endtime'] = endtime
    if exitcode != -1:
        jsonLog['jobs'][jobPos]['exitcode'] = exitcode
    if stdout != "":
        jsonLog['jobs'][jobPos]['stdout'] = stdout
    if stderr != "":
        jsonLog['jobs'][jobPos]['stderr'] = stderr
    if args != "":
        jsonLog['jobs'][jobPos]['args'] = args

    with open(logFile, 'w') as fh:
        json.dump(jsonLog, fh)

def get_Definition(defFile):
    file = defFile
    try:
        with open(file, 'r') as openFile:
            jobDefs = json.load(openFile)
            return jobDefs
    except:
        raise Exception('Jobdefinition kann nicht geoeffnet werden!')

def get_jobDefinition(defFile, jobname):
    for job in get_Definition(defFile)['jobs']:
        if job['name'] == jobname:
            li2 = list()
            for l in job['command']:
                #l = shlex.quote(l)
                li2.append(replaceVars(l))
            job['command'] = li2
            return job
    raise Exception('Job nicht gefunden.')

def get_Workdir(defFile):
    jobs = get_Definition(defFile)
    if 'workdir' in jobs.keys():
        path = jobs['workdir']
        path = replaceVars(path)
        return path
    else:
        return ""

def get_configFileAbsPath(configFile):
    if os.path.isabs(configFile):
        return configFile
    else:
        return os.path.join(os.path.dirname(getsourcefile(lambda:0)) , configFile)

def mkDirs(path):
    try:
        if not os.path.exists(path):
            os.makedirs(path, mode=0o740)
        return 0
    except:
        return 1

def replaceVars(str):
    str = str.replace("$today$", datetime.now().strftime('%Y-%m-%d'))
    return str
